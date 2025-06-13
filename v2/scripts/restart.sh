#!/bin/bash

# Restart the Legal Document Automation v2 Server
echo "Restarting Mallon Legal Document Automation v2 Server..."
echo "---------------------------------------------------------"

# Define color codes for status messages
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory of this script
SCRIPT_DIR="$(dirname "$0")"

# Step 1: Stop the server
echo "Step 1: Stopping the server..."
"$SCRIPT_DIR/stop.sh"
STOP_EXIT_CODE=$?

if [ $STOP_EXIT_CODE -ne 0 ]; then
    echo -e "${YELLOW}Warning: Stop script reported issues. Continuing with restart anyway...${NC}"
    # Give a moment for resources to be released
    sleep 2
fi

# Step 2: Start the server
echo "Step 2: Starting the server..."
"$SCRIPT_DIR/start.sh"
START_EXIT_CODE=$?

if [ $START_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}Mallon Legal Assistant v2 successfully restarted!${NC}"
    echo "Version: v2.5.28 (Go SSR + HTMX)"
    echo "Features: Dynamic document processing (Task 8), document editing"
    echo "Access at: http://localhost:8080"
else
    echo "Failed to restart the v2 server. Please check error messages above."
    exit $START_EXIT_CODE
fi
