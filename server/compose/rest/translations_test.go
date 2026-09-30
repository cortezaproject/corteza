package rest

import (
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/stretchr/testify/require"
)

func Test_checkTranslatedResources(t *testing.T) {
	var (
		mod     = &types.Module{ID: 2, NamespaceID: 1}
		field   = &types.ModuleField{ID: 3, ModuleID: 2, NamespaceID: 1}
		foreign = &types.Module{ID: 9, NamespaceID: 8}

		allowed = []string{mod.ResourceTranslation(), field.ResourceTranslation()}
	)

	require.NoError(t, checkTranslatedResources(nil, allowed...))

	require.NoError(t, checkTranslatedResources(locale.ResourceTranslationSet{
		{Resource: mod.ResourceTranslation(), Key: "name", Msg: "Module"},
		{Resource: field.ResourceTranslation(), Key: "label", Msg: "Field"},
	}, allowed...))

	require.Error(t, checkTranslatedResources(locale.ResourceTranslationSet{
		{Resource: mod.ResourceTranslation(), Key: "name", Msg: "Module"},
		{Resource: foreign.ResourceTranslation(), Key: "name", Msg: "defaced"},
	}, allowed...))

	require.Error(t, checkTranslatedResources(locale.ResourceTranslationSet{nil}, allowed...))

	// layouts of the page are translated together with the page
	var (
		page          = &types.Page{ID: 5, NamespaceID: 1}
		layout        = &types.PageLayout{ID: 6, PageID: 5, NamespaceID: 1}
		foreignLayout = &types.PageLayout{ID: 7, PageID: 55, NamespaceID: 1}

		pageAllowed = []string{page.ResourceTranslation(), types.PageLayoutResourceTranslation(1, 5, 0)[:len(types.PageLayoutResourceTranslation(1, 5, 0))-1]}
	)

	require.NoError(t, checkTranslatedResources(locale.ResourceTranslationSet{
		{Resource: page.ResourceTranslation()},
		{Resource: layout.ResourceTranslation()},
	}, pageAllowed...))

	require.Error(t, checkTranslatedResources(locale.ResourceTranslationSet{
		{Resource: foreignLayout.ResourceTranslation()},
	}, pageAllowed...))

	// prefix alone is not a resource
	require.Error(t, checkTranslatedResources(locale.ResourceTranslationSet{{Resource: pageAllowed[1]}}, pageAllowed...))
}
