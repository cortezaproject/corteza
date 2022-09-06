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
				{{- if .primaryKey }}PrimaryKey: true, {{ end }}
				{{- if .sortable }}Sortable: true,     {{ end }}
				{{- if .filterable }}Filterable: true, {{ end }}
				Type: &{{ .dal.fqType }}{
					{{- if .dal.nullable }}Nullable: true,{{ end }}
					{{- if .dal.quotedDefault }}
						DefaultValue: {{ printf "%q" .dal.quotedDefault }},
					{{- else if .dal.default }}
						DefaultValue: {{.dal.default }},
					{{- else if .dal.defaultEmptyObject }}
						DefaultValue: "{}",
					{{- else if .dal.defaultCurrentTimestamp }}
						DefaultCurrentTimestamp: true,
					{{ end -}}
					{{- if .dal.timezone }}  Timezone:  {{ .dal.timezone }},  {{ end }}
					{{- if .dal.precision }} Precision: {{ .dal.precision }}, {{ end }}
					{{- if .dal.scale }}     Scale:     {{ .dal.scale }},     {{ end }}
					{{- if .dal.length }}    Length:    {{ .dal.length }},    {{ end }}
					{{- if eq .dal.type "Ref" }}
						RefAttribute: {{ printf "%q" .dal.attribute }},
						RefModel: &dal.ModelRef{
							ResourceType: {{ printf "%q" .dal.refModelResType }},
						},
					{{- end }}
				},
				Store: &dal.CodecAlias{Ident: {{ printf "%q" .storeIdent }}},
			},
		{{ end -}}
		},

		Indexes: dal.IndexSet{
		{{- range .indexes }}
			&dal.Index{
				Ident: {{ printf "%q" .ident }},
				Type: {{ printf "%q" .type }},
				{{ if .unique }}Unique: true,{{ end }}
				{{ if .predicate }}Predicate: {{ printf "%q" .predicate }},{{ end }}
				Fields:    []*dal.IndexField{
				{{- range .fields }}
					{
						AttributeIdent: {{ printf "%q" .attribute }},
						{{- if .modifiers }}
						Modifiers: []dal.IndexFieldModifier{ {{- range .modifiers }}{{ printf "%q" . }},{{- end }}  },
						{{- end }}
						{{- if eq .sort  "ASC"   }}Sort:  dal.IndexFieldSortAsc,   {{ end }}
						{{- if eq .sort  "DESC"  }}Sort:  dal.IndexFieldSortDesc,  {{ end }}
						{{- if eq .nulls "FIRST" }}Nulls: dal.IndexFieldNullsFirst,{{ end }}
						{{- if eq .nulls "LAST"  }}Nulls: dal.IndexFieldNullsLast, {{ end }}
					},
				{{ end -}}
				},
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
