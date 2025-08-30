/*
 * Copyright 2024 CloudWeGo Authors
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
 */

package cronjob

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cloudwego/cwgo/config"
	"github.com/stretchr/testify/assert"
)

func TestCronJobValidCreation(t *testing.T) {
	args := &config.CronJobArgument{
		JobName:       []string{"cronjob1"},
		GoMod:         "github.com/cloudwego/cwgo",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
		JobFile:       "job.go",
	}

	err := Cronjob(args)
	assert.NoError(t, err)

	checkFiles := []string{
		"test_out/cmd/main.go",
		"test_out/internal/schedule.go",
		"test_out/scripts/run.sh",
		"test_out/internal/job/job.go",
	}

	for _, file := range checkFiles {
		_, err := os.Stat(file)
		assert.NoError(t, err)
	}

	contentChecks := map[string]string{
		"test_out/internal/schedule.go": "job.Cronjob1(ctx)",
	}

	for file, content := range contentChecks {
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		assert.Contains(t, string(data), content)
	}

	err = os.RemoveAll(args.OutDir)
	assert.NoError(t, err)
}

func TestAddScheduleNewJobs(t *testing.T) {
	original := `

package schedule

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"test/internal/job"
)

func Init(ctx context.Context, c *cron.Cron) {
	var err error

	_, err = c.AddFunc("* * * * *", func() {
		select {
		case <-ctx.Done():
			log.Println("JobOne terminated.")
			return
		default:
			job.JobOne(ctx)
		}
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	_, err = c.AddFunc("* * * * *", func() {
		select {
		case <-ctx.Done():
			log.Println("JobTwo terminated.")
			return
		default:
			job.JobTwo(ctx)
		}
	})
	if err != nil {
		log.Fatalf("Error adding new cron job: %v", err)
	}

}


`

	jobs := []JobInfo{
		{JobName: "JobOne", PackagePrefix: "github.com/cloudwego/cwgo", GoModule: "github.com/cloudwego/cwgo"},
		{JobName: "JobThree", PackagePrefix: "github.com/cloudwego/cwgo", GoModule: "github.com/cloudwego/cwgo"},
	}

	expected := `package schedule

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"test/internal/job"
)

func Init(ctx context.Context, c *cron.Cron) {
	var err error

	_, err = c.AddFunc("* * * * *", func() {
		select {
		case <-ctx.Done():
			log.Println("JobOne terminated.")
			return
		default:
			job.JobOne(ctx)
		}
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	_, err = c.AddFunc("* * * * *", func() {
		select {
		case <-ctx.Done():
			log.Println("JobTwo terminated.")
			return
		default:
			job.JobTwo(ctx)
		}
	})
	if err != nil {
		log.Fatalf("Error adding new cron job: %v", err)
	}
	
	_, err = c.AddFunc("* * * * *", func() {
		select {
		case <-ctx.Done():
			log.Println("JobThree terminated.")
			return
		default:
			job.JobThree(ctx)
		}
	})
	if err != nil {
		log.Fatalf("Error adding new cron job: %v", err)
	}

}
`

	result, err := addScheduleNewJobs(original, jobs)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestAddNewJobs(t *testing.T) {
	original := `

package job

import (
	"context"
)

func JobOne(ctx context.Context) {
	// TODO: fill with your own logic
}

func JobTwo(ctx context.Context) {
	// TODO: fill with your own logic
}

`

	jobs := []JobInfo{
		{JobName: "JobOne", PackagePrefix: "github.com/cloudwego/cwgo", GoModule: "github.com/cloudwego/cwgo"},
		{JobName: "JobThree", PackagePrefix: "github.com/cloudwego/cwgo", GoModule: "github.com/cloudwego/cwgo"},
	}

	expected := `package job

import (
	"context"
)

func JobOne(ctx context.Context) {
	// TODO: fill with your own logic
}

func JobTwo(ctx context.Context) {
	// TODO: fill with your own logic
}

func JobThree(ctx context.Context) {
	// TODO: fill with your own logic
}
`

	result, err := addNewJobs(original, jobs)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

//func TestCreateCronRunCall(t *testing.T) {
//	jobName := "cronjob1"
//	expected := `wg.Add(1)
//go func() {
//	defer wg.Done()
//	cronjob1.Run()
//}()`
//
//	result, err := addNewJobs(jobName)
//	assert.NoError(t, err)
//	var buf bytes.Buffer
//	err := printer.Fprint(&buf, token.NewFileSet(), result)
//	assert.NoError(t, err)
//	assert.Equal(t, expected, buf.String())
//}

func TestCronJobMissingJobName(t *testing.T) {
	args := &config.CronJobArgument{
		JobName:       []string{""},
		GoMod:         "github.com/cloudwego/cwgo",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
		JobFile:       "job.go",
	}

	if os.Getenv("BE_CRASHER") == "1" {
		err := Cronjob(args)
		assert.Error(t, err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCronJobMissingJobName")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	defer func() {
		_ = os.RemoveAll(args.OutDir)
	}()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Logf("process ran with err %v, want exit status 1", err)
}

func TestCronJobWithDifferentModule(t *testing.T) {
	args := &config.CronJobArgument{
		JobName:       []string{"cronjob2"},
		GoMod:         "github.com/cloudwego/another_test",
		PackagePrefix: "github.com/cloudwego/another_test",
		OutDir:        "./another_test_out",
		JobFile:       "job.go",
	}

	if os.Getenv("BE_CRASHER") == "1" {
		err := Cronjob(args)
		assert.Error(t, err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCronJobWithDifferentModule")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Logf("process ran with err %v, want exit status 1", err)
}

func TestCronJobWithEmptyModule(t *testing.T) {
	args := &config.CronJobArgument{
		JobName:       []string{"cronjob3"},
		GoMod:         "",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
		JobFile:       "job.go",
	}

	if os.Getenv("BE_CRASHER") == "1" {
		err := Cronjob(args)
		assert.Error(t, err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCronJobWithEmptyModule")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Logf("process ran with err %v, want exit status 1", err)
}

func TestCronJobWithoutGoModInGOPATH(t *testing.T) {
	args := &config.CronJobArgument{
		JobName:       []string{"cronjob3"},
		GoMod:         "",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
		JobFile:       "job.go",
	}

	// Setup a temporary GOPATH
	tmpGopath, err := os.MkdirTemp("", "gopath")
	if err != nil {
		t.Fatalf("Failed to create temp GOPATH: %v", err)
	}
	defer os.RemoveAll(tmpGopath)

	// Create a directory structure under GOPATH/src
	projectPath := filepath.Join(tmpGopath, "src", "github.com", "cloudwego", "cwgo")
	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Move test file to the project directory
	currentTestFile := filepath.Join(projectPath, "cronjob_test.go")
	if err := os.Rename(os.Args[0], currentTestFile); err != nil {
		t.Fatalf("Failed to move test file to project directory: %v", err)
	}

	// Set the GOPATH environment variable
	oldGopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", tmpGopath)
	defer os.Setenv("GOPATH", oldGopath)

	if os.Getenv("BE_CRASHER") == "1" {
		err := Cronjob(args)
		assert.Error(t, err)
		return
	}

	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestCronJobWithoutGoModInGOPATH")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("Expected process to exit with an error, got nil. Output: %s", output)
	}

	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
		t.Logf("Expected process to exit with status 0, got %d. Output: %s", exitErr.ExitCode(), output)
	}
}
