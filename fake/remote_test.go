package fake

import (
	"context"
	"strings"
	"sync"
	"testing"

	accountpb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/account"
	backuppb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/backup"
	challengepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/challenge"
	credentialspb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/credentials"
	devicepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/device"
	keyspb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/keys"
	messagespb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/messages"
	profilepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/profile"
)

func TestRemoteRecordsRepresentativeFlowsWithoutSecrets(t *testing.T) {
	ctx := context.Background()
	r := NewRemote()

	_, _ = r.LookupUsernameHash(ctx, &accountpb.LookupUsernameHashRequest{UsernameHash: []byte("secret-username-hash")})
	_, _ = r.LookupUsernameLink(ctx, &accountpb.LookupUsernameLinkRequest{UsernameLinkHandle: []byte("secret-link-handle")})
	_, _ = r.GetDevices(ctx, &devicepb.GetDevicesRequest{})
	_, _ = r.GetPreKeys(ctx, &keyspb.GetPreKeysAnonymousRequest{})
	_, _ = r.SendMessage(ctx, &messagespb.SendAuthenticatedSenderMessageRequest{})
	_, _ = r.GetVersionedProfile(ctx, &profilepb.GetVersionedProfileAnonymousRequest{})
	_, _ = r.GetBackupAuthCredentials(ctx, &backuppb.GetBackupAuthCredentialsRequest{})
	_, _ = r.AnswerChallenge(ctx, &challengepb.AnswerChallengeRequest{})
	_, _ = r.GetExternalServiceCredentials(ctx, &credentialspb.GetExternalServiceCredentialsRequest{})

	records := r.Records()
	if got, want := len(records), 9; got != want {
		t.Fatalf("records = %d, want %d", got, want)
	}
	var joined strings.Builder
	for _, rec := range records {
		if rec.Operation == "" || rec.RequestType == "" {
			t.Fatalf("incomplete record: %+v", rec)
		}
		joined.WriteString(rec.Operation)
		joined.WriteString(rec.RequestType)
	}
	for _, secret := range []string{"secret-username-hash", "secret-link-handle"} {
		if strings.Contains(joined.String(), secret) {
			t.Fatalf("record leaked secret %q: %s", secret, joined.String())
		}
	}
}

func TestRemoteRecordsConcurrently(t *testing.T) {
	ctx := context.Background()
	r := NewRemote()
	var wg sync.WaitGroup
	for range 32 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = r.GetDevices(ctx, &devicepb.GetDevicesRequest{})
		}()
	}
	wg.Wait()
	if got, want := len(r.Records()), 32; got != want {
		t.Fatalf("records = %d, want %d", got, want)
	}
}
