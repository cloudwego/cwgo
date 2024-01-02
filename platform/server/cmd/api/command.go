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

package api

import (
	"github.com/cloudwego/cwgo/platform/server/shared/args"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	opt := args.NewApiArgs()
	cmd := &cobra.Command{
		Use:   "api",
		Short: "api service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := run(opt); err != nil {
				return err
			}
			return nil
		},
	}
	opt.AddFlags(cmd)
	return cmd
}
