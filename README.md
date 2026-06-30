# libsignal-service-go

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)

Go service-boundary packages for Signal-compatible Workflow integrations.

This repository is a community compatibility project for selected
`signalapp/libsignal` service wire artifacts. It is not affiliated with,
endorsed by, or maintained by Signal Messenger LLC.

The first development milestone is intentionally limited to generated
protobuf/gRPC types, deterministic fake transports, and policy gates that keep
live official Signal service access disabled.

## Upstream Baseline

The copied protobuf sources are pinned to
[`signalapp/libsignal` `v0.96.4`](https://github.com/signalapp/libsignal/tree/v0.96.4).
`internal/upstream/manifest.json` records every upstream source path, upstream
blob SHA, local SHA-256, and the generated descriptor checksum.

Copied upstream files preserve their original license/header comments. Signal
Messenger sources are AGPL-3.0-only; bundled Google RPC protos retain their
Apache-2.0 notices.

## Phase 2B Boundary

This package does not log in to, register with, send to, receive from, or
otherwise contact the official Signal service. Code paths that would create a
live service transport return `servicepolicy.ErrLiveServiceDisabled`.

