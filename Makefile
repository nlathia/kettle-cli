
.PHONY: install
install:
	find . | grep -E "(.DS_Store)" | xargs rm -rf
	rm -f templates/templates.go
	go-bindata -pkg templates -ignore=\\*.DS_Store -o templates/templates.go templates/...
	go install .