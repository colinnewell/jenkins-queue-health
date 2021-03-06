package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
)

var password string
var project string
var url string
var user string
var build string
var successToo bool

func main() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&password, "password", "", "Token password")
	flag.StringVar(&project, "project", "", "Jenkins project")
	flag.StringVar(&build, "build", "", "Jenkins build")
	flag.StringVar(&url, "url", "http://localhost:8080", "Jenkins url")
	flag.BoolVar(&successToo, "success-too", false, "Successful builds too")
	flag.Parse()

	client := resty.New()
	if user != "" && password != "" {
		client.SetBasicAuth(user, password)
		client.SetDisableWarn(true)
	}
	j := &jenkins.API{
		Client:     client,
		JenkinsURL: url,
	}

	var builds []jenkins.BuildInfo
	if build != "" {
		build, err := j.GetBuildInfo(j.BuildURL(project, build))
		if err != nil {
			log.Fatal(err)
		}
		err = j.ConsoleLog(&build)
		if err != nil {
			log.Fatal(err)
		}
		builds = []jenkins.BuildInfo{build}
	} else {
		var err error
		builds, err = j.BuildsForProject(project, successToo)
		if err != nil {
			log.Fatal(err)
		}
	}
	bytes, err := json.Marshal(builds)
	if err != nil {
		log.Fatal(err)
	}
	// FIXME: potentially output in different formats
	fmt.Println(string(bytes))
}
