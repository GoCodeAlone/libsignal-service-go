# libsignal-service-go

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)

Go service-boundary packages for Signal-compatible Workflow integrations.

This repository is a community compatibility project for selected
`signalapp/libsignal` service wire artifacts. It is not affiliated with,
endorsed by, or maintained by Signal Messenger LLC.

The current public surface includes generated protobuf/gRPC types, typed
operation envelopes, deterministic fake and sandbox transports, approval-aware
policy gates, and a local operator fixture. It still does not ship an official
Signal production endpoint client.

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

## Service Boundary

This package does not log in to, register with, send to, receive from, or
otherwise contact the official Signal service by default. Code paths that
hard-disable live mode or lack a complete machine-checkable approval package
return `servicepolicy.ErrLiveServiceDisabled`.

`servicepolicy.EvaluateCompliance` reports which live official-service actions
remain blocked and which approval identifiers are required for the requested
mode/actions. Approval-package completeness is checked by
`servicepolicy.ValidateApprovalPackage`.

## Service Client Contracts

`serviceclient.Client` defines deterministic official-service request contracts
for registration, linked devices, send, receive, username reservation, backup
upload/download, and challenge responses. These contracts are for test doubles
and host-supplied policy-gated transports; this package still ships no live
official Signal endpoint client.

Every request carries an account ref, device ref, idempotency key, request
timestamp, consent evidence ref, audit ref, credential ref, and optional
non-exportable key handle ref. Responses return request/status metadata,
challenge refs, and secret refs. Private key material, credentials, backup
keys, and challenge responses remain host-managed secrets and must not be
returned as ordinary output values.

`servicepolicy.Policy.AllowsLiveTransport` is approval-package aware. It returns
true only for live mode with a machine-checkable approval package that includes
operator approval, supported service authorization, account-owner consent,
custody, abuse/rate-limit, egress allowlist, idempotency, and audit policy
metadata. Human/operator evidence is recorded, but narrative evidence alone does
not enable live transport. No official Signal endpoint constant exists in this
repository.

## Operation Envelopes

The `service` package exposes typed operation envelopes for registration,
linked-device preparation, send, receive, challenge, username, backup, and SVR
request flows. Envelopes require `operation_id`, `idempotency_key`,
`account_ref`, and `requested_at` before any transport can submit them.

Linked-device envelopes additionally require a display name, consent reference,
consent expiry, revocation URI, and unlink proof reference. Audit metadata stores
a redacted account hash and rejects message-body and phone-number fields.

## Operation Adapters

`service.NewAdapter` wraps an operation transport with mode-specific validation.
Fake and sandbox modes are deterministic test paths. Live mode remains
approval-gated and requires a machine-checkable approval package, account-owner
consent, custody policy, abuse/rate-limit policy, audit policy, and endpoint
allowlist entry before any operation can be submitted.

The `fake` package accepts every named operation and records redacted audit
metadata. The `sandbox` package requires an explicit sandbox endpoint and
rejects official-looking endpoints unless a test-only override is supplied. The
`operatorfixture` package exercises the same live adapter path against a local
allowlisted endpoint; it is a conformance fixture, not an official Signal service
client.

This repository still does not compile official production endpoint constants or
perform official Signal service egress.
