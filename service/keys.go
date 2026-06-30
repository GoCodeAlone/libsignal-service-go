package service

import (
	"context"

	keyspb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/keys"
)

type KeyLookup interface {
	GetPreKeys(context.Context, *keyspb.GetPreKeysAnonymousRequest) (*keyspb.GetPreKeysAnonymousResponse, error)
}
