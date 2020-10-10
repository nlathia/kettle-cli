VIRTUALENV_NAME=$(shell pwd | rev | cut -d '/' -f 1 | rev)-dev

.PHONY: install uninstall test clean

uninstall:
	@./bin/remove_pyenv.sh

install: uninstall
	@./bin/setup_pyenv.sh $(VIRTUALENV_NAME)

test:
	@pytest

clean:
	@./bin/cleanup.sh

localhost:
	@./bin/launch.sh
