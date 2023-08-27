/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package service

import (
	"context"
	"errors"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/agent/agentservice"
	"github.com/cloudwego/kitex/client"
	"time"
)

type BuiltinService struct {
	Id             string
	LastUpdateTime time.Time
	RpcClient      agentservice.Client
}

func NewBuiltinService(serviceId, address string) (*BuiltinService, error) {
	rpcClient, err := agentservice.NewClient(serviceId, client.WithHostPorts(address))
	if err != nil {
		return nil, err
	}

	return &BuiltinService{
		Id:             serviceId,
		LastUpdateTime: time.Now(),
		RpcClient:      rpcClient,
	}, nil
}

func (s *BuiltinService) GenerateCode(ctx context.Context, idlId int64) error {
	rpcRes, err := s.RpcClient.GenerateCode(ctx, &agent.GenerateCodeReq{IdlId: idlId})
	if err != nil {
		return err
	}
	if rpcRes.Code != 0 {
		return errors.New(rpcRes.Msg)
	}

	return nil
}
