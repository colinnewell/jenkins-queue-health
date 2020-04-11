package yath

import (
	"regexp"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

type Analyser struct {
}

func (yath *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	failSummaryPattern := `(?m)The following test jobs failed:(?:\s+\[[-0-9A-F]+\] (\d+): (.*)$)+`
	r := regexp.MustCompile(failSummaryPattern)
	matches := r.FindAllStringSubmatch(an.BuildInfo.ConsoleLog, -1)
	for _, v := range matches {
		for j, submatch := range v[1:] {
			if j%2 == 1 {
				an.FailureSummary = append(an.FailureSummary, submatch)
			}
		}
	}
	return nil
}
