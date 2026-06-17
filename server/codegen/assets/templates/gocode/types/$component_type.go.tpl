package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

{{ if or .needsTime .needsLabel .imports .idLists .jsonTypes }}
import (
	{{ if .jsonTypes }}"database/sql/driver"
	{{ end }}{{ if or .idLists .jsonTypes }}"encoding/json"
	{{ end }}{{ if .idLists }}"strconv"
	{{ end }}{{ if .needsTime }}"time"
	{{ end }}{{ if .jsonTypes }}"github.com/crusttech/human/server/pkg/sql"
	{{ end }}{{ if .needsLabel }}labelTypes "github.com/crusttech/human/server/pkg/label/types"
	{{ end }}{{ range .imports }}"{{ . }}"
	{{ end }}
)
{{ end }}

type {{ .expIdent }} struct {
{{ range .fields }}	{{ .expIdent }} {{ .goType }} `{{ .jsonTag }}`
{{ end }}{{ if .needsLabel }}	Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`
{{ end }}{{ if .flags }}	Flags []string `json:"flags,omitempty"`
{{ end }}}
{{ range .idLists }}
// {{ . }} is a []uint64 that serializes each element as a JSON string.
type {{ . }} []uint64

func (ll {{ . }}) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(ll))
	for i, id := range ll {
		ss[i] = strconv.FormatUint(id, 10)
	}
	return json.Marshal(ss)
}

func (ll *{{ . }}) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ll = make({{ . }}, 0, len(raw))
	for _, r := range raw {
		s := string(r)
		if len(s) >= 2 && s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*ll = append(*ll, id)
	}
	return nil
}
{{ end }}
{{ range .jsonTypes }}
func (m *{{ .name }}) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m {{ .name }}) Value() (driver.Value, error) { return json.Marshal(m) }

{{ if .ptr }}func Parse{{ .name }}(ss []string) (p *{{ .name }}, err error) {
	p = &{{ .name }}{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}{{ else }}func Parse{{ .name }}(ss []string) (p {{ .name }}, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}{{ end }}
{{ end }}
