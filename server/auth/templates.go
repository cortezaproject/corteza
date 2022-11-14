package auth

import (
	"embed"
	"html/template"
	"io"
)

//go:embed assets/templates/*.tpl
var Templates embed.FS

type (
	templateLoader func(tpls *template.Template) (tpl *template.Template, err error)

	templates struct {
		reload bool
	}

	templateExecutor interface {
		ExecuteTemplate(io.Writer, string, interface{}, ...func(*template.Template) *template.Template) error
	}

	templateReloader struct {
		base   *template.Template
		loader templateLoader
	}

	templateStatic struct {
		base *template.Template
	}
)

func NewReloadableTemplates(base *template.Template, loader templateLoader) *templateReloader {
	return &templateReloader{
		base:   base,
		loader: loader,
	}
}

// Reloads templates before every execution
//
// This is great for local development
func (t templateReloader) ExecuteTemplate(w io.Writer, name string, data interface{}, oo ...func(*template.Template) *template.Template) error {
	tpl, err := t.loader(applyTemplateExecOpt(t.base, oo...))
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(w, name, data)
}

func NewStaticTemplates(base *template.Template, loader templateLoader) (s *templateStatic, err error) {
	s = &templateStatic{}
	s.base, err = loader(base)
	return
}

// ExecuteTemplate executes preloaded templates
func (t templateStatic) ExecuteTemplate(w io.Writer, name string, data interface{}, oo ...func(*template.Template) *template.Template) error {
	return applyTemplateExecOpt(t.base, oo...).ExecuteTemplate(w, name, data)
}

// EmbeddedTemplates returns embedded templates.
func EmbeddedTemplates(t *template.Template) (tpl *template.Template, err error) {
	return t.ParseFS(Templates, "assets/templates/*.html.tpl")
}

func applyTemplateExecOpt(t *template.Template, oo ...func(*template.Template) *template.Template) *template.Template {
	for _, fn := range oo {
		t = fn(t)
	}

	return t
}
