# Makefile command: `make clean` 
# This script removes all of the pycache files from the current directory

echo "\n ⏱  Finding and removing pycache files..."

find . | grep -E "(__pycache__|\.pyc|\.pyo$)" | xargs rm -rf

echo "\n ✅  Done."