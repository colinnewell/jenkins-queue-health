all: jenkins-queue-health jenkins-queue-health-analysis

jenkins-queue-health-analysis: analysis/main.go jenkins/jenkins.go
	go build -o jenkins-queue-health-analysis analysis/main.go

jenkins-queue-health: main.go jenkins/jenkins.go
	go build -o jenkins-queue-health main.go

test:
	go test

install:
	cp jenkins-queue-health /usr/local/bin
