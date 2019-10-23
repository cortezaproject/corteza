package service

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	internalSettings "github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/repository"
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

		LoadAuthSettings() (*AuthSettings, error)
		LoadSystemSettings() (*SystemSettings, error)
		UpdateAuthSettings(*AuthSettings) error
		UpdateSystemSettings(*SystemSettings) error
		AutoDiscovery() error
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

func (svc settings) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

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
func (svc settings) LoadAuthSettings() (*AuthSettings, error) {
	as := &AuthSettings{}
	return as, svc.UpdateAuthSettings(as)
}

// Loads system.% settings, initializes & fills system settings struct
func (svc settings) LoadSystemSettings() (*SystemSettings, error) {
	ss := &SystemSettings{}
	return ss, svc.UpdateSystemSettings(ss)
}

func (svc settings) UpdateSystemSettings(ss *SystemSettings) error {
	vv, err := svc.internalSettings.FindByPrefix("system.")
	if err != nil {
		return err
	}

	return ss.ReadKV(vv.KV())
}

func (svc settings) UpdateAuthSettings(as *AuthSettings) error {
	vv, err := svc.internalSettings.FindByPrefix("auth.")
	if err != nil {
		return err
	}

	return as.ReadKV(vv.KV())
}

// AutoDiscovery orchestrates settings auto discovery
func (svc settings) AutoDiscovery() (err error) {
	var (
		current, discovered internalSettings.ValueSet
	)

	current, err = svc.internalSettings.FindByPrefix("")
	if err != nil {
		return
	}

	discovered, err = authSettingsAutoDiscovery(svc.logger, current)
	if err != nil || len(discovered) == 0 {
		return
	}

	return svc.internalSettings.BulkSet(discovered)
}
