package internal

import (
	"bytes"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"go/format"
	"io"
	"os"
	"text/template"
)

func WriteFormattedTo(tpl *template.Template, payload interface{}, tplName, dst string) {
	var output io.WriteCloser
	buf := bytes.Buffer{}

	if err := tpl.ExecuteTemplate(&buf, tplName, payload); err != nil {
		cli.HandleError(err)
	} else {
		fmtsrc, err := format.Source(buf.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "fmt warn: %v", err)
			fmtsrc = buf.Bytes()
		}

		if dst == "" || dst == "-" {
			output = os.Stdout
		} else {
			// cli.HandleError(os.Remove(dst))
			if output, err = os.Create(dst); err != nil {
				cli.HandleError(err)
			}

			defer output.Close()
		}

		if _, err := output.Write(fmtsrc); err != nil {
			cli.HandleError(err)
		}
	}
}

func WritePlainTo(tpl *template.Template, payload interface{}, tplName, dst string) {
	var output io.WriteCloser
	buf := bytes.Buffer{}

	if err := tpl.ExecuteTemplate(&buf, tplName, payload); err != nil {
		cli.HandleError(err)
	} else {
		if dst == "" || dst == "-" {
			output = os.Stdout
		} else {
			// cli.HandleError(os.Remove(dst))
			if output, err = os.Create(dst); err != nil {
				cli.HandleError(err)
			}

			defer output.Close()
		}

		if _, err := output.Write(buf.Bytes()); err != nil {
			cli.HandleError(err)
		}
	}
}
