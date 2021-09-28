package service

// func (svc resourceTranslationsManager) loadReport(ctx context.Context, s store.Storer, reportID uint64) (*types.Report, error) {
// 	return store.LookupReportByID(ctx, s, reportID)
// }

// func (svc resourceTranslationsManager) reportExtended(ctx context.Context, res *types.Report) (out locale.ResourceTranslationSet, er error) {
// 	var (
// 		k types.LocaleKey
// 	)

// 	for _, tag := range svc.locale.Tags() {
// 		for i, block := range res.Blocks {
// 			blockID := locale.ContentID(block.BlockID, i)
// 			rpl := strings.NewReplacer(
// 				"{{blockID}}", strconv.FormatUint(blockID, 10),
// 			)

// 			// base stuff
// 			k = types.LocaleKeyReportBlockTitle
// 			out = append(out, &locale.ResourceTranslation{
// 				Resource: res.ResourceTranslation(),
// 				Lang:     tag.String(),
// 				Key:      rpl.Replace(k.Path),
// 				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
// 			})

// 			k = types.LocaleKeyReportBlockDescription
// 			out = append(out, &locale.ResourceTranslation{
// 				Resource: res.ResourceTranslation(),
// 				Lang:     tag.String(),
// 				Key:      rpl.Replace(k.Path),
// 				Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), rpl.Replace(k.Path)),
// 			})
// 		}
// 	}

// 	return
// }
