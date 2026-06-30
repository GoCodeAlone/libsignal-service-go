package serviceclient

import (
	"context"
	"testing"
	"time"
)

func TestClientContractsRequireRequestMetadata(t *testing.T) {
	now := time.Unix(1_782_847_600, 0).UTC()
	metadata := RequestMetadata{
		AccountRef:       "secret://signal/account/alice",
		DeviceRef:        "secret://signal/device/alice-primary",
		IdempotencyKey:   "request-123",
		RequestedAt:      now,
		ConsentRef:       "audit://consent/alice",
		AuditRef:         "audit://signal/request-123",
		CredentialRef:    "secret://signal/credential/alice",
		NonExportableKey: "kms://signal/alice/identity",
	}

	requests := []RequestMetadata{
		RegisterRequest{Metadata: metadata, Username: "alice.01"}.Metadata,
		LinkDeviceRequest{Metadata: metadata, LinkCodeRef: "secret://signal/link-code"}.Metadata,
		SendRequest{Metadata: metadata, RecipientRef: "signal://recipient/bob", PayloadRef: "ciphertext://payload/1"}.Metadata,
		ReceiveRequest{Metadata: metadata, CursorRef: "cursor://inbox/alice"}.Metadata,
		ReserveUsernameRequest{Metadata: metadata, Username: "alice.01"}.Metadata,
		UploadBackupRequest{Metadata: metadata, BackupRef: "backup://alice/1"}.Metadata,
		DownloadBackupRequest{Metadata: metadata, BackupID: "backup-1"}.Metadata,
		RespondToChallengeRequest{Metadata: metadata, ChallengeRef: "challenge://alice/1", ResponseRef: "secret://challenge/response"}.Metadata,
	}
	for _, got := range requests {
		if got != metadata {
			t.Fatalf("request metadata = %#v, want %#v", got, metadata)
		}
	}
}

func TestClientInterfaceCoversOfficialServiceTestDoubleDomains(t *testing.T) {
	var client Client = recordingClient{}
	ctx := context.Background()
	metadata := RequestMetadata{
		AccountRef:     "secret://signal/account/alice",
		DeviceRef:      "secret://signal/device/alice-primary",
		IdempotencyKey: "request-123",
		RequestedAt:    time.Unix(1_782_847_600, 0).UTC(),
	}

	register, err := client.Register(ctx, RegisterRequest{Metadata: metadata})
	checkResponse(t, register, err)
	linkDevice, err := client.LinkDevice(ctx, LinkDeviceRequest{Metadata: metadata})
	checkResponse(t, linkDevice, err)
	send, err := client.Send(ctx, SendRequest{Metadata: metadata})
	checkResponse(t, send, err)
	receive, err := client.Receive(ctx, ReceiveRequest{Metadata: metadata})
	checkResponse(t, receive, err)
	username, err := client.ReserveUsername(ctx, ReserveUsernameRequest{Metadata: metadata})
	checkResponse(t, username, err)
	upload, err := client.UploadBackup(ctx, UploadBackupRequest{Metadata: metadata})
	checkResponse(t, upload, err)
	download, err := client.DownloadBackup(ctx, DownloadBackupRequest{Metadata: metadata})
	checkResponse(t, download, err)
	challenge, err := client.RespondToChallenge(ctx, RespondToChallengeRequest{Metadata: metadata})
	checkResponse(t, challenge, err)
}

func checkResponse[T interface{ Common() ResponseMetadata }](t *testing.T, response T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	meta := response.Common()
	if meta.RequestID == "" || meta.Status == "" {
		t.Fatalf("response metadata = %#v, want request id and status", meta)
	}
}

type recordingClient struct{}

func (recordingClient) Register(context.Context, RegisterRequest) (RegisterResponse, error) {
	return RegisterResponse{Metadata: okResponse()}, nil
}

func (recordingClient) LinkDevice(context.Context, LinkDeviceRequest) (LinkDeviceResponse, error) {
	return LinkDeviceResponse{Metadata: okResponse()}, nil
}

func (recordingClient) Send(context.Context, SendRequest) (SendResponse, error) {
	return SendResponse{Metadata: okResponse()}, nil
}

func (recordingClient) Receive(context.Context, ReceiveRequest) (ReceiveResponse, error) {
	return ReceiveResponse{Metadata: okResponse()}, nil
}

func (recordingClient) ReserveUsername(context.Context, ReserveUsernameRequest) (ReserveUsernameResponse, error) {
	return ReserveUsernameResponse{Metadata: okResponse()}, nil
}

func (recordingClient) UploadBackup(context.Context, UploadBackupRequest) (UploadBackupResponse, error) {
	return UploadBackupResponse{Metadata: okResponse()}, nil
}

func (recordingClient) DownloadBackup(context.Context, DownloadBackupRequest) (DownloadBackupResponse, error) {
	return DownloadBackupResponse{Metadata: okResponse()}, nil
}

func (recordingClient) RespondToChallenge(context.Context, RespondToChallengeRequest) (RespondToChallengeResponse, error) {
	return RespondToChallengeResponse{Metadata: okResponse()}, nil
}

func okResponse() ResponseMetadata {
	return ResponseMetadata{
		RequestID:  "request-123",
		Status:     "accepted",
		SecretRefs: map[string]string{"credential": "secret://signal/credential/alice"},
	}
}
