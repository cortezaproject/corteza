package cred_registry

import (
	"context"
	"net/http"
)

// BearerCredential holds a static bearer token.
type BearerCredential struct {
	connID uint64
	Token  string
}

func NewBearerCredential(connID uint64, token string) *BearerCredential {
	return &BearerCredential{connID: connID, Token: token}
}

func (c *BearerCredential) ConnectionID() uint64        { return c.connID }
func (c *BearerCredential) AuthType() string             { return "bearer" }
func (c *BearerCredential) GetAccessToken() string       { return c.Token }
func (c *BearerCredential) NeedsRefresh() bool           { return false }
func (c *BearerCredential) MarshalState() map[string]any { return nil }

func (c *BearerCredential) Refresh(_ context.Context, _ *http.Client) error {
	return nil
}
