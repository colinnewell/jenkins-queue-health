all: jenkins-queue-health test

jenkins-queue-health: main.go jenkins/jenkins.go
	go build -o jenkins-queue-health main.go

test:
	go test

install:
	cp jenkins-queue-health /usr/local/bin
