#!/usr/bin/env bash

set -euo pipefail

GO_MODULE_NAME="github.com/lcavajani/gojo"
GOLANGCI_LINT_CONFIG_FILE=""

for arg in "$@"; do
  case $arg in
    --golangci-lint-config=*)
    GOLANGCI_LINT_CONFIG_FILE="-c ${arg#*=}"
    shift
    ;;
  esac
done

echo ">>> Check"

echo ">>> Executing golangci-lint"
# shellcheck disable=SC2068,SC2086
golangci-lint run $GOLANGCI_LINT_CONFIG_FILE $@

echo ">>> Executing go vet"
# shellcheck disable=SC2068
go vet -mod=vendor $@

echo ">>> Executing gofmt"
folders=()
# shellcheck disable=SC2068
for f in $@; do
  folders+=( "$(echo $f | sed 's/\(.*\)\/\.\.\./\1/')" )
done
# shellcheck disable=SC2086
unformatted_files="$(goimports -l -local="$GO_MODULE_NAME" ${folders[*]})"
if [[ "$unformatted_files" ]]; then
  echo ">>> Unformatted files detected:"
  echo "$unformatted_files"
  exit 1
fi

echo ">>> Executing goconst"
goconst cmd/ pkg/

echo ">>> Executing shellcheck"
shellcheck hack/*.sh

echo ">>> All checks successful"
