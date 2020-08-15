package codegen

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
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
		NoIdField bool `yaml:"noIdField"`
	}
)

func procTypes() ([]*typesDef, error) {
	var (
		dd = make([]*typesDef, 0)
	)

	mm, err := filepath.Glob(filepath.Join("*", "*", "types.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob failed: %w", err)
	}

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
func genTypes(tpl *template.Template, dd []*typesDef) (err error) {
	var (
		typeGen     = tpl.Lookup("type_set.gen.go.tpl")
		typeGenTest = tpl.Lookup("type_set.gen_test.go.tpl")
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
	}

	return nil
}
