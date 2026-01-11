#!/usr/bin/env bash
set -euo pipefail

PID_FILE=".backend_pids"

if [[ ! -f "$PID_FILE" ]]; then
    echo "No PID file found"
    exit 0
fi

echo "Stopping backend servers..."

while read -r pid; do
    if kill -0 "$pid" 2>/dev/null; then
        echo "  stopping pid=$pid"
        kill "$pid"

        timeout=10
        while kill -0 "$pid" 2>/dev/null && [ $timeout -gt 0 ]; do
            sleep 0.5
            timeout=$((timeout - 1))
        done

        if kill -0 "$pid" 2>/dev/null; then
            echo "  forcing kill pid=$pid"
            kill -9 "$pid"
        fi
    fi
done < "$PID_FILE"

rm -f "$PID_FILE"

echo "Done."
