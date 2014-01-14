build: deps
	go build

deps:
	go get ./...

install: build
	go install
