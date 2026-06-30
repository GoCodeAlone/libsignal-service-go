package service

import (
	"context"

	devicepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/device"
)

type Devices interface {
	GetDevices(context.Context, *devicepb.GetDevicesRequest) (*devicepb.GetDevicesResponse, error)
}
