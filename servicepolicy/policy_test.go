package servicepolicy

import (
	"errors"
	"testing"
	"time"
)

func TestPolicyValidateAllowsDisabledAndTestDouble(t *testing.T) {
	for _, mode := range []Mode{"", ModeDisabled, ModeTestDouble} {
		if err := (Policy{Mode: mode}).Validate(); err != nil {
			t.Fatalf("mode %q: %v", mode, err)
		}
	}
}

func TestPolicyValidateRejectsLiveMode(t *testing.T) {
	err := (Policy{Mode: ModeLive}).Validate()
	if !errors.Is(err, ErrLiveServiceDisabled) {
		t.Fatalf("live mode error = %v, want %v", err, ErrLiveServiceDisabled)
	}
}

func TestPolicyNeverAllowsLiveTransport(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	for _, mode := range []Mode{"", ModeDisabled, ModeTestDouble, ModeLive, "other"} {
		if (Policy{Mode: mode}).AllowsLiveTransport(ApprovalPackage{}, now) {
			t.Fatalf("mode %q unexpectedly allowed live transport", mode)
		}
	}
}

func TestPolicyAllowsLiveTransportOnlyWithApprovalPackage(t *testing.T) {
	now := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	approval := validApprovalPackage(now)
	if !(Policy{Mode: ModeLive}).AllowsLiveTransport(approval, now, ActionSend) {
		t.Fatal("live mode with complete approval package should be policy-ready for a future live transport")
	}
	approval.ServiceAuthorization.EvidenceRef = ""
	if (Policy{Mode: ModeLive}).AllowsLiveTransport(approval, now, ActionSend) {
		t.Fatal("live mode with incomplete approval package allowed live transport")
	}
	if (Policy{Mode: ModeTestDouble}).AllowsLiveTransport(validApprovalPackage(now), now, ActionSend) {
		t.Fatal("test-double mode must not allow live transport even with approval package")
	}
}

func TestComplianceBlocksLiveActions(t *testing.T) {
	actions := []Action{
		ActionRegister,
		ActionLogin,
		ActionLinkedDevice,
		ActionSend,
		ActionReceive,
		ActionBackupUpload,
		ActionBackupDownload,
		ActionUsernameReserve,
		ActionProductionEgress,
	}
	for _, mode := range []Mode{"", ModeDisabled, ModeTestDouble, ModeLive} {
		report := EvaluateCompliance(ComplianceRequest{Mode: mode, RequestedActions: actions})
		if report.Approved {
			t.Fatalf("mode %q approved live actions: %+v", mode, report)
		}
		if !report.LiveServiceDisabled {
			t.Fatalf("mode %q did not report live service disabled", mode)
		}
		if len(report.BlockedActions) != len(actions) {
			t.Fatalf("mode %q blocked actions = %v, want %v", mode, report.BlockedActions, actions)
		}
		for _, approval := range []string{
			"operator_live_service_approval",
			"legal_tos_review",
			"account_owner_consent",
			"abuse_rate_limit_plan",
			"credential_custody_plan",
			"audit_retention_plan",
			"egress_allowlist",
		} {
			if !containsString(report.RequiredApprovals, approval) {
				t.Fatalf("mode %q missing approval %q in %v", mode, approval, report.RequiredApprovals)
			}
		}
	}
}

func TestComplianceAllowsNoActionsForDisabledBoundary(t *testing.T) {
	report := EvaluateCompliance(ComplianceRequest{Mode: ModeDisabled})
	if !report.Approved {
		t.Fatalf("disabled/no-action report unexpectedly denied: %+v", report)
	}
	if !report.LiveServiceDisabled {
		t.Fatal("disabled/no-action report must still state live service is disabled")
	}
	if len(report.BlockedActions) != 0 {
		t.Fatalf("blocked actions = %v, want empty", report.BlockedActions)
	}
}

func TestComplianceRejectsUnsupportedMode(t *testing.T) {
	report := EvaluateCompliance(ComplianceRequest{Mode: Mode("unsupported")})
	if report.Approved {
		t.Fatalf("unsupported mode approved: %+v", report)
	}
	if !report.LiveServiceDisabled {
		t.Fatal("unsupported mode must still keep live service disabled")
	}
}

func containsString(values []string, want string) bool {
	for _, got := range values {
		if got == want {
			return true
		}
	}
	return false
}
