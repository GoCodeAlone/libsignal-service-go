package serviceclient_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/fake"
	"github.com/GoCodeAlone/libsignal-service-go/serviceclient"
	"github.com/GoCodeAlone/libsignal-service-go/servicepolicy"
)

func TestTransportFakeRequiresIdempotencyThroughSharedClient(t *testing.T) {
	transport, err := serviceclient.NewTransport(fake.NewServiceClient(), serviceclient.TransportConfig{Mode: serviceclient.TransportModeFake})
	if err != nil {
		t.Fatal(err)
	}
	_, err = transport.Register(context.Background(), serviceclient.RegisterRequest{Username: "alice.01"})
	if !errors.Is(err, fake.ErrMissingIdempotencyKey) {
		t.Fatalf("register without idempotency error = %v, want %v", err, fake.ErrMissingIdempotencyKey)
	}
	_, err = transport.Register(context.Background(), serviceclient.RegisterRequest{
		Metadata: requestMetadata("register-1"),
		Username: "alice.01",
	})
	if err != nil {
		t.Fatalf("register with idempotency: %v", err)
	}
}

func TestTransportLiveRejectsMissingApprovalPackage(t *testing.T) {
	_, err := serviceclient.NewTransport(fake.NewServiceClient(), serviceclient.TransportConfig{
		Mode:             serviceclient.TransportModeLive,
		ApprovalTime:     time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		RequestedActions: []servicepolicy.Action{servicepolicy.ActionSend},
	})
	if !errors.Is(err, servicepolicy.ErrLiveServiceDisabled) {
		t.Fatalf("live transport error = %v, want %v", err, servicepolicy.ErrLiveServiceDisabled)
	}
}

func TestTransportSandboxUsesExplicitEndpointWithoutProductionApproval(t *testing.T) {
	transport, err := serviceclient.NewTransport(fake.NewServiceClient(), serviceclient.TransportConfig{
		Mode:            serviceclient.TransportModeSandbox,
		SandboxEndpoint: "signal-sandbox.example.invalid",
	})
	if err != nil {
		t.Fatal(err)
	}
	if transport.Mode() != serviceclient.TransportModeSandbox {
		t.Fatalf("mode = %q, want %q", transport.Mode(), serviceclient.TransportModeSandbox)
	}
	if transport.Endpoint() != "signal-sandbox.example.invalid" {
		t.Fatalf("endpoint = %q", transport.Endpoint())
	}
	_, err = transport.Send(context.Background(), serviceclient.SendRequest{
		Metadata:     requestMetadata("send-1"),
		RecipientRef: "signal://recipient/bob",
		PayloadRef:   "ciphertext://payload/1",
	})
	if err != nil {
		t.Fatalf("sandbox send: %v", err)
	}
}

func TestTransportFakeCoversRepresentativeServiceDomains(t *testing.T) {
	client := fake.NewServiceClient()
	transport, err := serviceclient.NewTransport(client, serviceclient.TransportConfig{Mode: serviceclient.TransportModeFake})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	calls := []struct {
		name string
		run  func() error
	}{
		{
			name: "register",
			run: func() error {
				_, err := transport.Register(ctx, serviceclient.RegisterRequest{Metadata: requestMetadata("register-1")})
				return err
			},
		},
		{
			name: "link",
			run: func() error {
				_, err := transport.LinkDevice(ctx, serviceclient.LinkDeviceRequest{Metadata: requestMetadata("link-1"), LinkCodeRef: "secret://link"})
				return err
			},
		},
		{
			name: "send",
			run: func() error {
				_, err := transport.Send(ctx, serviceclient.SendRequest{Metadata: requestMetadata("send-1"), RecipientRef: "signal://recipient/bob", PayloadRef: "ciphertext://payload/1"})
				return err
			},
		},
		{
			name: "receive",
			run: func() error {
				_, err := transport.Receive(ctx, serviceclient.ReceiveRequest{Metadata: requestMetadata("receive-1"), CursorRef: "cursor://inbox"})
				return err
			},
		},
		{
			name: "username",
			run: func() error {
				_, err := transport.ReserveUsername(ctx, serviceclient.ReserveUsernameRequest{Metadata: requestMetadata("username-1"), Username: "alice.01"})
				return err
			},
		},
		{
			name: "backup-upload",
			run: func() error {
				_, err := transport.UploadBackup(ctx, serviceclient.UploadBackupRequest{Metadata: requestMetadata("backup-upload-1"), BackupRef: "backup://alice/1"})
				return err
			},
		},
		{
			name: "backup-download",
			run: func() error {
				_, err := transport.DownloadBackup(ctx, serviceclient.DownloadBackupRequest{Metadata: requestMetadata("backup-download-1"), BackupID: "backup-1"})
				return err
			},
		},
		{
			name: "challenge",
			run: func() error {
				_, err := transport.RespondToChallenge(ctx, serviceclient.RespondToChallengeRequest{Metadata: requestMetadata("challenge-1"), ChallengeRef: "challenge://1", ResponseRef: "secret://challenge/1"})
				return err
			},
		},
	}
	for _, call := range calls {
		if err := call.run(); err != nil {
			t.Fatalf("%s: %v", call.name, err)
		}
	}
	if got, want := len(client.Records()), len(calls); got != want {
		t.Fatalf("records = %d, want %d", got, want)
	}
}

func requestMetadata(idempotencyKey string) serviceclient.RequestMetadata {
	return serviceclient.RequestMetadata{
		AccountRef:     "secret://signal/account/alice",
		DeviceRef:      "secret://signal/device/alice-primary",
		IdempotencyKey: idempotencyKey,
		RequestedAt:    time.Unix(1_782_847_600, 0).UTC(),
		ConsentRef:     "audit://consent/alice",
		AuditRef:       "audit://signal/" + idempotencyKey,
		CredentialRef:  "secret://signal/credential/alice",
	}
}
