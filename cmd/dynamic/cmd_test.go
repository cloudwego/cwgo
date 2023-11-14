/*
 * Copyright 2022 CloudWeGo Authors
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
 * MIT License
 *
 * Copyright (c) 2018 Alec Aivazis
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package dynamic

import (
	"io"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/cloudwego/cwgo/config"
	pseudotty "github.com/creack/pty"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

type expectConsole interface {
	ExpectString(string)
	ExpectEOF()
	SendLine(string)
	Send(string)
}

type consoleWithErrorHandling struct {
	console *expect.Console
	t       *testing.T
}

func (c *consoleWithErrorHandling) ExpectString(s string) {
	if _, err := c.console.ExpectString(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("ExpectString(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) SendLine(s string) {
	if _, err := c.console.SendLine(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("SendLine(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) Send(s string) {
	if _, err := c.console.Send(s); err != nil {
		c.t.Helper()
		c.t.Fatalf("Send(%q) = %v", s, err)
	}
}

func (c *consoleWithErrorHandling) ExpectEOF() {
	if _, err := c.console.ExpectEOF(); err != nil {
		c.t.Helper()
		c.t.Fatalf("ExpectEOF() = %v", err)
	}
}

type PromptTest struct {
	name      string
	prompt    []*survey.Question
	procedure func(expectConsole)
	expected  interface{}
}

func RunTest(t *testing.T, procedure func(expectConsole), test func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error) {
	t.Helper()
	t.Parallel()

	pty, tty, err := pseudotty.Open()
	if err != nil {
		t.Fatalf("failed to open pseudotty: %v", err)
	}

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	if err != nil {
		t.Fatalf("failed to create console: %v", err)
	}
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		procedure(&consoleWithErrorHandling{console: c, t: t})
	}()

	if err := test(c.Tty(), c.Tty(), c.Tty()); err != nil {
		t.Error(err)
	}

	if err := c.Tty().Close(); err != nil {
		t.Errorf("error closing Tty: %v", err)
	}
	<-donec
}

func TestTerminal(t *testing.T) {
	var answer []string
	test := PromptTest{
		name:   "Terminal",
		prompt: generateType,
		procedure: func(c expectConsole) {
			c.Send(string(terminal.KeyArrowDown))
			c.Send(" ")
			// Select Wednesday.
			c.Send(string(terminal.KeyArrowDown))
			c.Send(string(terminal.KeyArrowDown))
			c.SendLine(" ")
			c.ExpectEOF()
		},
		expected: []string{"server", "client"},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestCommon(t *testing.T) {
	var answer config.CommonParam
	test := PromptTest{
		name:   "Common",
		prompt: commonQuestion(),
		procedure: func(c expectConsole) {
			c.ExpectString("Select service type")
			c.SendLine("HTTP")

			c.ExpectString("Please input service name")
			c.SendLine("test")

			c.ExpectString("Please input module")
			c.SendLine("cwgo/test")

			c.ExpectString("Please input idlpath")
			c.SendLine("./idl")

			c.ExpectEOF()
		},
		expected: config.CommonParam{
			Type:    "HTTP",
			Service: "test",
			GoMod:   "cwgo/test",
			IdlPath: "./idl",
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestProtoSearch(t *testing.T) {
	var answer config.SliceParam
	test := PromptTest{
		name:   "ProtoSearch",
		prompt: protoSearch(),
		procedure: func(c expectConsole) {
			c.ExpectString("Please input proto search path if exists, space as separator")
			c.SendLine("./aaa ./bbb ./ccc /home/ddd")

			c.ExpectEOF()
		},
		expected: config.SliceParam{
			ProtoSearchPath: []string{"./aaa", "./bbb", "./ccc", "/home/ddd"},
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestDefaultConfig(t *testing.T) {
	var answer dfConfig
	test := PromptTest{
		name:   "DefaultConfig",
		prompt: defaultConfig(),
		procedure: func(c expectConsole) {
			c.ExpectString("Whether use default config to generate project")
			c.SendLine("Y")

			c.ExpectEOF()
		},
		expected: dfConfig{
			DefaultConfig: true,
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestRegistryConfig(t *testing.T) {
	var answer config.ServerArgument
	test := PromptTest{
		name:   "RegistryConfig",
		prompt: registryConfig(),
		procedure: func(c expectConsole) {
			c.ExpectString("Please select a registry")
			c.SendLine("ZK")

			c.ExpectEOF()
		},
		expected: config.ServerArgument{
			Registry: "ZK",
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestCustomConfig(t *testing.T) {
	var answer config.SliceParam
	test := PromptTest{
		name:   "CustomConfig",
		prompt: customConfig(),
		procedure: func(c expectConsole) {
			c.ExpectString("Please input custom param")
			c.SendLine("param1 param2 param3")

			c.ExpectEOF()
		},
		expected: config.SliceParam{
			Pass: []string{"param1", "param2", "param3"},
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestDBConfig(t *testing.T) {
	answer := config.NewModelArgument()
	test := PromptTest{
		name:   "DBConfig",
		prompt: dbConfig,
		procedure: func(c expectConsole) {
			c.ExpectString("Select db type")
			c.SendLine("SQLite")

			c.ExpectString("Please input db DSN")
			c.SendLine("this is a dsn")

			c.ExpectEOF()
		},
		expected: &config.ModelArgument{
			DSN:     "this is a dsn",
			Type:    "sqlite",
			OutPath: "biz/dal/query",
			OutFile: "gen.go",
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}

func TestClientNum(t *testing.T) {
	answer := cNum{}
	test := PromptTest{
		name:   "ClientNum",
		prompt: clientNum,
		procedure: func(c expectConsole) {
			c.ExpectString("Please enter the number of generated clients")
			c.SendLine("6")

			c.ExpectEOF()
		},
		expected: cNum{
			ClientNum: "6",
		},
	}
	RunTest(t, test.procedure, func(in terminal.FileReader, out terminal.FileWriter, stderr io.Writer) error {
		return survey.Ask(test.prompt, &answer, survey.WithStdio(in, out, stderr))
	})
	require.Equal(t, test.expected, answer)
}
