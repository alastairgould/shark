#!/usr/bin/env bash
#
# Starts the Pack Calculator server.
#
# Edit the values below to change the configuration, then run ./run.sh
# Any variable already set in your shell takes precedence over these defaults.
#
set -euo pipefail

# Port the HTTP server listens on.
export PORT="${PORT:-8080}"

# Comma-separated pack sizes, largest to smallest is not required.
export PACK_SIZES="${PACK_SIZES:-250,500,1000,2000,5000}"

# Maximum quantity a single request may ask for.
export MAX_QUANTITY="${MAX_QUANTITY:-1000000}"

echo "Building..."
go build -o gymshark .

echo "Starting Pack Calculator"
echo "  PORT         = $PORT"
echo "  PACK_SIZES   = $PACK_SIZES"
echo "  MAX_QUANTITY = $MAX_QUANTITY"
echo "  UI           = http://localhost:$PORT"
echo

exec ./gymshark
