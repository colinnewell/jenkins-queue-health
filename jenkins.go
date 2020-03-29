package jenkins_health

import (
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type JenkinsAPI struct {
	Client     resty.Client
	JenkinsURL string
}

type build struct {
	Url string `json:"url"`
}

type runs struct {
	Builds []build `json:"build"`
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
	if r.IsError() {
		return nil, fmt.Errorf("Error getting job runs: %s - %s", r.Status(), r)
	}
	// unmarshal the json
	var v runs
	err = json.Unmarshal(r.Body(), &v)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(v.Builds))
	for i := 0; i < len(v.Builds); i++ {
		urls[i] = v.Builds[i].Url
	}
	return urls, nil
}

func (jenkins *JenkinsAPI) ConsoleLog(jobName string, buildUrl string) (string, error) {
	url := fmt.Sprintf("%s/logText/progressiveText?start=0", buildUrl)
	return "", nil
}
