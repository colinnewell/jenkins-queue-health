package human

import (
	"time"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

// Analyser for yath test output (a Perl test runner)
type Analyser struct {
}

// AnalyseBuild fills in the analysis fields after examining the ConsoleLog
func (a *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	// set DurationReadable and TimeRunReadable
	d := time.Duration(an.Duration) * time.Millisecond
	an.DurationReadable = d.String()

	t := time.Unix(an.Timestamp/1000, an.Timestamp%1000)
	an.TimeRunReadable = t.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	return nil
}
