package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	Password string
	Project  string
	URL      string
	User     string
	Pause    int64
	Triggers map[string]string
}

var c config

func main() {
	viper.SetConfigName("jenkins-watcher.toml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	// FIXME: set some defaults
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	pflag.String("password", "", "Token password")
	pflag.Int64("pause", 60, "Polling interval")
	pflag.String("project", "", "Jenkins project")
	pflag.String("url", "http://localhost:8080", "Jenkins url")
	pflag.String("user", "", "Username")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}

	client := resty.New()
	if c.User != "" && c.Password != "" {
		client.SetBasicAuth(c.User, c.Password)
		client.SetDisableWarn(true)
	}
	j := &jenkins.API{
		Client:     client,
		JenkinsURL: c.URL,
	}

	builds, err := j.BuildsForProject(c.Project)
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
					err := j.MonitorLog(&build, c.Pause,
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
		time.Sleep(time.Duration(c.Pause) * time.Second)
		builds, err = j.BuildsForProject(c.Project)
		if err != nil {
			log.Fatal(err)
		}
	}
	// find the builds that are still on going
	// and watch their consoles.
}
