#!/usr/bin/env bash

set -eu
SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)

CLI_VERSION=$(cat "$SCRIPT_DIR/../go.mod" | grep "github.com/conduktor/ctl" | awk '{print $2}')
echo "Using CLI version: $CLI_VERSION"

CONSOLE_URL=${CDK_BASE_URL}
CONSOLE_USER=${CDK_ADMIN_EMAIL}
CONSOLE_PASSWORD=${CDK_ADMIN_PASSWORD}

GW_URL=${CDK_GATEWAY_BASE_URL}
GW_USER=${CDK_GATEWAY_USER}
GW_PASSWORD=${CDK_GATEWAY_PASSWORD}

# re-export as CDK_USER and CDK_PASSWORD for CLI
printenv | grep CDK_ | grep -v CDK_LICENSE

# unset all CDK_ variables
unset CDK_BASE_URL
unset CDK_ADMIN_EMAIL
unset CDK_ADMIN_PASSWORD
unset CDK_GATEWAY_BASE_URL
unset CDK_GATEWAY_USER
unset CDK_GATEWAY_PASSWORD

echo "Logging in Console and applying setup"
export CDK_BASE_URL=${CONSOLE_URL}
export CDK_USER=${CONSOLE_USER}
export CDK_PASSWORD=${CONSOLE_PASSWORD}
CDK_DEBUG=false go run github.com/conduktor/ctl@${CLI_VERSION} login # disable debug logs for the login
go run github.com/conduktor/ctl@${CLI_VERSION} apply -f "${SCRIPT_DIR}"/../testdata/init/init_console.yaml
if [[ "${CONDUKTOR_CONSOLE_IMAGE}" != *"1.26.0"* ]];then # only applying some resources for newer console versions
	go run github.com/conduktor/ctl@${CLI_VERSION} apply -f "${SCRIPT_DIR}"/../testdata/init/init_console_2.yaml
fi

echo "Logging in Gateway and applying setup"
unset CDK_BASE_URL
unset CDK_USER
unset CDK_PASSWORD
export CDK_GATEWAY_BASE_URL=${GW_URL}
export CDK_GATEWAY_USER=${GW_USER}
export CDK_GATEWAY_PASSWORD=${GW_PASSWORD}
go run github.com/conduktor/ctl@${CLI_VERSION} apply -f "${SCRIPT_DIR}"/../testdata/init/init_gateway.yaml
