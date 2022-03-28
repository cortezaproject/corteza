package types

import (
	"time"
)

type (
	AuthSession struct {
		ID         string
		Data       []byte
		ExpiresAt  time.Time
		CreatedAt  time.Time
		UserID     uint64
		RemoteAddr string
		UserAgent  string
	}

	AuthSessionFilter struct {
		UserID uint64
		Limit  uint
	}
)
