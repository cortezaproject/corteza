package gsheets

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
	apidal "github.com/cortezaproject/corteza/server/store/adapters/api/dal"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cast"
)

const (
	SCHEMA        = "gsheets"
	sheetsAPIBase = "https://sheets.googleapis.com/v4/spreadsheets"
)

func init() {
	dal.RegisterConnector(dalConnector, SCHEMA)
}

func dalConnector(ctx context.Context, dsn string) (_ dal.Connection, err error) {
	spew.Dump("hello==!=!=!=!=!=!", dsn)
	parsed, err := dal.ParseDSN(dsn)
	if err != nil {
		return
	}

	// spreadsheetID from DSN path or arbitrary params
	spreadsheetID := parsed.Path
	// if spreadsheetID == "" && parsed.Arbitrary != nil {
	// 	if sid, ok := parsed.Arbitrary["spreadsheetID"]; ok {
	// 		spreadsheetID = cast.ToString(sid)
	// 	}
	// }
	spreadsheetID = "1OfWlosXxacBz7b-EEemeIXNoZz09SS4byTZDiA35SVM"

	if spreadsheetID == "" {
		return nil, fmt.Errorf("spreadsheetID is required (use gsheets://<spreadsheetID> or set spreadsheetID param)")
	}

	baseURL := fmt.Sprintf("%s/%s", sheetsAPIBase, spreadsheetID)

	wrapper := newWrapper(baseURL, parsed.ConnectionID)
	dl := Dialect(spreadsheetID)

	// Register credential in the registry
	spew.Dump("xo", parsed)
	if parsed.ConnectionID > 0 {
		authType := parsed.AuthType
		switch authType {
		case "", "none", "jwt-bearer", "jwt_bearer":
			authType = "google_service_account"
		}

		// Scopes from arbitrary params if provided
		var scopes []string
		if parsed.Arbitrary != nil {
			if s, ok := parsed.Arbitrary["scopes"]; ok {
				if ss, ok := s.([]any); ok {
					for _, v := range ss {
						scopes = append(scopes, cast.ToString(v))
					}
				}
			}
		}

		cred, err := cred_registry.NewCredential(cred_registry.CredentialConfig{
			ConnectionID:        parsed.ConnectionID,
			AuthType:            authType,
			ServiceAccountEmail: parsed.Username,
			PrivateKey:          parsed.Token,
			Scopes:              scopes,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create credential: %w", err)
		}

		if err := cred_registry.Default().Store(cred); err != nil {
			return nil, fmt.Errorf("failed to store credentials: %w", err)
		}
	}

	return apidal.Connection(wrapper, dl), nil
}
