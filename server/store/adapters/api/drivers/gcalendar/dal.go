package gcalendar

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/store/adapters/api/cred_registry"
	apidal "github.com/cortezaproject/corteza/server/store/adapters/api/dal"
	"github.com/spf13/cast"
)

const (
	SCHEMA          = "gcalendar"
	calendarAPIBase = "https://www.googleapis.com/calendar/v3/calendars"
)

func init() {
	dal.RegisterConnector(dalConnector, SCHEMA)
}

func dalConnector(ctx context.Context, dsn string) (_ dal.Connection, err error) {
	parsed, err := dal.ParseDSN(dsn)
	if err != nil {
		return
	}

	calendarID := parsed.Path
	if calendarID == "" && parsed.Arbitrary != nil {
		if cid, ok := parsed.Arbitrary["calendarID"]; ok {
			calendarID = cast.ToString(cid)
		}
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	baseURL := fmt.Sprintf("%s/%s", calendarAPIBase, calendarID)
	wrapper := newWrapper(baseURL, parsed.ConnectionID)
	dl := Dialect(calendarID)

	if parsed.ConnectionID > 0 {
		authType := parsed.AuthType
		switch authType {
		case "", "none", "jwt-bearer", "jwt_bearer":
			authType = "google_service_account"
		}

		var scopes []string
		var subject string
		if parsed.Arbitrary != nil {
			if s, ok := parsed.Arbitrary["scopes"]; ok {
				if ss, ok := s.([]any); ok {
					for _, v := range ss {
						scopes = append(scopes, cast.ToString(v))
					}
				}
			}
			if sub, ok := parsed.Arbitrary["subject"]; ok {
				subject = cast.ToString(sub)
			}
		}

		cred, err := cred_registry.NewCredential(cred_registry.CredentialConfig{
			ConnectionID:        parsed.ConnectionID,
			AuthType:            authType,
			ServiceAccountEmail: parsed.Username,
			PrivateKey:          parsed.Token,
			Scopes:              scopes,
			Subject:             subject,
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
