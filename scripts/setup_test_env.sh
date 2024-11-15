#!/usr/bin/env bash

set -eu
SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)

CLI_VERSION=$(cat "$SCRIPT_DIR/../go.mod" | grep "github.com/conduktor/ctl" | awk '{print $2}')
echo "Using CLI version: $CLI_VERSION"

# re-export as CDK_USER and CDK_PASSWORD for CLI
export CDK_USER=${CDK_ADMIN_EMAIL}
export CDK_PASSWORD=${CDK_ADMIN_PASSWORD}
printenv | grep CDK_ | grep -v CDK_LICENSE


echo "Logging in and applying setup"
CDK_DEBUG=false go run github.com/conduktor/ctl@${CLI_VERSION} login # disable debug logs for the login
go run github.com/conduktor/ctl@${CLI_VERSION} apply -f "${SCRIPT_DIR}"/../testdata/init/*.yaml
