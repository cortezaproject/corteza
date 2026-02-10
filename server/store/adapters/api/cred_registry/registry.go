package cred_registry

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

type (
	// dalConnectionStore defines the interface for persisting credential updates
	dalConnectionStore interface {
		LookupDalConnectionByID(ctx context.Context, id uint64) (*types.DalConnection, error)
		UpdateDalConnection(ctx context.Context, conn ...*types.DalConnection) error
	}

	registry struct {
		mux         sync.RWMutex
		credentials map[uint64]*Credential
		httpClient  *http.Client
		refresher   *refresher
		store       dalConnectionStore
		logger      *zap.Logger
	}
)

var (
	defaultRegistry *registry
)

func New(store dalConnectionStore, logger *zap.Logger) (*registry, error) {
	reg := &registry{
		credentials: make(map[uint64]*Credential),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		store:  store,
		logger: logger.Named("cred_registry"),
	}

	reg.refresher = newRefresher(reg, logger.Named("cred_refresher"))
	reg.refresher.start()

	return reg, nil
}

// Default returns the default registry instance
func Default() *registry {
	if defaultRegistry == nil {
		panic("credential registry not initialized")
	}
	return defaultRegistry
}

// SetDefault sets the default registry instance
func SetDefault(reg *registry) {
	defaultRegistry = reg
}

func (r *registry) Close() error {
	if r.refresher != nil {
		r.refresher.stop()
	}
	return nil
}

func (r *registry) Store(cred *Credential) error {
	if cred.ConnectionID == 0 {
		return fmt.Errorf("connection ID is required")
	}

	r.mux.Lock()
	defer r.mux.Unlock()
	r.credentials[cred.ConnectionID] = cred
	return nil
}

func (r *registry) Get(connectionID uint64) (*Credential, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	cred, ok := r.credentials[connectionID]
	if !ok {
		return nil, fmt.Errorf("credential not found for connection %d", connectionID)
	}

	return cred, nil
}

func (r *registry) Delete(connectionID uint64) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	delete(r.credentials, connectionID)
	return nil
}

func (r *registry) GetAccessToken(ctx context.Context, connectionID uint64) (string, error) {
	cred, err := r.Get(connectionID)
	if err != nil {
		return "", err
	}

	// If OAuth2 and needs refresh, do it now
	if cred.AuthType == "oauth2_client_credentials" && cred.NeedsRefresh() {
		if err := r.refreshToken(ctx, cred); err != nil {
			return "", fmt.Errorf("failed to refresh token: %w", err)
		}
	}

	return cred.GetAccessToken(), nil
}

func (r *registry) GetAllCredentials() []*Credential {
	r.mux.RLock()
	defer r.mux.RUnlock()

	creds := make([]*Credential, 0, len(r.credentials))
	for _, cred := range r.credentials {
		creds = append(creds, cred)
	}
	return creds
}

// persistTokenUpdate saves the updated OAuth2 token back to the database
func (r *registry) persistTokenUpdate(ctx context.Context, cred *Credential) error {
	if r.store == nil {
		r.logger.Warn("store not configured, skipping token persistence")
		return nil
	}

	conn, err := r.store.LookupDalConnectionByID(ctx, cred.ConnectionID)
	if err != nil {
		return fmt.Errorf("failed to load connection: %w", err)
	}

	if conn.Config.DAL == nil || conn.Config.DAL.Params == nil {
		return fmt.Errorf("connection has no DAL config")
	}

	// Update the auth params with new token
	if auth, ok := conn.Config.DAL.Params["auth"].(map[string]any); ok {
		if params, ok := auth["params"].(map[string]any); ok {
			params["accessToken"] = cred.AccessToken
			params["expiresAt"] = cred.ExpiresAt.Format(time.RFC3339)
		}
	}

	return r.store.UpdateDalConnection(ctx, conn)
}
