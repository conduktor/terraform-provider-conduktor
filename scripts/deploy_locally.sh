#!/usr/bin/env bash

set -eu
SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
ARCH=$(uname -m | sed -e 's/x86_64/amd64/' -e 's/arm64/arm64/')
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH="${OS}_${ARCH}"

VERSION=${1:-0.0.1}
COMMIT=$(git rev-parse HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%S%Z)

cd "${SCRIPT_DIR}/.."

echo "Publish locally with version ${VERSION}"

echo "Build with : go build -ldflags=\"-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}\""
go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

echo "Install in ${HOME}/.terraform.d/plugins/terraform.local/conduktor/conduktor/${VERSION}/${ARCH} with name terraform-provider-conduktor_v${VERSION}"

# need to configure .terraformrc provider_installation.filesystem_mirror to point to ~/.terraform.d/plugins
mkdir -p ${HOME}/.terraform.d/plugins/terraform.local/conduktor/conduktor/${VERSION}/${ARCH}
cp terraform-provider-conduktor ${HOME}/.terraform.d/plugins/terraform.local/conduktor/conduktor/${VERSION}/${ARCH}/terraform-provider-conduktor_v${VERSION}

echo "Provider installed !"
echo
echo "Update provider inside your project with: 'terraform init -upgrade'"
