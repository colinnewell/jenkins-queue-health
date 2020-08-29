package stages

import (
	"strings"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

// Analyser for decoding stages
type Analyser struct {
}

// AnalyseBuild fills in the analysis fields after examining the ConsoleLog
func (a *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	// rip through the log and turn it into stages
	stages := strings.Split(an.ConsoleLog, "\r\n// [Pipeline] stage\r\n")
	for _, stage := range stages {
		name := ""
		// find first
		//[Pipeline] { (Test)
		// then look for first and last timestamps
		if name != "" {
			an.Stages = append(an.Stages, analysis.BuildStage{
				Name: name,
				Log:  stage,
			})
		}
	}
	return nil
}
