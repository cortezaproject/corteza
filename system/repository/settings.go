package repository

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

type (
	settings struct {
		*repository
	}

	Settings interface {
		With(ctx context.Context, db *factory.DB) Settings

		Get(name string, value interface{}) (bool, error)
		Set(name string, value interface{}) error
	}
)

func NewSettings(ctx context.Context, db *factory.DB) Settings {
	return (&settings{}).With(ctx, db)
}

func (r *settings) With(ctx context.Context, db *factory.DB) Settings {
	return &settings{
		repository: r.repository.With(ctx, db),
	}
}

func (r *settings) Set(name string, value interface{}) error {
	if jsonValue, err := json.Marshal(value); err != nil {
		return errors.Wrap(err, "Error marshaling settings value")
	} else {
		return r.db().Replace("settings", struct {
			Key string          `db:"name"`
			Val json.RawMessage `db:"value"`
		}{name, jsonValue})
	}
}

func (r *settings) Get(name string, value interface{}) (bool, error) {
	sql := "SELECT value FROM settings WHERE name = ?"

	var stored json.RawMessage

	if err := r.db().Get(&stored, sql, name); err != nil {
		return false, errors.Wrap(err, "Error reading settings from the database")
	} else if stored == nil {
		return false, nil
	} else if err := json.Unmarshal(stored, value); err != nil {
		return false, errors.Wrap(err, "Error unmarshaling settings value")
	}

	return true, nil
}
