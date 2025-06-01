#!/bin/bash

# Usage check
if [ $# -eq 0 ]; then
  echo "Usage: $0 [salloc options] (e.g. --nodes=2 --ntasks-per-node=2 --time=01:00:00)"
  exit 1
fi

# Temp file to capture output
TMP_LOG=$(mktemp)

# Launch salloc in background with --no-shell and capture output
salloc --no-shell --quiet "$@" >"$TMP_LOG" 2>&1 &
SALLOC_PID=$!

# Give salloc a moment to print its output
sleep 1

# Extract Job ID from log file
JOB_ID=$(grep -oP 'Granted job allocation \K[0-9]+' "$TMP_LOG")

# If not found, try detecting recent pending job for this user with no-shell name
if [ -z "$JOB_ID" ]; then
  JOB_ID=$(squeue --user="$USER" --name=no-shell --states=PD,R --format=%A --noheader | sort -n | tail -n1)
fi

# Kill the background salloc process if still running (avoid hanging)
if ps -p $SALLOC_PID > /dev/null 2>&1; then
  kill $SALLOC_PID >/dev/null 2>&1
  wait $SALLOC_PID 2>/dev/null
fi

# Clean up log
rm -f "$TMP_LOG"

# Validate job ID
if [ -z "$JOB_ID" ]; then
  echo "Failed to detect Job ID from salloc or squeue."
  exit 1
fi

# Query job info
JOB_INFO=$(scontrol show job "$JOB_ID")

# Extract fields safely
JOB_STATE=$(echo "$JOB_INFO" | grep -oP '\bJobState=\K\S+')
NODE_LIST=$(echo "$JOB_INFO" | grep -oP '(^|\s)NodeList=\K\S+')
BATCH_HOST=$(echo "$JOB_INFO" | grep -oP '(^|\s)BatchHost=\K\S+')

# Output
echo "Job ID: $JOB_ID"
echo "State: ${JOB_STATE:-N/A}"
echo "Node List: ${NODE_LIST:-(not allocated yet)}"
echo "BatchHost: ${BATCH_HOST:-(not assigned yet)}"
