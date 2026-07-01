package service

import (
	"errors"
	"time"
)

type Operation string

const (
	OperationRegister   Operation = "register"
	OperationLinkDevice Operation = "link_device"
	OperationSend       Operation = "send"
	OperationReceive    Operation = "receive"
	OperationChallenge  Operation = "challenge"
	OperationUsername   Operation = "username"
	OperationBackup     Operation = "backup"
	OperationSVR        Operation = "svr"
)

var (
	ErrMissingOperationID       = errors.New("missing operation_id")
	ErrMissingOperation         = errors.New("missing operation")
	ErrUnsupportedOperation     = errors.New("unsupported operation")
	ErrMissingIdempotencyKey    = errors.New("missing idempotency_key")
	ErrMissingAccountRef        = errors.New("missing account_ref")
	ErrMissingRequestedAt       = errors.New("missing requested_at")
	ErrMissingLinkedDevice      = errors.New("missing linked-device envelope")
	ErrMissingDeviceDisplayName = errors.New("missing device_display_name")
	ErrMissingConsentRef        = errors.New("missing consent_ref")
	ErrMissingConsentExpiry     = errors.New("missing consent_expires_at")
	ErrMissingRevocationURI     = errors.New("missing revocation_uri")
	ErrMissingUnlinkProofRef    = errors.New("missing unlink_proof_ref")
)

type OperationEnvelope struct {
	OperationID    string
	Operation      Operation
	IdempotencyKey string
	AccountRef     string
	RequestedAt    time.Time
	LinkedDevice   *LinkedDeviceEnvelope
	Audit          AuditMetadata
}

type LinkedDeviceEnvelope struct {
	DeviceDisplayName string
	ConsentRef        string
	ConsentExpiresAt  time.Time
	RevocationURI     string
	UnlinkProofRef    string
}

func (e OperationEnvelope) Validate() error {
	if e.OperationID == "" {
		return ErrMissingOperationID
	}
	if e.Operation == "" {
		return ErrMissingOperation
	}
	if !e.Operation.Valid() {
		return ErrUnsupportedOperation
	}
	if e.IdempotencyKey == "" {
		return ErrMissingIdempotencyKey
	}
	if e.AccountRef == "" {
		return ErrMissingAccountRef
	}
	if e.RequestedAt.IsZero() {
		return ErrMissingRequestedAt
	}
	if e.Operation == OperationLinkDevice {
		if e.LinkedDevice == nil {
			return ErrMissingLinkedDevice
		}
		if err := e.LinkedDevice.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (e LinkedDeviceEnvelope) Validate() error {
	if e.DeviceDisplayName == "" {
		return ErrMissingDeviceDisplayName
	}
	if e.ConsentRef == "" {
		return ErrMissingConsentRef
	}
	if e.ConsentExpiresAt.IsZero() {
		return ErrMissingConsentExpiry
	}
	if e.RevocationURI == "" {
		return ErrMissingRevocationURI
	}
	if e.UnlinkProofRef == "" {
		return ErrMissingUnlinkProofRef
	}
	return nil
}

func (o Operation) Valid() bool {
	switch o {
	case OperationRegister,
		OperationLinkDevice,
		OperationSend,
		OperationReceive,
		OperationChallenge,
		OperationUsername,
		OperationBackup,
		OperationSVR:
		return true
	default:
		return false
	}
}
