package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

func NewIdempotencyKey(operation Operation, accountRef string, requestedAt time.Time) string {
	h := sha256.New()
	h.Write([]byte(operation))
	h.Write([]byte{0})
	h.Write([]byte(accountRef))
	h.Write([]byte{0})
	h.Write([]byte(strconv.FormatInt(requestedAt.UTC().UnixNano(), 10)))
	sum := h.Sum(nil)
	return string(operation) + ":" + hex.EncodeToString(sum[:16])
}
