// Package servicemetadata exposes release-safe Signal service compatibility
// metadata for Workflow plugins and downstream compatibility checks.
package servicemetadata

// Baseline describes the upstream service wire baseline carried by this module.
type Baseline struct {
	UpstreamTag        string
	DescriptorChecksum string
	ManifestDigest     string
	SelectedDomains    []string
	BlockedLiveActions []string
}

var selectedDomains = []string{
	"account",
	"device",
	"messages",
	"profile",
	"keys",
	"backups_metadata",
	"challenge",
	"credentials",
}

var blockedLiveActions = []string{
	"register",
	"login",
	"linked_device",
	"send",
	"receive",
	"backup_upload",
	"backup_download",
	"username_reserve",
	"production_egress",
}

// Current returns the service compatibility baseline for this release.
func Current() Baseline {
	return Baseline{
		UpstreamTag:        "v0.96.4",
		DescriptorChecksum: "9f647ca4a75f581514cbe080c792871e10d7dbd7b22bd6faf2832e15d447e484",
		ManifestDigest:     "5973a782b6ce4836f4588c3e797a1070eaac04669f758cb82572b967cfcc0b60",
		SelectedDomains:    append([]string(nil), selectedDomains...),
		BlockedLiveActions: append([]string(nil), blockedLiveActions...),
	}
}

// HasDomain reports whether the baseline includes the named service domain.
func (b Baseline) HasDomain(name string) bool {
	for _, domain := range b.SelectedDomains {
		if domain == name {
			return true
		}
	}
	return false
}
