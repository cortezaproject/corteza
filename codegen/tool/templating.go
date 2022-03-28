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

func baseTemplate() *template.Template {
	return template.New("").
		Funcs(sprig.TxtFuncMap())
}

func loadTemplates(rTpl *template.Template, rootDir string) (*template.Template, error) {
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

func writeFormattedGo(dst string, tpl *template.Template, payload interface{}) error {
	return write(dst, tpl, payload, func(in io.ReadWriter) (out io.ReadWriter, err error) {
		var (
			org, bb []byte
		)

		if org, err = ioutil.ReadAll(in); err != nil {
			return
		}

		cp := bytes.NewBuffer(org)

		if bb, err = format.Source(cp.Bytes()); err != nil {
			// output error and return un-formatted source
			_, _ = fmt.Fprintf(os.Stderr, "%s fmt warn: %v\n", dst, err)
			return cp, nil
		}

		return bytes.NewBuffer(bb), nil
	})
}

func write(dst string, tpl *template.Template, payload interface{}, pp ...func(io.ReadWriter) (io.ReadWriter, error)) (err error) {
	var (
		output io.WriteCloser
		buf    io.ReadWriter
	)

	if tpl == nil {
		return fmt.Errorf("could not find template for %s", dst)
	}

	buf = &bytes.Buffer{}
	if err := tpl.Execute(buf, payload); err != nil {
		return err
	}

	for _, proc := range pp {
		if buf, err = proc(buf); err != nil {
			return
		}
	}

	if dst == "" || dst == "-" {
		output = os.Stdout
	} else {
		if output, err = os.Create(dst); err != nil {
			return err
		}

		defer output.Close()
	}

	if _, err = io.Copy(output, buf); err != nil {
		return err
	}

	return nil
}
