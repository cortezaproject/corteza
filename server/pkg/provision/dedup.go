package provision

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/store"
	"go.uber.org/zap"
)

func invalidateDedupRules(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	var ll types.ModuleSet

	if ll, _, err = s.SearchComposeModules(ctx, types.ModuleFilter{}); err != nil {
		return
	}

	// find only the invalid ones and empty the rules, leave the valid ones
	// as they were
	ll, _ = ll.Filter(func(m *types.Module) (bool, error) {
		return m.Config.RecordDeDup.Rules.Validate() != nil, nil
	})

	ll.Walk(func(m *types.Module) error {
		m.Config.RecordDeDup.Rules = types.DeDupRuleSet{}
		return nil
	})

	err = s.UpdateComposeModule(ctx, ll...)

	return
}
