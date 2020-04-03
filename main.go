package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
)

func main() {
	// talk to jenkins and grab jobs.
	// job number, pass/fail, failure type, time to run, machine
	// console log
	// can I get timing info?  Mathew mentioned that was somewhere
	client := resty.New()
	client.SetBasicAuth("admin", "119f8713bc75a829dbc4df57170ed8f5a3")
	client.SetDisableWarn(true)
	j := &jenkins.API{
		Client:     client,
		JenkinsURL: "http://localhost:8080",
	}

	urls, err := j.Runs("test")
	if err != nil {
		log.Fatal(err)
	}
	var builds []jenkins.BuildInfo
	for _, url := range urls {
		build, err := j.BuildInfo(url)
		if err != nil {
			log.Fatal(err)
		}
		if build.Result != "SUCCESS" {
			err := j.ConsoleLog(&build)
			// check what the deal is with the log
			if err != nil {
				// FIXME: Fatal is a bit lame
				log.Fatal(err)
			}
			builds = append(builds, build)
		}
	}
	bytes, err := json.Marshal(builds)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
