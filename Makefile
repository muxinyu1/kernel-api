.PHONY: all linux windows macos arm clean

BINARY_NAME=kernel-api

all: linux windows macos

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 .

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 .

arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o bin/$(BINARY_NAME)-linux-arm .

clean:
	rm -rf bin/
	mkdir -p bin/