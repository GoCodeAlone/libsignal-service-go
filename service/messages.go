package service

import (
	"context"

	messagespb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/messages"
)

type Messages interface {
	SendMessage(context.Context, *messagespb.SendAuthenticatedSenderMessageRequest) (*messagespb.SendMessageAuthenticatedSenderResponse, error)
}
