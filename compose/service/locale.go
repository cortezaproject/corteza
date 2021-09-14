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

func (svc resourceTranslation) moduleExtended(ctx context.Context, res *types.Module) (out locale.ResourceTranslationSet, err error) {
	var (
		k types.LocaleKey
	)

	for _, tag := range svc.locale.Tags() {
		for _, f := range res.Fields {
			k = types.LocaleKeyModuleFieldLabel
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      k.Path,
				Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), k.Path),
			})

			// Extra field bits
			converted, err := svc.moduleFieldValidatorErrorHandler(ctx, tag, f, k.Path)
			if err != nil {
				return nil, err
			}
			out = append(out, converted...)
		}
	}

	return out, nil
}

func (svc resourceTranslation) moduleFieldValidatorErrorHandler(ctx context.Context, tag language.Tag, f *types.ModuleField, k string) (locale.ResourceTranslationSet, error) {
	out := make(locale.ResourceTranslationSet, 0, 10)

	for i, v := range f.Expressions.Validators {
		vContentID := locale.ContentID(v.ValidatorID, i)
		rpl := strings.NewReplacer(
			"{{validatorID}}", strconv.FormatUint(vContentID, 10),
		)
		tKey := rpl.Replace(k)

		out = append(out, &locale.ResourceTranslation{
			Resource: f.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      tKey,
			Msg:      svc.locale.TRFor(tag, f.ResourceTranslation(), tKey),
		})
	}

	return out, nil
}

func (svc resourceTranslation) pageExtended(ctx context.Context, res *types.Page) (out locale.ResourceTranslationSet, err error) {
	var (
		k types.LocaleKey
	)

	for _, tag := range svc.locale.Tags() {
		for i, block := range res.Blocks {
			pbContentID := locale.ContentID(block.BlockID, i)
			rpl := strings.NewReplacer(
				"{{blockID}}", strconv.FormatUint(pbContentID, 10),
			)

			// base stuff
			k = types.LocaleKeyPageBlockTitle
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      rpl.Replace(k.Path),
				Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
			})

			k = types.LocaleKeyPageBlockDescription
			out = append(out, &locale.ResourceTranslation{
				Resource: res.ResourceTranslation(),
				Lang:     tag.String(),
				Key:      rpl.Replace(k.Path),
				Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
			})

			switch block.Kind {
			case "Automation":
				aux, err := svc.pageExtendedAutomatinBlock(tag, res, block, pbContentID, k)
				if err != nil {
					return nil, err
				}

				out = append(out, aux...)
			}
		}
	}

	return
}

func (svc resourceTranslation) pageExtendedAutomatinBlock(tag language.Tag, res *types.Page, block types.PageBlock, blockID uint64, k types.LocaleKey) (locale.ResourceTranslationSet, error) {
	out := make(locale.ResourceTranslationSet, 0, 10)

	bb, _ := block.Options["buttons"].([]interface{})
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
			Msg:      svc.locale.TRFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
		})
	}

	return out, nil
}

// Helper loaders

func (svc resourceTranslation) loadModule(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (m *types.Module, err error) {
	return loadModule(ctx, s, moduleID)
}

func (svc resourceTranslation) loadNamespace(ctx context.Context, s store.Storer, namespaceID uint64) (m *types.Namespace, err error) {
	return loadNamespace(ctx, s, namespaceID)
}

func (svc resourceTranslation) loadPage(ctx context.Context, s store.Storer, namespaceID, pageID uint64) (m *types.Page, err error) {
	_, m, err = loadPage(ctx, s, namespaceID, pageID)
	return m, err
}
