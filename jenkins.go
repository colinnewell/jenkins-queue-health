package jenkins_health

import (
	"fmt"
	"net/http"
)

type JenkinsAPI struct {
	Client http.Client
}

func (jenkins *JenkinsAPI) Runs(jobName string) ([]int, error) {
	return nil, nil
}

func (jenkins *JenkinsAPI) ConsoleLog(jobName string, build int) (string, error) {
	// get a list of recent builds for the job
	// grab the logs
	fmt.Println("vim-go")
	return "", nil
}
