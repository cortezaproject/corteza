package codegen

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/automation/types"
	. "github.com/cortezaproject/corteza-server/pkg/y7s"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
	"text/template"
)

type (
	// definitions are in one file
	aFuncDefs struct {
		Package   string
		Name      string
		Source    string
		Prefix    string
		outputDir string

		// List of imports
		// Used only by generated file and not pre-generated-user-file
		Imports []string

		Functions aFunctionSet
	}

	aFunctionSet []aFuncDef

	aFuncDef struct {
		Name    string
		Meta    *types.FunctionMeta
		Params  aFuncParamSet
		Results aFuncResultSet
	}

	aFuncParamSet  []*aFuncParamDef
	aFuncResultSet []*aFuncResultDef

	aFuncParamDef struct {
		Name     string
		Required bool
		SetOf    bool
		Types    []*aFuncParamTypeVarDef
		Meta     *types.ParamMeta
	}

	aFuncParamTypeVarDef struct {
		WorkflowType string `yaml:"wf"`
		GoType       string `yaml:"go"`
		Suffix       string
	}

	aFuncResultDef struct {
		Name         string
		Required     bool
		SetOf        bool
		WorkflowType string `yaml:"wf"`
		GoType       string `yaml:"go"`
		Meta         *types.ParamMeta
	}
)

func procAutomationFunctions(mm ...string) (dd []*aFuncDefs, err error) {
	for _, m := range mm {
		f, err := os.Open(m)
		if err != nil {
			return nil, fmt.Errorf("%s read failed: %w", m, err)
		}

		defer f.Close()

		var (
			d = &aFuncDefs{
				Package:   "automation",
				Source:    m,
				Name:      path.Base(m),
				outputDir: path.Dir(m),
			}
		)

		d.Name = d.Name[:len(d.Name)-13]

		if err := yaml.NewDecoder(f).Decode(d); err != nil {
			return nil, err
		}

		dd = append(dd, d)
	}

	return
}

func (set *aFunctionSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k *yaml.Node, v *yaml.Node) (err error) {
		def := aFuncDef{Name: k.Value}

		if err = v.Decode(&def); err != nil {
			return err
		}

		*set = append(*set, def)

		return nil
	})
}

func (set *aFuncParamSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k *yaml.Node, v *yaml.Node) (err error) {
		def := aFuncParamDef{}
		if k != nil {
			def.Name = k.Value
		}

		*set = append(*set, &def)
		return v.Decode(&def)
	})
}

func (set *aFuncResultSet) UnmarshalYAML(n *yaml.Node) error {
	return Each(n, func(k *yaml.Node, v *yaml.Node) (err error) {
		def := aFuncResultDef{}
		if k != nil {
			def.Name = k.Value
		}

		*set = append(*set, &def)
		return v.Decode(&def)
	})
}

func genAutomationFunctions(tpl *template.Template, dd ...*aFuncDefs) (err error) {
	var (
		// Will only be generated if file does not exist previously
		tplAFuncGen = tpl.Lookup("afunc.gen.go.tpl")

		dst string
	)

	for _, d := range dd {
		// Generic code, actions for every resource goes to a separated file
		dst = path.Join(d.outputDir, path.Base(d.Source)[:strings.LastIndex(path.Base(d.Source), ".")]+".gen.go")
		json.NewEncoder(os.Stdout).SetIndent("", "  ")
		err = goTemplate(dst, tplAFuncGen, d)
		if err != nil {
			return
		}

	}

	return nil
}
