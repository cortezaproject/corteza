package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"strings"
{{- range .imports }}
    {{ . }}
{{- end }}
)


// Parse generates resource setting logic for each resource
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ParseRule(res string) (string, *Ref, []*Ref, error) {
	if res == "" {
		return "", nil, nil, fmt.Errorf("empty resource")
	}

	sp := "/"

	res = strings.TrimSpace(res)
	res = strings.TrimRight(res, sp)
	rr := strings.Split(res, sp)

	// only service defined (corteza::system, corteza::compose, ...)
	if len(rr) == 1 {
		return res, nil, nil, nil
	}

	// full thing
	resourceType, path := rr[0], rr[1:]

	for p := 1; p < len(path); p++ {
		if path[p] != "*" && path[p-1] == "*" {
			return "", nil, nil, fmt.Errorf("invalid path wildcard combination for '%s'", res)
		}
	}

	// @todo this is to support res. tr. resources also.
	//       Split it into a separate function and remove this.
	if !strings.HasPrefix(resourceType, "corteza::") {
		resourceType = "corteza::" + resourceType
	}

	// make the resource provide the slice of parent resources we should nest under
	switch resourceType {
	{{- range .resources }}
	case {{ .typeConst }}:
		if len(path) != {{ len .references }} {
			return "", nil, nil, fmt.Errorf("expecting {{ len .references }} reference components in path, got %d", len(path))
		}
		{{- if gt (len .references) 0 }}
		ref, pp, err := {{ .rbacRefFunc }}(
			{{- range $i, $r := .references }}
				path[{{ $i }}],
			{{- end }}
		)
		return resourceType, ref, pp, err
		{{ else }}

		// Component resource, no path
		return resourceType, nil, nil, nil
		{{- end }}
	{{- end}}
	}

	// return unhandled resource as-is
	return resourceType, nil, nil, nil
}
