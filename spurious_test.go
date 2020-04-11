package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/colinnewell/jenkins-queue-health/analysis"
	. "github.com/colinnewell/jenkins-queue-health/analysis/spurious"
	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

var _ = Describe("Analysis/Spurious", func() {

	var _ = Describe("With a spurious test failure", func() {
		consoleSnippet := `
Triggered by Gerrit: https://gerrit/c/repo/+/62501
Running as SYSTEM
Building remotely on jenkins-build01 (auto-scm branch-e2e branch-tests) in workspace /home/jenkins/
workspace/gerrit
( PASSED )  job  4    t/unit/bin/free.t
( PASSED )  job  2    t/unit/api/euro/v1.t
[  FAIL  ]  job 378  + Unexpected warning: AFTER MUNGE: $VAR1 = {
[  FAIL  ]  job 378  +           'type' => [],
[  FAIL  ]  job 378  +           'county' => [],
[  FAIL  ]  job 378  +           'language' => [],
Solr request failed - Timed out while waiting for socket to become ready for reading
[  FAIL  ]  job 378  +           'stem_phrases' => '1',
[  FAIL  ]  job 378  +           'none' => '',
[  FAIL  ]  job 378  +           'basic' => '0',
[  FAIL  ]  job 378  +           'exact' => '',
[  FAIL  ]  job 378  +           'industry' => [],
[  FAIL  ]  job 378  +           'min_match' => '0',
( PASSED )  job 374    t/unit/web/admin/email_reports.t
( PASSED )  job 370    t/unit/web/admin/verification.t
[  FAIL  ]  job 378  +           'keywords' => ''
[  FAIL  ]  job 378  +         };
[  DEBUG ]  job 378    t/unit/web/blah/some_search.t line 832
(  DIAG  )  job 378    +--------------+------------------+---------+
( STDERR )  job 378              'type' => [],
( STDERR )  job 378              'county' => [],
( STDERR )  job 378              'salary' => [],
( STDERR )  job 378              'language' => [],
( STDERR )  job 378              'language_type' => 'any',
( STDERR )  job 378              'stem_phrases' => '1'
( FAILED )  job 378    t/unit/web/blah/some_search.t
< REASON >  job 378    Test script returned error (Err: 124)
< REASON >  job 378    Assertion failures were encountered (Count: 124)
< REASON >  job 378    Subtest failures were encountered (Count: 17)
( PASSED )  job 409    t/unit/web/www/results.t

================================================================================

Run ID: CBE22A3E-75CD-11EA-9047-B53E7E167E9F

The following test jobs failed:
  [CBEDFF26-75CD-11EA-9047-B53E7E167E9F] 378: t/unit/web/blah/some_search.t
  [CBEDFF26-75CD-11EA-9047-B53E7E167E9F] 379: t/unit/web/blah/another_search.t

ESC[0m
40.0690s on wallclock (9.34 usr 1.35 sys + 161.03 cusr 71.00 csys = 242.72 CPU)

Build step 'Execute shell' marked build as failure
The following test jobs failed:
  [CBEDFF26-75CD-11EA-9047-B53E7E167E0F] 378: t/unit/web/blah/some_search.t
  [CBEDFF26-75CD-11EA-9047-B53E7E167E0F] 379: t/unit/web/blah/another_search.t

ESC[0m
`

		It("Should set the spurious flag", func() {
			build := analysis.AnalysedBuild{BuildInfo: jenkins.BuildInfo{ConsoleLog: consoleSnippet}}
			var a Analyser
			err := a.AnalyseBuild(&build)

			Expect(build.SpuriousFail).To(Equal(true))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
