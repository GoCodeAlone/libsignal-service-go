package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

var ErrForbiddenAuditField = errors.New("forbidden audit field")

type AuditMetadata struct {
	AccountHash string
	fields      map[string]string
}

func NewAuditMetadata(accountRef string, fields map[string]string) (AuditMetadata, error) {
	for key := range fields {
		if forbiddenAuditField(key) {
			return AuditMetadata{}, ErrForbiddenAuditField
		}
	}
	hash := sha256.Sum256([]byte(accountRef))
	return AuditMetadata{
		AccountHash: "sha256:" + hex.EncodeToString(hash[:16]),
		fields:      cloneStringMap(fields),
	}, nil
}

func (m AuditMetadata) CopyFields() map[string]string {
	return cloneStringMap(m.fields)
}

func forbiddenAuditField(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
	switch normalized {
	case "message_body", "body", "phone_number", "phone", "e164":
		return true
	default:
		return false
	}
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
