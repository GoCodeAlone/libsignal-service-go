#!/usr/bin/env bash
set -euo pipefail

tag="${1:-v0.96.4}"
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

files=(
  rust/net/grpc/proto/TextSecure.proto
  rust/net/grpc/proto/google/rpc/error_details.proto
  rust/net/grpc/proto/google/rpc/status.proto
  rust/net/grpc/proto/org/signal/chat/account.proto
  rust/net/grpc/proto/org/signal/chat/backups.proto
  rust/net/grpc/proto/org/signal/chat/challenge.proto
  rust/net/grpc/proto/org/signal/chat/common.proto
  rust/net/grpc/proto/org/signal/chat/credentials.proto
  rust/net/grpc/proto/org/signal/chat/device.proto
  rust/net/grpc/proto/org/signal/chat/donations.proto
  rust/net/grpc/proto/org/signal/chat/errors.proto
  rust/net/grpc/proto/org/signal/chat/keys.proto
  rust/net/grpc/proto/org/signal/chat/messages.proto
  rust/net/grpc/proto/org/signal/chat/profile.proto
  rust/net/grpc/proto/org/signal/chat/require.proto
  rust/net/grpc/proto/org/signal/chat/subscriptions.proto
  rust/net/grpc/proto/org/signal/chat/tag.proto
  rust/net/src/proto/chat_websocket.proto
)

rm -rf proto

for upstream in "${files[@]}"; do
  case "$upstream" in
    rust/net/grpc/proto/*) local_path="proto/signal/net/grpc/${upstream#rust/net/grpc/proto/}" ;;
    rust/net/src/proto/*) local_path="proto/signal/net/src/${upstream#rust/net/src/proto/}" ;;
    *) echo "unsupported upstream path: $upstream" >&2; exit 1 ;;
  esac

  mkdir -p "$(dirname "$local_path")"
  gh api "repos/signalapp/libsignal/contents/${upstream}?ref=${tag}" --jq .content | base64 -d > "$local_path"
done

go run ./internal/upstream/cmd/manifestgen --tag "$tag" --out internal/upstream/manifest.json
