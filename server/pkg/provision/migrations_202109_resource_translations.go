package provision

import (
	"context"
	"fmt"
	"strings"

	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/store"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

// Migrates resource translations from the resource
// struct to the dedicated store (table)
//
// While doing this, we also modify some resource substructure:
//  - page blocks (assign page block IDs
//  - automation buttons on page blocks (assign automation button IDs)
//  - field expressions (assign expression IDs)
//
// Note: we will migrate all translations to current default language
// If you do not like that, shut down Corteza after migrations and fix this directly in the store
func migrateResourceTranslations(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	log.Info("migrating resource locales")

	var (
		migrated = make(map[string]bool)
	)

	set, _, err := store.SearchResourceTranslations(ctx, s, sysTypes.ResourceTranslationFilter{})
	set.Walk(func(r *sysTypes.ResourceTranslation) error {
		var pos = strings.Index(r.Resource, "/")
		if pos < 0 {
			return nil
		}

		migrated[r.Resource[0:pos]] = true
		return nil
	})

	if !migrated[cmpTypes.NamespaceResourceTranslationType] {
		if err = migrateComposeNamespaceResourceTranslations(ctx, log, s); err != nil {
			return
		}
	}

	if !migrated[cmpTypes.ModuleResourceTranslationType] {
		// @todo migrate module locales
		if err = migrateComposeModuleResourceTranslations(ctx, log, s); err != nil {
			return
		}
	}

	if !migrated[cmpTypes.PageResourceTranslationType] {
		if err = migrateComposePageResourceTranslations(ctx, log, s); err != nil {
			return
		}
	}

	return
}

// helper fn that creates ResourceTranslation from the given resource, key and message
// and default language
func makeResourceTranslation(lr interface{ ResourceTranslation() string }, k, m string) *sysTypes.ResourceTranslation {
	return &sysTypes.ResourceTranslation{
		ID:        id.Next(),
		Lang:      sysTypes.Lang{Tag: locale.Global().Default().Tag},
		Resource:  lr.ResourceTranslation(),
		K:         k,
		Message:   m,
		CreatedAt: *now(),
	}
}

// migrate resource translations for compose namespace
func migrateComposeNamespaceResourceTranslations(ctx context.Context, log *zap.Logger, s store.Storer) error {
	set, _, err := store.SearchComposeNamespaces(ctx, s, cmpTypes.NamespaceFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return err
	}

	log.Info("migrating compose namespaces", zap.Int("count", len(set)))

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return set.Walk(func(res *cmpTypes.Namespace) error {
			return store.CreateResourceTranslation(ctx, s,
				makeResourceTranslation(res, "name", res.Name),
				makeResourceTranslation(res, "subtitle", res.Meta.Subtitle),
				makeResourceTranslation(res, "description", res.Meta.Description),
			)
		})
	})
}

// migrate resource translations for compose module
func migrateComposeModuleResourceTranslations(ctx context.Context, log *zap.Logger, s store.Storer) error {
	set, _, err := store.SearchComposeModules(ctx, s, cmpTypes.ModuleFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return err
	}

	log.Info("migrating compose modules", zap.Int("count", len(set)))

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return set.Walk(func(res *cmpTypes.Module) (err error) {
			if err = migrateComposeModuleFieldResourceTranslations(ctx, s, res.NamespaceID, res.ID); err != nil {
				return err
			}

			err = store.CreateResourceTranslation(ctx, s,
				makeResourceTranslation(res, "name", res.Name),
			)

			if err != nil {
				return err
			}

			return store.UpdateComposeModule(ctx, s, res)
		})
	})
}

// migrate module fields translations for compose module field
//
// Adds validatorID on field expressions
func migrateComposeModuleFieldResourceTranslations(ctx context.Context, s store.Storer, namespaceID uint64, moduleIDs ...uint64) error {
	set, _, err := store.SearchComposeModuleFields(ctx, s, cmpTypes.ModuleFieldFilter{ModuleID: moduleIDs})
	if err != nil {
		return err
	}

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return set.Walk(func(res *cmpTypes.ModuleField) (err error) {
			tt := sysTypes.ResourceTranslationSet{}

			// store does not contain this info,
			// but we need it to generate  resource id
			res.NamespaceID = namespaceID
			tt = append(tt, makeResourceTranslation(res, "label", res.Label))

			var update = len(res.Expressions.Validators) > 0

			for i := range res.Expressions.Validators {
				validatorID := i + 1
				res.Expressions.Validators[i].ValidatorID = uint64(validatorID)
				tt = append(tt, makeResourceTranslation(res, fmt.Sprintf("expression.validator.%d.error", validatorID), res.Expressions.Validators[i].Error))
			}

			if err = store.CreateResourceTranslation(ctx, s, tt...); err != nil {
				return
			}

			if update {
				if err = store.UpdateComposeModuleField(ctx, s, res); err != nil {
					return
				}
			}

			return nil
		})
	})
}

// migrate resource translations for compose module
func migrateComposePageResourceTranslations(ctx context.Context, log *zap.Logger, s store.Storer) error {
	var tt sysTypes.ResourceTranslationSet
	set, _, err := store.SearchComposePages(ctx, s, cmpTypes.PageFilter{Deleted: filter.StateInclusive})
	if err != nil {
		return err
	}

	// @todo migrate page locales
	// @todo page block IDs
	// @todo automation button IDs

	log.Info("migrating compose pages", zap.Int("count", len(set)))

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return set.Walk(func(res *cmpTypes.Page) (err error) {
			tt = sysTypes.ResourceTranslationSet{
				makeResourceTranslation(res, "title", res.Title),
				makeResourceTranslation(res, "description", res.Description),
			}

			if pbtt, err := convertComposePageBlockTranslations(res); err != nil {
				return err
			} else {
				tt = append(tt, pbtt...)
			}

			if err = store.CreateResourceTranslation(ctx, s, tt...); err != nil {
				return
			}

			if err = store.UpdateComposePage(ctx, s, res); err != nil {
				return
			}

			return
		})
	})
}

// collects translations for compose page blocks and alters them (adding block ID)
func convertComposePageBlockTranslations(res *cmpTypes.Page) (sysTypes.ResourceTranslationSet, error) {
	var tt sysTypes.ResourceTranslationSet

	for i, b := range res.Blocks {
		blockID := i + 1
		res.Blocks[i].BlockID = uint64(blockID)
		pfx := fmt.Sprintf("pageBlock.%d.", blockID)

		tt = append(tt,
			makeResourceTranslation(res, pfx+"title", b.Title),
			makeResourceTranslation(res, pfx+"description", b.Description),
		)

		switch b.Kind {
		case "Automation":
			att, bb, err := convertComposeAutomationPageBlockTranslations(res, pfx, res.Blocks[i].Options)
			if err != nil {
				return nil, err
			}

			res.Blocks[i].Options["buttons"] = bb
			tt = append(tt, att...)
		}
	}

	return tt, nil
}

// takes automation options for page block and returns translations + fixed set of buttons (containing content ID)
func convertComposeAutomationPageBlockTranslations(res *cmpTypes.Page, pfx string, opt map[string]interface{}) (tt sysTypes.ResourceTranslationSet, bb []interface{}, err error) {
	if tmp, ok := opt["buttons"]; !ok {
		return nil, nil, nil
	} else if bb, ok = tmp.([]interface{}); !ok {
		return nil, nil, nil
	}

	for i := range bb {
		b, is := bb[i].(map[string]interface{})
		if !is {
			continue
		}

		buttonID := i + 1

		b["buttonID"] = buttonID
		bb[i] = b

		var label string
		if _, has := b["label"]; has {
			label = b["label"].(string)
		}

		tt = append(tt, makeResourceTranslation(res, fmt.Sprintf("%sbuttons.%d.label", pfx, buttonID), label))
	}

	return
}
