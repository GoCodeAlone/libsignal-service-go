package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/servicepolicy"
)

func TestAdapterLiveRejectsMissingApprovalPackage(t *testing.T) {
	_, err := NewAdapter(&recordingOperationTransport{}, AdapterConfig{
		Mode:         AdapterModeLive,
		Endpoint:     "signal-test.example.invalid",
		ApprovalTime: testAdapterTime(),
	})
	if !errors.Is(err, servicepolicy.ErrLiveServiceDisabled) {
		t.Fatalf("error = %v, want %v", err, servicepolicy.ErrLiveServiceDisabled)
	}
}

func TestAdapterLiveRejectsIncompleteApprovalPackage(t *testing.T) {
	for _, tc := range []struct {
		name string
		mut  func(*servicepolicy.ApprovalPackage)
	}{
		{name: "expired", mut: func(p *servicepolicy.ApprovalPackage) {
			p.OperatorApproval.ExpiresAt = testAdapterTime().Add(-time.Minute)
		}},
		{name: "consent", mut: func(p *servicepolicy.ApprovalPackage) { p.AccountConsent = servicepolicy.AccountConsent{} }},
		{name: "custody", mut: func(p *servicepolicy.ApprovalPackage) { p.CustodyPolicy = servicepolicy.CustodyPolicy{} }},
		{name: "abuse", mut: func(p *servicepolicy.ApprovalPackage) { p.AbusePolicy = servicepolicy.AbusePolicy{} }},
		{name: "audit", mut: func(p *servicepolicy.ApprovalPackage) { p.AuditPolicy = servicepolicy.AuditPolicy{} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			approval := testApprovalPackage(testAdapterTime(), "signal-test.example.invalid")
			tc.mut(&approval)
			_, err := NewAdapter(&recordingOperationTransport{}, AdapterConfig{
				Mode:             AdapterModeLive,
				Endpoint:         "signal-test.example.invalid",
				ApprovalPackage:  approval,
				ApprovalTime:     testAdapterTime(),
				RequestedActions: []servicepolicy.Action{servicepolicy.ActionSend},
			})
			if !errors.Is(err, servicepolicy.ErrLiveServiceDisabled) {
				t.Fatalf("error = %v, want %v", err, servicepolicy.ErrLiveServiceDisabled)
			}
		})
	}
}

func TestAdapterLiveRejectsEndpointOutsideAllowlist(t *testing.T) {
	_, err := NewAdapter(&recordingOperationTransport{}, AdapterConfig{
		Mode:             AdapterModeLive,
		Endpoint:         "outside.example.invalid",
		ApprovalPackage:  testApprovalPackage(testAdapterTime(), "signal-test.example.invalid"),
		ApprovalTime:     testAdapterTime(),
		RequestedActions: []servicepolicy.Action{servicepolicy.ActionSend},
	})
	if !errors.Is(err, ErrEndpointNotAllowed) {
		t.Fatalf("error = %v, want %v", err, ErrEndpointNotAllowed)
	}
}

func TestAdapterSubmitValidatesEnvelopeAndDelegates(t *testing.T) {
	transport := &recordingOperationTransport{}
	adapter, err := NewAdapter(transport, AdapterConfig{Mode: AdapterModeFake})
	if err != nil {
		t.Fatal(err)
	}
	env := testOperationEnvelope(OperationSend, "send-1")
	result, err := adapter.SubmitOperation(context.Background(), env)
	if err != nil {
		t.Fatal(err)
	}
	if result.OperationID != env.OperationID || result.Status != "accepted" {
		t.Fatalf("result = %#v", result)
	}
	if transport.last.OperationID != env.OperationID {
		t.Fatalf("delegated envelope = %#v", transport.last)
	}
	_, err = adapter.SubmitOperation(context.Background(), OperationEnvelope{Operation: OperationSend})
	if !errors.Is(err, ErrMissingOperationID) {
		t.Fatalf("invalid envelope error = %v, want %v", err, ErrMissingOperationID)
	}
}

type recordingOperationTransport struct {
	last OperationEnvelope
}

func (r *recordingOperationTransport) SubmitOperation(_ context.Context, env OperationEnvelope) (OperationResult, error) {
	r.last = env
	return OperationResult{OperationID: env.OperationID, Status: "accepted", Audit: env.Audit}, nil
}

func testOperationEnvelope(operation Operation, id string) OperationEnvelope {
	now := testAdapterTime()
	env := OperationEnvelope{
		OperationID:    id,
		Operation:      operation,
		IdempotencyKey: NewIdempotencyKey(operation, "secret://signal/account/alice", now),
		AccountRef:     "secret://signal/account/alice",
		RequestedAt:    now,
	}
	if operation == OperationLinkDevice {
		env.LinkedDevice = &LinkedDeviceEnvelope{
			DeviceDisplayName: "Alice laptop",
			ConsentRef:        "audit://consent/alice/link-device",
			ConsentExpiresAt:  now.Add(time.Hour),
			RevocationURI:     "https://operator.example.invalid/signal/devices/revoke/alice-laptop",
			UnlinkProofRef:    "proof://signal/unlink/alice-laptop",
		}
	}
	return env
}

func testAdapterTime() time.Time {
	return time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
}

func testApprovalPackage(now time.Time, endpoint string) servicepolicy.ApprovalPackage {
	return servicepolicy.ApprovalPackage{
		OperatorApproval: servicepolicy.OperatorApproval{
			ID:        "approval-123",
			Scope:     "signal-live-fixture",
			GrantedAt: now.Add(-time.Hour),
			ExpiresAt: now.Add(time.Hour),
		},
		ServiceAuthorization: servicepolicy.ServiceAuthorization{
			Type:        servicepolicy.ServiceAuthorizationOfficialTestEndpoint,
			EvidenceRef: "evidence://signal/test-endpoint",
			ExpiresAt:   now.Add(time.Hour),
		},
		AccountConsent: servicepolicy.AccountConsent{
			AccountRef:  "aci:test",
			EvidenceRef: "consent://owner",
			ExpiresAt:   now.Add(time.Hour),
		},
		CustodyPolicy: servicepolicy.CustodyPolicy{
			Backend:      "host-secret-file",
			KeyHandleRef: "key://signal/account",
			BackupRef:    "backup://signal/account",
			RotationRef:  "rotation://quarterly",
		},
		AbusePolicy: servicepolicy.AbusePolicy{
			IdempotencyRequired:   true,
			RateLimitRef:          "rate://one-per-minute",
			RecipientAllowlistRef: "allowlist://owned-test-recipient",
			ChallengePolicyRef:    "challenge://respond",
			BackoffPolicyRef:      "backoff://exponential",
		},
		EgressPolicy: servicepolicy.EgressPolicy{
			EndpointAllowlist: []string{endpoint},
			TLSPolicyRef:      "tls://pinned",
		},
		AuditPolicy: servicepolicy.AuditPolicy{
			AuditRef:     "audit://signal-live",
			RetentionRef: "retention://30d",
			RedactionRef: "redaction://pii",
		},
	}
}
