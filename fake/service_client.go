package fake

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"

	"github.com/GoCodeAlone/libsignal-service-go/serviceclient"
)

var (
	ErrMissingIdempotencyKey = errors.New("missing idempotency key")
	ErrIdempotencyConflict   = errors.New("idempotency key reused with different request")
)

type ServiceClientOption func(*ServiceClient)

func WithChallenge(operation, challengeRef string) ServiceClientOption {
	return func(c *ServiceClient) {
		c.challenges[operation] = challengeRef
	}
}

func WithExpiredCredentials() ServiceClientOption {
	return func(c *ServiceClient) {
		c.expiredCredentials = true
	}
}

func WithLedger(ledger map[string]LedgerRecord) ServiceClientOption {
	return func(c *ServiceClient) {
		c.ledger = cloneLedger(ledger)
	}
}

type ServiceClient struct {
	mu                 sync.Mutex
	ledger             map[string]LedgerRecord
	records            []ServiceClientRecord
	challenges         map[string]string
	expiredCredentials bool
}

type LedgerRecord struct {
	Operation   string
	Fingerprint string
	Response    serviceclient.ResponseMetadata
}

type ServiceClientRecord struct {
	Operation      string
	IdempotencyKey string
	RequestID      string
	Status         string
}

func NewServiceClient(opts ...ServiceClientOption) *ServiceClient {
	client := &ServiceClient{
		ledger:     map[string]LedgerRecord{},
		challenges: map[string]string{},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func (c *ServiceClient) Records() []ServiceClientRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]ServiceClientRecord, len(c.records))
	copy(out, c.records)
	return out
}

func (c *ServiceClient) Ledger() map[string]LedgerRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	return cloneLedger(c.ledger)
}

func (c *ServiceClient) Register(_ context.Context, req serviceclient.RegisterRequest) (serviceclient.RegisterResponse, error) {
	meta, err := c.handle("register", req.Metadata, req)
	return serviceclient.RegisterResponse{Metadata: meta}, err
}

func (c *ServiceClient) LinkDevice(_ context.Context, req serviceclient.LinkDeviceRequest) (serviceclient.LinkDeviceResponse, error) {
	meta, err := c.handle("linked_device", req.Metadata, req)
	return serviceclient.LinkDeviceResponse{Metadata: meta}, err
}

func (c *ServiceClient) Send(_ context.Context, req serviceclient.SendRequest) (serviceclient.SendResponse, error) {
	meta, err := c.handle("send", req.Metadata, req)
	return serviceclient.SendResponse{Metadata: meta}, err
}

func (c *ServiceClient) Receive(_ context.Context, req serviceclient.ReceiveRequest) (serviceclient.ReceiveResponse, error) {
	meta, err := c.handle("receive", req.Metadata, req)
	return serviceclient.ReceiveResponse{Metadata: meta}, err
}

func (c *ServiceClient) ReserveUsername(_ context.Context, req serviceclient.ReserveUsernameRequest) (serviceclient.ReserveUsernameResponse, error) {
	meta, err := c.handle("username_reserve", req.Metadata, req)
	return serviceclient.ReserveUsernameResponse{Metadata: meta}, err
}

func (c *ServiceClient) UploadBackup(_ context.Context, req serviceclient.UploadBackupRequest) (serviceclient.UploadBackupResponse, error) {
	meta, err := c.handle("backup_upload", req.Metadata, req)
	return serviceclient.UploadBackupResponse{Metadata: meta}, err
}

func (c *ServiceClient) DownloadBackup(_ context.Context, req serviceclient.DownloadBackupRequest) (serviceclient.DownloadBackupResponse, error) {
	meta, err := c.handle("backup_download", req.Metadata, req)
	return serviceclient.DownloadBackupResponse{Metadata: meta}, err
}

func (c *ServiceClient) RespondToChallenge(_ context.Context, req serviceclient.RespondToChallengeRequest) (serviceclient.RespondToChallengeResponse, error) {
	meta, err := c.handle("challenge_response", req.Metadata, req)
	return serviceclient.RespondToChallengeResponse{Metadata: meta}, err
}

func (c *ServiceClient) handle(operation string, metadata serviceclient.RequestMetadata, request any) (serviceclient.ResponseMetadata, error) {
	if metadata.IdempotencyKey == "" {
		return serviceclient.ResponseMetadata{}, ErrMissingIdempotencyKey
	}
	fingerprint, err := requestFingerprint(request)
	if err != nil {
		return serviceclient.ResponseMetadata{}, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if existing, ok := c.ledger[metadata.IdempotencyKey]; ok {
		if existing.Operation != operation || existing.Fingerprint != fingerprint {
			return serviceclient.ResponseMetadata{}, ErrIdempotencyConflict
		}
		return cloneResponse(existing.Response), nil
	}

	response := serviceclient.ResponseMetadata{
		RequestID:  metadata.IdempotencyKey,
		Status:     "accepted",
		SecretRefs: map[string]string{},
	}
	if c.expiredCredentials {
		response.Status = "credential_expired"
	}
	if challengeRef := c.challenges[operation]; challengeRef != "" {
		response.Status = "challenge_required"
		response.ChallengeRef = challengeRef
	}
	if metadata.CredentialRef != "" {
		response.SecretRefs["credential"] = metadata.CredentialRef
	}
	c.ledger[metadata.IdempotencyKey] = LedgerRecord{
		Operation:   operation,
		Fingerprint: fingerprint,
		Response:    cloneResponse(response),
	}
	c.records = append(c.records, ServiceClientRecord{
		Operation:      operation,
		IdempotencyKey: metadata.IdempotencyKey,
		RequestID:      response.RequestID,
		Status:         response.Status,
	})
	return cloneResponse(response), nil
}

func requestFingerprint(request any) (string, error) {
	encoded, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

func cloneLedger(in map[string]LedgerRecord) map[string]LedgerRecord {
	out := make(map[string]LedgerRecord, len(in))
	for key, record := range in {
		record.Response = cloneResponse(record.Response)
		out[key] = record
	}
	return out
}

func cloneResponse(in serviceclient.ResponseMetadata) serviceclient.ResponseMetadata {
	out := in
	if in.SecretRefs != nil {
		out.SecretRefs = make(map[string]string, len(in.SecretRefs))
		for key, value := range in.SecretRefs {
			out.SecretRefs[key] = value
		}
	}
	return out
}
