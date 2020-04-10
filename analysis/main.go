package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

func main() {
	err := processFiles(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func processFiles(files []string) error {
	var builds []jenkins.BuildInfo
	if len(files) == 0 {
		dat, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Failed to read from stdin", err)
		}
		builds, err = readBuild(dat)
		if err != nil {
			return fmt.Errorf("Failed to process %v", err)
		}
	} else {
		for _, f := range files {
			// FIXME: support - as a filename for stdin
			dat, err := ioutil.ReadFile(f)
			if err != nil {
				return fmt.Errorf("Failed to read %s - %v", f, err)
			}
			b, err := readBuild(dat)
			if err != nil {
				return fmt.Errorf("Failed to process %s - %v", f, err)
			}
			builds = append(builds, b...)
		}
	}
	// now process the build info
	// then output stuff
	bytes, err := json.Marshal(builds)
	if err != nil {
		log.Fatal(err)
	}
	// FIXME: potentially output in different formats
	fmt.Println(string(bytes))
	return nil
}

func readBuild(fileContents []byte) ([]jenkins.BuildInfo, error) {
	var builds []jenkins.BuildInfo

	err := json.Unmarshal(fileContents, &builds)
	return builds, err
}
