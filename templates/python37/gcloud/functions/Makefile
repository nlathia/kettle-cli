VIRTUALENV_NAME=gcloud-function.$(shell pwd | rev | cut -d '/' -f 1 | rev)

.PHONY: install uninstall test clean localhost

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
