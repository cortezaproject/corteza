package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/def"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/gen"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v3/internal/tpl"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

var _ = spew.Dump

func main() {
	dd, err := loadDefinitions(os.Args[1])
	cli.HandleError(err)

	tpls, err := tpl.LoadTemplates(tpl.BaseTemplate(), "./pkg/codegen-v3/assets/templates/gocode")
	if err != nil {
		cli.HandleError(fmt.Errorf("could not load templates: %w", err))
	}

	cli.HandleError(gen.List{
		"RBAC":  gen.RBAC,
		"Envoy": gen.Envoy,
	}.Generate(tpls, dd))
}

func loadDefinition(r io.Reader) (*def.Document, error) {
	doc := &def.Document{
		Envoy: true,
	}

	return doc, yaml.NewDecoder(r).Decode(doc)
}

func loadDefinitions(path string) (dd []*def.Document, err error) {
	var (
		fh    *os.File
		doc   *def.Document
		files []string
	)

	files, err = filepath.Glob(path + "/*.yaml")
	if err != nil {
		return nil, fmt.Errorf("could not load ddefinitions form path '%s': %w", path, err)
	}

	for _, file := range files {
		fh, err = os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("could not load definiton file '%s': %w", file, err)
		}

		doc, err = loadDefinition(fh)
		if err != nil {
			return nil, fmt.Errorf("could not load definiton from '%s': %w", file, err)
		}

		if doc.Skip {
			continue
		}

		if err = doc.Proc(filepath.Base(file)); err != nil {
			return nil, fmt.Errorf("failed to preprocess definitions from '%s': %w", file, err)
		}

		dd = append(dd, doc)
	}

	return
}
