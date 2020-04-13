package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
)

var password string
var project string
var url string
var user string

func main() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&password, "password", "", "Token password")
	flag.StringVar(&project, "project", "", "Jenkins project")
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

	builds, err := j.BuildsForProject(project)
	if err != nil {
		log.Fatal(err)
	}

	checking := map[string]bool{}
	for {
		for _, build := range builds {
			_, ok := checking[build.ID]
			if !ok {
				checking[build.ID] = true
				go func(build jenkins.BuildInfo) {
					if build.Result != "" {
						return
					}
					fmt.Printf("Monitoring %s\n", build.FullDisplayName)
					errorIndex := 0
					err := j.MonitorLog(&build,
						func(build *jenkins.BuildInfo, moreToCome bool) error {
							if foundAt := strings.Index(
								build.ConsoleLog[errorIndex:], "FAILED",
							); foundAt > 0 {
								for foundAt > 0 {
									errorIndex = errorIndex + foundAt + 1
									if errorIndex >= len(build.ConsoleLog) {
										break
									}
									foundAt = strings.Index(build.ConsoleLog[errorIndex:], "FAILED")
								}
								fmt.Printf("Failure on build %s - %sconsole\n",
									build.FullDisplayName,
									build.URL,
								)
							}
							return nil
						})
					if err != nil {
						log.Fatal(err)
					}
					build, err = j.GetBuildInfo(build.URL)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Printf("Done monitoring %s - %s\n", build.FullDisplayName, build.Result)
					if err != nil {
						log.Fatal(err)
					}
				}(build)
			}
		}
		time.Sleep(5 * time.Second)
		builds, err = j.BuildsForProject(project)
		if err != nil {
			log.Fatal(err)
		}
	}
	// find the builds that are still on going
	// and watch their consoles.
}
