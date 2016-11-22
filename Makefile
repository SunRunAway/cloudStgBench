goinstall:
	go get github.com/aws/aws-sdk-go

build: goinstall
	GOOS=linux GOARCH=amd64 go build -o ./bin/cloudStgBench
