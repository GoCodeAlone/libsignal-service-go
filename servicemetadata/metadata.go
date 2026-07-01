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
	"donations_metadata",
	"subscriptions_metadata",
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
		DescriptorChecksum: "203ff86981ca5b249cd0a373296f1754303e070fd5341225a3b4d8995f6c2286",
		ManifestDigest:     "e6d117566fe76ce537709b45bbdfb08a148d89b4c6e7273e5401f7ea1f72ca08",
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
