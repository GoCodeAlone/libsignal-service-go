package servicepolicy

import (
	"slices"
	"testing"
	"time"
)

func TestValidateApprovalPackageRejectsExpiredOperatorApproval(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	approval := validApprovalPackage(now)
	approval.OperatorApproval.ExpiresAt = now.Add(-time.Minute)

	report := ValidateApprovalPackage(approval, now, []Action{ActionProductionEgress})
	if report.LiveAllowed {
		t.Fatalf("expired approval allowed live transport: %+v", report)
	}
	if !slices.Contains(report.DeniedReasons, "operator_approval_expired") {
		t.Fatalf("denied reasons = %v, want operator_approval_expired", report.DeniedReasons)
	}
}

func TestValidateApprovalPackageRejectsMissingOrUnknownAuthorization(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	for _, authType := range []ServiceAuthorizationType{
		ServiceAuthorizationNone,
		ServiceAuthorizationType("browser-scrape"),
		"",
	} {
		approval := validApprovalPackage(now)
		approval.ServiceAuthorization.Type = authType

		report := ValidateApprovalPackage(approval, now, []Action{ActionProductionEgress})
		if report.LiveAllowed {
			t.Fatalf("authorization %q allowed live transport: %+v", authType, report)
		}
		if !slices.Contains(report.DeniedReasons, "service_authorization_unsupported") {
			t.Fatalf("authorization %q denied reasons = %v, want service_authorization_unsupported", authType, report.DeniedReasons)
		}
	}
}

func TestValidateApprovalPackageRequiresSendAbuseControls(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	approval := validApprovalPackage(now)
	approval.AbusePolicy.IdempotencyRequired = false
	approval.AbusePolicy.RateLimitRef = ""
	approval.AbusePolicy.RecipientAllowlistRef = ""
	approval.AbusePolicy.DeclaredAudienceRef = ""

	report := ValidateApprovalPackage(approval, now, []Action{ActionSend})
	if report.LiveAllowed {
		t.Fatalf("send without abuse controls allowed live transport: %+v", report)
	}
	for _, reason := range []string{
		"abuse_idempotency_required",
		"abuse_rate_limit_missing",
		"abuse_audience_missing",
	} {
		if !slices.Contains(report.DeniedReasons, reason) {
			t.Fatalf("denied reasons = %v, want %s", report.DeniedReasons, reason)
		}
	}
}

func TestValidateApprovalPackageAllowsCompleteApproval(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	report := ValidateApprovalPackage(validApprovalPackage(now), now, []Action{ActionSend})
	if !report.LiveAllowed {
		t.Fatalf("complete approval denied live transport: %+v", report)
	}
	if len(report.DeniedReasons) != 0 {
		t.Fatalf("denied reasons = %v, want empty", report.DeniedReasons)
	}
}

func TestPolicyAllowsApprovedLiveTransportOnlyForLiveMode(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	approval := validApprovalPackage(now)
	if !((Policy{Mode: ModeLive}).AllowsApprovedLiveTransport(approval, now, ActionSend)) {
		t.Fatal("live mode with complete approval should allow approved live transport")
	}
	for _, mode := range []Mode{"", ModeDisabled, ModeTestDouble, "unsupported"} {
		if (Policy{Mode: mode}).AllowsApprovedLiveTransport(approval, now, ActionSend) {
			t.Fatalf("mode %q unexpectedly allowed approved live transport", mode)
		}
	}
}

func validApprovalPackage(now time.Time) ApprovalPackage {
	return ApprovalPackage{
		OperatorApproval: OperatorApproval{
			ID:        "approval-123",
			Scope:     "signal-live-send-test",
			GrantedAt: now.Add(-time.Hour),
			ExpiresAt: now.Add(time.Hour),
		},
		ServiceAuthorization: ServiceAuthorization{
			Type:        ServiceAuthorizationOfficialTestEndpoint,
			EvidenceRef: "evidence://signal/test-endpoint",
			ExpiresAt:   now.Add(time.Hour),
		},
		AccountConsent: AccountConsent{
			AccountRef:  "aci:test",
			EvidenceRef: "consent://owner",
			ExpiresAt:   now.Add(time.Hour),
		},
		CustodyPolicy: CustodyPolicy{
			Backend:      "host-secret-file",
			KeyHandleRef: "key://signal/account",
			BackupRef:    "backup://signal/account",
			RotationRef:  "rotation://quarterly",
		},
		AbusePolicy: AbusePolicy{
			IdempotencyRequired:   true,
			RateLimitRef:          "rate://one-per-minute",
			RecipientAllowlistRef: "allowlist://owned-test-recipient",
			ChallengePolicyRef:    "challenge://respond",
			BackoffPolicyRef:      "backoff://exponential",
		},
		EgressPolicy: EgressPolicy{
			EndpointAllowlist: []string{"signal-test.example.invalid"},
			TLSPolicyRef:      "tls://pinned",
		},
		AuditPolicy: AuditPolicy{
			AuditRef:     "audit://signal-live",
			RetentionRef: "retention://30d",
			RedactionRef: "redaction://pii",
		},
	}
}
