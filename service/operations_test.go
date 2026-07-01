package service

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestOperationEnvelopeRequiresCommonMetadata(t *testing.T) {
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	operations := []Operation{
		OperationRegister,
		OperationLinkDevice,
		OperationSend,
		OperationReceive,
		OperationChallenge,
		OperationUsername,
		OperationBackup,
		OperationSVR,
	}
	for _, operation := range operations {
		t.Run(string(operation), func(t *testing.T) {
			env := OperationEnvelope{
				OperationID:    "op-" + string(operation),
				Operation:      operation,
				IdempotencyKey: NewIdempotencyKey(operation, "account-ref", now),
				AccountRef:     "secret://signal/account/alice",
				RequestedAt:    now,
			}
			if operation == OperationLinkDevice {
				env.LinkedDevice = validLinkedDeviceEnvelope(now)
			}
			if err := env.Validate(); err != nil {
				t.Fatalf("valid envelope: %v", err)
			}
			for _, tc := range []struct {
				name string
				mut  func(*OperationEnvelope)
				err  error
			}{
				{name: "operation_id", mut: func(e *OperationEnvelope) { e.OperationID = "" }, err: ErrMissingOperationID},
				{name: "idempotency_key", mut: func(e *OperationEnvelope) { e.IdempotencyKey = "" }, err: ErrMissingIdempotencyKey},
				{name: "account_ref", mut: func(e *OperationEnvelope) { e.AccountRef = "" }, err: ErrMissingAccountRef},
				{name: "requested_at", mut: func(e *OperationEnvelope) { e.RequestedAt = time.Time{} }, err: ErrMissingRequestedAt},
			} {
				t.Run(tc.name, func(t *testing.T) {
					broken := env
					tc.mut(&broken)
					if err := broken.Validate(); !errors.Is(err, tc.err) {
						t.Fatalf("error = %v, want %v", err, tc.err)
					}
				})
			}
		})
	}
}

func TestLinkedDeviceEnvelopeRequiresConsentAndRevocationMetadata(t *testing.T) {
	now := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	linked := *validLinkedDeviceEnvelope(now)
	env := OperationEnvelope{
		OperationID:    "op-link",
		Operation:      OperationLinkDevice,
		IdempotencyKey: NewIdempotencyKey(OperationLinkDevice, "account-ref", now),
		AccountRef:     "secret://signal/account/alice",
		RequestedAt:    now,
		LinkedDevice:   &linked,
	}
	if err := env.Validate(); err != nil {
		t.Fatalf("valid linked-device envelope: %v", err)
	}
	for _, tc := range []struct {
		name string
		mut  func(*LinkedDeviceEnvelope)
		err  error
	}{
		{name: "device_display_name", mut: func(e *LinkedDeviceEnvelope) { e.DeviceDisplayName = "" }, err: ErrMissingDeviceDisplayName},
		{name: "consent_ref", mut: func(e *LinkedDeviceEnvelope) { e.ConsentRef = "" }, err: ErrMissingConsentRef},
		{name: "consent_expires_at", mut: func(e *LinkedDeviceEnvelope) { e.ConsentExpiresAt = time.Time{} }, err: ErrMissingConsentExpiry},
		{name: "revocation_uri", mut: func(e *LinkedDeviceEnvelope) { e.RevocationURI = "" }, err: ErrMissingRevocationURI},
		{name: "unlink_proof_ref", mut: func(e *LinkedDeviceEnvelope) { e.UnlinkProofRef = "" }, err: ErrMissingUnlinkProofRef},
	} {
		t.Run(tc.name, func(t *testing.T) {
			brokenLinked := linked
			tc.mut(&brokenLinked)
			broken := env
			broken.LinkedDevice = &brokenLinked
			if err := broken.Validate(); !errors.Is(err, tc.err) {
				t.Fatalf("error = %v, want %v", err, tc.err)
			}
		})
	}
}

func validLinkedDeviceEnvelope(now time.Time) *LinkedDeviceEnvelope {
	return &LinkedDeviceEnvelope{
		DeviceDisplayName: "Alice laptop",
		ConsentRef:        "audit://consent/alice/link-device",
		ConsentExpiresAt:  now.Add(time.Hour),
		RevocationURI:     "https://operator.example.invalid/signal/devices/revoke/alice-laptop",
		UnlinkProofRef:    "proof://signal/unlink/alice-laptop",
	}
}

func TestAuditMetadataRedactsAccountIdentifiersAndRejectsSensitiveFields(t *testing.T) {
	audit, err := NewAuditMetadata("secret://signal/account/alice", map[string]string{
		"operation_id": "op-send",
		"result":       "accepted",
	})
	if err != nil {
		t.Fatal(err)
	}
	if audit.AccountHash == "" {
		t.Fatal("account hash is empty")
	}
	if strings.Contains(audit.AccountHash, "alice") || strings.Contains(audit.AccountHash, "secret://") {
		t.Fatalf("account hash leaks account identifier: %q", audit.AccountHash)
	}
	fields := audit.CopyFields()
	if fields["operation_id"] != "op-send" {
		t.Fatalf("audit fields = %#v", fields)
	}
	fields["operation_id"] = "mutated"
	if audit.CopyFields()["operation_id"] == "mutated" {
		t.Fatal("audit metadata exposes mutable field backing storage")
	}
	for _, field := range []string{"message_body", "body", "phone_number", "e164"} {
		t.Run(field, func(t *testing.T) {
			_, err := NewAuditMetadata("secret://signal/account/alice", map[string]string{field: "sensitive"})
			if !errors.Is(err, ErrForbiddenAuditField) {
				t.Fatalf("error = %v, want %v", err, ErrForbiddenAuditField)
			}
		})
	}
}
