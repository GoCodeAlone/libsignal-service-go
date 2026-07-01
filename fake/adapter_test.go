package fake

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/GoCodeAlone/libsignal-service-go/service"
)

func TestAdapterAcceptsEveryNamedOperationAndRecordsRedactedAudit(t *testing.T) {
	adapter := NewAdapter()
	for _, operation := range []service.Operation{
		service.OperationRegister,
		service.OperationLinkDevice,
		service.OperationSend,
		service.OperationReceive,
		service.OperationChallenge,
		service.OperationUsername,
		service.OperationBackup,
		service.OperationSVR,
	} {
		env := testEnvelope(operation, string(operation)+"-1")
		if _, err := adapter.SubmitOperation(context.Background(), env); err != nil {
			t.Fatalf("%s: %v", operation, err)
		}
	}
	records := adapter.Records()
	if got, want := len(records), 8; got != want {
		t.Fatalf("records = %d, want %d", got, want)
	}
	for _, record := range records {
		if record.Audit.AccountHash == "" {
			t.Fatalf("record missing audit hash: %#v", record)
		}
		if strings.Contains(record.Audit.AccountHash, "alice") || strings.Contains(record.Audit.AccountHash, "secret://") {
			t.Fatalf("record leaked account identifier: %#v", record)
		}
	}
}

func testTime() time.Time {
	return time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
}

func testDuration() time.Duration {
	return time.Hour
}

func testEnvelope(operation service.Operation, id string) service.OperationEnvelope {
	now := testTime()
	env := service.OperationEnvelope{
		OperationID:    id,
		Operation:      operation,
		IdempotencyKey: service.NewIdempotencyKey(operation, "secret://signal/account/alice", now),
		AccountRef:     "secret://signal/account/alice",
		RequestedAt:    now,
	}
	if operation == service.OperationLinkDevice {
		env.LinkedDevice = &service.LinkedDeviceEnvelope{
			DeviceDisplayName: "Alice laptop",
			ConsentRef:        "audit://consent/alice/link-device",
			ConsentExpiresAt:  now.Add(testDuration()),
			RevocationURI:     "https://operator.example.invalid/signal/devices/revoke/alice-laptop",
			UnlinkProofRef:    "proof://signal/unlink/alice-laptop",
		}
	}
	return env
}
