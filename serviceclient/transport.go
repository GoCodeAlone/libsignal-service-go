package serviceclient

import (
	"errors"
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/servicepolicy"
)

var (
	ErrMissingTransportClient   = errors.New("missing Signal service transport client")
	ErrMissingSandboxEndpoint   = errors.New("missing Signal service sandbox endpoint")
	ErrUnsupportedTransportMode = errors.New("unsupported Signal service transport mode")
)

type TransportMode string

const (
	TransportModeFake    TransportMode = "fake"
	TransportModeSandbox TransportMode = "sandbox"
	TransportModeLive    TransportMode = "live"
)

type TransportConfig struct {
	Mode             TransportMode
	SandboxEndpoint  string
	ApprovalPackage  servicepolicy.ApprovalPackage
	ApprovalTime     time.Time
	RequestedActions []servicepolicy.Action
}

type Transport struct {
	Client
	mode     TransportMode
	endpoint string
}

func NewTransport(client Client, cfg TransportConfig) (*Transport, error) {
	if client == nil {
		return nil, ErrMissingTransportClient
	}
	mode := cfg.Mode
	if mode == "" {
		mode = TransportModeFake
	}
	switch mode {
	case TransportModeFake:
		return &Transport{Client: client, mode: mode}, nil
	case TransportModeSandbox:
		if cfg.SandboxEndpoint == "" {
			return nil, ErrMissingSandboxEndpoint
		}
		return &Transport{Client: client, mode: mode, endpoint: cfg.SandboxEndpoint}, nil
	case TransportModeLive:
		if !(servicepolicy.Policy{Mode: servicepolicy.ModeLive}).AllowsLiveTransport(
			cfg.ApprovalPackage,
			cfg.ApprovalTime,
			cfg.RequestedActions...,
		) {
			return nil, servicepolicy.ErrLiveServiceDisabled
		}
		return &Transport{Client: client, mode: mode}, nil
	default:
		return nil, ErrUnsupportedTransportMode
	}
}

func (t *Transport) Mode() TransportMode {
	if t == nil {
		return ""
	}
	return t.mode
}

func (t *Transport) Endpoint() string {
	if t == nil {
		return ""
	}
	return t.endpoint
}
