package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	settings struct {
		store         store.Settings
		accessControl accessController
		logger        *zap.Logger

		// Holds reference to the "current" settings that
		// are used by the services
		current interface{}
	}

	accessController interface {
		CanReadSettings(ctx context.Context) bool
		CanManageSettings(ctx context.Context) bool
	}
)

var (
	ErrNoReadPermission   = fmt.Errorf("not allowed to read settings")
	ErrNoManagePermission = fmt.Errorf("not allowed to manage settings")
)

func Settings(s store.Settings, log *zap.Logger, ac accessController, current interface{}) *settings {
	svc := &settings{
		store:         s,
		accessControl: ac,
		logger:        log.Named("settings"),
		current:       current,
	}

	return svc
}

func (svc settings) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc settings) FindByPrefix(ctx context.Context, pp ...string) (types.SettingValueSet, error) {
	if !svc.accessControl.CanReadSettings(ctx) {
		return nil, ErrNoReadPermission
	}

	return svc.findByPrefix(ctx, pp...)
}

func (svc settings) findByPrefix(ctx context.Context, pp ...string) (types.SettingValueSet, error) {
	var (
		f = types.SettingsFilter{
			Prefix: strings.Join(pp, "."),
			Check:  func(*types.SettingValue) (bool, error) { return true, nil },
		}
	)

	vv, _, err := store.SearchSettings(ctx, svc.store, f)
	return vv, err
}

func (svc settings) Get(ctx context.Context, name string, ownedBy uint64) (out *types.SettingValue, err error) {
	if !svc.accessControl.CanReadSettings(ctx) {
		return nil, ErrNoReadPermission
	}

	out, err = store.LookupSettingByNameOwnedBy(ctx, svc.store, name, ownedBy)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	return out, nil
}

func (svc settings) UpdateCurrent(ctx context.Context) error {
	if vv, err := svc.findByPrefix(ctx); err != nil {
		return err
	} else {
		return svc.updateCurrent(ctx, vv)
	}
}

func (svc settings) updateCurrent(ctx context.Context, vv types.SettingValueSet) (err error) {
	// update current settings with new values
	if err = vv.KV().Decode(svc.current); err != nil {
		return
	}

	return
}

func (svc settings) Set(ctx context.Context, v *types.SettingValue) (err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	var current *types.SettingValue
	current, err = store.LookupSettingByNameOwnedBy(ctx, svc.store, v.Name, v.OwnedBy)
	if errors.IsNotFound(err) {
		v.UpdatedAt = *now()
		err = store.CreateSetting(ctx, svc.store, v)
	} else if err != nil {
		return err
	}

	if !current.Eq(v) {
		v.UpdatedAt = *now()
		err = store.UpdateSetting(ctx, svc.store, v)
	}

	if err != nil || current.Eq(v) {
		// Return on error or when there is nothing to update (same value)
		return
	}

	svc.logChange(ctx, types.SettingValueSet{v})
	return svc.updateCurrent(ctx, types.SettingValueSet{v})
}

func (svc settings) BulkSet(ctx context.Context, vv types.SettingValueSet) (err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	// Load current settings and get changed values
	var current, old, new types.SettingValueSet
	if current, err = svc.FindByPrefix(ctx); err != nil {
		return
	} else {
		vv = current.Changed(vv)
		_ = vv.Walk(func(v *types.SettingValue) error {
			v.UpdatedAt = *now()
			return nil
		})

		old = current.Old(vv)
		new = current.New(vv)
	}

	if err = store.UpdateSetting(ctx, svc.store, old...); err != nil {
		return
	}

	if err = store.CreateSetting(ctx, svc.store, new...); err != nil {
		return
	}

	svc.logChange(ctx, vv)

	return svc.updateCurrent(ctx, vv)
}

func (svc settings) logChange(ctx context.Context, vv types.SettingValueSet) {
	for _, v := range vv {
		svc.log(ctx,
			zap.String("name", v.Name),
			zap.Uint64("owned-by", v.OwnedBy),
			zap.Stringer("value", v.Value)).
			WithOptions(zap.AddCallerSkip(1)).
			Debug("setting value updated")
	}
}

func (svc settings) Delete(ctx context.Context, name string, ownedBy uint64) error {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	current, err := store.LookupSettingByNameOwnedBy(ctx, svc.store, name, ownedBy)
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}

	err = store.DeleteSetting(ctx, svc.store, current)
	if err != nil {
		return err
	}

	vv := types.SettingValueSet{current}

	svc.log(ctx,
		zap.String("name", name),
		zap.Uint64("owned-by", ownedBy)).Info("setting value removed")

	return svc.updateCurrent(ctx, vv)
}
