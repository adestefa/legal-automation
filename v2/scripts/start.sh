#!/bin/bash

# Start the Legal Document Automation v2 (Go SSR + HTMX)
echo "Starting Mallon Legal Document Automation v2..."
echo "-----------------------------------------------"

# Define color codes for status messages
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Navigate to the project root directory (where main.go is located)
cd "$(dirname "$0")/.."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH. Please install Go to run this prototype.${NC}"
    exit 1
fi

# Check if port 8080 is already in use
PORT_CHECK_COMMAND=""
if command -v lsof &> /dev/null; then
    PORT_CHECK_COMMAND="lsof -i:8080"
elif command -v netstat &> /dev/null; then
    PORT_CHECK_COMMAND="netstat -tuln | grep ':8080'"
elif command -v ss &> /dev/null; then
    PORT_CHECK_COMMAND="ss -tuln | grep ':8080'"
fi

if [ -n "$PORT_CHECK_COMMAND" ] && eval "$PORT_CHECK_COMMAND" &>/dev/null; then
    echo -e "${YELLOW}Warning: Port 8080 is already in use.${NC}"
    echo "Please stop any running instances with ./stop.sh before starting a new one."
    exit 1
fi

# Build and run the v2 backend with Go SSR + HTMX
echo "Building and starting the v2 server (Go SSR + HTMX)..."
echo "This may take a moment..."

# Check dependencies first
go mod tidy
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to check Go dependencies${NC}"
    exit 1
fi

# Build the v2 binary first for better error handling
echo "Building mallon-v2 binary..."
go build -o mallon-v2 main.go
if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to build v2 server${NC}"
    exit 1
fi

# Start the server in the background and save PID
PID_FILE=".server.pid"
./mallon-v2 &
SERVER_PID=$!

# Store the PID for the stop script
echo $SERVER_PID > "$PID_FILE"

# Give the server a moment to start and check if it's running
sleep 2
if ! kill -0 $SERVER_PID 2>/dev/null; then
    echo -e "${RED}Error: Server failed to start${NC}"
    rm -f "$PID_FILE"
    exit 1
fi

echo -e "${GREEN}Mallon Legal Assistant v2 started successfully!${NC}"
echo "Process ID: $SERVER_PID"
echo "Version: v2.5.28 (Go SSR + HTMX)"
echo "Features: Dynamic document processing (Task 8), document editing, server-side rendering"
echo "Open your browser and navigate to: http://localhost:8080"
echo "Press Ctrl+C to stop the server, or run ./stop.sh from another terminal"
echo ""
echo "Server output:"

# Function to handle cleanup on exit
cleanup() {
    echo ""
    echo "Stopping server..."
    kill $SERVER_PID 2>/dev/null
    rm -f "$PID_FILE"
    echo "Server stopped."
}

# Set trap to handle Ctrl+C
trap cleanup SIGINT SIGTERM

# Wait for the server process
wait $SERVER_PID

# Clean up PID file when the server exits
rm -f "$PID_FILE"