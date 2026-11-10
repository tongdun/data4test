#!/bin/sh
# biz initializes config/i18n relative to the repository root, before TestMain.
# Usage: go test -vet=off -exec "$PWD/build/test-exec.sh" ./...
set -eu
cd "$(dirname "$0")/.."
exec "$@"
