#!/usr/bin/env bash

set -euo pipefail

echo ">>> Format"

# shellcheck disable=SC2068
goimports -l -w -local=github.com/lcavajani/gojo $@
