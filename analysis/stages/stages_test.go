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

func TestCoupleOfStages(t *testing.T) {
	build := analysis.AnalysedBuild{}
	build.ConsoleLog = "foo\r\n[Pipeline] stage\r\n[Pipeline] { (Git)\r\nsomething\r\n[Pipeline] }\r\n[Pipeline] // stage\r\n[Pipeline] stage\r\n[Pipeline] { (Test)\r\n[2020-08-27T09:16:54.990Z]  > git --version # timeout=10\r\n[Pipeline] }\r\n[Pipeline] // stage\r\n[Pipeline] End of Pipeline\r\nFinished: NOT_BUILT"
	a := stages.Analyser{}

	err := a.AnalyseBuild(&build)
	if err != nil {
		t.Error(err)
	}

	expectedStages := []analysis.BuildStage{
		analysis.BuildStage{
			Log:  "foo\r\n",
			Name: "",
		},
		analysis.BuildStage{
			Log:  "[Pipeline] { (Git)\r\nsomething\r\n[Pipeline] }\r\n",
			Name: "Git",
		},
		analysis.BuildStage{
			Log:  "[Pipeline] { (Test)\r\n[2020-08-27T09:16:54.990Z]  > git --version # timeout=10\r\n[Pipeline] }\r\n",
			Name: "Test",
		},
		analysis.BuildStage{
			Log: "[Pipeline] End of Pipeline\r\nFinished: NOT_BUILT",
		},
	}
	if diff := cmp.Diff(build.Stages, expectedStages); diff != "" {
		t.Errorf("Different json (-got +expected):\n%s\n", diff)
	}
}
