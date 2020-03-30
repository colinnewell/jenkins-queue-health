package main

import (
	"fmt"
	"log"

	jenkins "github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
)

func main() {
	// talk to jenkins and grab jobs.
	// job number, pass/fail, failure type, time to run, machine
	// console log
	// can I get timing info?  Mathew mentioned that was somewhere
	client := resty.New()
	client.SetBasicAuth("admin", "119f8713bc75a829dbc4df57170ed8f5a3")
	j := &jenkins.JenkinsAPI{
		Client:     client,
		JenkinsURL: "http://localhost:8080",
	}

	urls, err := j.Runs("test")
	//text, err := j.ConsoleLog("cvl-gerrit", 1003)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", urls)
	//fmt.Println(text)
}
