#!/usr/bin/env bash

set -euo pipefail

pushd "$GOENV_ROOT"
printf '\nUpdating goenv to %s...\n' "${GOENV_TOOL_VERSION}"
git pull origin "${GOENV_TOOL_VERSION}"
popd

go_version=$(goenv local 2>/dev/null)
echo "Local Go version: ${go_version}"

if [ -z "$GO_VERSION" ]; then
  echo "\$GO_VERSION: ${GO_VERSION}"
  go_version=$GO_VERSION
fi

echo "Installing Go version $go_version..."
goenv install --skip-existing "$go_version" && goenv rehash
