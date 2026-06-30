package service

import (
	"context"

	challengepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/challenge"
)

type Challenge interface {
	AnswerChallenge(context.Context, *challengepb.AnswerChallengeRequest) (*challengepb.AnswerChallengeResponse, error)
}
