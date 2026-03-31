package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

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
// connectionID is used to register credentials in the cred_registry so that
// OAuth2 token fetching works; pass 0 to skip registration.
func RunnerFromConnection(conn *systemTypes.Connection, connectionID uint64) (*restAPIWrapper, error) {
	baseURL, err := url.Parse(conn.Service.BaseURL.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid baseURL %q: %w", conn.Service.BaseURL.Value, err)
	}

	headers := make(map[string][]string, len(conn.Service.Headers))
	for k, tpl := range conn.Service.Headers {
		headers[k] = []string{tpl.Value}
	}

	dsn := dal.DSN{
		Scheme:       baseURL.Scheme,
		Host:         baseURL.Hostname(),
		Port:         baseURL.Port(),
		Path:         baseURL.Path,
		AuthType:     conn.Service.Auth.Method,
		Headers:      headers,
		ConnectionID: connectionID,
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
		case "oauth2ClientCredentials", "serviceAccountJSON":
			// A JSON block containing either standard OAuth2 client credentials
			// or a Google Service Account key. We must extract the fields.
			var data map[string]string
			if err := json.Unmarshal([]byte(tpl.Value), &data); err == nil {
				// Google SA fields
				if email, ok := data["client_email"]; ok {
					dsn.Username = email
				}
				if pk, ok := data["private_key"]; ok {
					dsn.Token = pk
				}
				// Standard OAuth2 fields
				if id, ok := data["client_id"]; ok {
					dsn.ClientID = id
				}
				if secret, ok := data["client_secret"]; ok {
					dsn.ClientSecret = secret
				}
				if tokenURL, ok := data["token_uri"]; ok {
					dsn.TokenURL = tokenURL
				}
			} else {
				// Sometimes private_key PEM blocks contain literal newlines that break standard unmarshal
				fallbackJSON := strings.ReplaceAll(tpl.Value, "\n", `\n`)
				if err2 := json.Unmarshal([]byte(fallbackJSON), &data); err2 == nil {
					if email, ok := data["client_email"]; ok {
						dsn.Username = email
					}
					if pk, ok := data["private_key"]; ok {
						dsn.Token = pk
					}
					if id, ok := data["client_id"]; ok {
						dsn.ClientID = id
					}
					if secret, ok := data["client_secret"]; ok {
						dsn.ClientSecret = secret
					}
					if tokenURL, ok := data["token_uri"]; ok {
						dsn.TokenURL = tokenURL
					}
				}
			}
		}
	}

	// Register credentials from the resolved config so that OAuth2 token
	// fetching (via cred_registry) works without a database round-trip.
	if connectionID > 0 && dsn.AuthType != "" {
		// If the JSON block contained an email and private key, it's actually
		// a Google Service Account disguised as oauth2_client_credentials (or similar).
		authType := dsn.AuthType
		if dsn.Username != "" && strings.Contains(dsn.Token, "PRIVATE KEY") {
			authType = "google_service_account"
		}

		var scopes []string
		if authType == "google_service_account" {
			baseURL := conn.Service.BaseURL.Value
			switch {
			case strings.Contains(baseURL, "googleapis.com/calendar"):
				scopes = []string{"https://www.googleapis.com/auth/calendar"}
			case strings.Contains(baseURL, "sheets.googleapis.com"):
				scopes = []string{"https://www.googleapis.com/auth/spreadsheets"}
			case strings.Contains(baseURL, "googleapis.com/drive"):
				scopes = []string{"https://www.googleapis.com/auth/drive"}
			case strings.Contains(baseURL, "tasks.googleapis.com"):
				scopes = []string{"https://www.googleapis.com/auth/tasks"}
			}
		}

		cred, credErr := cred_registry.NewCredential(cred_registry.CredentialConfig{
			ConnectionID:        connectionID,
			AuthType:            authType,
			Token:               dsn.Token,
			APIKey:              dsn.APIKey,
			ClientID:            dsn.ClientID,
			ClientSecret:        dsn.ClientSecret,
			TokenURL:            dsn.TokenURL,
			ServiceAccountEmail: dsn.Username,
			PrivateKey:          dsn.Token,
			Scopes:              scopes,
		})
		if credErr == nil {
			_ = cred_registry.Default().Store(cred)
		}
	}

	return newWrapperFromParsedDSN(dsn), nil
}
