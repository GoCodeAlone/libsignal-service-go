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

Downstream packages should use `servicemetadata.Current()` for the public
compatibility baseline. The internal manifest remains an implementation detail;
the public baseline exposes only the upstream tag, descriptor checksum,
manifest digest, selected service domains, and live actions blocked by policy.

Copied upstream files preserve their original license/header comments. Signal
Messenger sources are AGPL-3.0-only; bundled Google RPC protos retain their
Apache-2.0 notices.

## Phase 2B Boundary

This package does not log in to, register with, send to, receive from, or
otherwise contact the official Signal service. Code paths that would create a
live service transport return `servicepolicy.ErrLiveServiceDisabled`.

`servicepolicy.EvaluateCompliance` reports which live official-service actions
remain blocked and which approvals a future design would need before any live
transport can be enabled. Phase 2C keeps every live transport disabled.

## Phase 2D Service Client Contracts

`serviceclient.Client` defines deterministic official-service request contracts
for registration, linked devices, send, receive, username reservation, backup
upload/download, and challenge responses. These contracts are for test doubles
and future policy-gated transports; this package still ships no live official
Signal endpoint client.

Every request carries an account ref, device ref, idempotency key, request
timestamp, consent evidence ref, audit ref, credential ref, and optional
non-exportable key handle ref. Responses return request/status metadata,
challenge refs, and secret refs. Private key material, credentials, backup
keys, and challenge responses remain host-managed secrets and must not be
returned as ordinary output values.

`servicepolicy.Policy.AllowsLiveTransport` is approval-aware. It returns true
only for live mode with every required approval identifier present, but no live
transport implementation exists in this phase.
