package codegen

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"text/template"
)

type (
	exprTypesDef struct {
		// source file path
		Source string

		// outputDir
		// dir where the source file is
		outputDir string

		Imports []string
		Package string            `yaml:"package"`
		Types   map[string]string `yaml:"types"`
	}
)

func procExprTypes(mm ...string) (dd []*exprTypesDef, err error) {
	dd = make([]*exprTypesDef, 0)

	for _, m := range mm {
		var (
			d = &exprTypesDef{
				Source:    m,
				outputDir: path.Dir(m),

				Package: "types",
				Types:   make(map[string]string),
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
func genExprTypes(tpl *template.Template, dd ...*exprTypesDef) (err error) {
	var (
		typeGen = tpl.Lookup("expr_types.gen.go.tpl")
	)

	for _, d := range dd {
		err = goTemplate(path.Join(d.outputDir, "expr_types.gen.go"), typeGen, d)

		if err != nil {
			return
		}
	}

	return nil
}
