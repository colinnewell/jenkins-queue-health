all: jenkins-queue-health jenkins-queue-health-analysis jenkins-watcher \
		jenkins-queue-summary

jenkins-queue-health-analysis: analysis/cli/main.go jenkins/jenkins.go \
	analysis/build.go analysis/*/*.go
	go build -o jenkins-queue-health-analysis analysis/cli/main.go

jenkins-queue-summary: summary/cli/main.go jenkins/jenkins.go \
	analysis/build.go analysis/*/*.go summary/build.go
	go build -o jenkins-queue-summary summary/cli/main.go

jenkins-queue-health: main.go jenkins/jenkins.go
	go build -o jenkins-queue-health main.go

jenkins-watcher: watcher/cli/main.go jenkins/jenkins.go
	go build -o jenkins-watcher watcher/cli/main.go

lint:
	golangci-lint run
	golint ./...

test:
	go test ./...

install:
	cp jenkins-queue-health jenkins-queue-health-analysis jenkins-watcher /usr/local/bin
