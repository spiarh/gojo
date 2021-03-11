#!/usr/bin/env bash

set -euo pipefail

GO111MODULE=on go mod vendor
GO111MODULE=on go mod tidy
