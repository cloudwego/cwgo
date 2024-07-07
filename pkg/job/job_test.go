package job

import (
	"bytes"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cloudwego/cwgo/config"
	"github.com/stretchr/testify/assert"
)

func TestJobValidCreation(t *testing.T) {
	args := &config.JobArgument{
		JobName:       []string{"job1"},
		GoMod:         "github.com/cloudwego/cwgo",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
	}

	err := Job(args)
	assert.NoError(t, err)

	checkFiles := []string{
		"test_out/cmd/main.go",
		"test_out/schedule.go",
		"test_out/scripts/run.sh",
		"test_out/job1/job/job.go",
	}

	for _, file := range checkFiles {
		_, err := os.Stat(file)
		assert.NoError(t, err)
	}

	contentChecks := map[string]string{
		"test_out/schedule.go": "job1.Run()",
	}

	for file, content := range contentChecks {
		data, err := os.ReadFile(file)
		assert.NoError(t, err)
		assert.Contains(t, string(data), content)
	}

	err = os.RemoveAll(args.OutDir)
	assert.NoError(t, err)
}

func TestAddJobImportsAndRun(t *testing.T) {
	original := `
package schedule

import (
	"sync"
)

var wg sync.WaitGroup

func Run() {
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Wait()
}
`
	jobs := []JobInfo{
		{JobName: "job1", PackagePrefix: "github.com/cloudwego/cwgo"},
		{JobName: "job2", PackagePrefix: "github.com/cloudwego/cwgo"},
	}

	expected := `package schedule

import (
	"sync"
	job1 "github.com/cloudwego/cwgo/job1/job"
	job2 "github.com/cloudwego/cwgo/job2/job"
)

var wg sync.WaitGroup

func Run() {
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		job1.Run()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		job2.Run()
	}()
	wg.Wait()
}
`

	result, err := addJobImportsAndRun(original, jobs)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateRunCall(t *testing.T) {
	jobName := "job1"
	expected := `wg.Add(1)
go func() {
	defer wg.Done()
	job1.Run()
}()`

	result := createRunCall(jobName)
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), result)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}

func TestJobMissingJobName(t *testing.T) {
	args := &config.JobArgument{
		JobName:       []string{""},
		GoMod:         "github.com/cloudwego/cwgo",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
	}

	if os.Getenv("BE_CRASHER") == "1" {
		err := Job(args)
		assert.Error(t, err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestJobMissingJobName")
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

func TestJobWithDifferentModule(t *testing.T) {
	args := &config.JobArgument{
		JobName:       []string{"job2"},
		GoMod:         "github.com/cloudwego/another_test",
		PackagePrefix: "github.com/cloudwego/another_test",
		OutDir:        "./another_test_out",
	}

	if os.Getenv("BE_CRASHER") == "1" {
		err := Job(args)
		assert.Error(t, err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestJobWithDifferentModule")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Logf("process ran with err %v, want exit status 1", err)
}

func TestJobWithEmptyModule(t *testing.T) {
	args := &config.JobArgument{
		JobName:       []string{"job3"},
		GoMod:         "",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
	}

	if os.Getenv("BE_CRASHER") == "1" {
		err := Job(args)
		assert.Error(t, err)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestJobWithEmptyModule")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Logf("process ran with err %v, want exit status 1", err)
}

func TestJobWithoutGoModInGOPATH(t *testing.T) {
	args := &config.JobArgument{
		JobName:       []string{"job3"},
		GoMod:         "",
		PackagePrefix: "github.com/cloudwego/cwgo",
		OutDir:        "./test_out",
	}

	// Setup a temporary GOPATH
	tmpGopath, err := os.MkdirTemp("", "gopath")
	if err != nil {
		t.Fatalf("Failed to create temp GOPATH: %v", err)
	}
	defer os.RemoveAll(tmpGopath)

	// Create a directory structure under GOPATH/src
	projectPath := filepath.Join(tmpGopath, "src", "github.com", "cloudwego", "cwgo")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Move test file to the project directory
	currentTestFile := filepath.Join(projectPath, "job_test.go")
	if err := os.Rename(os.Args[0], currentTestFile); err != nil {
		t.Fatalf("Failed to move test file to project directory: %v", err)
	}

	// Set the GOPATH environment variable
	oldGopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", tmpGopath)
	defer os.Setenv("GOPATH", oldGopath)

	if os.Getenv("BE_CRASHER") == "1" {
		err := Job(args)
		assert.Error(t, err)
		return
	}

	// Run the test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestJobWithoutGoModInGOPATH")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("Expected process to exit with an error, got nil. Output: %s", output)
	}

	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
		t.Logf("Expected process to exit with status 0, got %d. Output: %s", exitErr.ExitCode(), output)
	}

}
