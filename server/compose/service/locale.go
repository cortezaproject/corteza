package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/store"
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

	// We need to get page layouts to include them also
	// @todo refactor the base logic to simplify this; it's suboptimal
	layouts, _, err := store.SearchComposePageLayouts(ctx, svc.store, types.PageLayoutFilter{
		NamespaceID: res.NamespaceID,
		PageID:      res.ID,
	})
	if err != nil {
		return
	}

	for _, tag := range svc.locale.Tags() {
		for _, l := range layouts {
			aux, err = svc.PageLayout(ctx, res.NamespaceID, res.ID, l.ID)
			if err != nil {
				return
			}
			out = append(out, aux...)
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
			case "RecordList":
				aux, err = svc.pageExtendedRecordListBlock(tag, res, block, pbContentID)
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

func (svc resourceTranslationsManager) pageLayoutExtended(ctx context.Context, res *types.PageLayout) (out locale.ResourceTranslationSet, err error) {
	var (
		k types.LocaleKey
	)

	for _, tag := range svc.locale.Tags() {
		// Standard buttons
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      types.LocaleKeyPageLayoutConfigButtonsNewLabel.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageLayoutConfigButtonsNewLabel.Path),
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      types.LocaleKeyPageLayoutConfigButtonsEditLabel.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageLayoutConfigButtonsEditLabel.Path),
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      types.LocaleKeyPageLayoutConfigButtonsSubmitLabel.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageLayoutConfigButtonsSubmitLabel.Path),
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      types.LocaleKeyPageLayoutConfigButtonsDeleteLabel.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageLayoutConfigButtonsDeleteLabel.Path),
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      types.LocaleKeyPageLayoutConfigButtonsCloneLabel.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageLayoutConfigButtonsCloneLabel.Path),
		})

		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      types.LocaleKeyPageLayoutConfigButtonsBackLabel.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), types.LocaleKeyPageLayoutConfigButtonsBackLabel.Path),
		})

		// Actions
		for i, action := range res.Config.Actions {
			acContentID := locale.ContentID(action.ActionID, i)
			rpl := strings.NewReplacer(
				"{{actionID}}", strconv.FormatUint(acContentID, 10),
			)

			k = types.LocaleKeyPageLayoutConfigActionsActionIDMetaLabel
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

func (svc resourceTranslationsManager) chartExtended(_ context.Context, res *types.Chart) (out locale.ResourceTranslationSet, err error) {
	var (
		yAxisLabelK   = types.LocaleKeyChartYAxisLabel
		metricLabelK  = types.LocaleKeyChartMetricsMetricIDLabel
		dimStepLabelK = types.LocaleKeyChartDimensionsDimensionIDMetaStepsStepIDLabel
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

		report.WalkMetrics(func(metricID string, _ map[string]interface{}) {
			mpl := strings.NewReplacer(
				"{{metricID}}", metricID,
			)

			for _, tag := range svc.locale.Tags() {
				out = append(out, &locale.ResourceTranslation{
					Resource: res.ResourceTranslation(),
					Lang:     tag.String(),
					Key:      mpl.Replace(metricLabelK.Path),
					Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), mpl.Replace(metricLabelK.Path)),
				})
			}
		})

		report.WalkDimensionSteps(func(dimensionID string, stepID string, _ map[string]interface{}) {
			mpl := strings.NewReplacer(
				"{{dimensionID}}", dimensionID,
				"{{stepID}}", stepID,
			)

			for _, tag := range svc.locale.Tags() {
				out = append(out, &locale.ResourceTranslation{
					Resource: res.ResourceTranslation(),
					Lang:     tag.String(),
					Key:      mpl.Replace(dimStepLabelK.Path),
					Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), mpl.Replace(dimStepLabelK.Path)),
				})
			}
		})
	}

	return
}

func (svc resourceTranslationsManager) pageExtendedAutomationBlock(tag language.Tag, res *types.Page, block types.PageBlock, blockID uint64) (locale.ResourceTranslationSet, error) {
	var (
		bb, _ = block.Options["buttons"].([]interface{})
	)

	return svc.pageBlockButtons(tag, res, blockID, bb)
}

func (svc resourceTranslationsManager) pageExtendedRecordListBlock(tag language.Tag, res *types.Page, block types.PageBlock, blockID uint64) (locale.ResourceTranslationSet, error) {
	var (
		bb, _ = block.Options["selectionButtons"].([]interface{})
	)

	return svc.pageBlockButtons(tag, res, blockID, bb)
}

func (svc resourceTranslationsManager) pageBlockButtons(tag language.Tag, res *types.Page, blockID uint64, bb []interface{}) (locale.ResourceTranslationSet, error) {
	var (
		k   = types.LocaleKeyPagePageBlockBlockIDButtonButtonIDLabel
		out = make(locale.ResourceTranslationSet, 0, 10)
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

func (svc resourceTranslationsManager) loadModule(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (*types.Module, error) {
	return loadModule(ctx, s, namespaceID, moduleID)
}

func (svc resourceTranslationsManager) loadNamespace(ctx context.Context, s store.Storer, namespaceID uint64) (*types.Namespace, error) {
	return loadNamespace(ctx, s, namespaceID)
}

func (svc resourceTranslationsManager) loadPage(ctx context.Context, s store.Storer, namespaceID, pageID uint64) (*types.Page, error) {
	return loadPage(ctx, s, namespaceID, pageID)
}

func (svc resourceTranslationsManager) loadPageLayout(ctx context.Context, s store.Storer, namespaceID, pageID, pageLayoutID uint64) (res *types.PageLayout, err error) {
	return loadPageLayout(ctx, s, namespaceID, pageID, pageLayoutID)
}

func (svc resourceTranslationsManager) loadChart(ctx context.Context, s store.Storer, namespaceID, chartID uint64) (*types.Chart, error) {
	return loadChart(ctx, s, namespaceID, chartID)
}
