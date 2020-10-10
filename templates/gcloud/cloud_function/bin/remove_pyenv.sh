# Makefile command: make uninstall
# If a .python-version file exists, this script will
# uninstall the virtual environment that is defined in it

set -e

source $(dirname $0)/_config.sh

if [[ -f ".python-version" ]]; then
    echo "\n ⏱  Force removing: $VIRTUALENV_NAME"
    pyenv uninstall -f $VIRTUALENV_NAME
    rm .python-version

    echo "\n ✅  Done."
fi

