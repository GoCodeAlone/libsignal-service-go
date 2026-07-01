package servicemetadata

import "testing"

func TestCurrentBaseline(t *testing.T) {
	b := Current()
	if b.UpstreamTag != "v0.96.4" {
		t.Fatalf("upstream tag = %q, want v0.96.4", b.UpstreamTag)
	}
	if b.DescriptorChecksum != "203ff86981ca5b249cd0a373296f1754303e070fd5341225a3b4d8995f6c2286" {
		t.Fatalf("descriptor checksum = %q", b.DescriptorChecksum)
	}
	if b.ManifestDigest != "e6d117566fe76ce537709b45bbdfb08a148d89b4c6e7273e5401f7ea1f72ca08" {
		t.Fatalf("manifest digest = %q", b.ManifestDigest)
	}
	for _, domain := range []string{"account", "device", "messages", "profile", "keys", "backups_metadata", "challenge", "credentials", "donations_metadata", "subscriptions_metadata"} {
		if !b.HasDomain(domain) {
			t.Fatalf("baseline missing domain %q in %v", domain, b.SelectedDomains)
		}
	}
	for _, action := range []string{"register", "login", "linked_device", "send", "receive", "backup_upload", "backup_download", "username_reserve", "production_egress"} {
		if !contains(b.BlockedLiveActions, action) {
			t.Fatalf("baseline missing blocked action %q in %v", action, b.BlockedLiveActions)
		}
	}
}

func TestCurrentReturnsIndependentSlices(t *testing.T) {
	b := Current()
	b.SelectedDomains[0] = "mutated"
	b.BlockedLiveActions[0] = "mutated"

	next := Current()
	if next.SelectedDomains[0] == "mutated" {
		t.Fatal("SelectedDomains shares backing storage")
	}
	if next.BlockedLiveActions[0] == "mutated" {
		t.Fatal("BlockedLiveActions shares backing storage")
	}
}

func contains(values []string, want string) bool {
	for _, got := range values {
		if got == want {
			return true
		}
	}
	return false
}
