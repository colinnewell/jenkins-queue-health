package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/colinnewell/jenkins-queue-health/analysis"
	"github.com/colinnewell/jenkins-queue-health/analysis/human"
	"github.com/colinnewell/jenkins-queue-health/analysis/spurious"
	"github.com/colinnewell/jenkins-queue-health/analysis/yath"
	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

func main() {
	err := processFiles(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

func processFiles(files []string) error {
	var builds []analysis.AnalysedBuild
	if len(files) == 0 {
		dat, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Failed to read from stdin - %v", err)
		}
		builds, err = readBuild(dat)
		if err != nil {
			return fmt.Errorf("Failed to process - %v", err)
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

func readBuild(fileContents []byte) ([]analysis.AnalysedBuild, error) {
	var builds []jenkins.BuildInfo
	var analysed []analysis.AnalysedBuild

	err := json.Unmarshal(fileContents, &builds)
	if err != nil {
		return analysed, err
	}
	analysed = make([]analysis.AnalysedBuild, len(builds))
	var yath yath.Analyser
	var human human.Analyser
	var spurious spurious.Analyser
	for i, b := range builds {
		// FIXME: do I want to allow some optimisation based on job?
		// perhaps have some kind of analyser object?
		analysed[i] = analysis.AnalyseBuild(b)

		if err := yath.AnalyseBuild(&analysed[i]); err != nil {
			return analysed, err
		}
		if err := human.AnalyseBuild(&analysed[i]); err != nil {
			return analysed, err
		}
		if err := spurious.AnalyseBuild(&analysed[i]); err != nil {
			return analysed, err
		}
	}
	return analysed, err
}
