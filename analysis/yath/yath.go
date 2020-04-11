package yath

import (
	"regexp"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

// Analyser for yath test output (a Perl test runner)
type Analyser struct {
}

// AnalyseBuild fills in the analysis fields after examining the ConsoleLog
func (a *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	failSummaryPattern := `(?m)The following test jobs failed:(?:\s+\[[-0-9A-F]+\] \d+: .*$)+`
	r := regexp.MustCompile(failSummaryPattern)
	matches := r.FindAllString(an.BuildInfo.ConsoleLog, -1)
	for _, v := range matches {
		an.FailureSummary = an.FailureSummary + v + "\n"
	}
	return nil
}
