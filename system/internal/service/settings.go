package service

import (
	"context"
	"errors"

	"github.com/titpetric/factory"
	"go.uber.org/zap"

	internalSettings "github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
)

type (
	// Wrapper service for system around internal settings service
	settings struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac               settingsAccessController
		internalSettings internalSettings.Service
	}

	settingsAccessController interface {
		CanReadSettings(ctx context.Context) bool
		CanManageSettings(ctx context.Context) bool
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
		ac:               DefaultAccessControl,
		logger:           DefaultLogger.Named("settings"),
	}).With(ctx)
}

func (svc settings) With(ctx context.Context) SettingsService {
	db := repository.DB(ctx)

	return &settings{
		ctx:    ctx,
		db:     db,
		ac:     svc.ac,
		logger: svc.logger,

		internalSettings: svc.internalSettings.With(ctx),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc settings) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc settings) FindByPrefix(prefix string) (vv internalSettings.ValueSet, err error) {
	if !svc.ac.CanReadSettings(svc.ctx) {
		return nil, errors.New("not allowed to read settings")
	}

	return svc.internalSettings.FindByPrefix(prefix)
}

func (svc settings) Set(v *internalSettings.Value) (err error) {
	if !svc.ac.CanManageSettings(svc.ctx) {
		return errors.New("not allowed to manage settings")
	}

	return svc.internalSettings.Set(v)
}

func (svc settings) BulkSet(vv internalSettings.ValueSet) (err error) {
	if !svc.ac.CanManageSettings(svc.ctx) {
		return errors.New("not allowed to manage settings")
	}

	return svc.internalSettings.BulkSet(vv)
}

func (svc settings) Get(name string, ownedBy uint64) (out *internalSettings.Value, err error) {
	if !svc.ac.CanReadSettings(svc.ctx) {
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
