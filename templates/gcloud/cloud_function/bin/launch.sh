# Makefile command: `make localhost` 
# This script launches the docker container locally

echo "\n ⏱  Launching locally (use ctrl+c to exit)..."

echo "\n 🎯 Replace the URL and input in the following command to test:"
echo " $ curl -X POST http://0.0.0.0:8080 -d '{\"key\": \"value\"}'\n\n"

functions-framework --target={{ .FunctionName }}
