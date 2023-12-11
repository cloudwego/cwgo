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

package errx

import (
	"reflect"
)

const (
	PkgName = "github.com/cloudwego/cwgo/platform/server/shared/errx"
)

type (
	Error interface {
		Code() int32
		Error() string
	}
)

type myErr struct {
	c int32
	s string
}

func New(code int32, text string) Error {
	return &myErr{
		c: code,
		s: text,
	}
}

func (e *myErr) Code() int32 {
	return e.c
}

func (e *myErr) Error() string {
	return e.s
}

func GetCode(err error) int32 {
	if err == nil {
		return -1
	}

	rv := reflect.ValueOf(err)

	rvk := rv.Kind()
	for rvk == reflect.Ptr {
		rv = rv.Elem()
		rvk = rv.Kind()
	}

	rvt := rv.Type()

	if rvt == reflect.TypeOf(myErr{}) {
		nf := rv.NumField()
		for i := 0; i < nf; i++ {
			rvtf := rvt.Field(i)
			rvf := rv.Field(i)

			if rvtf.PkgPath == PkgName && rvtf.Name == "c" && rvf.Kind() == reflect.Int32 {
				if rvf.CanInt() {
					return int32(rvf.Int())
				}
				return -1
			}
		}
	}

	return -1
}
