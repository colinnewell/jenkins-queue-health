all: jenkins-queue-health jenkins-queue-health-analysis jenkins-watcher

jenkins-queue-health-analysis: analysis/cli/main.go jenkins/jenkins.go \
	analysis/build.go analysis/human/human.go analysis/yath/yath.go \
	analysis/spurious/spurious.go
	go build -o jenkins-queue-health-analysis analysis/cli/main.go

jenkins-queue-health: main.go jenkins/jenkins.go
	go build -o jenkins-queue-health main.go

jenkins-watcher: watcher/cli/main.go jenkins/jenkins.go
	go build -o jenkins-watcher watcher/cli/main.go

lint:
	golangci-lint run
	golint ./...

test:
	go test

install:
	cp jenkins-queue-health jenkins-queue-health-analysis jenkins-watcher /usr/local/bin
