/*
 *
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	taskmodel "github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/task"
	"github.com/cloudwego/cwgo/platform/server/shared/log"
	"github.com/cloudwego/cwgo/platform/server/shared/task"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (m *Manager) tryPromoteApiToMasterSrv() {
	ctx := context.Background()
	for {
		log.Debug("trying to promote to master...")
		if ok, err := m.trySetMaster(ctx); err != nil {
			log.Error("failed to promote to master", zap.Error(err))
		} else if ok {
			log.Info("the [api] service is master now")
			m.isMasterApi = true

			// Regularly update the key TTL
			maintainMasterWithRetry(ctx, m)
		}

		time.Sleep(10 * time.Second)
	}
}

// trySetMaster is used to make this api service to master api service
func (m *Manager) trySetMaster(ctx context.Context) (bool, error) {
	ok, err := m.rdb.SetNX(ctx, consts.RdbKeyApiMaster, m.apiID, 30*time.Second).Result()
	if err != nil {
		log.Error(consts.ErrMsgDatabaseRedisSetNX,
			zap.Error(err),
			zap.String("key", consts.RdbKeyApiMaster),
			zap.String("value", m.apiID),
		)
	}
	return ok, err
}

// maintainMasterWithRetry is used to keep this api service as the master service
func maintainMasterWithRetry(ctx context.Context, manager *Manager) {
	setApiMasterKeyExpiration := func() error {
		err := manager.rdb.Expire(ctx, consts.RdbKeyApiMaster, 30*time.Second).Err()
		if err != nil {
			log.Error(consts.ErrMsgDatabaseRedisExpire,
				zap.Error(err),
				zap.String("key", consts.RdbKeyApiMaster),
			)
			return err
		}
		return nil
	}

	for {
		opts := []retry.Option{
			retry.LastErrorOnly(true),
			retry.Attempts(6),
			retry.Delay(3 * time.Second),
		}
		if err := retry.Do(setApiMasterKeyExpiration, opts...); err != nil {
			log.Error("trying to maintain the master identity failed", zap.Error(err))
			manager.isMasterApi = false
			return
		}
		time.Sleep(10 * time.Second)
	}
}

func (m *Manager) watchTaskUpdate() {
	ctx := context.Background()

	subscribe := m.rdb.Subscribe(ctx, consts.RdbKeyTask)
	defer func() {
		err := subscribe.Close()
		if err != nil {
			log.Error(consts.ErrMsgDatabaseRedisPubSubClose,
				zap.Error(err),
				zap.String("key", consts.RdbKeyTask),
			)
		}
	}()

	channel := subscribe.Channel()
	for msg := range channel {
		var taskMessage task.Message

		// get task message
		err := sonic.UnmarshalString(msg.Payload, &taskMessage)
		if err != nil {
			log.Error(consts.ErrMsgCommonJsonUnmarshal,
				zap.Error(err),
				zap.String("value", msg.Payload),
			)
		}

		log.Debug("get message from redis task channel",
			zap.Reflect("message", taskMessage),
		)

		switch taskMessage.Command {
		case task.AddTask:
			log.Debug("add task from redis task channel",
				zap.Reflect("task", taskMessage.Task),
			)
			_ = m.AddTask(&taskMessage.Task, false)

		case task.DeleteTask:
			log.Debug("del task from redis task channel",
				zap.String("task_id", taskMessage.Task.ID),
			)
			switch taskMessage.Task.Type {
			case consts.Sync:
				_ = m.DeleteTask(consts.Sync, taskMessage.Task.IdlID, false)
			}

		default:
			log.Warn("invalid task command type",
				zap.String("type", taskMessage.Command),
			)
		}
	}
}

func (m *Manager) syncTaskFromDB() {
	var page int32 = 1
	activeModel := model.IDL{
		Status: consts.IdlStatusNumActive,
	}
	log.Info("acquiring all sync task from database")
	for {
		idlModelsFromDB, total, err := m.daoManager.Idl.GetIDLList(context.Background(), activeModel, page, consts.DefaultPageSize, consts.OrderNumDec, "update_time")
		if err != nil {
			log.Error("get idl list failed", zap.Error(err))
			continue
		}

		tasksInMemory := m.dispatcher.GetAllTasks()
		syncIdlIDMap := make(map[int64]struct{})
		for _, t := range tasksInMemory {
			if t.Type == consts.Sync {
				syncIdlIDMap[t.IdlID] = struct{}{}
			}
		}

		for _, idlModel := range idlModelsFromDB {
			if _, ok := syncIdlIDMap[idlModel.Id]; !ok {
				err = m.AddTask(
					task.NewTask(consts.Sync, m.syncIdlInterval.String(), idlModel.Id), false,
				)
				if err != nil {
					log.Error("fail to add sync task", zap.Error(err))
				}
			}
			if int64(page)*consts.DefaultPageSize >= total {
				break
			}
			page++
		}
		time.Sleep(m.syncIdlInterval)
		log.Info("acquire all sync task complete")
	}
}

func (m *Manager) AddTask(t *model.Task, isNotify bool) error {
	switch t.Type {
	case consts.Sync:
		if m.syncIdlInterval == 0 {
			return nil
		}
		t.ScheduleTime = m.syncIdlInterval.String()
	}
	err := m.dispatcher.AddTask(t)
	if err != nil {
		return fmt.Errorf("add task to dispatcher failed, err: %v", err)
	}

	m.Lock()
	m.currentUpdateTaskTime = time.Now()
	m.Unlock()

	return nil
}

func (m *Manager) DeleteTask(taskType model.TaskType, payloadID int64, isNotify bool) error {
	switch taskType {
	case consts.Sync:
		err := m.dispatcher.RemoveTaskByIdlID(payloadID)
		if err != nil {
			log.Error("delete task fail", zap.Error(err))
			return fmt.Errorf("delete task by idl id in dispatcher failed, err: %v", err)
		}
	}

	m.Lock()
	m.currentUpdateTaskTime = time.Now()
	m.Unlock()

	if isNotify {
		// notify redis channel key
		// so that other api services can update task
		if err := m.notifyRedisChannel(taskType, payloadID); err != nil {
			log.Error("failed to notify Redis channel", zap.Error(err))
			return err
		}
	}

	return nil
}

func (m *Manager) notifyRedisChannel(taskType model.TaskType, idlID int64) error {
	ctx := context.Background()
	msg, err := sonic.MarshalString(task.Message{
		Command: task.DeleteTask,
		Task: model.Task{
			Type:  taskType,
			IdlID: idlID,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return retry.Do(func() error {
		if err := m.rdb.Publish(ctx, consts.RdbKeyTask, msg).Err(); err != nil {
			return fmt.Errorf("failed to publish to Redis channel: %w", err)
		}
		return nil
	}, retry.LastErrorOnly(true), retry.Context(ctx), retry.Attempts(6), retry.Delay(3*time.Second))
}

func (m *Manager) updateAgentTasks() {
	log.Debug("start update agentSrv tasks")
	var group errgroup.Group

	// Launch worker goroutines to update tasks for each agentSrv
	for _, agentSrv := range m.agents {
		agentSrv := agentSrv
		group.Go(func() error {
			return m.updateTasksForAgent(agentSrv.ID)
		})
	}
	if err := group.Wait(); err != nil {
		log.Error("update agentSrv tasks failed", zap.Error(err))
	}
}

func (m *Manager) updateTasksForAgent(serviceID string) error {
	kxClient, err := m.GetAgentClient()
	if err != nil {
		log.Error("get agent client failed", zap.Error(err))
		return err
	}

	tasks := m.dispatcher.GetTasksByServiceID(serviceID)
	rpcRes, err := kxClient.UpdateTask(context.Background(), &taskmodel.UpdateTaskReq{
		Tasks: tasks,
	})

	log.Debug("update tasks to agent service", zap.Reflect("rpcRes", rpcRes))

	if err != nil {
		log.Error("update tasks to RPC client failed", zap.Error(err))
		return err
	}

	if rpcRes.Code != 0 {
		log.Error("update tasks failed", zap.String("err", rpcRes.Msg))
		return err
	}

	log.Debug("update tasks to agent service successfully",
		zap.String("service_id", serviceID),
		zap.Reflect("tasks", tasks),
	)
	return nil
}

// 该方法通过两个 Goroutine 实现了定时同步服务和定时更新代理任务的功能。其中，定时同步服务的 Goroutine 以较长的时间间隔运行，而定时更新代理任务的 Goroutine 则以较短的时间间隔运行，监测任务变化并在变化时进行更新
func (m *Manager) startUpdate() {
	// 同步服务
	go func() {
		for {
			time.Sleep(m.syncAgentInterval)
			m.SyncService()
		}
	}()

	// 其实还是同步服务
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			m.Lock()
			if m.lastUpdateTaskTime != m.currentUpdateTaskTime {
				if m.currentUpdateTaskTime.Add(m.updateTaskInterval).After(time.Now()) {
					log.Debug("task changed, start update agent tasks")
					m.updateAgentTasks()
					m.lastUpdateTaskTime = m.currentUpdateTaskTime
				}
			}
			m.Unlock()
		}
	}()
}

// SyncService
// sync service from registry
func (m *Manager) SyncService() {
	log.Debug("start sync service  ")

	agentMetas, err := m.registry.GetAgents()
	if err != nil {
		log.Error("get agentMetas fail", zap.Error(err))
		return
	}

	metaInMemory := make(map[string]struct{})
	metaInRedis := make(map[string]struct{})

	for _, agentSrv := range m.agents {
		metaInMemory[agentSrv.ID] = struct{}{}
	}
	for _, agentSrv := range agentMetas {
		metaInRedis[agentSrv.ID] = struct{}{}
	}

	var needAddAgentIDs, needDelAgentIDs []string

	for serviceID := range metaInMemory {
		if _, ok := metaInRedis[serviceID]; !ok {
			needDelAgentIDs = append(needDelAgentIDs, serviceID)
		}
	}

	for serviceID := range metaInRedis {
		if _, ok := metaInMemory[serviceID]; !ok {
			needAddAgentIDs = append(needAddAgentIDs, serviceID)
		}
	}

	for _, serviceID := range needAddAgentIDs {
		err = m.dispatcher.AddService(serviceID)
		if err != nil {
			log.Error("add service fail", zap.Error(err))
			continue
		}
	}

	for _, serviceId := range needDelAgentIDs {
		err = m.dispatcher.DelService(serviceId)
		if err != nil {
			log.Error("delete service fail", zap.Error(err))
			continue
		}
	}

	m.agents = agentMetas

	log.Debug("sync service complete", zap.Reflect("agentMetas", agentMetas))

	if len(needAddAgentIDs) != 0 || len(needDelAgentIDs) != 0 {
		// service changed, update cron
		m.updateAgentTasks()
	}
}
