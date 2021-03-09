
.PHONY: get install full-install

full-install: get rebuild install

get:
	go get ./...

install:
	go install .

