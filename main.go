package main

import (
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
	j := &jenkins.API{
		Client:     client,
		JenkinsURL: "http://localhost:8080",
	}

	urls, err := j.Runs("test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", urls)
	for _, url := range urls {
		// grab the build info.

		build, err := j.BuildInfo(url)
		if err != nil {
			log.Fatal(err)
		}
		if build.Result != "SUCCESS" {
			err := j.ConsoleLog(&build)
			if err != nil {
				// FIXME: Fatal is a bit lame
				log.Fatal(err)
			}
			log.Printf("%s: %s\n", url, build.ConsoleLog)
		}
	}
}
