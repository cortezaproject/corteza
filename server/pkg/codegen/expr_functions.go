package codegen

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type (
	exprFunctionsDef struct {
		// source file path
		Source string

		// dir where the source file is
		outputDir string

		Package string               `yaml:"package"`
		Groups  []*exprFunctionGroup `yaml:"groups"`
	}

	exprFunctionGroup struct {
		Name string `yaml:"name"`
		// RegisteredIn names the package that registers the group's functions
		// when it is not the definition's own; those stay out of Names.
		RegisteredIn string             `yaml:"registeredIn"`
		Intro        string             `yaml:"intro"`
		Functions    []*exprFunctionDef `yaml:"functions"`
	}

	exprFunctionDef struct {
		Name        string               `yaml:"name"`
		Params      []*exprFunctionParam `yaml:"params"`
		Returns     string               `yaml:"returns"`
		Description string               `yaml:"description"`
		Caution     string               `yaml:"caution"`
		Example     exprFunctionExample  `yaml:"example"`
	}

	exprFunctionParam struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Variadic bool   `yaml:"variadic"`
		Optional bool   `yaml:"optional"`
	}

	exprFunctionExample struct {
		Expr   string `yaml:"expr"`
		Scope  string `yaml:"scope"`
		Result string `yaml:"result"`
		Note   string `yaml:"note"`
	}
)

func procExprFunctions(mm ...string) (dd []*exprFunctionsDef, err error) {
	dd = make([]*exprFunctionsDef, 0, len(mm))

	for _, m := range mm {
		d := &exprFunctionsDef{
			Source:    m,
			outputDir: path.Dir(m),
		}

		raw, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("%s read failed: %w", m, err)
		}

		if err = yaml.Unmarshal(raw, d); err != nil {
			return nil, fmt.Errorf("%s decode failed: %w", m, err)
		}

		dd = append(dd, d)
	}

	return dd, nil
}

// Names returns the function names the definition's own package registers, sorted.
func (d exprFunctionsDef) Names() []string {
	var nn []string
	for _, g := range d.Groups {
		if g.RegisteredIn != "" && g.RegisteredIn != d.Package {
			continue
		}

		for _, f := range g.Functions {
			nn = append(nn, f.Name)
		}
	}

	sort.Strings(nn)
	return nn
}

// Signature renders the call form, e.g. `join(list: Array, sep: String): String`.
func (f exprFunctionDef) Signature() string {
	pp := make([]string, len(f.Params))
	for i, p := range f.Params {
		switch {
		case p.Variadic:
			pp[i] = "..." + p.Name + ": " + p.Type
		case p.Optional:
			pp[i] = p.Name + "?: " + p.Type
		default:
			pp[i] = p.Name + ": " + p.Type
		}
	}

	return fmt.Sprintf("%s(%s): %s", f.Name, strings.Join(pp, ", "), f.Returns)
}

// genExprFunctions writes func_names.gen.go next to each definition file
func genExprFunctions(tpl *template.Template, dd ...*exprFunctionsDef) (err error) {
	namesGen := tpl.Lookup("expr_functions.gen.go.tpl")

	for _, d := range dd {
		if err = goTemplate(path.Join(d.outputDir, "func_names.gen.go"), namesGen, d); err != nil {
			return
		}
	}

	return nil
}

// genExprFunctionDocs writes the expression functions reference page
func genExprFunctionDocs(tpl *template.Template, docsPath string, dd ...*exprFunctionsDef) (err error) {
	return plainTemplate(
		path.Join(docsPath, "functions.gen.md"),
		tpl.Lookup("expr_functions.gen.md.tpl"),
		map[string]interface{}{"Definitions": dd},
	)
}
