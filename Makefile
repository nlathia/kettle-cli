
.PHONY: get install

get:
	go get ./...

install:
	go build -o $GOPATH/bin/kettle

