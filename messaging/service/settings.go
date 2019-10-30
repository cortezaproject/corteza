package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	internalSettings "github.com/cortezaproject/corteza-server/pkg/settings"
)

type (
	// Wrapper service for messaging around internal settings service
	settings struct {
		ctx    context.Context
		db     *factory.DB
		logger *zap.Logger

		ac               settingsAccessController
		internalSettings internalSettings.Service
		current          *types.Settings
	}

	settingsAccessController interface {
		CanReadSettings(ctx context.Context) bool
		CanManageSettings(ctx context.Context) bool
	}

	SettingsService interface {
		With(ctx context.Context) *settings
		FindByPrefix(prefix string) (vv internalSettings.ValueSet, err error)
		Set(v *internalSettings.Value) (err error)
		BulkSet(vv internalSettings.ValueSet) (err error)
		Get(name string, ownedBy uint64) (out *internalSettings.Value, err error)
	}
)

func Settings(ctx context.Context, intSet internalSettings.Service, current *types.Settings) *settings {
	return (&settings{
		internalSettings: intSet,
		ac:               DefaultAccessControl,
		logger:           DefaultLogger.Named("settings"),
		current:          current,
	}).With(ctx)
}

func (svc settings) With(ctx context.Context) *settings {
	db := repository.DB(ctx)

	return &settings{
		ctx:    ctx,
		db:     db,
		ac:     svc.ac,
		logger: svc.logger,

		internalSettings: svc.internalSettings.With(ctx),

		current: svc.current,
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

// UpdateCurrent loads settings values from storage and updates current settings variable
//
// It accesses internal settings directly because
// we do not want any security checks for this
func (svc settings) UpdateCurrent() error {
	if vv, err := svc.internalSettings.FindByPrefix(""); err != nil {
		return err
	} else {
		return svc.updateCurrent(vv.KV())
	}
}

func (svc settings) Set(v *internalSettings.Value) (err error) {
	if !svc.ac.CanManageSettings(svc.ctx) {
		return errors.New("not allowed to manage settings")
	}

	if err = svc.internalSettings.Set(v); err != nil {
		return
	}

	return svc.updateCurrent(internalSettings.KV{v.Name: v.Value})
}

func (svc settings) BulkSet(vv internalSettings.ValueSet) (err error) {
	if !svc.ac.CanManageSettings(svc.ctx) {
		return errors.New("not allowed to manage settings")
	}

	var old internalSettings.ValueSet
	if old, err = svc.internalSettings.FindByPrefix(""); err != nil {
		return
	} else {
		vv = old.Changed(vv)
	}

	if err = svc.internalSettings.BulkSet(vv); err != nil {
		return
	}

	for _, v := range vv {
		svc.log(svc.ctx,
			zap.String("name", v.Name),
			zap.Stringer("value", v.Value)).Info("settings changed")
	}

	return svc.updateCurrent(vv.KV())
}

func (svc settings) updateCurrent(kv internalSettings.KV) (err error) {
	// update current settings with new values
	if err = kv.Decode(svc.current); err != nil {
		return
	}

	svc.log(svc.ctx).Info("current settings updated")
	return
}

func (svc settings) Get(name string, ownedBy uint64) (out *internalSettings.Value, err error) {
	if !svc.ac.CanReadSettings(svc.ctx) {
		return nil, errors.New("not allowed to read settings")
	}

	return svc.internalSettings.Get(name, ownedBy)
}
