package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	systemTypes "github.com/cortezaproject/corteza/server/system/types"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
	apidal "github.com/cortezaproject/corteza/server/store/adapters/api/dal"
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
		cred, err := cred_registry.NewCredential(cred_registry.CredentialConfig{
			ConnectionID: parsed.ConnectionID,
			AuthType:     parsed.AuthType,
			Token:        parsed.Token,
			APIKey:       parsed.APIKey,
			ClientID:     parsed.ClientID,
			ClientSecret: parsed.ClientSecret,
			TokenURL:     parsed.TokenURL,
		})
		if err != nil {
			return nil, err
		}

		if err := cred_registry.Default().Store(cred); err != nil {
			return nil, fmt.Errorf("failed to store credentials: %w", err)
		}
	}

	return apidal.Connection(newWrapperFromParsedDSN(parsed), dl), nil
}

// newWrapperFromParsedDSN builds a restAPIWrapper from an already-parsed DSN.
// Shared by dalConnector (DSN string path) and RunnerFromConnection (Connection struct path).
func newWrapperFromParsedDSN(parsed dal.DSN) *restAPIWrapper {
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
		BaseURL:             u,
		DSN:                 parsed,
		Timeout:             parsed.Timeout,
		MaxIdleConns:        parsed.MaxIdleConns,
		MaxIdleConnsPerHost: parsed.MaxIdleConnsPerHost,
		IdleConnTimeout:     parsed.IdleConnTimeout,
		Headers:             parsed.Headers,
	})

	return &restAPIWrapper{
		client:       c,
		dsn:          parsed,
		connectionID: parsed.ConnectionID,
	}
}

// RunnerFromConnection builds a transient restAPIWrapper from a resolved Connection
// without persisting anything to the DAL. Used for pre-enable health checks.
func RunnerFromConnection(conn *systemTypes.Connection) (*restAPIWrapper, error) {
	baseURL, err := url.Parse(conn.Service.BaseURL.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL %q: %w", conn.Service.BaseURL.Value, err)
	}

	headers := make(map[string][]string, len(conn.Service.Headers))
	for k, tpl := range conn.Service.Headers {
		headers[k] = []string{tpl.Value}
	}

	dsn := dal.DSN{
		Scheme:   baseURL.Scheme,
		Host:     baseURL.Hostname(),
		Port:     baseURL.Port(),
		Path:     baseURL.Path,
		AuthType: conn.Service.Auth.Method,
		Headers:  headers,
	}

	for k, tpl := range conn.Service.Auth.Params {
		switch k {
		case "apikey":
			dsn.APIKey = tpl.Value
		case "APIKeyHeader":
			dsn.APIKeyHeader = tpl.Value
		case "token":
			dsn.Token = tpl.Value
		case "username":
			dsn.Username = tpl.Value
		case "password":
			dsn.Password = tpl.Value
		case "clientID":
			dsn.ClientID = tpl.Value
		case "clientSecret":
			dsn.ClientSecret = tpl.Value
		case "tokenURL":
			dsn.TokenURL = tpl.Value
		}
	}

	return newWrapperFromParsedDSN(dsn), nil
}
