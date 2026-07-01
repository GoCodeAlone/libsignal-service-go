package service

import "context"

type OperationTransport interface {
	SubmitOperation(context.Context, OperationEnvelope) (OperationResult, error)
}

type OperationResult struct {
	OperationID  string
	Status       string
	ChallengeRef string
	SecretRefs   map[string]string
	Audit        AuditMetadata
}
