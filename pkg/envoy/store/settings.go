package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	settingsState struct {
		cfg *EncoderConfig

		res types.SettingValueSet
		exs types.SettingValueSet
	}
)

var (
	// gSettingsState will aggregate all of the setting resources.
	gSettingsState *settingsState = nil
)

func NewSettingsState(res *resource.Settings, cfg *EncoderConfig) resourceState {
	return &settingsState{
		cfg: mergeConfig(cfg, res.Config()),

		res: res.Res,
	}
}

func (n *settingsState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Init global state
	if gSettingsState == nil {
		gSettingsState = &settingsState{
			cfg: n.cfg,
		}

		// Preload settings
		gSettingsState.exs, _, err = store.SearchSettings(ctx, s, types.SettingsFilter{})
		if err == store.ErrNotFound {
			gSettingsState.exs = make(types.SettingValueSet, 0, len(n.res))
		} else if err != nil {
			return err
		}
	}

	// Default values
	for _, s := range n.res {
		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = *now()
		}

		gSettingsState.res = append(gSettingsState.res, s)
	}

	return nil
}

func (n *settingsState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	ss := make(types.SettingValueSet, 0, len(n.res))

	for _, ns := range n.res {
		os := n.exs.First(ns.Name)
		if os != nil {
			// Update existing setting
			switch n.cfg.OnExisting {
			case resource.Skip,
				resource.MergeLeft:
				ss = append(ss, os)

			case resource.Replace,
				resource.MergeRight:
				ss = append(ss, ns)
			}
		} else {
			// Create fresh setting
			ss = append(ss, ns)
		}
	}

	return store.UpsertSetting(ctx, s, ss...)
}
