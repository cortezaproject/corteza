package rest

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	"github.com/cortezaproject/corteza-server/system/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"strconv"
)

type (
	Settings struct {
		svc struct {
			settings settings.Service
			att      service.AttachmentService
		}
	}
)

func (Settings) New() *Settings {
	ctrl := &Settings{}
	ctrl.svc.settings = service.DefaultSettings
	ctrl.svc.att = service.DefaultAttachment

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

func (ctrl *Settings) Set(ctx context.Context, r *request.SettingsSet) (interface{}, error) {
	if r.Upload != nil {
		file, err := r.Upload.Open()

		if err != nil {
			return nil, err
		}

		defer file.Close()

		// @todo this whole attachment + settings logic must be moved to settings service
		//       this can be done when we generalize attachment handling
		//       and move that our of sys/msg/cmp to pkg
		att, err := ctrl.svc.att.With(ctx).CreateSettingsAttachment(
			r.Upload.Filename,
			r.Upload.Size,
			file,
			map[string]string{"key": r.Key, "ownedBy": strconv.FormatUint(r.OwnerID, 10)},
		)

		s := &settings.Value{Name: r.Key, OwnedBy: r.OwnerID}
		if err = s.SetValue(fmt.Sprintf("attachment:%d", att.ID)); err != nil {
			return nil, err
		}

		return s, ctrl.svc.settings.Set(ctx, s)
	}

	return nil, nil
}

// Current settings, structured
//
// This is available to all authenticated users
//
// @todo selectively apply subset of user's own settings (like ui.one.panes.*)
func (ctrl *Settings) Current(ctx context.Context, r *request.SettingsCurrent) (interface{}, error) {
	return service.CurrentSettings, nil
}
