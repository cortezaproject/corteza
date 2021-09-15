package yaml

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/pkg/y7s"
	"github.com/cortezaproject/corteza-server/system/types"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func decodeLocale(n *yaml.Node) (resourceTranslationSet, error) {
	var (
		rr = make(resourceTranslationSet, 0, 20)
	)

	err := y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		if k.Value != "locale" {
			return nil
		}

		return y7s.EachMap(v, func(k, v *yaml.Node) error {
			lang := ""
			if err = y7s.DecodeScalar(k, "language", &lang); err != nil {
				return err
			}

			rr, err = rr.decodeLocale(lang, v)
			if err != nil {
				return err
			}
			return nil
		})
	})

	return rr, err
}

func (rr resourceTranslationSet) decodeLocale(lang string, locRes *yaml.Node) (oo resourceTranslationSet, err error) {
	if rr == nil {
		oo = make(resourceTranslationSet, 0, 10)
	} else {
		oo = rr
	}

	parseOps := func(km *yaml.Node, lang, resource string) error {
		return y7s.EachMap(km, func(k, msg *yaml.Node) error {
			lr := &resourceTranslation{
				locales: types.ResourceTranslationSet{{
					Lang:     types.Lang{Tag: language.Make(lang)},
					Resource: resource,
					K:        k.Value,
					Message:  msg.Value,
				}},
			}

			if resource != "" {
				if err = lr.setResourceTranslation(resource); err != nil {
					return fmt.Errorf("failed to decode locale resource: %w", err)
				}
			}

			oo = append(oo, lr)
			return nil
		})
	}

	// If its a mapping node, keys represent resources
	if locRes.Content[1].Kind == yaml.MappingNode {
		return oo, y7s.EachMap(locRes, func(k, v *yaml.Node) error {
			return parseOps(v, lang, k.Value)
		})
	} else {
		err = parseOps(locRes, lang, "")
		if err != nil {
			return nil, err
		}
	}

	return
}

func (rr resourceTranslationSet) bindResource(resI resource.Interface) resourceTranslationSet {
	res, ref, pp := resI.(resource.LocaleInterface).ResourceTranslationParts()

	rtr := make(resourceTranslationSet, 0, len(rr))
	for _, r := range rr {

		r.refResourceTranslation = res
		r.refLocaleRes = ref
		r.refPathRes = pp
		rtr = append(rtr, r)
	}

	return rtr
}

func (rr resourceTranslationSet) MarshalEnvoy() ([]resource.Interface, error) {
	rr = rr.groupByResourceTranslation()

	var nn = make([]resource.Interface, 0, len(rr))

	for _, r := range rr {
		nn = append(nn, resource.NewResourceTranslation(r.locales, r.refResourceTranslation, r.refLocaleRes, r.refPathRes...))
	}
	return nn, nil
}

func (r *resourceTranslation) setResourceTranslation(res string) error {
	_, ref, pp, err := resource.ParseResourceTranslation(res)
	if err != nil {
		return err
	}

	r.refResourceTranslation = res
	r.refLocaleRes = ref
	r.refPathRes = pp
	return nil
}
