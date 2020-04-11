package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/colinnewell/jenkins-queue-health/analysis"
	"github.com/colinnewell/jenkins-queue-health/analysis/human"
	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

var _ = Describe("Analysis/Human/Human", func() {

	var _ = Describe("With a build", func() {

		It("Should produce human readable time info", func() {
			build := analysis.AnalysedBuild{BuildInfo: jenkins.BuildInfo{Timestamp: 1585940109411, Duration: 1311955}}
			var a human.Analyser
			err := a.AnalyseBuild(&build)

			Expect(build.TimeRunReadable).To(Equal("Fri Apr 3 19:55:09 +0100 BST 2020"))
			Expect(build.DurationReadable).To(Equal("21m51.955s"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
