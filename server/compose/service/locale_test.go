package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/locale"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// Most of a module's translatable strings hang off its fields — labels,
// descriptions, hints, select options, bool labels. The manager reads them off
// the loaded module, so a module loaded without its fields yields a set holding
// nothing but the module name and every field-level translator in the webapp
// opens empty.
func TestResourceTranslationsManagerModuleCoversFields(t *testing.T) {
	var (
		req = require.New(t)
		ctx = context.Background()

		ns = &types.Namespace{ID: nextID(), Name: "translations", Slug: "translations"}
	)

	s, err := sqlite.ConnectInMemoryWithDebug(ctx)
	req.NoError(err)
	req.NoError(store.Upgrade(ctx, zap.NewNop(), s))
	req.NoError(store.TruncateComposeNamespaces(ctx, s))
	req.NoError(store.TruncateComposeModules(ctx, s))
	req.NoError(store.TruncateComposeModuleFields(ctx, s))

	mod := &types.Module{ID: nextID(), NamespaceID: ns.ID, Name: "Probe", Handle: "probe"}

	selectField := &types.ModuleField{
		ID:          nextID(),
		ModuleID:    mod.ID,
		NamespaceID: ns.ID,
		Name:        "status",
		Label:       "Status",
		Kind:        "Select",
		Options: types.ModuleFieldOptions{
			"options": []any{
				map[string]any{"value": "new", "text": "New"},
			},
		},
	}
	boolField := &types.ModuleField{
		ID:          nextID(),
		ModuleID:    mod.ID,
		NamespaceID: ns.ID,
		Name:        "flag",
		Label:       "Flag",
		Kind:        "Bool",
		Options:     types.ModuleFieldOptions{"trueLabel": "Yes", "falseLabel": "No"},
	}

	req.NoError(store.CreateComposeNamespace(ctx, s, ns))
	req.NoError(store.CreateComposeModule(ctx, s, mod))
	req.NoError(store.CreateComposeModuleField(ctx, s, selectField, boolField))

	svc := resourceTranslationsManager{
		store:  s,
		locale: locale.Static(&locale.Language{Tag: language.English}),
	}

	set, err := svc.Module(ctx, ns.ID, mod.ID)
	req.NoError(err)

	keys := map[string][]string{}
	for _, tr := range set {
		keys[tr.Resource] = append(keys[tr.Resource], tr.Key)
	}

	req.Contains(keys[mod.ResourceTranslation()], "name")

	req.Subset(keys[selectField.ResourceTranslation()], []string{
		"label",
		"meta.description.view",
		"meta.hint.edit",
		"meta.options.new.text",
	})

	req.Subset(keys[boolField.ResourceTranslation()], []string{
		"label",
		"meta.bool.true.label",
		"meta.bool.false.label",
	})
}
