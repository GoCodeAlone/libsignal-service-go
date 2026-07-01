package servicepolicy

import (
	"slices"
	"time"
)

type ServiceAuthorizationType string

const (
	ServiceAuthorizationNone                 ServiceAuthorizationType = "none"
	ServiceAuthorizationSignalPermission     ServiceAuthorizationType = "signal-written-permission"
	ServiceAuthorizationOfficialTestEndpoint ServiceAuthorizationType = "official-test-endpoint"
	ServiceAuthorizationThrowawayAccount     ServiceAuthorizationType = "throwaway-owned-account"
)

var supportedServiceAuthorizations = []ServiceAuthorizationType{
	ServiceAuthorizationSignalPermission,
	ServiceAuthorizationOfficialTestEndpoint,
	ServiceAuthorizationThrowawayAccount,
}

type ApprovalPackage struct {
	OperatorApproval     OperatorApproval
	ServiceAuthorization ServiceAuthorization
	AccountConsent       AccountConsent
	CustodyPolicy        CustodyPolicy
	AbusePolicy          AbusePolicy
	EgressPolicy         EgressPolicy
	AuditPolicy          AuditPolicy
}

type OperatorApproval struct {
	ID        string
	Scope     string
	GrantedAt time.Time
	ExpiresAt time.Time
}

type ServiceAuthorization struct {
	Type        ServiceAuthorizationType
	EvidenceRef string
	ExpiresAt   time.Time
}

type AccountConsent struct {
	AccountRef  string
	EvidenceRef string
	ExpiresAt   time.Time
}

type CustodyPolicy struct {
	Backend      string
	KeyHandleRef string
	BackupRef    string
	RotationRef  string
}

type AbusePolicy struct {
	IdempotencyRequired   bool
	RateLimitRef          string
	RecipientAllowlistRef string
	DeclaredAudienceRef   string
	ChallengePolicyRef    string
	BackoffPolicyRef      string
}

type EgressPolicy struct {
	EndpointAllowlist []string
	TLSPolicyRef      string
	DryRun            bool
}

type AuditPolicy struct {
	AuditRef     string
	RetentionRef string
	RedactionRef string
}

type ApprovalValidationReport struct {
	LiveAllowed   bool
	DeniedReasons []string
}

func ValidateApprovalPackage(pkg ApprovalPackage, now time.Time, requested []Action) ApprovalValidationReport {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	var reasons []string
	reasons = validateOperatorApproval(reasons, pkg.OperatorApproval, now)
	reasons = validateServiceAuthorization(reasons, pkg.ServiceAuthorization, now)
	reasons = validateAccountConsent(reasons, pkg.AccountConsent, now)
	reasons = validateCustodyPolicy(reasons, pkg.CustodyPolicy)
	reasons = validateAbusePolicy(reasons, pkg.AbusePolicy, requested)
	reasons = validateEgressPolicy(reasons, pkg.EgressPolicy)
	reasons = validateAuditPolicy(reasons, pkg.AuditPolicy)
	return ApprovalValidationReport{
		LiveAllowed:   len(reasons) == 0,
		DeniedReasons: reasons,
	}
}

func (p Policy) AllowsApprovedLiveTransport(pkg ApprovalPackage, now time.Time, requested ...Action) bool {
	return p.AllowsLiveTransport(pkg, now, requested...)
}

func validateOperatorApproval(reasons []string, approval OperatorApproval, now time.Time) []string {
	if approval.ID == "" {
		reasons = append(reasons, "operator_approval_id_missing")
	}
	if approval.Scope == "" {
		reasons = append(reasons, "operator_approval_scope_missing")
	}
	if approval.ExpiresAt.IsZero() {
		return append(reasons, "operator_approval_expiry_missing")
	}
	if !approval.ExpiresAt.After(now) {
		reasons = append(reasons, "operator_approval_expired")
	}
	return reasons
}

func validateServiceAuthorization(reasons []string, auth ServiceAuthorization, now time.Time) []string {
	if !slices.Contains(supportedServiceAuthorizations, auth.Type) {
		reasons = append(reasons, "service_authorization_unsupported")
	}
	if auth.EvidenceRef == "" {
		reasons = append(reasons, "service_authorization_evidence_missing")
	}
	if auth.ExpiresAt.IsZero() {
		return append(reasons, "service_authorization_expiry_missing")
	}
	if !auth.ExpiresAt.After(now) {
		reasons = append(reasons, "service_authorization_expired")
	}
	return reasons
}

func validateAccountConsent(reasons []string, consent AccountConsent, now time.Time) []string {
	if consent.AccountRef == "" {
		reasons = append(reasons, "account_consent_account_missing")
	}
	if consent.EvidenceRef == "" {
		reasons = append(reasons, "account_consent_evidence_missing")
	}
	if consent.ExpiresAt.IsZero() {
		return append(reasons, "account_consent_expiry_missing")
	}
	if !consent.ExpiresAt.After(now) {
		reasons = append(reasons, "account_consent_expired")
	}
	return reasons
}

func validateCustodyPolicy(reasons []string, custody CustodyPolicy) []string {
	if custody.Backend == "" {
		reasons = append(reasons, "custody_backend_missing")
	}
	if custody.KeyHandleRef == "" {
		reasons = append(reasons, "custody_key_handle_missing")
	}
	if custody.BackupRef == "" {
		reasons = append(reasons, "custody_backup_missing")
	}
	if custody.RotationRef == "" {
		reasons = append(reasons, "custody_rotation_missing")
	}
	return reasons
}

func validateAbusePolicy(reasons []string, abuse AbusePolicy, requested []Action) []string {
	if !abuse.IdempotencyRequired {
		reasons = append(reasons, "abuse_idempotency_required")
	}
	if abuse.RateLimitRef == "" {
		reasons = append(reasons, "abuse_rate_limit_missing")
	}
	if abuse.ChallengePolicyRef == "" {
		reasons = append(reasons, "abuse_challenge_policy_missing")
	}
	if abuse.BackoffPolicyRef == "" {
		reasons = append(reasons, "abuse_backoff_policy_missing")
	}
	if needsAudiencePolicy(requested) && abuse.RecipientAllowlistRef == "" && abuse.DeclaredAudienceRef == "" {
		reasons = append(reasons, "abuse_audience_missing")
	}
	return reasons
}

func validateEgressPolicy(reasons []string, egress EgressPolicy) []string {
	if len(egress.EndpointAllowlist) == 0 {
		reasons = append(reasons, "egress_allowlist_missing")
	}
	if egress.TLSPolicyRef == "" {
		reasons = append(reasons, "egress_tls_policy_missing")
	}
	return reasons
}

func validateAuditPolicy(reasons []string, audit AuditPolicy) []string {
	if audit.AuditRef == "" {
		reasons = append(reasons, "audit_ref_missing")
	}
	if audit.RetentionRef == "" {
		reasons = append(reasons, "audit_retention_missing")
	}
	if audit.RedactionRef == "" {
		reasons = append(reasons, "audit_redaction_missing")
	}
	return reasons
}

func needsAudiencePolicy(requested []Action) bool {
	for _, action := range requested {
		switch action {
		case ActionSend, ActionReceive, ActionRegister, ActionLogin, ActionLinkedDevice, ActionProductionEgress:
			return true
		}
	}
	return false
}
