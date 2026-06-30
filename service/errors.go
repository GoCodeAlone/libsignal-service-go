package service

import (
	"errors"

	"google.golang.org/grpc/status"
)

type RPCError struct {
	Code    string
	Message string
}

func NormalizeError(err error) RPCError {
	if err == nil {
		return RPCError{}
	}
	st, ok := status.FromError(err)
	if !ok {
		return RPCError{Code: "unknown", Message: err.Error()}
	}
	return RPCError{Code: st.Code().String(), Message: st.Message()}
}

func IsRPCError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := status.FromError(err)
	return ok || !errors.Is(err, nil)
}
