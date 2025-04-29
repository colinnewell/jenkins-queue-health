package jenkins

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

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
	ID              string `json:"id"`
	Result          string `json:"result"`
	Timestamp       int64  `json:"timestamp"`
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
			"error getting job runs: %s - %s", r.Status(), r,
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

// BuildURL construct the url for the build in Jenkins
func (jenkins *API) BuildURL(jobName string, build string) string {
	return fmt.Sprintf("%s/job/%s/%s/",
		jenkins.JenkinsURL,
		jobName,
		build,
	)
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
			"error getting console: %s - %s", r.Status(), r,
		)
	}
	build.ConsoleLog = r.String()
	return nil
}

// GetLogChunk pull the a chunk of the console log and return the response so
// that it can be checked to see if there is more.
func (jenkins *API) GetLogChunk(build *BuildInfo, offset uint64) (*resty.Response, error) {
	url := fmt.Sprintf("%slogText/progressiveText?start=%d", build.URL, offset)
	r, err := jenkins.Client.R().Get(url)
	if err != nil {
		return r, err
	}
	if r.IsError() {
		return r, fmt.Errorf(
			"error getting console: %s - %s", r.Status(), r,
		)
	}
	build.ConsoleLog = build.ConsoleLog + r.String()
	return r, nil
}

// GetBuildInfo gather the headline build info from the build url
func (jenkins *API) GetBuildInfo(buildURL string) (BuildInfo, error) {
	b := BuildInfo{URL: buildURL}

	url := fmt.Sprintf("%sapi/json?tree=id,fullDisplayName,result,timestamp,builtOn,changeSet,duration", buildURL)
	r, err := jenkins.Client.R().Get(url)
	if err != nil {
		return b, err
	}
	if r.IsError() {
		return b, fmt.Errorf(
			"error getting build info: %s - %s", r.Status(), r,
		)
	}

	err = json.Unmarshal(r.Body(), &b)
	return b, err
}

// BuildsForProject download a list of builds for the project filling in the
// ConsoleLog for any builds that aren't marked as success.
func (jenkins *API) BuildsForProject(project string, successToo bool) ([]BuildInfo, error) {
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
		if successToo || build.Result != "SUCCESS" {
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

// MonitorLog keep following the log and call callback until it's all
// retrieved.
func (jenkins *API) MonitorLog(build *BuildInfo, pause int64,
	callback func(*BuildInfo, bool) error) error {

	var offset uint64
	r, err := jenkins.GetLogChunk(build, offset)
	if err != nil {
		return err
	}
	moreToCome := true
	for moreToCome {
		headers := r.Header()
		moreData, ok := headers["X-More-Data"]
		moreToCome = ok && moreData[0] == "true"
		if r.Size() > 0 {
			// only call the callback if we have new data
			err := callback(build, moreToCome)
			if err != nil {
				return err
			}
		}
		if moreToCome {
			if r.Size() == 0 {
				// pause to wait for more content
				time.Sleep(time.Duration(pause) * time.Second)
			}
			offsetString, ok := headers["X-Text-Size"]
			if !ok {
				return fmt.Errorf("no X-Text-Size returned despite indicating there is more data to read")
			}
			offset, err := strconv.ParseUint(offsetString[0], 10, 64)
			if err != nil {
				return err
			}
			r, err = jenkins.GetLogChunk(build, offset)
			if err != nil {
				return err
			}
		}
	}
	return err
}
