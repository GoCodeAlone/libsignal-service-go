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
		if (Policy{Mode: mode}).AllowsLiveTransport() {
			t.Fatalf("mode %q unexpectedly allowed live transport", mode)
		}
	}
}

