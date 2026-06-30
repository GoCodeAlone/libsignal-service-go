package service

import (
	"context"

	accountpb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/account"
)

type AccountLookup interface {
	LookupUsernameHash(context.Context, *accountpb.LookupUsernameHashRequest) (*accountpb.LookupUsernameHashResponse, error)
	LookupUsernameLink(context.Context, *accountpb.LookupUsernameLinkRequest) (*accountpb.LookupUsernameLinkResponse, error)
}
