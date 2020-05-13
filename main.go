package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
	resty "github.com/go-resty/resty/v2"
)

var password string
var project string
var url string
var user string
var build string
var db string

func main() {
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&password, "password", "", "Token password")
	flag.StringVar(&project, "project", "", "Jenkins project")
	flag.StringVar(&build, "build", "", "Jenkins build")
	flag.StringVar(&url, "url", "http://localhost:8080", "Jenkins url")
	flag.StringVar(&db, "db", "", "SQLite database to store results")
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
		builds, err = j.BuildsForProject(project)
		if err != nil {
			log.Fatal(err)
		}
	}
	if db != "" {
		// FIXME: create db if it doesn't already exist.
		db, err := sql.Open("sqlite3", db)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		stmt, err := db.Prepare(`
			INSERT INTO builds
				(url, console_log, built_on, duration, displayName,
				timeStamp, result)
			VALUES (?, ?, ?, ?, ?, ?, ?)`)
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()
		for _, b := range builds {
			_, err := stmt.Exec(b.URL, b.ConsoleLog, b.BuiltOn,
				b.Duration, b.FullDisplayName, b.Timestamp, b.Timestamp,
				b.Result)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		bytes, err := json.Marshal(builds)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bytes))
	}
}
