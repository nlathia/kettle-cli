VIRTUALENV_NAME=aws-lambda.$(shell pwd | rev | cut -d '/' -f 1 | rev)

.PHONY: install uninstall test clean invoke-local invoke

clean:
	@./bin/cleanup.sh

uninstall: clean
	@./bin/remove_pyenv.sh

install: uninstall
	@./bin/setup_pyenv.sh $(VIRTUALENV_NAME)

test:
	@pytest

invoke-local:
	@./bin/invoke_local.sh 

invoke:
	@./bin/invoke.sh 
