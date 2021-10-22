GOPATH:=$(shell go env GOPATH)

all: build

.PHONY: test docker build clean

build:
		#go build -gcflags '-N -l'  -o darwin-app-api
		go build -o myapp

abuild:
		go build  -v -o myapp main.go

test:
		go test -v ./... -cover

docker:
		docker build . -t myapp:latest

clean:
		rm my-app

