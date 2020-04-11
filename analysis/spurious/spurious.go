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
	spuriousBuild := `Solr request failed - Timed out while waiting for socket to become ready for reading`
	an.SpuriousFail = strings.Contains(an.ConsoleLog, spuriousBuild)
	return nil
}
