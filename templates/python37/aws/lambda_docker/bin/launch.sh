# Makefile command: `make localhost` 
# This script launches the function locally

echo "\n ⏱  Launching locally (use ctrl+c to exit)..."

echo "\n 🎯 Example test command:"
echo " $ curl -X POST http://localhost:9000/2015-03-31/functions/function/invocations -d '{\"payload\":\"hello world!\"}'\n"

docker run -p 9000:8080 $1
