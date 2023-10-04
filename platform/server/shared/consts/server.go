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

package consts

type ServerType uint32

const (
	ProjectName = "cwgo"
)

const (
	ServerTypeNumApi ServerType = iota + 1
	ServerTypeNumAgent

	ServerTypeApi   = "api"
	ServerTypeAgent = "agent"

	ServerTypeEnvName = "CWGO_SERVER_TYPE"
)

var (
	ServerTypeMapToStr = map[ServerType]string{
		ServerTypeNumApi:   ServerTypeApi,
		ServerTypeNumAgent: ServerTypeAgent,
	}
)

type ServerMode uint32

const (
	ServerModeNumDev ServerMode = iota + 1
	ServerModeNumPro

	ServerModeDev = "dev"
	ServerModePro = "pro"
)

var (
	ServerModeMapToStr = map[ServerMode]string{
		ServerModeNumDev: ServerModeDev,
		ServerModeNumPro: ServerModePro,
	}
	ServerModeMapToNum = map[string]ServerMode{
		ServerModeDev: ServerModeNumDev,
		ServerModePro: ServerModeNumPro,
	}
)
