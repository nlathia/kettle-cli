# Makefile command: make install
# This script will create a pyenv virtual environment 
# And install all of the dependencies defined in the
# requirements file(s).

set -e

source $(dirname $0)/_config.sh

if [[ $(pyenv versions | grep -L $PYTHON_VERSION) ]]; then
    echo "\n ⏱  Installing Python $PYTHON_VERSION"
    pyenv install $PYTHON_VERSION
fi

echo "\n ⏱  Creating a $PYTHON_VERSION environment: $1"
env PYTHON_CONFIGURE_OPTS="--enable-framework CC=clang" \
    pyenv virtualenv \
        --force $PYTHON_VERSION \
        $1
pyenv local $1

echo "\n ✅  Done."

source $(dirname $0)/install_dependencies.sh
