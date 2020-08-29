package analysis

import (
	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

// Analyser interface for processing the build info, making the raw information
// more easily digested.
type Analyser interface {
	AnalyseBuild(AnalysedBuild) error
}

// AnalysedBuild struct with extra information about builds
type AnalysedBuild struct {
	SpuriousFail bool `json:"spuriousFail"`
	// divide up the console log into parts.
	// prelude
	// logs after
	// failed test list
	// passed test list
	// provide clearer info from the exisitng field
	// More human readable version of when the build was run
	FailureSummary   string       `json:"failureSummary"`
	TimeRunReadable  string       `json:"timeReadable"`
	DurationReadable string       `json:"durationReadable"`
	Stages           []BuildStage `json:"stages,omitempty"`

	jenkins.BuildInfo
}

type BuildStage struct {
	// FIXME: add timestamps for start/stop
	Name        string       `json:"name"`
	Log         string       `json:"log"`
	LockedSteps []LockedStep `json:"lockedSteps,omitempty"`
}

type LockedStep struct {
	Name string `json:"name"`
	Log  string `json:"log"`
}

// AnalyseBuild turn the BuildInfo object into an AnalysedBuild one with extra
// info to be filled in during analysis
func AnalyseBuild(build jenkins.BuildInfo) AnalysedBuild {
	b := AnalysedBuild{BuildInfo: build}
	return b
}
