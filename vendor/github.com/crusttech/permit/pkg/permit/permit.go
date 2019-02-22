package permit

import (
	"time"

	"github.com/pkg/errors"
)

type (
	Permit struct {
		Version    uint           `json:"version"`
		Key        string         `json:"key"`
		Domain     string         `json:"domain"`
		Expires    *time.Time     `json:"expires,omitempty"`
		Valid      bool           `json:"valid"`
		Attributes map[string]int `json:"attributes"`
	}
)

const (
	// KeyLength
	KeyLength = 64
)

var (
	PermitNotFound = errors.New("permit not found")
)

func (p Permit) IsValid() bool {
	return p.Valid && !p.Expired()
}

func (p Permit) Expired() bool {
	return p.Expires != nil && p.Expires.Before(time.Now())
}
