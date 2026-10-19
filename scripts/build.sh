#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
export PATH="${PATH}:/usr/local/go/bin"
export GOPROXY="${GOPROXY:-https://goproxy.cn,direct}"
CGO_ENABLED=0 go build -o "$ROOT/gitlab-code-reader-mcp" ./cmd/server/
echo "Built $ROOT/gitlab-code-reader-mcp"
