package rest

import (
	"context"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
)

var _ = errors.Wrap

type (
	Settings struct {
		svc struct {
			settings service.SettingsService
			auth service.AuthService
		}
	}
)

func (Settings) New() *Settings {
	ctrl := &Settings{}
	ctrl.svc.settings = service.DefaultSettings
	ctrl.svc.auth = service.DefaultAuth

	return ctrl
}

func (ctrl *Settings) List(ctx context.Context, r *request.SettingsList) (interface{}, error) {
	if vv, err := ctrl.svc.settings.With(ctx).FindByPrefix(r.Prefix); err != nil {
		return nil, err
	} else {
		return vv, err
	}
}

func (ctrl *Settings) Update(ctx context.Context, r *request.SettingsUpdate) (interface{}, error) {
	if err := ctrl.svc.settings.With(ctx).BulkSet(r.Values); err != nil {
		return nil, err
	} else {
		return true, nil
	}
}

func (ctrl *Settings) Get(ctx context.Context, r *request.SettingsGet) (interface{}, error) {
	if v, err := ctrl.svc.settings.With(ctx).Get(r.Key, r.OwnerID); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}

func (ctrl *Settings) Current(ctx context.Context, r *request.SettingsCurrent) (interface{}, error) {
	if err := ctrl.svc.auth.With(ctx).CanRegister(); err != nil {
		// Append signup warning to current settings
		service.CurrentSettings.Auth.SignupWarning = err.Error()
	}

	return service.CurrentSettings, nil
}
