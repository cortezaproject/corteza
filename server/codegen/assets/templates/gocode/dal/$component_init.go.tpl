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
	models []*dal.Model
)

func Models() dal.ModelSet {
	return models
}
