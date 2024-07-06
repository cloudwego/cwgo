package job

const jobTemplate = `package job

// Run is an example function job
func Run() {
	// TODO: fill with your own logic

}
`

const jobMainTemplate = `package main

import (
    "log"
    schedule "{{.PackagePrefix}}"
)

func main() {
    err := schedule.Run()
    if err != nil {
        log.Fatalf("job failed: %v", err)
    }
}

`

const jobScheduleTemplate = `package schedule

import (
	"sync"
    {{- range .JobInfos }}
	{{.JobName}} "{{.PackagePrefix}}/{{.JobName}}/job"
	{{- end }}
)

func Run() error {
	var wg sync.WaitGroup

	{{- range .JobInfos }}
	wg.Add(1)
	go func() {
		defer wg.Done()
		{{.JobName}}.Run()
	}()
	{{- end }}

	wg.Wait()
	return nil
}

`

const scriptTemplate = `#!/bin/bash

echo "Building job binary..."
go build -o job ../cmd/main.go

if [ $? -ne 0 ]; then
    echo "Error: Failed to build job."
    exit 1
fi

echo "Running job..."
./job

if [ $? -ne 0 ]; then
    echo "Error: job execution failed."
    exit 1
fi

echo "job Done."

rm job

if [ $? -ne 0 ]; then
    echo "Error: Failed to remove job binary."
    exit 1
fi

`
