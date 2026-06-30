package fake

import (
	"context"
	"testing"

	devicepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/device"
)

func TestRemoteReturnsClonedResponses(t *testing.T) {
	ctx := context.Background()
	r := NewRemote()
	r.DevicesResponse = &devicepb.GetDevicesResponse{
		Devices: []*devicepb.GetDevicesResponse_LinkedDevice{{Id: 2, Name: []byte("fixture-device")}},
	}

	got, err := r.GetDevices(ctx, &devicepb.GetDevicesRequest{})
	if err != nil {
		t.Fatalf("GetDevices: %v", err)
	}
	got.Devices[0].Name = []byte("mutated")

	again, err := r.GetDevices(ctx, &devicepb.GetDevicesRequest{})
	if err != nil {
		t.Fatalf("GetDevices second call: %v", err)
	}
	if string(again.Devices[0].Name) != "fixture-device" {
		t.Fatalf("stored fixture mutated: %q", string(again.Devices[0].Name))
	}
}
