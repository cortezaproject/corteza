package cred_registry

import (
	"context"
	"net/http"
)

// APIKeyCredential holds a static API key.
type APIKeyCredential struct {
	connID uint64
	Key    string
	// default: X-API-Key
	HeaderName string
	// optional prefix
	Prefix string
}

func NewAPIKeyCredential(connID uint64, key, headerName, prefix string) *APIKeyCredential {
	if headerName == "" {
		headerName = "X-API-Key"
	}
	return &APIKeyCredential{
		connID:     connID,
		Key:        key,
		HeaderName: headerName,
		Prefix:     prefix,
	}
}

func (c *APIKeyCredential) ConnectionID() uint64         { return c.connID }
func (c *APIKeyCredential) AuthType() string             { return "apikey" }
func (c *APIKeyCredential) GetAccessToken() string       { return c.Key }
func (c *APIKeyCredential) NeedsRefresh() bool           { return false }
func (c *APIKeyCredential) MarshalState() map[string]any { return nil }

func (c *APIKeyCredential) Refresh(_ context.Context, _ *http.Client) error {
	return nil
}
