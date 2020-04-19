package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/pkg/settings"
)

type (
	Settings struct {
		svc struct {
			settings settings.Service
		}
	}
)

func (Settings) New() *Settings {
	ctrl := &Settings{}
	ctrl.svc.settings = service.DefaultSettings

	return ctrl
}

func (ctrl *Settings) List(ctx context.Context, r *request.SettingsList) (interface{}, error) {
	if vv, err := ctrl.svc.settings.FindByPrefix(ctx, r.Prefix); err != nil {
		return nil, err
	} else {
		return vv, err
	}
}

func (ctrl *Settings) Update(ctx context.Context, r *request.SettingsUpdate) (interface{}, error) {
	if err := ctrl.svc.settings.BulkSet(ctx, r.Values); err != nil {
		return nil, err
	} else {
		return true, nil
	}
}

func (ctrl *Settings) Get(ctx context.Context, r *request.SettingsGet) (interface{}, error) {
	if v, err := ctrl.svc.settings.Get(ctx, r.Key, r.OwnerID); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}

// Current settings, structured
//
// This is available to all authenticated users
func (ctrl *Settings) Current(ctx context.Context, r *request.SettingsCurrent) (interface{}, error) {
	return service.CurrentSettings, nil
}
