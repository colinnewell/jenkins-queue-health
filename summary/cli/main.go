package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/colinnewell/jenkins-queue-health/analysis"
	"github.com/colinnewell/jenkins-queue-health/summary"
)

func main() {
	err := processFiles(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func processFiles(files []string) error {
	var builds []summary.SummarisedBuild
	if len(files) == 0 {
		dat, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin - %v", err)
		}
		builds, err = readBuild(dat)
		if err != nil {
			return fmt.Errorf("failed to process - %v", err)
		}
	} else {
		for _, f := range files {
			// FIXME: support - as a filename for stdin
			dat, err := os.ReadFile(f)
			if err != nil {
				return fmt.Errorf("failed to read %s - %v", f, err)
			}
			b, err := readBuild(dat)
			if err != nil {
				return fmt.Errorf("failed to process %s - %v", f, err)
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

func readBuild(fileContents []byte) ([]summary.SummarisedBuild, error) {
	var analysed []analysis.AnalysedBuild
	var summarised []summary.SummarisedBuild

	err := json.Unmarshal(fileContents, &analysed)
	if err != nil {
		return summarised, err
	}
	return summarised, nil
}
