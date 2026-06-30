package fake

import (
	"context"
	"sync"

	accountpb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/account"
	backuppb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/backup"
	challengepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/challenge"
	credentialspb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/credentials"
	devicepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/device"
	keyspb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/keys"
	messagespb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/messages"
	profilepb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/profile"
	"github.com/GoCodeAlone/libsignal-service-go/service"
	"google.golang.org/protobuf/proto"
)

var (
	_ service.AccountLookup  = (*Remote)(nil)
	_ service.Devices        = (*Remote)(nil)
	_ service.KeyLookup      = (*Remote)(nil)
	_ service.Messages       = (*Remote)(nil)
	_ service.ProfileLookup  = (*Remote)(nil)
	_ service.BackupMetadata = (*Remote)(nil)
	_ service.Challenge      = (*Remote)(nil)
	_ service.Credentials    = (*Remote)(nil)
)

type RequestRecord struct {
	Operation    string
	RequestType  string
	RequestBytes int
}

type Remote struct {
	mu      sync.Mutex
	records []RequestRecord

	UsernameHashResponse       *accountpb.LookupUsernameHashResponse
	UsernameLinkResponse       *accountpb.LookupUsernameLinkResponse
	DevicesResponse            *devicepb.GetDevicesResponse
	PreKeysResponse            *keyspb.GetPreKeysAnonymousResponse
	MessageResponse            *messagespb.SendMessageAuthenticatedSenderResponse
	ProfileResponse            *profilepb.GetVersionedProfileAnonymousResponse
	BackupAuthCredentials      *backuppb.GetBackupAuthCredentialsResponse
	ChallengeResponse          *challengepb.AnswerChallengeResponse
	ExternalCredentialResponse *credentialspb.GetExternalServiceCredentialsResponse
}

func NewRemote() *Remote {
	return &Remote{}
}

func (r *Remote) Records() []RequestRecord {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]RequestRecord, len(r.records))
	copy(out, r.records)
	return out
}

func (r *Remote) LookupUsernameHash(_ context.Context, req *accountpb.LookupUsernameHashRequest) (*accountpb.LookupUsernameHashResponse, error) {
	r.record("account.LookupUsernameHash", req)
	if r.UsernameHashResponse != nil {
		return proto.Clone(r.UsernameHashResponse).(*accountpb.LookupUsernameHashResponse), nil
	}
	return &accountpb.LookupUsernameHashResponse{}, nil
}

func (r *Remote) LookupUsernameLink(_ context.Context, req *accountpb.LookupUsernameLinkRequest) (*accountpb.LookupUsernameLinkResponse, error) {
	r.record("account.LookupUsernameLink", req)
	if r.UsernameLinkResponse != nil {
		return proto.Clone(r.UsernameLinkResponse).(*accountpb.LookupUsernameLinkResponse), nil
	}
	return &accountpb.LookupUsernameLinkResponse{}, nil
}

func (r *Remote) GetDevices(_ context.Context, req *devicepb.GetDevicesRequest) (*devicepb.GetDevicesResponse, error) {
	r.record("device.GetDevices", req)
	if r.DevicesResponse != nil {
		return proto.Clone(r.DevicesResponse).(*devicepb.GetDevicesResponse), nil
	}
	return &devicepb.GetDevicesResponse{}, nil
}

func (r *Remote) GetPreKeys(_ context.Context, req *keyspb.GetPreKeysAnonymousRequest) (*keyspb.GetPreKeysAnonymousResponse, error) {
	r.record("keys.GetPreKeys", req)
	if r.PreKeysResponse != nil {
		return proto.Clone(r.PreKeysResponse).(*keyspb.GetPreKeysAnonymousResponse), nil
	}
	return &keyspb.GetPreKeysAnonymousResponse{}, nil
}

func (r *Remote) SendMessage(_ context.Context, req *messagespb.SendAuthenticatedSenderMessageRequest) (*messagespb.SendMessageAuthenticatedSenderResponse, error) {
	r.record("messages.SendMessage", req)
	if r.MessageResponse != nil {
		return proto.Clone(r.MessageResponse).(*messagespb.SendMessageAuthenticatedSenderResponse), nil
	}
	return &messagespb.SendMessageAuthenticatedSenderResponse{}, nil
}

func (r *Remote) GetVersionedProfile(_ context.Context, req *profilepb.GetVersionedProfileAnonymousRequest) (*profilepb.GetVersionedProfileAnonymousResponse, error) {
	r.record("profile.GetVersionedProfile", req)
	if r.ProfileResponse != nil {
		return proto.Clone(r.ProfileResponse).(*profilepb.GetVersionedProfileAnonymousResponse), nil
	}
	return &profilepb.GetVersionedProfileAnonymousResponse{}, nil
}

func (r *Remote) GetBackupAuthCredentials(_ context.Context, req *backuppb.GetBackupAuthCredentialsRequest) (*backuppb.GetBackupAuthCredentialsResponse, error) {
	r.record("backup.GetBackupAuthCredentials", req)
	if r.BackupAuthCredentials != nil {
		return proto.Clone(r.BackupAuthCredentials).(*backuppb.GetBackupAuthCredentialsResponse), nil
	}
	return &backuppb.GetBackupAuthCredentialsResponse{}, nil
}

func (r *Remote) AnswerChallenge(_ context.Context, req *challengepb.AnswerChallengeRequest) (*challengepb.AnswerChallengeResponse, error) {
	r.record("challenge.AnswerChallenge", req)
	if r.ChallengeResponse != nil {
		return proto.Clone(r.ChallengeResponse).(*challengepb.AnswerChallengeResponse), nil
	}
	return &challengepb.AnswerChallengeResponse{}, nil
}

func (r *Remote) GetExternalServiceCredentials(_ context.Context, req *credentialspb.GetExternalServiceCredentialsRequest) (*credentialspb.GetExternalServiceCredentialsResponse, error) {
	r.record("credentials.GetExternalServiceCredentials", req)
	if r.ExternalCredentialResponse != nil {
		return proto.Clone(r.ExternalCredentialResponse).(*credentialspb.GetExternalServiceCredentialsResponse), nil
	}
	return &credentialspb.GetExternalServiceCredentialsResponse{}, nil
}

func (r *Remote) record(operation string, msg proto.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records = append(r.records, RequestRecord{
		Operation:    operation,
		RequestType:  string(msg.ProtoReflect().Descriptor().FullName()),
		RequestBytes: proto.Size(msg),
	})
}
