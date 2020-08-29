package stages

import (
	"regexp"
	"strings"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

// Analyser for decoding stages
type Analyser struct {
}

// AnalyseBuild fills in the analysis fields after examining the ConsoleLog
func (a *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	// rip through the log and turn it into stages
	stages := strings.Split(an.ConsoleLog, "\r\n[Pipeline] stage\r\n")
	if len(stages) < 2 {
		return nil
	}
	for _, stage := range stages {
		pipelineName := `(?m)\[Pipeline\] { \(([^)]+)\)`
		r := regexp.MustCompile(pipelineName)
		matches := r.FindStringSubmatch(stage)
		var name string
		if len(matches) > 1 {
			name = matches[1]
		}
		an.Stages = append(an.Stages, analysis.BuildStage{
			Name: name,
			Log:  stage,
		})
	}
	return nil
}
