package service

import (
	"context"

	backuppb "github.com/GoCodeAlone/libsignal-service-go/gen/org/signal/chat/backup"
)

type BackupMetadata interface {
	GetBackupAuthCredentials(context.Context, *backuppb.GetBackupAuthCredentialsRequest) (*backuppb.GetBackupAuthCredentialsResponse, error)
}
