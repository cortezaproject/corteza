package codegen

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type (
	optionsDef struct {
		Source    string
		outputDir string

		Name string

		// List of imports
		// Used only by generated file and not pre-generated-user-file
		Imports []string `yaml:"imports"`

		Properties optionsPropSet `yaml:"props"`
	}

	optionsPropSet []*optionsProp

	optionsProp struct {
		Name        string
		Type        string
		Env         string
		Default     *optionsPropDefault
		Description string
	}

	optionsPropDefault string
)

// Processes multiple options defenitions
func procOptions(mm ...string) (dd []*optionsDef, err error) {
	var (
		f io.ReadCloser
		d *optionsDef
	)

	dd = make([]*optionsDef, 0)
	for _, m := range mm {
		err = func() error {
			if f, err = os.Open(m); err != nil {
				return err
			}

			defer f.Close()

			d = &optionsDef{}

			if err := yaml.NewDecoder(f).Decode(d); err != nil {
				return err
			}

			for _, j := range d.Properties {

				if j.Type == "" {
					j.Type = "string"
				}

				if j.Env == "" {
					j.Env = strings.ToUpper(d.Name + "_" + cc2underscore(j.Name))
				}

			}

			d.Source = m
			d.outputDir = path.Dir(m)

			dd = append(dd, d)

			return nil
		}()
		if err != nil {
			return nil, fmt.Errorf("could not process %s: %w", m, err)
		}
	}
	return dd, nil
}

// Custom UnmarshalYAML function for
func (pd *optionsPropDefault) UnmarshalYAML(n *yaml.Node) error {

	val := n.Value

	if n.Style == yaml.DoubleQuotedStyle {
		val = "\"" + val + "\""
	}

	*pd = optionsPropDefault(val)

	return nil
}

// Gets package name from file path
func (o optionsDef) Package() string {
	return path.Base(path.Dir(o.Source))
}

func genOptions(tpl *template.Template, dd ...*optionsDef) (err error) {
	var (
		tplOptionsGen = tpl.Lookup("options.gen.go.tpl")

		tplOptionsAdoc = tpl.Lookup("options.gen.adoc.tpl")

		dst string
	)

	for _, d := range dd {
		dst = path.Join(d.outputDir, path.Base(d.Source)[:strings.LastIndex(path.Base(d.Source), ".")]+".gen.go")
		err = goTemplate(dst, tplOptionsGen, d)
		if err != nil {
			return
		}
		dst = path.Join(d.outputDir, path.Base(d.Source)[:strings.LastIndex(path.Base(d.Source), ".")]+".adoc")
		err = goTemplate(dst, tplOptionsAdoc, d)
		if err != nil {
			return
		}
	}

	// for _, d := range dd {
	// 	dst = path.Join(d.outputDir, path.Base(d.Source)[:strings.LastIndex(path.Base(d.Source), ".")]+".adoc")
	// 	err = goTemplate(dst, tplOptionsAdoc, d)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	return nil
}
