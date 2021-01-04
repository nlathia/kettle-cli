VIRTUALENV_NAME=$(shell pwd | rev | cut -d '/' -f 1 | rev)-dev

.PHONY: install uninstall test clean localhost

clean:
	@./bin/cleanup.sh

uninstall: clean
	@./bin/remove_pyenv.sh

install: uninstall
	@./bin/setup_pyenv.sh $(VIRTUALENV_NAME)

test:
	@pytest

localhost:
	@./bin/launch.sh 
