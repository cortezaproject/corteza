package codegen

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"syscall"
	"text/template"
)

type (
	typesDef struct {
		// source file path
		Source string

		// outputDir
		// dir where the source file is
		outputDir string

		Imports []string
		Package string             `yaml:"package"`
		Types   map[string]typeDef `yaml:"types"`
	}

	typeDef struct {
		NoIdField         bool   `yaml:"noIdField"`
		LabelResourceType string `yaml:"labelResourceType"`
	}
)

func procTypes(mm ...string) (dd []*typesDef, err error) {
	dd = make([]*typesDef, 0)

	for _, m := range mm {
		var (
			d = &typesDef{
				Source:    m,
				outputDir: path.Dir(m),

				Package: "types",
				Types:   map[string]typeDef{},
			}
		)

		f, err := os.Open(m)
		if err != nil {
			return nil, fmt.Errorf("%s read failed: %w", m, err)
		}

		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(d); err != nil {
			return nil, fmt.Errorf("%s decode failed: %w", m, err)
		}

		dd = append(dd, d)
	}

	return dd, nil
}

// Generates all type set files & accompanying tests
//
// generates 2 files per type definition
func genTypes(tpl *template.Template, dd ...*typesDef) (err error) {
	var (
		typeGen     = tpl.Lookup("type_set.gen.go.tpl")
		typeGenTest = tpl.Lookup("type_set.gen_test.go.tpl")

		typeLabelsGen = tpl.Lookup("type_labels.gen.go.tpl")
	)

	for _, d := range dd {
		err = goTemplate(path.Join(d.outputDir, "type_set.gen.go"), typeGen, d)

		if err != nil {
			return
		}

		err = goTemplate(path.Join(d.outputDir, "type_set.gen_test.go"), typeGenTest, d)

		if err != nil {
			return
		}

		labelsOutput := path.Join(d.outputDir, "type_labels.gen.go")
		if d.HasLabels() {
			err = goTemplate(labelsOutput, typeLabelsGen, d)
		} else if err = syscall.Unlink(labelsOutput); os.IsNotExist(err) {
			err = nil
		}

		if err != nil {
			return
		}
	}

	return nil
}

func (d typesDef) HasLabels() bool {
	for _, t := range d.Types {
		if len(t.LabelResourceType) > 0 {
			return true
		}
	}
	return false
}
