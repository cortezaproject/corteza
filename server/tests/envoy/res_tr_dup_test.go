package envoy

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// Importing configs where the same (lang, resource, key) translation appears
// multiple times (across files or on repeated imports) must not trip the
// unique constraint on resource translations
func TestResTrImportDuplicates(t *testing.T) {
	var (
		ctx = context.Background()
		req = require.New(t)
	)

	cleanup(t)

	runImport := func(uri string) error {
		nodes, _, err := defaultEnvoy.Decode(ctx, envoyx.DecodeParams{
			Type: envoyx.DecodeTypeURI,
			Params: map[string]any{
				"uri": uri,
			},
		})
		if err != nil {
			return err
		}

		gg, err := defaultEnvoy.Bake(ctx, envoyx.EncodeParams{
			Type: envoyx.EncodeTypeStore,
			Params: map[string]any{
				"storer": defaultStore,
				"dal":    defaultDal,
			},
		}, nil, nodes...)
		if err != nil {
			return err
		}

		return defaultEnvoy.Encode(ctx, envoyx.EncodeParams{
			Type: envoyx.EncodeTypeStore,
			Params: map[string]any{
				"storer": defaultStore,
				"dal":    defaultDal,
			},
		}, gg)
	}

	t.Run("import with duplicated translations", func(t *testing.T) {
		req.NoError(runImport("file://testdata/locale_dup"))

		namespaces, _, err := store.SearchComposeNamespaces(ctx, defaultStore, types.NamespaceFilter{})
		req.NoError(err)
		req.Len(namespaces, 1)
		nsRes := fmt.Sprintf("corteza::compose:namespace/%d", namespaces[0].ID)

		ll, _, err := store.SearchResourceTranslations(ctx, defaultStore, systemTypes.ResourceTranslationFilter{})
		req.NoError(err)

		// Only one row per (lang, resource, key) may exist; the later
		// duplicate wins and may not be silently dropped
		en := ll.FilterLanguage(language.English)
		req.Len(en.FilterResource(nsRes).FilterKey("res_tr_1"), 1)
		req.Equal("res_tr_1_text_duplicate", en.FilterResource(nsRes).FilterKey("res_tr_1")[0].Message)
		req.Len(en.FilterResource(nsRes).FilterKey("res_tr_2"), 1)
	})

	t.Run("re-import over existing translations", func(t *testing.T) {
		req.NoError(runImport("file://testdata/locale_dup"))
	})
}
