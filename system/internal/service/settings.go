package service

import (
	"context"
	"errors"

	internalSettings "github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/repository"
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
	}
)

func Settings(ctx context.Context) SettingsService {
	return (&settings{}).With(ctx)
}

func (svc settings) With(ctx context.Context) SettingsService {
	db := repository.DB(ctx)
	return &settings{
		ctx:              ctx,
		prm:              Permissions(ctx),
		internalSettings: internalSettings.NewService(internalSettings.NewRepository(db, "sys_settings")).With(ctx),
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
