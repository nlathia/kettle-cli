#!/bin/bash
# Makefile command: `make invoke` 
# This script tests the function locally & interactively

RESULT="outputfile.txt"
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

aws lambda  invoke --function-name {{.Name}} \
    --payload "fileb://$FILE" $RESULT

echo "✅  Result written to: $RESULT:"
cat $RESULT
