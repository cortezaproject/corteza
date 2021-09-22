package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (svc resourceTranslationsManager) loadReport(ctx context.Context, s store.Storer, reportID uint64) (*types.Report, error) {
	return store.LookupReportByID(ctx, s, reportID)
}

func (svc resourceTranslationsManager) reportExtended(ctx context.Context, res *types.Report) (out locale.ResourceTranslationSet, er error) {
	var (
		k types.LocaleKey
	)

	for _, tag := range svc.locale.Tags() {
		for i, projection := range res.Projections {
			projectionID := locale.ContentID(projection.ProjectionID, i)
			rpl := strings.NewReplacer(
				"{{projectionID}}", strconv.FormatUint(projectionID, 10),
			)

			// base stuff
			k = types.LocaleKeyReportProjectionTitle
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      rpl.Replace(k.Path),
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
			})

			k = types.LocaleKeyReportProjectionDescription
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      rpl.Replace(k.Path),
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
			})
		}
	}

	return
}
