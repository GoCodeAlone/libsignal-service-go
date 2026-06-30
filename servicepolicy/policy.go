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

func (p Policy) AllowsLiveTransport() bool {
	return false
}

