package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

type (
{{- range .types }}
	{{ .expIdent }}Set []*{{ .expIdent }}
{{- end }}
)

{{- range .types }}

func (set {{ .expIdent }}Set) Walk(w func(*{{ .expIdent }}) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

func (set {{ .expIdent }}Set) Filter(f func(*{{ .expIdent }}) (bool, error)) (out {{ .expIdent }}Set, err error) {
	var ok bool
	out = {{ .expIdent }}Set{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}
{{- if not .noIdField }}

func (set {{ .expIdent }}Set) FindByID(ID uint64) *{{ .expIdent }} {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

func (set {{ .expIdent }}Set) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
{{- end }}
{{- end }}
