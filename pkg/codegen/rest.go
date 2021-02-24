package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type (
	// The following structure represents
	// legacy API definition (spec.json)
	//
	//

	// definitions are in one file
	restDef struct {
		App       string
		Source    string
		outputDir string

		Endpoints []*restEndpointDef
	}

	restEndpointDef struct {
		Title          string                `yaml:"title"`
		Path           string                `yaml:"path"`
		Entrypoint     string                `yaml:"entrypoint"`
		Authentication []interface{}         `yaml:"authentication,omitempty"`
		Apis           []*restEndpointApi    `yaml:"apis"`
		Imports        []string              `yaml:"imports"`
		Description    string                `yaml:"description,omitempty"`
		Params         restEndpointParamsDef `yaml:"parameters,omitempty"`
	}

	restEndpointApi struct {
		Name   string                `yaml:"name"`
		Method string                `yaml:"method"`
		Title  string                `yaml:"title"`
		Path   string                `yaml:"path"`
		Params restEndpointParamsDef `yaml:"parameters,omitempty"`
	}

	restEndpointParamsDef struct {
		Post []*restEndpointParamDef `yaml:"post"`
		Path []*restEndpointParamDef `yaml:"path"`
		Get  []*restEndpointParamDef `yaml:"get"`
	}

	restEndpointParamDef struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Required bool   `yaml:"required"`
		Title    string `yaml:"title"`
		Origin   string

		Sensitive bool `yaml:"sensitive"`

		DefinedParser string `yaml:"parser"`
	}
)

func procRest(mm ...string) (dd []*restDef, err error) {
	dd = make([]*restDef, 0)

	for _, m := range mm {
		err = func() error {
			f, err := os.Open(m)
			if err != nil {
				return fmt.Errorf("%s read failed: %w", m, err)
			}

			defer f.Close()

			var d = &restDef{}

			if err := yaml.NewDecoder(f).Decode(d); err != nil {
				return err
			}

			d.outputDir = path.Dir(m)

			// Append params from endpoit to all apis
			for _, e := range d.Endpoints {
				for _, a := range e.Apis {
					a.Params.Path = append(e.Params.Path, a.Params.Path...)
					a.Params.Post = append(e.Params.Post, a.Params.Post...)
					a.Params.Get = append(e.Params.Get, a.Params.Get...)
				}
			}

			dd = append(dd, d)
			return nil
		}()

		if err != nil {
			return nil, fmt.Errorf("failed to process %s: %w", m, err)
		}
	}

	return dd, nil
}

func genRest(tpl *template.Template, dd ...*restDef) (err error) {
	var (
		// Will only be generated if file does not exist previously
		tplHandler = tpl.Lookup("rest_handler.go.tpl")
		tplRequest = tpl.Lookup("rest_request.go.tpl")

		dst string
	)

	for _, d := range dd {
		for _, e := range d.Endpoints {

			// Generic code, every event goes into one file (per app)
			dst = path.Join(d.outputDir, "rest", "handlers", e.Entrypoint+".go")
			err = goTemplate(dst, tplHandler, map[string]interface{}{
				"Source":   d.Source,
				"Endpoint": e,
				"App":      path.Base(d.outputDir),
			})
			if err != nil {
				return
			}

			// Generic code, every event goes into one file (per app)
			dst = path.Join(d.outputDir, "rest", "request", e.Entrypoint+".go")
			err = goTemplate(dst, tplRequest, map[string]interface{}{
				"Source":   d.Source,
				"Endpoint": e,
				"Imports":  e.Imports,
			})
			if err != nil {
				return
			}
		}

	}

	return nil
}

func (d *restEndpointParamsDef) All() []*restEndpointParamDef {
	var pp = make([]*restEndpointParamDef, 0)

	for _, p := range d.Path {
		p.Origin = "PATH"
		pp = append(pp, p)
	}

	for _, p := range d.Get {
		p.Origin = "GET"
		pp = append(pp, p)
	}

	for _, p := range d.Post {
		p.Origin = "POST"
		pp = append(pp, p)
	}

	return pp
}

func (d *restEndpointParamDef) IsUpload() bool {
	switch d.Type {
	case "*multipart.FileHeader":
		return true
	}

	return false
}

func (d *restEndpointParamDef) IsSlice() bool {
	return strings.HasPrefix(d.Type, "[]") || strings.HasSuffix(d.Type, "Set")
}

func (d *restEndpointParamDef) IsString() bool {
	switch d.Type {
	case "string", "[]string", "[]*string":
		return true
	}

	return false
}

func (d *restEndpointParamDef) FieldTag() string {
	switch d.Type {
	case "uint64":
		return "`json:\",string\"`"
	}

	return ""
}

func (d *restEndpointParamDef) HasExplicitParser() bool {
	return d.DefinedParser != ""
}

func (d *restEndpointParamDef) Parser(arg string) string {
	if d.HasExplicitParser() {
		return fmt.Sprintf("%s(%s)", d.DefinedParser, arg)
	}

	switch d.Type {
	case "[]uint64":
		return fmt.Sprintf("payload.ParseUint64s(%s), nil", arg)
	case "[]uint":
		return fmt.Sprintf("payload.ParseUints(%s), nil", arg)
	case "time.Time":
		return fmt.Sprintf("payload.ParseISODateWithErr(%s)", arg)
	case "*time.Time":
		return fmt.Sprintf("payload.ParseISODatePtrWithErr(%s)", arg)
	case "sqlxTypes.JSONText":
		return fmt.Sprintf("payload.ParseJSONTextWithErr(%s)", arg)
	case "int", "uint", "uint64", "int64", "float", "float64", "bool":
		return fmt.Sprintf("payload.Parse%s(%s), nil", export(d.Type), arg)
	case "string", "[]string":
		return fmt.Sprintf("%s, nil", arg)
	default:
		return fmt.Sprintf("%s(%s), nil", d.Type, arg)
	}

}
