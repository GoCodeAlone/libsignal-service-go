package service

import (
	"context"

	credentialspb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/credentials"
)

type Credentials interface {
	GetExternalServiceCredentials(context.Context, *credentialspb.GetExternalServiceCredentialsRequest) (*credentialspb.GetExternalServiceCredentialsResponse, error)
}
