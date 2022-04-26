package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/spf13/cast"
	"golang.org/x/text/language"
)

func (svc resourceTranslationsManager) moduleExtended(ctx context.Context, res *types.Module) (out locale.ResourceTranslationSet, err error) {
	var (
		k   types.LocaleKey
		set locale.ResourceTranslationSet
	)

	for _, tag := range svc.locale.Tags() {
		for _, f := range res.Fields {
			k = types.LocaleKeyModuleFieldLabel
			out = append(out, &locale.ResourceTranslation{
				Resource: f.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      k.Path,
				Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), k.Path),
			})

			k = types.LocaleKeyModuleFieldMetaDescriptionView
			out = append(out, &locale.ResourceTranslation{
				Resource: f.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      k.Path,
				Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), k.Path),
			})
			k = types.LocaleKeyModuleFieldMetaDescriptionEdit
			out = append(out, &locale.ResourceTranslation{
				Resource: f.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      k.Path,
				Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), k.Path),
			})

			k = types.LocaleKeyModuleFieldMetaHintView
			out = append(out, &locale.ResourceTranslation{
				Resource: f.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      k.Path,
				Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), k.Path),
			})
			k = types.LocaleKeyModuleFieldMetaHintEdit
			out = append(out, &locale.ResourceTranslation{
				Resource: f.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      k.Path,
				Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), k.Path),
			})

			// Expressions
			set, err = svc.moduleFieldExpressionsHandler(ctx, tag, f)
			if err != nil {
				return nil, err
			}
			out = append(out, set...)

			// Extra field bits
			set, err = svc.moduleFieldOptionsHandler(ctx, tag, f)
			if err != nil {
				return nil, err
			}
			out = append(out, set...)

			set, err = svc.moduleFieldBoolHandler(ctx, tag, f)
			if err != nil {
				return nil, err
			}
			out = append(out, set...)
		}
	}

	return out, nil
}

func (svc resourceTranslationsManager) moduleFieldExpressionsHandler(ctx context.Context, tag language.Tag, f *types.ModuleField) (locale.ResourceTranslationSet, error) {
	out := make(locale.ResourceTranslationSet, 0, 10)

	for i, v := range f.Expressions.Validators {
		vContentID := locale.ContentID(v.ValidatorID, i)
		rpl := strings.NewReplacer(
			"{{validatorID}}", strconv.FormatUint(vContentID, 10),
		)

		tKey := rpl.Replace(types.LocaleKeyModuleFieldExpressionValidatorValidatorIDError.Path)

		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      tKey,
			Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), tKey),
		})
	}

	return out, nil
}

func (svc resourceTranslationsManager) moduleFieldOptionsHandler(ctx context.Context, tag language.Tag, f *types.ModuleField) (locale.ResourceTranslationSet, error) {
	out := make(locale.ResourceTranslationSet, 0, 10)

	optsUnknown, has := f.Options["options"]
	if !has {
		return nil, nil
	}

	optsSlice, is := optsUnknown.([]interface{})
	if !is {
		return nil, nil
	}

	for _, optUnknown := range optsSlice {
		var value string

		// what is this we're dealing with?
		// slice of strings (values) or map (value+text)
		switch opt := optUnknown.(type) {
		case string:
			value = opt

		case map[string]interface{}:
			value, is = opt["value"].(string)
			if !is {
				continue
			}
		}

		trKey := strings.NewReplacer("{{value}}", value).Replace(types.LocaleKeyModuleFieldMetaOptionsValueText.Path)

		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      trKey,
			Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), trKey),
		})
	}

	return out, nil
}

func (svc resourceTranslationsManager) moduleFieldBoolHandler(ctx context.Context, tag language.Tag, f *types.ModuleField) (locale.ResourceTranslationSet, error) {
	if f.Kind != "Bool" {
		return nil, nil
	}

	out := make(locale.ResourceTranslationSet, 0, 2)

	trKey := strings.NewReplacer("{{value}}", "true").Replace(types.LocaleKeyModuleFieldMetaBoolValueLabel.Path)
	out = append(out, &locale.ResourceTranslation{
		Resource: f.ResourceTranslation(),
		Lang:     tag.String(),
		Key:      trKey,
		Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), trKey),
	})

	trKey = strings.NewReplacer("{{value}}", "false").Replace(types.LocaleKeyModuleFieldMetaBoolValueLabel.Path)
	out = append(out, &locale.ResourceTranslation{
		Resource: f.ResourceTranslation(),
		Lang:     tag.String(),
		Key:      trKey,
		Msg:      svc.locale.TResourceFor(tag, f.ResourceTranslation(), trKey),
	})

	return out, nil
}

func (svc resourceTranslationsManager) pageExtended(ctx context.Context, res *types.Page) (out locale.ResourceTranslationSet, err error) {
	var (
		k   types.LocaleKey
		aux []*locale.ResourceTranslation
	)

	for _, tag := range svc.locale.Tags() {
		// Button translations
		if res.Config.Buttons != nil {
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      types.LocaleKeyPageRecordToolbarNewLabel.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageRecordToolbarNewLabel.Path),
			})

			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      types.LocaleKeyPageRecordToolbarEditLabel.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageRecordToolbarEditLabel.Path),
			})

			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      types.LocaleKeyPageRecordToolbarSubmitLabel.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageRecordToolbarSubmitLabel.Path),
			})

			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      types.LocaleKeyPageRecordToolbarDeleteLabel.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageRecordToolbarDeleteLabel.Path),
			})

			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      types.LocaleKeyPageRecordToolbarCloneLabel.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageRecordToolbarCloneLabel.Path),
			})

			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      types.LocaleKeyPageRecordToolbarBackLabel.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageRecordToolbarBackLabel.Path),
			})
		}

		for i, block := range res.Blocks {
			pbContentID := locale.ContentID(block.BlockID, i)
			rpl := strings.NewReplacer(
				"{{blockID}}", strconv.FormatUint(pbContentID, 10),
			)

			// base stuff
			k = types.LocaleKeyPagePageBlockBlockIDTitle
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      rpl.Replace(k.Path),
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
			})

			k = types.LocaleKeyPagePageBlockBlockIDDescription
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      rpl.Replace(k.Path),
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
			})

			switch block.Kind {
			case "Automation":
				aux, err = svc.pageExtendedAutomationBlock(tag, res, block, pbContentID)
				if err != nil {
					return nil, err
				}

				out = append(out, aux...)
			case "Content":
				k = types.LocaleKeyPagePageBlockBlockIDContentBody
				out = append(out, &locale.ResourceTranslation{
					Resource: res.ResourceTranslation(),
					Lang:     tag.String(),
					Key:      rpl.Replace(k.Path),
					Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
				})

			}
		}
	}

	return
}

func (svc resourceTranslationsManager) chartExtended(_ context.Context, res *types.Chart) (out locale.ResourceTranslationSet, err error) {
	var (
		yAxisLabelK  = types.LocaleKeyChartYAxisLabel
		metricLabelK = types.LocaleKeyChartMetricsMetricIDLabel
	)

	for _, report := range res.Config.Reports {
		for _, tag := range svc.locale.Tags() {
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      yAxisLabelK.Path,
				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), yAxisLabelK.Path),
			})
		}

		for _, metric := range report.Metrics {
			if _, ok := metric["metricID"]; ok {
				mpl := strings.NewReplacer(
					"{{metricID}}", metric["metricID"].(string),
				)

				for _, tag := range svc.locale.Tags() {
					out = append(out, &locale.ResourceTranslation{
						Resource: res.ResourceTranslation(),
						Lang:     tag.String(),
						Key:      mpl.Replace(metricLabelK.Path),
						Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), mpl.Replace(metricLabelK.Path)),
					})
				}
			}
		}
	}

	return
}

func (svc resourceTranslationsManager) pageExtendedAutomationBlock(tag language.Tag, res *types.Page, block types.PageBlock, blockID uint64) (locale.ResourceTranslationSet, error) {
	var (
		k     = types.LocaleKeyPagePageBlockBlockIDButtonButtonIDLabel
		out   = make(locale.ResourceTranslationSet, 0, 10)
		bb, _ = block.Options["buttons"].([]interface{})
	)

	for j, auxBtn := range bb {
		btn := auxBtn.(map[string]interface{})

		bContentID := uint64(0)
		if aux, ok := btn["buttonID"]; ok {
			bContentID = cast.ToUint64(aux)
		}

		bContentID = locale.ContentID(bContentID, j)

		rpl := strings.NewReplacer(
			"{{blockID}}", strconv.FormatUint(blockID, 10),
			"{{buttonID}}", strconv.FormatUint(bContentID, 10),
		)

		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      rpl.Replace(k.Path),
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
		})
	}

	return out, nil
}

// Helper loaders

func (svc resourceTranslationsManager) loadModule(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (m *types.Module, err error) {
	return loadModule(ctx, s, moduleID)
}

func (svc resourceTranslationsManager) loadNamespace(ctx context.Context, s store.Storer, namespaceID uint64) (m *types.Namespace, err error) {
	return loadNamespace(ctx, s, namespaceID)
}

func (svc resourceTranslationsManager) loadPage(ctx context.Context, s store.Storer, namespaceID, pageID uint64) (m *types.Page, err error) {
	_, m, err = loadPage(ctx, s, namespaceID, pageID)
	return m, err
}

func (svc resourceTranslationsManager) loadChart(ctx context.Context, s store.Storer, namespaceID, chartID uint64) (m *types.Chart, err error) {
	_, m, err = loadChart(ctx, s, namespaceID, chartID)
	return m, err
}
