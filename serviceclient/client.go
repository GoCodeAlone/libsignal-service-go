// Package serviceclient defines official Signal service client contracts used
// by deterministic test doubles and future policy-gated live transports.
package serviceclient

import (
	"context"
	"time"
)

type Client interface {
	Register(context.Context, RegisterRequest) (RegisterResponse, error)
	LinkDevice(context.Context, LinkDeviceRequest) (LinkDeviceResponse, error)
	Send(context.Context, SendRequest) (SendResponse, error)
	Receive(context.Context, ReceiveRequest) (ReceiveResponse, error)
	ReserveUsername(context.Context, ReserveUsernameRequest) (ReserveUsernameResponse, error)
	UploadBackup(context.Context, UploadBackupRequest) (UploadBackupResponse, error)
	DownloadBackup(context.Context, DownloadBackupRequest) (DownloadBackupResponse, error)
	RespondToChallenge(context.Context, RespondToChallengeRequest) (RespondToChallengeResponse, error)
}

type RequestMetadata struct {
	AccountRef       string
	DeviceRef        string
	IdempotencyKey   string
	RequestedAt      time.Time
	ConsentRef       string
	AuditRef         string
	CredentialRef    string
	NonExportableKey string
}

type ResponseMetadata struct {
	RequestID    string
	Status       string
	ChallengeRef string
	SecretRefs   map[string]string
}

type RegisterRequest struct {
	Metadata RequestMetadata
	Username string
}

type RegisterResponse struct {
	Metadata ResponseMetadata
}

func (r RegisterResponse) Common() ResponseMetadata { return r.Metadata }

type LinkDeviceRequest struct {
	Metadata    RequestMetadata
	LinkCodeRef string
}

type LinkDeviceResponse struct {
	Metadata ResponseMetadata
}

func (r LinkDeviceResponse) Common() ResponseMetadata { return r.Metadata }

type SendRequest struct {
	Metadata     RequestMetadata
	RecipientRef string
	PayloadRef   string
}

type SendResponse struct {
	Metadata ResponseMetadata
}

func (r SendResponse) Common() ResponseMetadata { return r.Metadata }

type ReceiveRequest struct {
	Metadata  RequestMetadata
	CursorRef string
}

type ReceiveResponse struct {
	Metadata ResponseMetadata
}

func (r ReceiveResponse) Common() ResponseMetadata { return r.Metadata }

type ReserveUsernameRequest struct {
	Metadata RequestMetadata
	Username string
}

type ReserveUsernameResponse struct {
	Metadata ResponseMetadata
}

func (r ReserveUsernameResponse) Common() ResponseMetadata { return r.Metadata }

type UploadBackupRequest struct {
	Metadata  RequestMetadata
	BackupRef string
}

type UploadBackupResponse struct {
	Metadata ResponseMetadata
}

func (r UploadBackupResponse) Common() ResponseMetadata { return r.Metadata }

type DownloadBackupRequest struct {
	Metadata RequestMetadata
	BackupID string
}

type DownloadBackupResponse struct {
	Metadata ResponseMetadata
}

func (r DownloadBackupResponse) Common() ResponseMetadata { return r.Metadata }

type RespondToChallengeRequest struct {
	Metadata     RequestMetadata
	ChallengeRef string
	ResponseRef  string
}

type RespondToChallengeResponse struct {
	Metadata ResponseMetadata
}

func (r RespondToChallengeResponse) Common() ResponseMetadata { return r.Metadata }
