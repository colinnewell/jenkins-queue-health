package stages_test

import (
	"testing"

	"github.com/colinnewell/jenkins-queue-health/analysis"
	"github.com/colinnewell/jenkins-queue-health/analysis/stages"

	"github.com/google/go-cmp/cmp"
)

func TestNoStagesPresent(t *testing.T) {
	build := analysis.AnalysedBuild{}
	build.ConsoleLog = "foo"
	a := stages.Analyser{}

	err := a.AnalyseBuild(&build)
	if err != nil {
		t.Error(err)
	}

	var expectedStages []analysis.BuildStage
	if diff := cmp.Diff(build.Stages, expectedStages); diff != "" {
		t.Errorf("Different json (-got +expected):\n%s\n", diff)
	}

}
