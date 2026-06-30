package servicepolicy

import "errors"

// ErrLiveServiceDisabled is returned by Phase 2B code paths that would create
// a live official Signal service transport.
var ErrLiveServiceDisabled = errors.New("live Signal service transport disabled")

type Mode string

const (
	ModeDisabled   Mode = "disabled"
	ModeTestDouble Mode = "test_double"
	ModeLive       Mode = "live"
)

type Policy struct {
	Mode Mode
}

type Action string

const (
	ActionRegister         Action = "register"
	ActionLogin            Action = "login"
	ActionLinkedDevice     Action = "linked_device"
	ActionSend             Action = "send"
	ActionReceive          Action = "receive"
	ActionBackupUpload     Action = "backup_upload"
	ActionBackupDownload   Action = "backup_download"
	ActionUsernameReserve  Action = "username_reserve"
	ActionProductionEgress Action = "production_egress"
)

type ComplianceRequest struct {
	Mode             Mode
	RequestedActions []Action
}

type ComplianceReport struct {
	Mode                Mode
	Approved            bool
	LiveServiceDisabled bool
	BlockedActions      []Action
	RequiredApprovals   []string
	DeferredDomains     []string
}

var requiredLiveApprovals = []string{
	"operator_live_service_approval",
	"legal_tos_review",
	"account_owner_consent",
	"abuse_rate_limit_plan",
	"credential_custody_plan",
	"audit_retention_plan",
	"egress_allowlist",
}

var deferredDomains = []string{
	"svr_svrb",
	"message_backup",
	"zkgroup",
	"zkcredential",
	"poksho",
	"keytrans",
}

func (p Policy) Validate() error {
	switch p.Mode {
	case "", ModeDisabled, ModeTestDouble:
		return nil
	case ModeLive:
		return ErrLiveServiceDisabled
	default:
		return errors.New("unsupported Signal service boundary mode")
	}
}

type LiveApprovalSet map[string]bool

func NewLiveApprovalSet(approvals ...string) LiveApprovalSet {
	set := LiveApprovalSet{}
	for _, approval := range approvals {
		if approval != "" {
			set[approval] = true
		}
	}
	return set
}

func RequiredLiveApprovals() []string {
	return append([]string(nil), requiredLiveApprovals...)
}

// AllowsLiveTransport reports whether policy metadata contains the approvals a
// future live transport would require. Phase 2D still ships no live transport
// implementation, and Validate continues to reject ModeLive.
func (p Policy) AllowsLiveTransport(sets ...LiveApprovalSet) bool {
	if p.Mode != ModeLive {
		return false
	}
	var approvals LiveApprovalSet
	if len(sets) > 0 {
		approvals = sets[0]
	}
	for _, required := range requiredLiveApprovals {
		if !approvals[required] {
			return false
		}
	}
	return true
}

func EvaluateCompliance(req ComplianceRequest) ComplianceReport {
	mode := req.Mode
	if mode == "" {
		mode = ModeDisabled
	}
	blocked := append([]Action(nil), req.RequestedActions...)
	validMode := mode == ModeDisabled || mode == ModeTestDouble || mode == ModeLive
	report := ComplianceReport{
		Mode:                mode,
		Approved:            validMode && len(blocked) == 0,
		LiveServiceDisabled: true,
		BlockedActions:      blocked,
		DeferredDomains:     append([]string(nil), deferredDomains...),
	}
	if len(blocked) > 0 || mode == ModeLive {
		report.Approved = false
		report.RequiredApprovals = append([]string(nil), requiredLiveApprovals...)
	}
	return report
}
