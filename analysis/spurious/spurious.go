package spurious

import (
	"strings"

	"github.com/colinnewell/jenkins-queue-health/analysis"
)

// Analyser for spurious failures
type Analyser struct {
}

// AnalyseBuild toggles the SpuriousFail field based on whether this looks like
// a spurious failure
func (a *Analyser) AnalyseBuild(an *analysis.AnalysedBuild) error {
	// it would be neat to make this configurable.
	spuriousBuild := []string{
		`Solr request failed - Timed out while waiting for socket to become ready for reading`,
		`Timed out waiting for Solr 7 to come up.`,
		// `Event timeout after 60 second(s) for job`, - this ones a maybe
	}
	for _, message := range spuriousBuild {
		an.SpuriousFail = strings.Contains(an.ConsoleLog, message)
		if an.SpuriousFail {
			break
		}
	}
	return nil
}
