
.PHONY: install
install:
	go get ./...
	go install .

.PHONY: rebuild
rebuild:
	find . | grep -E "(.DS_Store)" | xargs rm -rf
	rm -f templates/templates.go
	go-bindata -pkg templates -o templates/templates.go templates/...
	go get ./...
	go install .
