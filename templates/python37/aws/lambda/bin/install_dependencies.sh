# Installs the Python dependencies

set -e

pip install --upgrade pip setuptools wheel
for i in ./requirements*txt; do
    echo "\n ⏱  Installing requirements in: $i"
    pip install -r $i
    echo "\n ✅  Done."
done
