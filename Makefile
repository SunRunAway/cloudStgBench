goinstall:
	go get github.com/aws/aws-sdk-go
	go get qiniupkg.com/api.v7/kodo
	go get golang.org/x/net/context

build: goinstall
	GOOS=linux GOARCH=amd64 go build -o ./bin/cloudStgBench
