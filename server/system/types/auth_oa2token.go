package types

import (
	"time"

	sqlxTypes "github.com/jmoiron/sqlx/types"
)

type (
	AuthOa2token struct {
		ID         uint64
		Code       string
		Access     string
		Refresh    string
		ExpiresAt  time.Time
		CreatedAt  time.Time
		Data       sqlxTypes.JSONText
		ClientID   uint64
		UserID     uint64
		RemoteAddr string
		UserAgent  string
	}

	AuthOa2tokenFilter struct {
		UserID uint64
	}
)
