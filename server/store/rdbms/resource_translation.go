package rdbms

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/system/types"
	"golang.org/x/text/language"
)

func (s Store) convertResourceTranslationFilter(f types.ResourceTranslationFilter) (query squirrel.SelectBuilder, err error) {
	query = s.resourceTranslationsSelectBuilder()

	query = filter.StateCondition(query, "rt.deleted_at", f.Deleted)

	if len(f.TranslationID) > 0 {
		query = query.Where(squirrel.Eq{"rt.id": f.TranslationID})
	}

	if f.Lang != "" {
		query = query.Where(squirrel.Eq{"rt.lang": f.Lang})
	}

	if f.Resource != "" {
		query = query.Where(squirrel.Eq{"rt.resource": f.Resource})
	}

	if f.ResourceType != "" {
		query = query.Where(squirrel.Like{"rt.resource": f.ResourceType + "%"})
	}

	return
}

// TransformResource converts raw resource translations into the format used by
// the locale package.
//
// @todo this function knows too much (locale pkg), move it out of store
func (s *Store) TransformResource(ctx context.Context, lang language.Tag) (out map[string]map[string]*locale.ResourceTranslation, err error) {
	out = make(map[string]map[string]*locale.ResourceTranslation)
	var cc types.ResourceTranslationSet
	var ok bool

	cc, _, err = s.SearchResourceTranslations(ctx, types.ResourceTranslationFilter{
		Lang: lang.String(),
	})
	if err != nil {
		return
	}

	for _, c := range cc {
		if _, ok = out[c.Resource]; !ok {
			out[c.Resource] = make(map[string]*locale.ResourceTranslation)
		}

		out[c.Resource][c.K] = &locale.ResourceTranslation{
			Resource: c.Resource,
			Lang:     c.Lang.String(),
			Key:      c.K,
			Msg:      c.Message,
		}
	}

	return
}
