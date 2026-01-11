#!/usr/bin/env bash

set -euo pipefail

for n in {0..5}; do
  port="808${n}"
  pids=$(lsof -ti:"$port" || true)

  if [[ -n "$pids" ]]; then
    echo "Killing processes on port $port: $pids"
    kill -9 $pids
  else
    echo "No process on port $port"
  fi
done
