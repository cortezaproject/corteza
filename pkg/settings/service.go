package settings

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	service struct {
		repository    Repository
		accessControl accessController
		logger        *zap.Logger
		current       interface{}
	}

	Service interface {
		FindByPrefix(ctx context.Context, pp ...string) (vv ValueSet, err error)
		BulkSet(ctx context.Context, vv ValueSet) (err error)
		Set(ctx context.Context, v *Value) (err error)
		Get(ctx context.Context, name string, ownedBy uint64) (out *Value, err error)
		Delete(ctx context.Context, name string, ownedBy uint64) error
		UpdateCurrent(ctx context.Context) error
	}

	accessController interface {
		CanReadSettings(ctx context.Context) bool
		CanManageSettings(ctx context.Context) bool
	}
)

var (
	ErrNoReadPermission   = errors.New("not allowed to read settings")
	ErrNoManagePermission = errors.New("not allowed to manage settings")
)

func NewService(r Repository, log *zap.Logger, ac accessController, current interface{}) *service {
	svc := &service{
		repository:    r,
		accessControl: ac,
		logger:        log.Named("settings"),
		current:       current,
	}

	return svc
}

func (svc service) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc service) FindByPrefix(ctx context.Context, pp ...string) (ValueSet, error) {
	if !svc.accessControl.CanReadSettings(ctx) {
		return nil, ErrNoReadPermission
	}

	var (
		f = Filter{
			Prefix: strings.Join(pp, "."),
		}
	)

	return svc.repository.With(ctx).Find(f)
}

func (svc service) Get(ctx context.Context, name string, ownedBy uint64) (out *Value, err error) {
	if !svc.accessControl.CanReadSettings(ctx) {
		return nil, ErrNoReadPermission
	}

	return svc.repository.With(ctx).Get(name, ownedBy)
}

func (svc service) UpdateCurrent(ctx context.Context) error {
	if vv, err := svc.FindByPrefix(ctx); err != nil {
		return err
	} else {
		return svc.updateCurrent(ctx, vv)
	}
}

func (svc service) updateCurrent(ctx context.Context, vv ValueSet) (err error) {
	// update current settings with new values
	if err = vv.KV().Decode(svc.current); err != nil {
		return
	}

	return
}

func (svc service) Set(ctx context.Context, v *Value) (err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	var current *Value
	current, err = svc.repository.With(ctx).Get(v.Name, v.OwnedBy)
	if err != nil || current.Eq(v) {
		// Return on error or when there is nothing to update (same value)
		return
	}

	err = svc.repository.With(ctx).Set(v)
	if err != nil {
		return
	}

	vv := ValueSet{v}

	svc.logChange(ctx, vv)

	return svc.updateCurrent(ctx, vv)
}

func (svc service) BulkSet(ctx context.Context, vv ValueSet) (err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	// Load current settings from repository
	// and get changed values
	var current ValueSet
	if current, err = svc.FindByPrefix(ctx); err != nil {
		return
	} else {
		vv = current.Changed(vv)
	}

	err = svc.repository.With(ctx).BulkSet(vv)
	if err != nil {
		return
	}

	svc.logChange(ctx, vv)

	return svc.updateCurrent(ctx, vv)
}

func (svc service) logChange(ctx context.Context, vv ValueSet) {
	for _, v := range vv {
		svc.log(ctx,
			zap.String("name", v.Name),
			zap.Uint64("owned-by", v.OwnedBy),
			zap.Stringer("value", v.Value)).Info("setting value updated")
	}
}

func (svc service) Delete(ctx context.Context, name string, ownedBy uint64) (err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return ErrNoManagePermission
	}

	var current *Value
	current, err = svc.repository.With(ctx).Get(name, ownedBy)
	if err != nil || current == nil {
		// Return on error or when there is nothing to delete (there is no value)
		return
	}

	err = svc.repository.With(ctx).Delete(name, ownedBy)
	if err != nil {
		return
	}

	vv := ValueSet{&Value{Name: name, OwnedBy: ownedBy}}

	svc.log(ctx,
		zap.String("name", name),
		zap.Uint64("owned-by", ownedBy)).Info("setting value removed")

	return svc.updateCurrent(ctx, vv)
}
