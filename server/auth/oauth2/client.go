package oauth2

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/store"
	systemService "github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/spf13/cast"
)

type (
	// Wrapper for store to satisfy oauth2.clientStore interface
	clientStore struct {
		store store.AuthClients
		def   *types.AuthClient
	}
)

var (
	_ oauth2.ClientStore = &clientStore{}
)

func NewClientStore(s store.AuthClients, def *types.AuthClient) *clientStore {
	return &clientStore{s, def}
}

// GetByID pulls client directly from context
//
// This requires that client is put in context before oauth2 procedures are executed!
func (cs clientStore) GetByID(ctx context.Context, id string) (_ oauth2.ClientInfo, err error) {
	var c *types.AuthClient

	if id == "0" || cs.def != nil && cast.ToUint64(id) == cs.def.ID {
		if cs.def == nil {
			return nil, fmt.Errorf("could not provide default auth client")
		}

		c = cs.def
	} else if c, err = clientLookup(ctx, cs.store, id); err != nil {
		return nil, fmt.Errorf("failed to do auth client lookup (%q): %w", id, err)
	}

	m := &models.Client{
		ID:     strconv.FormatUint(c.ID, 10),
		Secret: c.Secret,
		Domain: c.RedirectURI,
	}

	return m, nil
}

func clientLookup(ctx context.Context, s store.AuthClients, identifier interface{}) (*types.AuthClient, error) {
	if id := cast.ToUint64(identifier); id > 0 {
		return store.LookupAuthClientByID(ctx, s, id)
	} else if h := cast.ToString(identifier); handle.IsValid(h) {
		return store.LookupAuthClientByHandle(ctx, s, h)
	} else {
		return nil, systemService.AuthClientErrInvalidID()
	}
}
