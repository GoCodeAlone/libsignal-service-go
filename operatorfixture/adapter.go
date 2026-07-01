package operatorfixture

import (
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/fake"
	"github.com/GoCodeAlone/libsignal-service-go/service"
	"github.com/GoCodeAlone/libsignal-service-go/servicepolicy"
)

type Config struct {
	Endpoint string
	Now      time.Time
}

func NewAdapter(cfg Config) (*service.Adapter, error) {
	now := cfg.Now
	if now.IsZero() {
		now = ApprovalTime()
	}
	endpoint := cfg.Endpoint
	if endpoint == "" {
		endpoint = "127.0.0.1:19091"
	}
	return service.NewAdapter(fake.NewAdapter(), service.AdapterConfig{
		Mode:             service.AdapterModeLive,
		Endpoint:         endpoint,
		ApprovalPackage:  ApprovalPackage(now, endpoint),
		ApprovalTime:     now,
		RequestedActions: []servicepolicy.Action{servicepolicy.ActionProductionEgress},
	})
}

func ApprovalTime() time.Time {
	return time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
}

func ApprovalPackage(now time.Time, endpoint string) servicepolicy.ApprovalPackage {
	return servicepolicy.ApprovalPackage{
		OperatorApproval: servicepolicy.OperatorApproval{
			ID:        "fixture-approval",
			Scope:     "local-operator-fixture",
			GrantedAt: now.Add(-time.Hour),
			ExpiresAt: now.Add(time.Hour),
		},
		ServiceAuthorization: servicepolicy.ServiceAuthorization{
			Type:        servicepolicy.ServiceAuthorizationOfficialTestEndpoint,
			EvidenceRef: "evidence://operator-fixture",
			ExpiresAt:   now.Add(time.Hour),
		},
		AccountConsent: servicepolicy.AccountConsent{
			AccountRef:  "aci:fixture",
			EvidenceRef: "consent://operator-fixture",
			ExpiresAt:   now.Add(time.Hour),
		},
		CustodyPolicy: servicepolicy.CustodyPolicy{
			Backend:      "host-secret-file",
			KeyHandleRef: "key://fixture/account",
			BackupRef:    "backup://fixture/account",
			RotationRef:  "rotation://fixture",
		},
		AbusePolicy: servicepolicy.AbusePolicy{
			IdempotencyRequired:   true,
			RateLimitRef:          "rate://fixture",
			RecipientAllowlistRef: "allowlist://fixture",
			ChallengePolicyRef:    "challenge://fixture",
			BackoffPolicyRef:      "backoff://fixture",
		},
		EgressPolicy: servicepolicy.EgressPolicy{
			EndpointAllowlist: []string{endpoint},
			TLSPolicyRef:      "tls://fixture",
		},
		AuditPolicy: servicepolicy.AuditPolicy{
			AuditRef:     "audit://operator-fixture",
			RetentionRef: "retention://fixture",
			RedactionRef: "redaction://fixture",
		},
	}
}
