package oauth2

import (
	"context"
	"strconv"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-oauth2/oauth2/v4"
)

type (
	// Wrapper for store to satisfy oauth2.ClientStore interface
	ContextClientStore struct{}
)

var _ oauth2.ClientStore = &ContextClientStore{}

// Pull client directly from context
//
// This requires that client is put in context before oauth2 procedures are executed!
func (s ContextClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return &clientInfo{ctx.Value(&ContextClientStore{}).(*types.AuthClient)}, nil
}

type (
	// Wrapper for client info object to satisfy oauth2.ClientInfo interface
	clientInfo struct{ *types.AuthClient }
)

var _ oauth2.ClientInfo = &clientInfo{}

func (c clientInfo) GetID() string {
	return strconv.FormatUint(c.ID, 10)
}

func (c clientInfo) GetSecret() string {
	return c.Secret
}

func (c clientInfo) GetDomain() string {
	return c.RedirectURI
}

func (c clientInfo) GetUserID() string {
	panic("implement me")
}
