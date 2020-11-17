package store

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	settingsState struct {
		cfg *EncoderConfig

		res *resource.Settings
		ss  types.SettingValueSet
	}
)

func NewSettingsState(res *resource.Settings, cfg *EncoderConfig) resourceState {
	return &settingsState{
		cfg: cfg,

		res: res,
	}
}

func (n *settingsState) Prepare(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	// Preload settings
	n.ss, _, err = store.SearchSettings(ctx, s, types.SettingsFilter{})
	if err == store.ErrNotFound {
		n.ss = make(types.SettingValueSet, 0, len(n.res.Res))
	} else if err != nil {
		return err
	}

	// Default values
	for _, s := range n.res.Res {
		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = time.Now()
		}
	}

	// Nothing else to do.
	// Settings can't conflict either.

	return nil
}

func (n *settingsState) Encode(ctx context.Context, s store.Storer, state *envoy.ResourceState) (err error) {
	ss := make(types.SettingValueSet, 0, len(n.res.Res))

	for _, ns := range n.res.Res {
		os := n.ss.First(ns.Name)
		if os != nil {
			// Update existing setting
			switch n.cfg.OnExisting {
			case Skip:
				ss = append(ss, os)

			case Replace:
				ss = append(ss, ns)

			case MergeLeft:
				ss = append(ss, mergeSettings(os, ns))

			case MergeRight:
				ss = append(ss, mergeSettings(ns, os))
			}
		} else {
			// Create fresh setting
			ss = append(ss, ns)
		}
	}

	return store.UpsertSetting(ctx, s, ss...)
}

// mergeSettings merges b into a, prioritising a
func mergeSettings(a, b *types.SettingValue) *types.SettingValue {
	c := *a

	if len(c.Value) <= 0 {
		c.Value = b.Value
	}

	return &c
}
