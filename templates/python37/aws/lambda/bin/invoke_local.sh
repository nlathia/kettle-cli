#!/bin/bash
# Makefile command: `make invoke-local` 
# This script tests the function locally & interactively

FILE="event.json"
if [ -f "$1" ]; then
    echo "✅  Using $1 as input."
    FILE="$1"
elif [ -f "$FILE" ]; then
    echo "✅  Using $FILE as input."
else
    read -p " 🎯 Enter a JSON payload to test: "
    echo "$REPLY" > "$FILE"
    validation=$(cat event.json | python -m json.tool)
    if [[ "$?" -ne 0 ]]; then
        echo " ❌  Input is not valid JSON, aborting."
        rm "$FILE"
        exit 1
    fi
fi

python-lambda-local -f {{.FunctionName}} main.py "$FILE"
