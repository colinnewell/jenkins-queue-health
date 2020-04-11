package analysis

import (
	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

type Analyser interface {
	AnalyseBuild(AnalysedBuild) error
}

type AnalysedBuild struct {
	SpuriousFail bool `json:"spuriosFail"`
	// divide up the console log into parts.
	// prelude
	// logs after
	// failed test list
	// passed test list
	// provide clearer info from the exisitng field
	// More human readable version of when the build was run
	FailureSummary   string `json:"failureSummary"`
	TimeRunReadable  string `json:"timeReadable"`
	DurationReadable string `json:"durationReadable"`

	jenkins.BuildInfo
}

type JobConfig struct {
	JobRegex string
}

func AnalyseBuild(build jenkins.BuildInfo) AnalysedBuild {
	b := AnalysedBuild{BuildInfo: build}
	return b
}
