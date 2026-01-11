#!/usr/bin/env bash
set -euo pipefail

CONFIG_FILE=$1
PID_FILE=".backend_pids"
LOGS_DIR="logs"

> "$PID_FILE"

echo "Starting backend servers from $CONFIG_FILE"

DB_MASTER_PORT=$(yq -r '.db.replicas[] | select(.mode == "master") | .port' "$CONFIG_FILE")
DB_SLAVE_PORT=$(yq -r '.db.replicas[] | select(.mode == "slave") | .port' "$CONFIG_FILE")

if [[ -z "$DB_MASTER_PORT" || -z "$DB_SLAVE_PORT" ]]; then
    echo "❌ DB master or slave port not found in config"
    exit 1
fi

yq -r '
.service.replicas[]
| .group as $group
| .pods[]
| "\($group) \(.port) \(.mode)"
' "$CONFIG_FILE" | while read -r group port mode; do
    echo "→ group=$group port=$port mode=$mode"

    if [[ "$mode" == "ro" ]]; then
        POSTGRES_PORT="$DB_SLAVE_PORT"
    else
        POSTGRES_PORT="$DB_MASTER_PORT"
    fi

    SERVER_CLUSTER="lab-backend" \
    SERVER_GROUP="$group" \
    SERVER_PORT="$port" \
    SERVER_MODE="$mode" \
    POSTGRES_PORT="$POSTGRES_PORT" \
    nohup go run cmd/service/main.go \
        > "${LOGS_DIR}/server_${group}_${port}.log" 2>&1 &

    pid=$!
    echo "$pid" >> "$PID_FILE"

    echo "  started pid=$pid (POSTGRES_PORT=$POSTGRES_PORT)"
done

echo "All servers started."
