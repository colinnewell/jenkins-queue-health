package jenkins

import (
	"encoding/json"
	"fmt"
	"log"

	resty "github.com/go-resty/resty/v2"
)

// API struct for downloading info from Jenkins
type API struct {
	Client     *resty.Client
	JenkinsURL string
}

type build struct {
	URL string `json:"url"`
}

type runs struct {
	Builds []build `json:"builds"`
}

// BuildInfo build structure
type BuildInfo struct {
	BuiltOn         string `json:"builtOn"`
	Duration        int64  `json:"duration"`
	FullDisplayName string `json:"fullDisplayName"`
	Id              string `json:"id"`
	Result          string `json:"result"`
	Timestamp       int    `json:"timestamp"`
	ConsoleLog      string `json:"log"`
	URL             string `json:"url"`
	// FIXME: need to find structure for this.
	//	ChangeSet       string `json:changeSet`
}

// Runs provide an array of urls for the builds assocaited with the job
func (jenkins *API) Runs(jobName string) ([]string, error) {
	url := fmt.Sprintf("%s/job/%s/api/json?tree=builds[url]",
		jenkins.JenkinsURL,
		jobName,
	)
	r, err := jenkins.Client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		return nil, fmt.Errorf(
			"Error getting job runs: %s - %s", r.Status(), r,
		)
	}
	var v runs
	err = json.Unmarshal(r.Body(), &v)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(v.Builds))
	for i := 0; i < len(v.Builds); i++ {
		urls[i] = v.Builds[i].URL
	}
	return urls, nil
}

// ConsoleLog fills in the builds ConsoleLog.
func (jenkins *API) ConsoleLog(build *BuildInfo) error {
	url := fmt.Sprintf("%slogText/progressiveText?start=0", build.URL)
	r, err := jenkins.Client.R().Get(url)
	if err != nil {
		return err
	}
	if r.IsError() {
		return fmt.Errorf(
			"Error getting console: %s - %s", r.Status(), r,
		)
	}
	build.ConsoleLog = r.String()
	return nil
}

// BuildInfo gather the headline build info from the build url
func (jenkins *API) GetBuildInfo(buildUrl string) (BuildInfo, error) {
	b := BuildInfo{URL: buildUrl}

	url := fmt.Sprintf("%sapi/json?tree=id,fullDisplayName,result,timestamp,builtOn,changeSet,duration", buildUrl)
	r, err := jenkins.Client.R().Get(url)
	if err != nil {
		return b, err
	}
	if r.IsError() {
		return b, fmt.Errorf(
			"Error getting build info: %s - %s", r.Status(), r,
		)
	}

	err = json.Unmarshal(r.Body(), &b)
	return b, err
}

// BuildsForProject download a list of builds for the project filling in the
// ConsoleLog for any builds that aren't marked as success.
func (jenkins *API) BuildsForProject(project string) ([]BuildInfo, error) {
	urls, err := jenkins.Runs(project)
	if err != nil {
		log.Fatal(err)
	}
	var builds []BuildInfo
	for _, url := range urls {
		build, err := jenkins.GetBuildInfo(url)
		if err != nil {
			return builds, err
		}
		if build.Result != "SUCCESS" {
			// grab log so we can examine reasons for failure
			err := jenkins.ConsoleLog(&build)
			if err != nil {
				return builds, err
			}
		}
		builds = append(builds, build)
	}
	return builds, nil
}
