package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
{{- range .imports }}
    {{ . }}
{{- end }}
)

type (
	modelReplacer interface  {
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
	}
)

var (
{{ range .models }}
	{{ .var }} = &dal.Model{
		Ident: {{ printf "%q" .ident }},
		ResourceType: {{ .resType }},

		Attributes: dal.AttributeSet{
		{{- range .attributes }}
			&dal.Attribute{
				Ident: {{ printf "%q" .expIdent }},
				{{ if .primaryKey }}PrimaryKey: true, {{ end -}}
				{{ if .sortable }}Sortable: true,     {{ end -}}
				{{ if .filterable }}Filterable: true, {{ end }}
				Type: &{{ .dal.fqType }}{
					{{ if .dal.nullable }}Nullable: true,{{ end -}}
					{{ if eq .dal.type "Ref" }}
						RefAttribute: {{ printf "%q" .dal.attribute }},
						RefModel: &dal.ModelRef{
							ResourceType: {{ printf "%q" .dal.refModelResType }},
						},
					{{ end -}}
				},
				Store: &dal.CodecAlias{Ident: {{ printf "%q" .storeIdent }}},
			},
		{{ end -}}
		},
	}
{{ end }}
)

func All() dal.ModelSet {
	return dal.ModelSet{
	{{- range .models }}
		{{ .var }},
	{{- end }}
	}
}


func Register(ctx context.Context, mr modelReplacer) (err error) {
{{- range .models }}
	if err = mr.ReplaceModel(ctx, {{ .var }}); err != nil {
		return
	}
{{ end }}

	return
}
