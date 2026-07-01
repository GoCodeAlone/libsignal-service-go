package service

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/servicepolicy"
)

var (
	ErrMissingAdapterTransport = errors.New("missing Signal service operation transport")
	ErrUnsupportedAdapterMode  = errors.New("unsupported Signal service adapter mode")
	ErrEndpointNotAllowed      = errors.New("endpoint outside approval allowlist")
)

type AdapterMode string

const (
	AdapterModeFake    AdapterMode = "fake"
	AdapterModeSandbox AdapterMode = "sandbox"
	AdapterModeLive    AdapterMode = "live"
)

type AdapterConfig struct {
	Mode             AdapterMode
	Endpoint         string
	ApprovalPackage  servicepolicy.ApprovalPackage
	ApprovalTime     time.Time
	RequestedActions []servicepolicy.Action
}

type Adapter struct {
	transport OperationTransport
	mode      AdapterMode
	endpoint  string
}

func NewAdapter(transport OperationTransport, cfg AdapterConfig) (*Adapter, error) {
	if transport == nil {
		return nil, ErrMissingAdapterTransport
	}
	mode := cfg.Mode
	if mode == "" {
		mode = AdapterModeFake
	}
	switch mode {
	case AdapterModeFake, AdapterModeSandbox:
		return &Adapter{transport: transport, mode: mode, endpoint: cfg.Endpoint}, nil
	case AdapterModeLive:
		if !(servicepolicy.Policy{Mode: servicepolicy.ModeLive}).AllowsApprovedLiveTransport(
			cfg.ApprovalPackage,
			cfg.ApprovalTime,
			cfg.RequestedActions...,
		) {
			return nil, servicepolicy.ErrLiveServiceDisabled
		}
		if cfg.Endpoint == "" || !slices.Contains(cfg.ApprovalPackage.EgressPolicy.EndpointAllowlist, cfg.Endpoint) {
			return nil, ErrEndpointNotAllowed
		}
		return &Adapter{transport: transport, mode: mode, endpoint: cfg.Endpoint}, nil
	default:
		return nil, ErrUnsupportedAdapterMode
	}
}

func (a *Adapter) SubmitOperation(ctx context.Context, env OperationEnvelope) (OperationResult, error) {
	if err := env.Validate(); err != nil {
		return OperationResult{}, err
	}
	return a.transport.SubmitOperation(ctx, env)
}

func (a *Adapter) Mode() AdapterMode {
	if a == nil {
		return ""
	}
	return a.mode
}

func (a *Adapter) Endpoint() string {
	if a == nil {
		return ""
	}
	return a.endpoint
}

func ActionForOperation(operation Operation) servicepolicy.Action {
	switch operation {
	case OperationRegister:
		return servicepolicy.ActionRegister
	case OperationLinkDevice:
		return servicepolicy.ActionLinkedDevice
	case OperationSend:
		return servicepolicy.ActionSend
	case OperationReceive:
		return servicepolicy.ActionReceive
	case OperationUsername:
		return servicepolicy.ActionUsernameReserve
	case OperationBackup:
		return servicepolicy.ActionBackupUpload
	default:
		return servicepolicy.ActionProductionEgress
	}
}
