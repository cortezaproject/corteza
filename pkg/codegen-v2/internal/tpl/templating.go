package tpl

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"go/format"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type (
	Wrap struct {
		Package string

		// will be set when grouping definitions by component
		Component string

		Imports []string
		Def     interface{}
	}
)

var nonIdentChars = regexp.MustCompile(`[\s\\/\-.]+`)

func Export(pp ...string) (out string) {
	for _, p := range pp {
		if len(p) > 1 {
			p = strings.ToUpper(p[:1]) + p[1:]
		}

		if ss := nonIdentChars.Split(p, -1); len(ss) > 1 {
			p = Export(ss...)
		}

		out = out + p
	}

	return out
}

func Unexport(pp ...string) (out string) {
	out = Export(pp...)
	if len(out) > 0 {
		return
	}

	return strings.ToLower(out[:1]) + out[1:]
}

func NormalizeImport(i string) string {
	if strings.Contains(i, " ") {
		p := strings.SplitN(i, " ", 2)
		return fmt.Sprintf(`%s "%s"`, p[0], strings.Trim(p[1], `"`))
	} else {
		return fmt.Sprintf(`"%s"`, strings.Trim(i, `"`))
	}
}

func BaseTemplate() *template.Template {
	return template.New("").
		Funcs(sprig.TxtFuncMap()).
		Funcs(map[string]interface{}{
			"export":          Export,
			"unexport":        Unexport,
			"normalizeImport": NormalizeImport,
		})
}

func GoTemplate(dst string, tpl *template.Template, payload Wrap) (err error) {
	var output io.WriteCloser
	buf := bytes.Buffer{}

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
