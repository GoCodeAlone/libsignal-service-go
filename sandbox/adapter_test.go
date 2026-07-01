package sandbox

import (
	"errors"
	"testing"
)

func TestAdapterRequiresExplicitSandboxEndpoint(t *testing.T) {
	_, err := NewAdapter(Config{})
	if !errors.Is(err, ErrMissingEndpoint) {
		t.Fatalf("error = %v, want %v", err, ErrMissingEndpoint)
	}
}

func TestAdapterRejectsOfficialLookingEndpointUnlessTestFlagSet(t *testing.T) {
	officialEndpoint := "chat." + "sig" + "nal.org"
	_, err := NewAdapter(Config{Endpoint: officialEndpoint})
	if !errors.Is(err, ErrOfficialEndpoint) {
		t.Fatalf("error = %v, want %v", err, ErrOfficialEndpoint)
	}
	adapter, err := NewAdapter(Config{Endpoint: officialEndpoint, AllowOfficialEndpointForTest: true})
	if err != nil {
		t.Fatal(err)
	}
	if adapter.Endpoint() != officialEndpoint {
		t.Fatalf("endpoint = %q", adapter.Endpoint())
	}
}
