package yaml

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"golang.org/x/text/language"
)

type (
	auxLocaleMarshalRes struct {
		res    string
		parts  []string
		keys   [][2]string
		nested []*auxLocaleMarshalRes
	}
	auxLocaleMarshalLang struct {
		lang string
		res  []*auxLocaleMarshalRes
	}
)

func (n *resourceTranslation) Prepare(ctx context.Context, state *envoy.ResourceState) (err error) {
	// For now we will only allow resource specific resource translations if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.
	//
	// @todo how can we improve this?
	res := state.Res.(*resource.ResourceTranslation)
	refRes := res.RefRes
	if refRes.ResourceType == composeTypes.ModuleFieldResourceType {
		for _, r := range state.ParentResources {
			if r.ResourceType() == composeTypes.ModuleResourceType && r.Identifiers().HasAny(res.RefPath[1].Identifiers) {
				return
			}
		}
	} else {
		for _, r := range state.ParentResources {
			if refRes.ResourceType == r.ResourceType() && r.Identifiers().HasAny(refRes.Identifiers) {
				return
			}
		}

		// We couldn't find it...
		return resource.ResourceTranslationErrNotFound(refRes.Identifiers)
	}

	return nil
}

func (r *resourceTranslation) Encode(ctx context.Context, doc *Document, state *envoy.ResourceState) (err error) {
	// @todo Improve locale placement
	//
	// In cases where a specific rule is created for a specific resource, nest the rule
	// under the related namespace.
	// For now all rules will be nested under a root node for simplicity sake.

	// @todo move out of the encoding logic
	if r == nil {
		return
	}

	refResource, err := r.makeResourceTranslationResource(state)
	if err != nil {
		return err
	}

	r.refResourceTranslation = refResource

	doc.addResourceTranslation(r)

	return nil
}

func (r *resourceTranslation) makeResourceTranslationResource(state *envoy.ResourceState) (string, error) {
	res := state.Res.(*resource.ResourceTranslation)

	p0ID := "*"
	p1ID := "*"
	p2ID := "*"

	switch r.refLocaleRes.ResourceType {
	case composeTypes.NamespaceResourceType:
		if res.RefRes != nil {
			p1 := resource.FindComposeNamespace(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Slug
		}

		return fmt.Sprintf(composeTypes.NamespaceResourceTranslationTpl(), composeTypes.NamespaceResourceTranslationType, p1ID), nil

	case composeTypes.ModuleResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefRes != nil {
			p1 := resource.FindComposeModule(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeModuleErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(composeTypes.ModuleResourceTranslationTpl(), composeTypes.ModuleResourceTranslationType, p0ID, p1ID), nil

	case composeTypes.PageResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefRes != nil {
			p1 := resource.FindComposePage(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposePageErrUnresolved(res.RefRes.Identifiers)
			}

			if p1.Handle == "" {
				p1ID = strconv.FormatUint(p1.ID, 10)
			} else {
				p1ID = p1.Handle
			}
		}

		return fmt.Sprintf(composeTypes.PageResourceTranslationTpl(), composeTypes.PageResourceTranslationType, p0ID, p1ID), nil

	case composeTypes.ModuleFieldResourceType:
		p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
		if p0 == nil {
			return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
		}
		p0ID = p0.Slug

		p1 := resource.FindComposeModule(state.ParentResources, res.RefPath[1].Identifiers)
		if p1 == nil {
			return "", resource.ComposeModuleErrUnresolved(res.RefPath[1].Identifiers)
		}
		p1ID = p1.Handle

		// field
		f := resource.FindComposeModuleField(state.ParentResources, res.RefPath[1].Identifiers, res.RefRes.Identifiers)
		if f == nil {
			return "", resource.ComposeModuleFieldErrUnresolved(res.RefRes.Identifiers)
		}
		p2ID = f.Name

		return fmt.Sprintf(composeTypes.ModuleFieldResourceTranslationTpl(), composeTypes.ModuleFieldResourceTranslationType, p0ID, p1ID, p2ID), nil

	case composeTypes.ChartResourceType:
		if len(res.RefPath) > 0 {
			p0 := resource.FindComposeNamespace(state.ParentResources, res.RefPath[0].Identifiers)
			if p0 == nil {
				return "", resource.ComposeNamespaceErrUnresolved(res.RefPath[0].Identifiers)
			}
			p0ID = p0.Slug
		}

		if res.RefRes != nil {
			p1 := resource.FindComposeChart(state.ParentResources, res.RefRes.Identifiers)
			if p1 == nil {
				return "", resource.ComposeChartErrUnresolved(res.RefRes.Identifiers)
			}
			p1ID = p1.Handle
		}

		return fmt.Sprintf(composeTypes.ChartResourceTranslationTpl(), composeTypes.ChartResourceTranslationType, p0ID, p1ID), nil

	// case automationTypes.WorkflowResourceType:
	// 	if res.RefRes != nil {
	// 		p0 := resource.FindAutomationWorkflow(state.ParentResources, res.RefRes.Identifiers)
	// 		if p0 == nil {
	// 			return "", resource.AutomationWorkflowErrUnresolved(res.RefRes.Identifiers)
	// 		}
	// 		p0ID = p0.Handle
	// 	}

	// 	return fmt.Sprintf(automationTypes.WorkflowResourceTranslationTpl(), automationTypes.WorkflowResourceTranslationType, p0ID), nil

	// case systemTypes.ReportResourceType:
	// 	if res.RefRes != nil {
	// 		p0 := resource.FindReport(state.ParentResources, res.RefRes.Identifiers)
	// 		if p0 == nil {
	// 			return "", resource.ReportErrUnresolved(res.RefRes.Identifiers)
	// 		}
	// 		p0ID = p0.Handle
	// 	}

	// 	return fmt.Sprintf(systemTypes.ReportResourceTranslationTpl(), systemTypes.ReportResourceTranslationType, p0ID), nil

	default:
		return "", fmt.Errorf("unsupported resource type '%s' for locale resource YAML encode", r.refLocaleRes.ResourceType)
	}
}

func (ll resourceTranslationSet) MarshalYAML() (interface{}, error) {

	// Flow outline
	// . group by language & trans. resource
	// . sort and slice based on the resource path
	// . nest resources based on their path parameters
	// . prepare YAML nodes

	var err error
	_ = err

	if ll == nil || len(ll) == 0 {
		return nil, nil
	}

	// . group by language & trans. resource
	//
	// We use an array+map combo to help assure consistent order when outputing.
	resourceList := make([]*auxLocaleMarshalLang, 0, 2)
	lngIndex := make(map[language.Tag]int)
	resIndex := make(map[string]int)
	for _, l := range ll {
		for _, lr := range l.locales {
			lang := lr.Lang.Tag

			i, ok := lngIndex[lang]
			if !ok {
				i = len(resourceList)
				lngIndex[lang] = i
				resourceList = append(resourceList, &auxLocaleMarshalLang{
					lang: lang.String(),
					res:  make([]*auxLocaleMarshalRes, 0, 4),
				})
			}

			key := fmt.Sprintf("%s/%s", lang.String(), l.refResourceTranslation)
			j, ok := resIndex[key]
			if !ok {
				j = len(resourceList[i].res)
				resIndex[key] = j
				resourceList[i].res = append(resourceList[i].res, &auxLocaleMarshalRes{res: l.refResourceTranslation, keys: make([][2]string, 0, 2)})
			}

			resourceList[i].res[j].keys = append(resourceList[i].res[j].keys, [2]string{lr.K, lr.Message})
		}
	}

	for li, l := range resourceList {
		// . sort and slice based on the resource path
		//
		// Firstly sort; have resources with less path items (slashes) at the top.
		sort.Slice(l.res, func(i, j int) bool {
			a := l.res[i]
			b := l.res[j]

			a.parts = strings.Split(a.res, "/")
			b.parts = strings.Split(b.res, "/")

			return len(a.parts) < len(b.parts)
		})
		//
		// Next chunk; go through the slice and group resources based on path items.
		// Since they are ordered, this is a simple traversal.
		bySize := make([][]*auxLocaleMarshalRes, 0)
		for i := len(l.res) - 1; i >= 0; i-- {
			if i == len(l.res)-1 {
				bySize = append(bySize, []*auxLocaleMarshalRes{l.res[i]})
				continue
			}

			if len(l.res[i].parts) != len(l.res[i+1].parts) {
				bySize = append(bySize, []*auxLocaleMarshalRes{l.res[i]})
			} else {
				bySize[len(bySize)-1] = append(bySize[len(bySize)-1], l.res[i])
			}
		}

		// . nest resources based on their path parameters.
		//
		// Going from most specific resources to less specific resources assures
		// that we will have a correct and minimal output.
		for i := 1; i < len(bySize); i++ {
			aa := bySize[i-1]
			bb := bySize[i]

			for _, b := range bb {
				for i, a := range aa {
					// We're not deleting rows, just setting to nil
					if a == nil {
						continue
					}

					if isAuxLocaleParent(a, b) {
						b.nested = append(b.nested, a)
						aa[i] = nil
					}
				}
			}
		}
		//
		// Lastly, remove the nested hierarchy for easier YAML node construction
		tmpRes := make([]*auxLocaleMarshalRes, 0)
		for i := len(bySize) - 1; i >= 0; i-- {
			bs := bySize[i]
			for _, b := range bs {
				tmpRes = append(tmpRes, unpackAuxLocale(b)...)
			}
		}

		l.res = tmpRes
		resourceList[li] = l
	}

	// . prepare YAML nodes
	localeNode, _ := makeMap()
	for _, byLang := range resourceList {
		resNode, _ := makeMap()

		for _, res := range byLang.res {
			kmNode, _ := makeMap()

			for _, km := range res.keys {
				kmNode, err = addMap(kmNode, km[0], km[1])
				if err != nil {
					return nil, err
				}
			}

			resNode, err = addMap(resNode, res.res, kmNode)
			if err != nil {
				return nil, err
			}

		}

		localeNode, err = addMap(localeNode, byLang.lang, resNode)
		if err != nil {
			return nil, err
		}
	}

	return localeNode, nil
}

func isAuxLocaleParent(main, sub *auxLocaleMarshalRes) bool {
	for i, p := range sub.parts {
		if i == 0 {
			continue
		}
		if main.parts[i] != p {
			return false
		}
	}

	return true
}

func unpackAuxLocale(l *auxLocaleMarshalRes) []*auxLocaleMarshalRes {
	if l == nil {
		return nil
	}

	out := make([]*auxLocaleMarshalRes, 0)
	nn := l.nested
	l.nested = nil
	out = append(out, l)
	for _, n := range nn {
		out = append(out, unpackAuxLocale(n)...)
	}
	return out
}
