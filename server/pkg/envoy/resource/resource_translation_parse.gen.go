package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"strings"
)

// ParseResourceTranslation generates resource setting logic for each resource
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ParseResourceTranslation(res string) (string, *Ref, []*Ref, error) {
	if res == "" {
		return "", nil, nil, fmt.Errorf("empty resource")
	}

	sp := "/"

	if strings.Index(res, "corteza::") == 0 {
		res = res[9:]
	}

	res = strings.TrimSpace(res)
	res = strings.TrimRight(res, sp)
	rr := strings.Split(res, sp)

	// only service defined (corteza::system, corteza::compose, ...)
	if len(rr) == 1 {
		return "", nil, nil, fmt.Errorf("only service defined: %s", res)
	}

	// full thing
	resourceType, path := rr[0], rr[1:]
	for p := 1; p < len(path); p++ {
		if path[p] == "*" {
			return "", nil, nil, fmt.Errorf("path wildcard not allowed for locale resources: '%s'", res)
		}
	}

	// make the resource provide the slice of parent resources we should nest under
	switch resourceType {
	case composeTypes.ChartResourceTranslationType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeChartResourceTranslationReferences(
			path[0],
			path[1],
		)
		return composeTypes.ChartResourceTranslationType, ref, pp, err

	case composeTypes.ModuleResourceTranslationType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeModuleResourceTranslationReferences(
			path[0],
			path[1],
		)
		return composeTypes.ModuleResourceTranslationType, ref, pp, err

	case composeTypes.ModuleFieldResourceTranslationType:
		if len(path) != 3 {
			return "", nil, nil, fmt.Errorf("expecting 3 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeModuleFieldResourceTranslationReferences(
			path[0],
			path[1],
			path[2],
		)
		return composeTypes.ModuleFieldResourceTranslationType, ref, pp, err

	case composeTypes.NamespaceResourceTranslationType:
		if len(path) != 1 {
			return "", nil, nil, fmt.Errorf("expecting 1 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposeNamespaceResourceTranslationReferences(
			path[0],
		)
		return composeTypes.NamespaceResourceTranslationType, ref, pp, err

	case composeTypes.PageResourceTranslationType:
		if len(path) != 2 {
			return "", nil, nil, fmt.Errorf("expecting 2 reference components in path, got %d", len(path))
		}
		ref, pp, err := ComposePageResourceTranslationReferences(
			path[0],
			path[1],
		)
		return composeTypes.PageResourceTranslationType, ref, pp, err

	}

	// return unhandled resource as-is
	return resourceType, nil, nil, nil
}
