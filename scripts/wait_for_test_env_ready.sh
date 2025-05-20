#!/usr/bin/env bash

set -eu

# get the container name
CONSOLE_CONTAINER_ID=$(docker container ls -a --filter "name=conduktor-console" --filter "status=running" --format "{{.ID}}")
if [ -z "$CONSOLE_CONTAINER_ID" ]; then
    echo "Conduktor Console container not found. Exiting."
    exit 1
fi

echo "Waiting for the Conduktor Console container $CONSOLE_CONTAINER_ID to be ready..."
until [ "$(docker inspect -f {{.State.Health.Status}} "$CONSOLE_CONTAINER_ID")" == "healthy" ]; do
    sleep 1;
    printf "."
done;
echo "Conduktor Console container is ready!"

# get the container name
GATEWAY_CONTAINER_ID=$(docker container ls -a --filter "name=conduktor-gateway" --filter "status=running" --format "{{.ID}}")
if [ -z "$GATEWAY_CONTAINER_ID" ]; then
    echo "Conduktor Gateway container not found. Exiting."
    exit 1
fi

echo "Waiting for the Conduktor Gateway container $GATEWAY_CONTAINER_ID to be ready..."
until [ "$(docker inspect -f {{.State.Health.Status}} "$GATEWAY_CONTAINER_ID")" == "healthy" ]; do
    sleep 1;
    printf "."
done;
echo "Conduktor Gateway container is ready!"
