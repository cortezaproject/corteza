package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
	"fmt"
	"strings"
{{- range .Imports }}
    {{ . }}
{{- end }}
)


// Parse generates resource setting logic for each resource
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
	{{- range .Def }}
	case {{ unexport .Component "types" }}.{{ export .Resource }}ResourceTranslationType:
		if len(path) != {{ len .Locale.Resource.References }} {
			return "", nil, nil, fmt.Errorf("expecting {{ len .Locale.Resource.References }} reference components in path, got %d", len(path))
		}
		{{- if gt (len .Locale.Resource.References) 0 }}
		ref, pp, err := {{ export .Component .Resource }}ResourceTranslationReferences(
			{{- range $i, $r := .Locale.Resource.References }}
				// {{ unexport $r.Resource }}
				path[{{ $i }}],
			{{ end }}
		)
		return {{ unexport .Component "types" }}.{{ export .Resource }}ResourceTranslationType, ref, pp, err
		{{ else }}

		// Component resource, no path
		return {{ unexport .Component "types" }}.{{ export .Resource }}ResourceTranslationType, nil, nil, nil
		{{- end }}
	{{- end}}
	}

	// return unhandled resource as-is
	return resourceType, nil, nil, nil
}
