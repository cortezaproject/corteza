package provision

import (
	"context"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/service"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

func migratePages(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	layouts, _, err := store.SearchComposePageLayouts(ctx, s, types.PageLayoutFilter{})
	if err != nil {
		return
	}

	// Probably already ran, no need to continue
	if len(layouts) > 0 {
		return
	}

	var (
		crs   *filter.PagingCursor
		pages types.PageSet
		auxf  types.PageFilter

		translations systemTypes.ResourceTranslationSet
	)

	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
		translations, err = getRelevantTranslations(ctx, s)
		if err != nil {
			return
		}

		for {
			pages, auxf, err = store.SearchComposePages(ctx, s, types.PageFilter{
				Deleted: filter.StateInclusive,
				Paging: filter.Paging{
					PageCursor: crs,
				},
			})
			if err != nil {
				return
			}
			if len(pages) == 0 {
				break
			}

			err = migratePageChunk(ctx, s, translations, pages)
			if err != nil {
				return
			}

			crs = auxf.Paging.PageCursor
			if crs == nil {
				break
			}
		}
		return nil
	})
}

func getRelevantTranslations(ctx context.Context, s store.Storer) (out systemTypes.ResourceTranslationSet, err error) {
	var (
		crs  *filter.PagingCursor
		ll   systemTypes.ResourceTranslationSet
		auxf systemTypes.ResourceTranslationFilter
	)

	for {
		ll, auxf, err = store.SearchResourceTranslations(ctx, s, systemTypes.ResourceTranslationFilter{
			Paging: filter.Paging{
				PageCursor: crs,
			},
		})
		if err != nil {
			return
		}
		if len(ll) == 0 {
			break
		}

		for _, l := range ll {
			if !strings.HasPrefix(l.K, "pageBlock.") {
				out = append(out, l)
			}
		}

		crs = auxf.Paging.PageCursor
		if crs == nil {
			break
		}
	}

	return
}

func migratePageChunk(ctx context.Context, s store.Storer, translations systemTypes.ResourceTranslationSet, pages types.PageSet) (err error) {
	n := now()
	for _, p := range pages {

		// Base layout
		ly := &types.PageLayout{
			ID: nextID(),

			NamespaceID: p.NamespaceID,
			PageID:      p.ID,
			Handle:      "primary",
			Primary:     true,

			Weight: 1,

			Meta: types.PageLayoutMeta{
				Title:       p.Title,
				Description: p.Description,
			},

			Config: types.PageLayoutConfig{
				Buttons: extractPageButtons(p),
			},

			CreatedAt: *n,
			UpdatedAt: n,
			DeletedAt: p.DeletedAt,
		}

		// Translations
		tt := translations.FilterResource(p.ResourceTranslation())

		// Button translations
		sr := strings.NewReplacer("recordToolbar", "config.buttons")
		for _, t := range tt {
			if !strings.HasPrefix(t.K, "recordToolbar.") {
				continue
			}

			t.K = sr.Replace(t.K)
			t.Resource = ly.ResourceTranslation()
		}

		// Title, description
		x := tt.FilterKey(types.LocaleKeyPageTitle.Path)
		for _, t := range x {
			tt = append(tt, &systemTypes.ResourceTranslation{
				ID:        nextID(),
				Resource:  ly.ResourceTranslation(),
				K:         types.LocaleKeyPageLayoutMetaTitle.Path,
				Lang:      t.Lang,
				Message:   t.Message,
				CreatedAt: *n,
			})
		}

		x = tt.FilterKey(types.LocaleKeyPageDescription.Path)
		for _, t := range x {
			tt = append(tt, &systemTypes.ResourceTranslation{
				ID:        nextID(),
				Resource:  ly.ResourceTranslation(),
				K:         types.LocaleKeyPageLayoutMetaDescription.Path,
				Lang:      t.Lang,
				Message:   t.Message,
				CreatedAt: *n,
			})
		}

		// Blocks
		for _, b := range p.Blocks {
			b.XYWH = adjustBlockScale(b.XYWH, 12, 48)

			ly.Blocks = append(ly.Blocks, types.PageLayoutBlock{
				BlockID: b.BlockID,
				XYWH:    b.XYWH,
			})
		}

		err = store.UpsertComposePage(ctx, s, p)
		if err != nil {
			return
		}

		err = store.UpsertComposePageLayout(ctx, s, ly)
		if err != nil {
			return
		}

		err = store.UpsertResourceTranslation(ctx, s, tt...)
		if err != nil {
			return
		}
	}

	return
}

func extractPageButtons(p *types.Page) (out types.PageLayoutButtonConfig) {
	ss := service.CurrentSettings.Compose.UI.RecordToolbar

	if p.ModuleID == 0 {
		return
	}

	out = types.PageLayoutButtonConfig{}
	out.New.Enabled = !ss.HideNew
	out.Edit.Enabled = !ss.HideEdit
	out.Submit.Enabled = !ss.HideSubmit
	out.Delete.Enabled = !ss.HideDelete
	out.Clone.Enabled = !ss.HideClone
	out.Back.Enabled = !ss.HideBack

	return
}

func adjustBlockScale(b [4]int, prev, new int) [4]int {
	for i, v := range b {
		b[i] = v * new / prev
	}

	return b
}
