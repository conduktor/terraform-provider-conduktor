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

# Gateway container may fail to start without a license (version > 3.15.0).
# In that case we just warn and continue — gateway tests will be skipped.
GATEWAY_CONTAINER_ID=$(docker container ls -a --filter "name=conduktor-gateway" --filter "status=running" --format "{{.ID}}")
if [ -z "$GATEWAY_CONTAINER_ID" ]; then
    echo "WARNING: Conduktor Gateway container is not running. Gateway tests will be skipped."
else
    echo "Waiting for the Conduktor Gateway container $GATEWAY_CONTAINER_ID to be ready..."
    RETRIES=0
    MAX_RETRIES=30
    while [ "$RETRIES" -lt "$MAX_RETRIES" ]; do
        STATUS=$(docker inspect -f {{.State.Health.Status}} "$GATEWAY_CONTAINER_ID" 2>/dev/null || echo "exited")
        if [ "$STATUS" = "healthy" ]; then
            echo "Conduktor Gateway container is ready!"
            break
        elif [ "$STATUS" = "unhealthy" ] || [ "$STATUS" = "exited" ]; then
            echo "WARNING: Conduktor Gateway container failed to start (status: $STATUS). Gateway tests will be skipped."
            break
        fi
        sleep 2;
        printf "."
        RETRIES=$((RETRIES + 1))
    done
    if [ "$RETRIES" -ge "$MAX_RETRIES" ]; then
        echo "WARNING: Conduktor Gateway container did not become healthy in time. Gateway tests will be skipped."
    fi
fi
