package main_test

import (
	resty "github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/colinnewell/jenkins-queue-health/jenkins"
)

var _ = Describe("Jenkins Client", func() {

	// this stuff should go in BeforeSuite etc. but I'm
	// struggling to get that stuff to fire.
	httpmock.Reset()
	client := resty.New()
	// add a default 404 handler
	client.SetTransport(httpmock.DefaultTransport)
	defer httpmock.DeactivateAndReset()

	var j *jenkins.API

	var _ = BeforeEach(func() {
		httpmock.RegisterResponder("GET", "http://test/job/test/api/json?tree=builds[url]", httpmock.NewStringResponder(200, `{
			  "_class": "hudson.model.FreeStyleProject",
			  "builds": [
				{
				  "_class": "hudson.model.FreeStyleBuild",
				  "url": "http://test/job/test/3/"
				},
				{
				  "_class": "hudson.model.FreeStyleBuild",
				  "url": "http://test/job/test/1/"
				}
			  ]
			}`))
		httpmock.RegisterResponder("GET", `=~^http://test/job/test/\d+/api/json`, httpmock.NewStringResponder(200, `{
  "_class": "hudson.model.FreeStyleBuild",
  "duration": 73,
  "fullDisplayName": "test #3",
  "id": "3",
  "result": "FAILURE",
  "timestamp": 1585940109411,
  "builtOn": "machine1",
  "changeSet": { "_class": "hudson.scm.EmptyChangeLogSet" }
			}`))
		httpmock.RegisterResponder("GET", `=~^http://test/job/test/\d+/logText`, httpmock.NewStringResponder(200, `FAIL`))
		dt := httpmock.DefaultTransport
		client.SetTransport(dt)

		j = &jenkins.API{
			Client:     client,
			JenkinsURL: "http://test",
		}
	})

	var _ = Describe("Get Builds", func() {

		It("Should retrieve a list of builds", func() {
			builds, err := j.BuildsForProject("test", false)

			expected := []jenkins.BuildInfo{
				{
					BuiltOn:         "machine1",
					Duration:        73,
					FullDisplayName: "test #3",
					ID:              "3",
					Result:          "FAILURE",
					Timestamp:       1585940109411,
					ConsoleLog:      "FAIL",
					URL:             "http://test/job/test/3/",
				},
				{
					BuiltOn:         "machine1",
					Duration:        73,
					FullDisplayName: "test #3",
					ID:              "3",
					Result:          "FAILURE",
					Timestamp:       1585940109411,
					ConsoleLog:      "FAIL",
					URL:             "http://test/job/test/1/",
				},
			}
			Expect(builds).To(Equal(expected))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	var _ = Describe("GetBuildInfo", func() {

		It("Should return headling info about a build", func() {
			build, err := j.GetBuildInfo("http://test/job/test/1/")

			expected := jenkins.BuildInfo{
				BuiltOn:         "machine1",
				Duration:        73,
				FullDisplayName: "test #3",
				ID:              "3",
				Result:          "FAILURE",
				Timestamp:       1585940109411,
				ConsoleLog:      "",
				URL:             "http://test/job/test/1/",
			}
			Expect(build).To(Equal(expected))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
