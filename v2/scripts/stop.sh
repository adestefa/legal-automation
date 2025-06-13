#!/bin/bash

# Stop the Legal Document Automation Server (v1 and v2)
echo "Stopping Mallon Legal Document Automation Server..."
echo "---------------------------------------------------"

# Define color codes for status messages
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Port the server is running on
SERVER_PORT=8080

# Check for both v1 and v2 server processes
SERVER_PIDS=""

# Check for v2 mallon_v2 processes
if command -v pgrep &> /dev/null; then
    V2_PIDS=$(pgrep -f "mallon_v2" 2>/dev/null)
    if [ -n "$V2_PIDS" ]; then
        SERVER_PIDS="$SERVER_PIDS $V2_PIDS"
        echo "Found v2 mallon_v2 processes: $V2_PIDS"
    fi
    
    # Check for v1 main processes (in backend directory)
    V1_PIDS=$(pgrep -f "main.go" 2>/dev/null)
    if [ -n "$V1_PIDS" ]; then
        SERVER_PIDS="$SERVER_PIDS $V1_PIDS"
        echo "Found v1 main.go processes: $V1_PIDS"
    fi
fi

# Also check by port
if command -v lsof &> /dev/null; then
    # Use lsof if available (macOS, many Linux distros)
    PORT_PIDS=$(lsof -ti:${SERVER_PORT} 2>/dev/null)
    if [ -n "$PORT_PIDS" ]; then
        # Add any port-specific PIDs that aren't already in our list
        for pid in $PORT_PIDS; do
            if [[ ! " $SERVER_PIDS " =~ " $pid " ]]; then
                SERVER_PIDS="$SERVER_PIDS $pid"
            fi
        done
    fi
elif command -v netstat &> /dev/null; then
    # Fallback to netstat (older Linux distros)
    PORT_PIDS=$(netstat -tulnp 2>/dev/null | grep ":${SERVER_PORT}" | awk '{print $7}' | cut -d'/' -f1 | sort -u)
    if [ -n "$PORT_PIDS" ]; then
        for pid in $PORT_PIDS; do
            if [[ ! " $SERVER_PIDS " =~ " $pid " ]]; then
                SERVER_PIDS="$SERVER_PIDS $pid"
            fi
        done
    fi
elif command -v ss &> /dev/null; then
    # Another fallback using ss (newer Linux distros)
    PORT_PIDS=$(ss -tulnp 2>/dev/null | grep ":${SERVER_PORT}" | awk '{print $7}' | cut -d',' -f2 | cut -d'=' -f2 | sort -u)
    if [ -n "$PORT_PIDS" ]; then
        for pid in $PORT_PIDS; do
            if [[ ! " $SERVER_PIDS " =~ " $pid " ]]; then
                SERVER_PIDS="$SERVER_PIDS $pid"
            fi
        done
    fi
fi

# Clean up SERVER_PIDS (remove extra spaces)
SERVER_PIDS=$(echo $SERVER_PIDS | tr -s ' ' | sed 's/^ *//' | sed 's/ *$//')

# Check if a PID file exists (enhancement - not yet implemented in start.sh)
PID_FILE="$(dirname "$0")/.server.pid"
if [ -f "$PID_FILE" ]; then
    PID_FROM_FILE=$(cat "$PID_FILE")
    if [ -n "$PID_FROM_FILE" ] && kill -0 "$PID_FROM_FILE" 2>/dev/null; then
        # Add the PID from file to our list if it's not already there
        if [[ ! " $SERVER_PIDS " =~ " $PID_FROM_FILE " ]]; then
            SERVER_PIDS="$SERVER_PIDS $PID_FROM_FILE"
        fi
    else
        echo -e "${YELLOW}Warning: PID file exists but process is not running. Removing stale PID file.${NC}"
        rm "$PID_FILE"
    fi
fi

# Check if we found any processes to kill
if [ -z "$SERVER_PIDS" ]; then
    echo -e "${YELLOW}No server process found running on port ${SERVER_PORT}.${NC}"
    echo "Either the server is not running or it's running on a different port."
    exit 0
fi

# Kill the processes
echo "Found server processes: $SERVER_PIDS"
echo "Sending termination signal..."

for PID in $SERVER_PIDS; do
    # First try graceful termination with SIGTERM
    if kill -15 "$PID" 2>/dev/null; then
        echo -e "Sent SIGTERM to process ${PID}."
    else
        echo -e "${YELLOW}Failed to send SIGTERM to process ${PID}. Process may have already terminated.${NC}"
    fi
done

# Wait a moment for processes to terminate
sleep 2

# Check if any processes are still running
REMAINING_PIDS=""
for PID in $SERVER_PIDS; do
    if kill -0 "$PID" 2>/dev/null; then
        REMAINING_PIDS="$REMAINING_PIDS $PID"
    fi
done

# If processes are still running, try SIGKILL
if [ -n "$REMAINING_PIDS" ]; then
    echo -e "${YELLOW}Some processes did not terminate gracefully. Forcing termination...${NC}"
    for PID in $REMAINING_PIDS; do
        if kill -9 "$PID" 2>/dev/null; then
            echo -e "Sent SIGKILL to process ${PID}."
        else
            echo -e "${RED}Failed to terminate process ${PID}.${NC}"
        fi
    done
    
    # Check again if all processes are terminated
    sleep 1
    for PID in $REMAINING_PIDS; do
        if kill -0 "$PID" 2>/dev/null; then
            echo -e "${RED}Process ${PID} is still running. Please terminate it manually.${NC}"
            TERMINATION_FAILED=true
        fi
    done
fi

# Check if the port is now free
sleep 1
PORT_CHECK_FAILED=false
if command -v lsof &> /dev/null; then
    if lsof -i:${SERVER_PORT} &>/dev/null; then
        PORT_CHECK_FAILED=true
    fi
elif command -v netstat &> /dev/null; then
    if netstat -tuln 2>/dev/null | grep ":${SERVER_PORT}" &>/dev/null; then
        PORT_CHECK_FAILED=true
    fi
elif command -v ss &> /dev/null; then
    if ss -tuln | grep ":${SERVER_PORT}" &>/dev/null; then
        PORT_CHECK_FAILED=true
    fi
fi

# Clean up PID file if it exists
if [ -f "$PID_FILE" ]; then
    rm "$PID_FILE"
fi

# Final status message
if [ "$PORT_CHECK_FAILED" = true ]; then
    echo -e "${RED}Port ${SERVER_PORT} is still in use. Not all processes were terminated.${NC}"
    echo "Please check running processes manually with: lsof -i:${SERVER_PORT} or similar command."
    exit 1
elif [ "$TERMINATION_FAILED" = true ]; then
    echo -e "${YELLOW}Warning: Some processes could not be terminated, but port ${SERVER_PORT} appears to be free.${NC}"
    echo "The server should be able to restart successfully."
    exit 0
else
    echo -e "${GREEN}Server successfully stopped!${NC}"
    echo "Port ${SERVER_PORT} is now free. You can start the server again with ./start.sh"
    exit 0
fi
