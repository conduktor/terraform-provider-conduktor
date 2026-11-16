#!/usr/bin/env bash

set -eu
SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)


echo "CONDUKTOR_CONSOLE_IMAGE: $CONDUKTOR_CONSOLE_IMAGE"
echo "CONDUKTOR_GATEWAY_IMAGE: $CONDUKTOR_GATEWAY_IMAGE"
echo "CONDUKTOR_CONSOLE_CORTEX_IMAGE: $CONDUKTOR_CONSOLE_CORTEX_IMAGE"
echo "CDK_BASE_URL: $CDK_BASE_URL"
echo "CDK_ADMIN_EMAIL: $CDK_ADMIN_EMAIL"
echo "CDK_ADMIN_PASSWORD: $CDK_ADMIN_PASSWORD"

if [ -n "${CDK_LICENSE:-}" ]; then
    echo "CDK_LICENSE is set"
else
    echo "CDK_LICENSE is not set"
fi

# Gateway 3.18 added the GATEWAY_LISTENER_* settings and 3.23 removed the
# GATEWAY_PORT_START-style ones, so older Gateways get the legacy settings.
gateway_tag=${CONDUKTOR_GATEWAY_IMAGE##*:}
if [[ $gateway_tag =~ ^[0-9]+\.[0-9]+ ]] && [ "$(printf '%s\n' 3.18.0 "$gateway_tag" | sort -V | head -1)" != "3.18.0" ]; then
    export CONDUKTOR_GATEWAY_LISTENER_ENV=legacy-listeners
fi
echo "CONDUKTOR_GATEWAY_LISTENER_ENV: ${CONDUKTOR_GATEWAY_LISTENER_ENV:-listeners}"

docker compose -f ${SCRIPT_DIR}/../docker-compose.yaml up -d
