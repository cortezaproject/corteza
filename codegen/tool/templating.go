package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

func BaseTemplate() *template.Template {
	return template.New("").
		Funcs(sprig.TxtFuncMap())
}

func LoadTemplates(rTpl *template.Template, rootDir string) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1

	return rTpl, filepath.Walk(cleanRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".tpl") || err != nil {
			return err
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := path[pfx:]
		rTpl, err = rTpl.New(name).Parse(string(b))

		return err
	})
}

func GoTemplate(dst string, tpl *template.Template, payload interface{}) (err error) {
	var output io.WriteCloser
	buf := bytes.Buffer{}

	if tpl == nil {
		return fmt.Errorf("could not find template for %s", dst)
	}

	if err := tpl.Execute(&buf, payload); err != nil {
		return err
	}

	fmtsrc, err := format.Source(buf.Bytes())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s fmt warn: %v\n", dst, err)

		err = nil
		fmtsrc = buf.Bytes()
	}

	if dst == "" || dst == "-" {
		output = os.Stdout
	} else {
		if output, err = os.Create(dst); err != nil {
			return err
		}

		defer output.Close()
	}

	if _, err = output.Write(fmtsrc); err != nil {
		return err
	}

	return nil
}
