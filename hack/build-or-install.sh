#!/usr/bin/env bash

set -euo pipefail

CURRENT_DIR=$(dirname "$0")
PROJECT_ROOT="${CURRENT_DIR}"/..

if [[ $EFFECTIVE_VERSION == "" ]]; then
  EFFECTIVE_VERSION=$(cat "$PROJECT_ROOT/VERSION")
fi

if [[ $GO_ACTION == "" ]]; then
  GO_ACTION="build"
fi

TREE_STATE=""
STATUS="$(git status --porcelain 2>/dev/null)"
if [[ "$STATUS" == "" ]]; then
    TREE_STATE="clean"
else
    TREE_STATE="dirty"
fi    

echo ">>> $GO_ACTION $EFFECTIVE_VERSION"

CGO_ENABLED=0 GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) GO111MODULE=on \
  go "$GO_ACTION" -mod=vendor \
  -ldflags "-X github.com/lcavajani/gojo/pkg/version.GitVersion=$EFFECTIVE_VERSION \
            -X github.com/lcavajani/gojo/pkg/version.gitTreeState=$TREE_STATE \
            -X github.com/lcavajani/gojo/pkg/version.gitCommit=$(git rev-parse --verify HEAD) \
            -X github.com/lcavajani/gojo/pkg/version.buildDate=$(date --rfc-3339=seconds | sed 's/ /T/')"
