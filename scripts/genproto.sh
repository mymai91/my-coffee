#!/usr/bin/env bash
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.." && pwd)

command -v buf >/dev/null 2>&1 || { echo "buf not found; please install from https://docs.buf.build/installation" >&2; exit 1; }

echo "Generating protos with buf..."
cd "$ROOT"
buf generate
echo "buf generate completed"
