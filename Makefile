deps:
	go get ./...

build: deps
	go build

install: build
	go install
