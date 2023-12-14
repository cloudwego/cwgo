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

package idl

import (
	"context"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/biz/model/idl"
	"github.com/cloudwego/cwgo/platform/server/cmd/api/internal/svc"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/model"
	"github.com/cloudwego/cwgo/platform/server/shared/logger"
	"github.com/cloudwego/cwgo/platform/server/shared/task"
	"go.uber.org/zap"
	"strconv"
)

const (
	successMsgUpdateIDL = "update idl successfully"
)

type UpdateIDLLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateIDLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateIDLLogic {
	return &UpdateIDLLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateIDLLogic) UpdateIDL(req *idl.UpdateIDLReq) (res *idl.UpdateIDLRes) {
	if req.Status != 0 {
		if _, ok := consts.IdlStatusNumMap[int(req.Status)]; !ok {
			return &idl.UpdateIDLRes{
				Code: consts.ErrNumParamIdlStatus,
				Msg:  consts.ErrMsgParamIdlStatus,
			}
		}
	}

	client, err := l.svcCtx.Manager.GetAgentClient()
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcGetClient, zap.Error(err))
		return &idl.UpdateIDLRes{
			Code: consts.ErrNumRpcGetClient,
			Msg:  consts.ErrMsgRpcGetClient,
		}
	}

	rpcRes, err := client.UpdateIDL(l.ctx, &agent.UpdateIDLReq{
		RepositoryId: req.ID,
		ServiceName:  req.ServiceName,
		Status:       req.Status,
	})
	if err != nil {
		logger.Logger.Error(consts.ErrMsgRpcConnectClient, zap.Error(err))
		return &idl.UpdateIDLRes{
			Code: consts.ErrNumRpcConnectClient,
			Msg:  consts.ErrMsgRpcConnectClient,
		}
	}
	if rpcRes.Code != 0 {
		return &idl.UpdateIDLRes{
			Code: rpcRes.Code,
			Msg:  rpcRes.Msg,
		}
	}

	switch req.Status {
	case consts.IdlStatusNumInactive:
		go func() {
			// delete task
			_ = l.svcCtx.Manager.DeleteTask(strconv.FormatInt(req.ID, 10))
		}()
	case consts.IdlStatusNumActive:
		go func() {
			// delete task
			_ = l.svcCtx.Manager.AddTask(
				task.NewTask(
					model.Type_sync_idl_data,
					l.svcCtx.Manager.SyncIdlInterval.String(),
					&model.Data{
						SyncIdlData: &model.SyncIdlData{
							IdlId: req.ID,
						},
					},
				),
			)
		}()
	}

	return &idl.UpdateIDLRes{
		Code: 0,
		Msg:  successMsgUpdateIDL,
	}
}
