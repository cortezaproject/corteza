package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type (
	settings struct {
		store         settingsStore
		accessControl accessController
		logger        *zap.Logger

		// Holds reference to the "current" settings that
		// are used by the services
		current interface{}
	}

	//Service interface {
	//	FindByPrefix(ctx context.Context, pp ...string) (vv types.SettingValueSet, err error)
	//	BulkSet(ctx context.Context, vv types.SettingValue) (err error)
	//	Set(ctx context.Context, v *types.SettingValue) (err error)
	//	Get(ctx context.Context, name string, ownedBy uint64) (out *types.SettingValue, err error)
	//	Delete(ctx context.Context, name string, ownedBy uint64) error
	//	UpdateCurrent(ctx context.Context) error
	//}

	accessController interface {
		CanReadSettings(ctx context.Context) bool
		CanManageSettings(ctx context.Context) bool
	}
)

var (
	ErrNoReadPermission   = fmt.Errorf("not allowed to read settings")
	ErrNoManagePermission = fmt.Errorf("not allowed to manage settings")
)

func Settings(s settingsStore, log *zap.Logger, ac accessController, current interface{}) *settings {
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
		}
	)

	vv, _, err := svc.store.SearchSettings(ctx, f)
	return vv, err
}

func (svc settings) Get(ctx context.Context, name string, ownedBy uint64) (out *types.SettingValue, err error) {
	if !svc.accessControl.CanReadSettings(ctx) {
		return nil, ErrNoReadPermission
	}

	out, err = svc.store.LookupSettingByNameOwnedBy(ctx, name, ownedBy)
	if err != nil && err != store.ErrNotFound {
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
	current, err = svc.store.LookupSettingByNameOwnedBy(ctx, v.Name, v.OwnedBy)
	if err == store.ErrNotFound {
		v.UpdatedAt = time.Now()
		err = svc.store.CreateSetting(ctx, v)
	} else if err != nil {
		return err
	}

	if !current.Eq(v) {
		v.UpdatedAt = time.Now()
		err = svc.store.UpdateSetting(ctx, v)
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

	// Load current settings from repository
	// and get changed values
	var current, old, new types.SettingValueSet
	if current, err = svc.FindByPrefix(ctx); err != nil {
		return
	} else {
		vv = current.Changed(vv)
		_ = vv.Walk(func(v *types.SettingValue) error {
			v.UpdatedAt = time.Now()
			return nil
		})

		old = current.Old(vv)
		new = current.New(vv)
	}

	if err = svc.store.UpdateSetting(ctx, old...); err != nil {
		return
	}

	if err = svc.store.CreateSetting(ctx, new...); err != nil {
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

func (svc settings) Delete(ctx context.Context, name string, ownedBy uint64) (err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	var current *types.SettingValue
	if current, err = svc.store.LookupSettingByNameOwnedBy(ctx, name, ownedBy); err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return
	}

	err = svc.store.RemoveSetting(ctx, current)
	if err != nil {
		return
	}

	vv := types.SettingValueSet{current}

	svc.log(ctx,
		zap.String("name", name),
		zap.Uint64("owned-by", ownedBy)).Info("setting value removed")

	return svc.updateCurrent(ctx, vv)
}
