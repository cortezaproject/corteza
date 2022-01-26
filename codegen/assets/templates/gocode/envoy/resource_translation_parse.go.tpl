package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"strings"
{{- range .imports }}
    {{ . }}
{{- end }}
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
	{{- range .resources }}
	case {{ .typeConst }}:
		if len(path) != {{ len .references }} {
			return "", nil, nil, fmt.Errorf("expecting {{ len .references }} reference components in path, got %d", len(path))
		}
		ref, pp, err := {{ .resTrRefFunc }}(
			{{- range $i, $r := .references }}
				path[{{ $i }}],
			{{- end }}
		)
		return {{ .typeConst }}, ref, pp, err
		{{ else }}

		// Component resource, no path
		return {{ .typeConst }}, nil, nil, nil
		{{- end }}
	}

	// return unhandled resource as-is
	return resourceType, nil, nil, nil
}
