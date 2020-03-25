package jenkins_health

import (
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type JenkinsAPI struct {
	Client     resty.Client
	JenkinsURL string
}

func (jenkins *JenkinsAPI) Runs(jobName string) ([]string, error) {
	url := fmt.Sprintf("%s/job/%s/api/json?tree=builds[url]",
		jenkins.JenkinsURL,
		jobName,
	)
	r, err := jenkins.Client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (jenkins *JenkinsAPI) ConsoleLog(jobName string, buildUrl string) (string, error) {
	url := fmt.Sprintf("%s/logText/progressiveText?start=0", buildUrl)
	return "", nil
}
