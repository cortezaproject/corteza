package store

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

var (
	gSettings map[string]bool
)

func NewSettingFromResource(res *resource.Setting, cfg *EncoderConfig) resourceState {
	return &setting{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,
		st:  &res.Res,
	}
}

func (n *setting) Prepare(ctx context.Context, pl *payload) (err error) {
	// Init global state
	if gSettings == nil {
		gSettings = make(map[string]bool)
		ss, _, err := store.SearchSettings(ctx, pl.s, types.SettingsFilter{})
		if err != store.ErrNotFound && err != nil {
			return err
		}
		for _, s := range ss {
			gSettings[s.Name] = true
		}
	}

	return nil
}

func (n *setting) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.st

	us, err := resolveUserstamps(ctx, pl.s, pl.state.ParentResources, n.res.Userstamps())
	if err != nil {
		return err
	}

	ts := n.res.Timestamps()
	if ts != nil {
		if ts.UpdatedAt != nil {
			res.UpdatedAt = *ts.UpdatedAt.T
		} else {
			res.UpdatedAt = *now()
		}
	} else {
		res.UpdatedAt = *now()
	}

	if us != nil {
		if us.OwnedBy != nil {
			res.OwnedBy = us.OwnedBy.UserID
		}
		if us.UpdatedBy != nil {
			res.UpdatedBy = us.UpdatedBy.UserID
		}
	}

	if _, exists := gSettings[res.Name]; !exists {
		return store.CreateSetting(ctx, pl.s, res)
	}

	// On existing setting, replace/merge right basically overwrites the existing value;
	// otherwise, the new value is ignored.
	switch n.cfg.OnExisting {
	case resource.Replace,
		resource.MergeRight:
		return store.UpdateSetting(ctx, pl.s, res)
	}

	return nil
}
