#!/bin/bash
set -e

echo "host replication replicator 0.0.0.0/0 trust" >> "$PGDATA/pg_hba.conf"
