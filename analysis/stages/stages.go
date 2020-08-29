package stages

import (
	"regexp"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

// Analyser for decoding stages
type Analyser struct {
}

// AnalyseBuild fills in the analysis fields after examining the ConsoleLog
func (a *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	// rip through the log and turn it into stages
	pipelineStage := regexp.MustCompile(`(?m)\[Pipeline\] (// )?stage\r\n`)
	stages := pipelineStage.Split(an.ConsoleLog, -1)
	if len(stages) < 2 {
		return nil
	}
	for _, stage := range stages {
		if stage == "" {
			continue
		}
		pipelineName := `(?m)\[Pipeline\] { \(([^)]+)\)`
		r := regexp.MustCompile(pipelineName)
		matches := r.FindStringSubmatch(stage)
		var name string
		if len(matches) > 1 {
			name = matches[1]
		}
		bs := analysis.BuildStage{
			Name: name,
			Log:  stage,
		}
		ExtractLockInfo(&bs)
		an.Stages = append(an.Stages, bs)
	}
	return nil
}

func ExtractLockInfo(stage *analysis.BuildStage) {
	lockText := regexp.MustCompile(`(?m)\[Pipeline\] (// )?lock\r\n`)
	lockChunks := lockText.Split(stage.Log, -1)
	if len(lockChunks) < 2 {
		return
	}
	for _, chunk := range lockChunks {
		if chunk == "" {
			continue
		}
		lockName := `(?m)Trying to acquire lock on \[([^\]]+)\]`
		r := regexp.MustCompile(lockName)
		matches := r.FindStringSubmatch(chunk)
		var name string
		if len(matches) > 1 {
			name = matches[1]
		}
		stage.LockedSteps = append(stage.LockedSteps, analysis.LockedStep{
			Name: name,
			Log:  chunk,
		})
	}
}
