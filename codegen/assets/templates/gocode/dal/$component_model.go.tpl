package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"github.com/cortezaproject/corteza-server/pkg/dal"
{{- range .imports }}
    {{ . }}
{{- end }}
)

{{ range .models }}
var {{ .var }} = &dal.Model{
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
				{{- if .dal.hasDefault }}HasDefault:  true,  {{ end }}
				{{- if .dal.quotedDefault }}
					DefaultValue: {{ printf "%q" .dal.quotedDefault }},
				{{- else if .dal.defaultEmptyObject }}
					DefaultValue: "{}",
				{{- else if .dal.defaultCurrentTimestamp }}
					DefaultCurrentTimestamp: true,
				{{- else if .dal.hasDefault }}
					DefaultValue: {{.dal.default }},
				{{ end -}}
				{{- if .dal.timezone }}  Timezone:  {{ .dal.timezone }},  {{ end }}
				{{- if .dal.timestampPrecision }} TimezonePrecision:  true,  {{ end }}
				{{- if .dal.precision }} Precision: {{ .dal.precision }}, {{ end }}
				{{- if .dal.scale }}     Scale:     {{ .dal.scale }},     {{ end }}
				{{- if .dal.length }}    Length:    {{ .dal.length }},    {{ end }}
				{{- if .dal.meta }}  Meta: {{ printf "%#v" .dal.meta }},  {{ end }}
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

func init () {
	models = append(
		models,
{{- range .models }}
		{{ .var }},
{{- end }}
	)
}

