
.PHONY: get rebuild install full-install

full-install: get rebuild install

get:
	go get ./...

rebuild:
	find . | grep -E "(.DS_Store)" | xargs rm -rf
	rm -f templates/templates.go
	go-bindata -pkg templates -o templates/templates.go templates/...

install:
	go install .

