package fake

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/serviceclient"
)

func TestServiceClientFakeReplaysDuplicateIdempotencyKey(t *testing.T) {
	client := NewServiceClient()
	req := serviceclient.RegisterRequest{Metadata: requestMetadata("request-1"), Username: "alice.01"}

	first, err := client.Register(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	second, err := client.Register(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if first.Metadata.RequestID != second.Metadata.RequestID || first.Metadata.Status != second.Metadata.Status {
		t.Fatalf("duplicate response = %#v, want replay of %#v", second.Metadata, first.Metadata)
	}
	if got, want := len(client.Records()), 1; got != want {
		t.Fatalf("records = %d, want %d", got, want)
	}
}

func TestServiceClientFakeRejectsConflictingDuplicatePayload(t *testing.T) {
	client := NewServiceClient()
	req := serviceclient.RegisterRequest{Metadata: requestMetadata("request-1"), Username: "alice.01"}
	if _, err := client.Register(context.Background(), req); err != nil {
		t.Fatal(err)
	}
	req.Username = "alice.02"
	if _, err := client.Register(context.Background(), req); !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("conflicting duplicate error = %v, want %v", err, ErrIdempotencyConflict)
	}
}

func TestServiceClientFakeSimulatesChallengeAndCredentialExpiry(t *testing.T) {
	challenged := NewServiceClient(WithChallenge("send", "challenge://alice/1"))
	send, err := challenged.Send(context.Background(), serviceclient.SendRequest{
		Metadata:     requestMetadata("send-1"),
		RecipientRef: "signal://recipient/bob",
		PayloadRef:   "ciphertext://payload/1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if send.Metadata.Status != "challenge_required" || send.Metadata.ChallengeRef != "challenge://alice/1" {
		t.Fatalf("challenge response = %#v", send.Metadata)
	}

	expired := NewServiceClient(WithExpiredCredentials())
	receive, err := expired.Receive(context.Background(), serviceclient.ReceiveRequest{Metadata: requestMetadata("receive-1")})
	if err != nil {
		t.Fatal(err)
	}
	if receive.Metadata.Status != "credential_expired" {
		t.Fatalf("expired response status = %q, want credential_expired", receive.Metadata.Status)
	}
}

func TestServiceClientFakeRestoresLedgerAfterRestart(t *testing.T) {
	client := NewServiceClient()
	req := serviceclient.LinkDeviceRequest{Metadata: requestMetadata("link-1"), LinkCodeRef: "secret://link-code"}
	first, err := client.LinkDevice(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	restarted := NewServiceClient(WithLedger(client.Ledger()))
	second, err := restarted.LinkDevice(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if second.Metadata.RequestID != first.Metadata.RequestID || second.Metadata.Status != first.Metadata.Status {
		t.Fatalf("restarted duplicate response = %#v, want %#v", second.Metadata, first.Metadata)
	}
	if got, want := len(restarted.Records()), 0; got != want {
		t.Fatalf("restart replay records = %d, want %d new records", got, want)
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
