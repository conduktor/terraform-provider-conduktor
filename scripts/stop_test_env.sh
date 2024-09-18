#!/usr/bin/env bash

set -eu
SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)

echo "Extract logs before stopping the containers"
mkdir -p ./logs
docker compose -f ${SCRIPT_DIR}/../docker-compose.yaml logs > ./logs/docker-compose.log

echo "Stopping the containers and removing volumes"
docker compose -f ${SCRIPT_DIR}/../docker-compose.yaml down -v
