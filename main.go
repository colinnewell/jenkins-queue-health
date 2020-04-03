package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
)

var user string
var password string
var url string

func main() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&password, "password", "", "Token password")
	flag.StringVar(&url, "url", "http://localhost:8080", "Jenkins url")
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
