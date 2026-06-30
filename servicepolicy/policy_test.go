package servicepolicy

import (
	"errors"
	"testing"
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
	for _, mode := range []Mode{"", ModeDisabled, ModeTestDouble, ModeLive, "other"} {
		if (Policy{Mode: mode}).AllowsLiveTransport(LiveApprovalSet{}) {
			t.Fatalf("mode %q unexpectedly allowed live transport", mode)
		}
	}
}

func TestPolicyAllowsLiveTransportOnlyWithEveryApproval(t *testing.T) {
	approvals := NewLiveApprovalSet(RequiredLiveApprovals()...)
	if !(Policy{Mode: ModeLive}).AllowsLiveTransport(approvals) {
		t.Fatal("live mode with every required approval should be policy-ready for a future live transport")
	}
	for _, missing := range RequiredLiveApprovals() {
		partial := NewLiveApprovalSet(RequiredLiveApprovals()...)
		delete(partial, missing)
		if (Policy{Mode: ModeLive}).AllowsLiveTransport(partial) {
			t.Fatalf("live mode with missing approval %q allowed live transport", missing)
		}
	}
	if (Policy{Mode: ModeTestDouble}).AllowsLiveTransport(approvals) {
		t.Fatal("test-double mode must not allow live transport even with approvals")
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
