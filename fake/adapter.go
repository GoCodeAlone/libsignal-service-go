package fake

import (
	"context"
	"sync"

	"github.com/GoCodeAlone/libsignal-service-go/service"
)

var _ service.OperationTransport = (*Adapter)(nil)

type Adapter struct {
	mu      sync.Mutex
	records []AdapterRecord
}

type AdapterRecord struct {
	Operation      service.Operation
	OperationID    string
	IdempotencyKey string
	Audit          service.AuditMetadata
}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) SubmitOperation(_ context.Context, env service.OperationEnvelope) (service.OperationResult, error) {
	if err := env.Validate(); err != nil {
		return service.OperationResult{}, err
	}
	audit := env.Audit
	if audit.AccountHash == "" {
		var err error
		audit, err = service.NewAuditMetadata(env.AccountRef, map[string]string{
			"operation_id": string(env.OperationID),
			"operation":    string(env.Operation),
		})
		if err != nil {
			return service.OperationResult{}, err
		}
	}
	a.mu.Lock()
	a.records = append(a.records, AdapterRecord{
		Operation:      env.Operation,
		OperationID:    env.OperationID,
		IdempotencyKey: env.IdempotencyKey,
		Audit:          audit,
	})
	a.mu.Unlock()
	return service.OperationResult{
		OperationID: env.OperationID,
		Status:      "accepted",
		Audit:       audit,
		SecretRefs:  map[string]string{},
	}, nil
}

func (a *Adapter) Records() []AdapterRecord {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]AdapterRecord, len(a.records))
	copy(out, a.records)
	return out
}
