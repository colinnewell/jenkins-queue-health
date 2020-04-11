package analysis

import (
	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

type AnalysedBuild struct {
	SpuriousFail bool
	// divide up the console log into parts.
	// prelude
	// logs after
	// failed test list
	// passed test list
	// provide clearer info from the exisitng field
	// More human readable version of when the build was run
	TimeRunReadable  string
	DurationReadable string

	jenkins.BuildInfo
}

func AnalyseBuild(build jenkins.BuildInfo) {
}
