package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJenkinsQueueHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JenkinsQueueHealth Suite")
}
