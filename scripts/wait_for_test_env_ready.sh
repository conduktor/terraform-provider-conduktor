#!/usr/bin/env bash

set -eu

# get the container name
CONSOLE_CONTAINER_ID=$(docker container ls -a --filter "name=conduktor-console" --filter "status=running" --format "{{.ID}}")

echo "Waiting for the Conduktor Console container $CONSOLE_CONTAINER_ID to be ready..."
until [ "$(docker inspect -f {{.State.Health.Status}} "$CONSOLE_CONTAINER_ID")" == "healthy" ]; do
    sleep 1;
    printf "."
done;
echo "Conduktor Console container is ready!"

# get the container name
GATEWAY_CONTAINER_ID=$(docker container ls -a --filter "name=conduktor-gateway" --filter "status=running" --format "{{.ID}}")

echo "Waiting for the Conduktor Console container $GATEWAY_CONTAINER_ID to be ready..."
until [ "$(docker inspect -f {{.State.Health.Status}} "$GATEWAY_CONTAINER_ID")" == "healthy" ]; do
    sleep 1;
    printf "."
done;
echo "Conduktor Gateway container is ready!"
