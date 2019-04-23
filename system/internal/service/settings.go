package service

import (
	"context"
	"errors"

	internalSettings "github.com/crusttech/crust/internal/settings"
)

type (
	// Wrapper service for system around internal settings service
	settings struct {
		db  db
		ctx context.Context

		prm              PermissionsService
		internalSettings internalSettings.Service
	}

	SettingsService interface {
		With(ctx context.Context) SettingsService
		FindByPrefix(prefix string) (vv internalSettings.ValueSet, err error)
		Set(v *internalSettings.Value) (err error)
		BulkSet(vv internalSettings.ValueSet) (err error)
		Get(name string, ownedBy uint64) (out *internalSettings.Value, err error)

		LoadAuthSettings() (authSettings, error)
	}
)

func Settings(ctx context.Context, intSet internalSettings.Service) SettingsService {
	return (&settings{
		internalSettings: intSet,
	}).With(ctx)
}

func (svc settings) With(ctx context.Context) SettingsService {
	return &settings{
		ctx:              ctx,
		prm:              Permissions(ctx),
		internalSettings: svc.internalSettings.With(ctx),
	}
}

func (svc settings) FindByPrefix(prefix string) (vv internalSettings.ValueSet, err error) {
	if !svc.prm.CanReadSettings() {
		return nil, errors.New("not allowed to read settings")
	}

	return svc.internalSettings.FindByPrefix(prefix)
}

func (svc settings) Set(v *internalSettings.Value) (err error) {
	if !svc.prm.CanManageSettings() {
		return errors.New("not allowed to manage settings")
	}

	return svc.internalSettings.Set(v)
}

func (svc settings) BulkSet(vv internalSettings.ValueSet) (err error) {
	if !svc.prm.CanManageSettings() {
		return errors.New("not allowed to manage settings")
	}

	return svc.internalSettings.BulkSet(vv)
}

func (svc settings) Get(name string, ownedBy uint64) (out *internalSettings.Value, err error) {
	if !svc.prm.CanReadSettings() {
		return nil, errors.New("not allowed to read settings")
	}

	return svc.internalSettings.Get(name, ownedBy)
}

// Loads auth.% settings, initializes & fills auth settings struct
func (svc settings) LoadAuthSettings() (authSettings, error) {
	vv, err := svc.internalSettings.FindByPrefix("auth.")
	if err != nil {
		return authSettings{}, err
	}
	return AuthSettings(vv.KV()), nil
}
