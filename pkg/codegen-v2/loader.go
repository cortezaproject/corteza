package main

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v2/internal/def"
	"github.com/cortezaproject/corteza-server/pkg/codegen-v2/internal/tpl"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

var _ = spew.Dump

func main() {
	def, err := loadDefinitions(os.Args[1])
	cli.HandleError(err)

	tpls, err := tpl.BaseTemplate().ParseGlob("./pkg/codegen-v2/assets/templates/gocode/*.tpl")
	if err != nil {
		cli.HandleError(fmt.Errorf("could not load templates: %w", err))
	}

	if err = rbacTypes(tpls, def); err != nil {
		cli.HandleError(fmt.Errorf("could not generate RBAC type code: %w", err))
	}

	if err = accessControlService(tpls, def); err != nil {
		cli.HandleError(fmt.Errorf("could not generate access control service code: %w", err))
	}
}

func loadDefinition(r io.Reader) (*def.Document, error) {
	doc := &def.Document{}
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
