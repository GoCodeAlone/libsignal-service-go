package service

import (
	"context"

	profilepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/profile"
)

type ProfileLookup interface {
	GetVersionedProfile(context.Context, *profilepb.GetVersionedProfileAnonymousRequest) (*profilepb.GetVersionedProfileAnonymousResponse, error)
}
