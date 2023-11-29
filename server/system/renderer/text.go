package renderer

import (
	"io"
	"text/template"

	"github.com/Masterminds/sprig"
)

func preprocPlainTemplate(tpl io.Reader, pp map[string]io.Reader) (*template.Template, error) {
	bb, err := io.ReadAll(tpl)
	if err != nil {
		return nil, err
	}

	gtpl := template.New("text/plain_render").
		Funcs(sprig.TxtFuncMap()).
		Funcs(template.FuncMap{
			"env": envGetter(),
		})

	// Prep the original template
	t, err := gtpl.Parse(string(bb))
	if err != nil {
		return nil, err
	}

	// Prep partials
	for _, p := range pp {
		bb, err = io.ReadAll(p)
		if err != nil {
			return nil, err
		}

		t, err = gtpl.Parse(string(bb))
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
