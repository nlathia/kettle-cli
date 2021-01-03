# Makefile command: `make localhost` 
# This script launches the function locally

echo "\n ⏱  Launching locally (use ctrl+c to exit)..."

echo "\n 🎯 Example test command:"
echo " $ curl -X POST http://localhost:9090/ -d '{\"key\": \"value\"}'\n\n"

PORT=8080 && docker run -p 9090:${PORT} -e PORT=${PORT} $1
