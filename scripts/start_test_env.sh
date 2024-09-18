#!/usr/bin/env bash

set -eu
SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)


echo "CONDUKTOR_CONSOLE_IMAGE: $CONDUKTOR_CONSOLE_IMAGE"
echo "CONDUKTOR_CONSOLE_CORTEX_IMAGE: $CONDUKTOR_CONSOLE_CORTEX_IMAGE"
echo "CDK_BASE_URL: $CDK_BASE_URL"
echo "CDK_ADMIN_EMAIL: $CDK_ADMIN_EMAIL"
echo "CDK_ADMIN_PASSWORD: $CDK_ADMIN_PASSWORD"

if [ -n "${CDK_LICENSE:-}" ]; then
    echo "CDK_LICENSE is set"
else
    echo "CDK_LICENSE is not set"
fi

docker compose -f ${SCRIPT_DIR}/../docker-compose.yaml up -d
