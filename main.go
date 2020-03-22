package jenkins_health

import (
	"fmt"
	"log"
)

func main() {
	// talk to jenkins and grab jobs.
	// job number, pass/fail, failure type, time to run, machine
	// console log
	// can I get timing info?  Mathew mentioned that was somewhere
	j := &JenkinsAPI{}

	text, err := j.ConsoleLog("cvl-gerrit", 1003)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(text)
}
