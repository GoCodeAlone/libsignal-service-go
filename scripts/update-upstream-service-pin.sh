#!/usr/bin/env bash
set -euo pipefail

check_only=0
if [[ "${1:-}" == "--check-only" ]]; then
  check_only=1
fi

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

latest="$(gh release view --repo signalapp/libsignal --json tagName --jq .tagName)"
pinned="$(go run ./internal/upstream/cmd/manifestgen --print-tag)"

if [[ "$latest" == "$pinned" ]]; then
  echo "upstream service pin current: ${pinned}"
  exit 0
fi

if [[ "$check_only" == "1" ]]; then
  echo "upstream service pin drift: pinned=${pinned} latest=${latest}" >&2
  exit 1
fi

bash scripts/fetch-upstream-protos.sh "$latest"
buf build
buf generate
go test ./... -count=1
