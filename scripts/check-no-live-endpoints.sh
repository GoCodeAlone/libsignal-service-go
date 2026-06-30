#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

if rg -n --glob '!docs/**' --glob '!README.md' --glob '!internal/upstream/manifest.json' \
  'textsecure-service\.whispersystems\.org|chat\.signal\.org|storage\.signal\.org|signal\.org/v1' .; then
  echo "live Signal endpoint constant found in runnable code" >&2
  exit 1
fi

