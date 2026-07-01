package operatorfixture

import (
	"context"
	"testing"

	"github.com/GoCodeAlone/libsignal-service-go/service"
)

func TestAdapterUsesAllowlistedLivePath(t *testing.T) {
	adapter, err := NewAdapter(Config{Endpoint: "127.0.0.1:19091"})
	if err != nil {
		t.Fatal(err)
	}
	env := service.OperationEnvelope{
		OperationID:    "send-1",
		Operation:      service.OperationSend,
		IdempotencyKey: service.NewIdempotencyKey(service.OperationSend, "secret://signal/account/alice", ApprovalTime()),
		AccountRef:     "secret://signal/account/alice",
		RequestedAt:    ApprovalTime(),
	}
	result, err := adapter.SubmitOperation(context.Background(), env)
	if err != nil {
		t.Fatal(err)
	}
	if result.OperationID != "send-1" || result.Status != "accepted" {
		t.Fatalf("result = %#v", result)
	}
}
