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

package options

import "github.com/spf13/cobra"

type GlobalOptions struct {
	ServerMode string
	ConfigType string
	ConfigPath string
}

func NewGlobalOptions() *GlobalOptions {
	return &GlobalOptions{}
}

func (o *GlobalOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()

	flagSet.StringVarP(&o.ServerMode, "server_mode", "", "", "server run mode (dev/pro)")
	flagSet.StringVarP(&o.ConfigType, "config_type", "", "", "config type (file)")
	flagSet.StringVarP(&o.ConfigPath, "config_path", "", "", "config file path")
}
