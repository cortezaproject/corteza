package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	apidal "github.com/cortezaproject/corteza/server/store/adapters/api/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
)

const (
	SCHEMA = "restapi"
)

func init() {
	dal.RegisterConnector(dalConnector, SCHEMA)
}

func dalConnector(ctx context.Context, dsn string) (_ dal.Connection, err error) {
	parsed, err := dal.ParseDSN(dsn)
	if err != nil {
		return
	}

	h := parsed.Host
	if parsed.Port != "" {
		h = fmt.Sprintf("%s:%s", h, parsed.Port)
	}

	u := url.URL{
		Scheme: parsed.Scheme,
		Host:   h,
		Path:   parsed.Path,
	}

	c := newClient(ClientConfig{
		BaseURL: u,
		DSN:     parsed,

		Timeout:             parsed.Timeout,
		MaxIdleConns:        parsed.MaxIdleConns,
		MaxIdleConnsPerHost: parsed.MaxIdleConnsPerHost,
		IdleConnTimeout:     parsed.IdleConnTimeout,
		Headers:             parsed.Headers,
	})

	dl := apiDialect{}

	if parsed.Arbitrary != nil {
		bb, err := json.Marshal(parsed.Arbitrary)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(bb, &dl.defaultOps)
	}

	// Load credentials into registry for all auth types
	if parsed.ConnectionID > 0 && parsed.AuthType != "" {
		cred := &cred_registry.Credential{
			ConnectionID: parsed.ConnectionID,
			AuthType:     parsed.AuthType,
			Token:        parsed.Token,
			APIKey:       parsed.APIKey,
			ClientID:     parsed.ClientID,
			ClientSecret: parsed.ClientSecret,
			TokenURL:     parsed.TokenURL,
		}

		// force immediate refresh on first use for oauth2
		if parsed.AuthType == "oauth2_client_credentials" {
			cred.ExpiresAt = time.Now().Add(-1 * time.Hour)
		}

		if err := cred_registry.Default().Store(cred); err != nil {
			return nil, fmt.Errorf("failed to store credentials: %w", err)
		}
	}

	return apidal.Connection(
		&restAPIWrapper{
			client:       c,
			dsn:          parsed,
			connectionID: parsed.ConnectionID,
		},
		dl,
	), nil
}
