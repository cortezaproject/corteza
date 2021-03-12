package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func goTemplate(dst string, tpl *template.Template, payload interface{}) (err error) {
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

	if _, err := output.Write(fmtsrc); err != nil {
		return err
	}

	return nil
}

func plainTemplate(dst string, tpl *template.Template, payload interface{}) (err error) {
	var output io.WriteCloser
	buf := bytes.Buffer{}

	if err := tpl.Execute(&buf, payload); err != nil {
		return err
	}

	if dst == "" || dst == "-" {
		output = os.Stdout
	} else {
		if output, err = os.Create(dst); err != nil {
			return err
		}

		defer output.Close()
	}

	if _, err := output.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func camelCase(pp ...string) (out string) {
	for i, p := range pp {
		if i > 0 && len(p) > 1 {
			p = strings.ToUpper(p[:1]) + p[1:]
		}

		out = out + p
	}

	return out
}

// PubIdent returns published identifier by uppercasing
// input, cammelcasing it and removing ident unfriendly characters
var nonIdentChars = regexp.MustCompile(`[\s\\/]+`)

func export(pp ...string) (out string) {
	for _, p := range pp {
		if len(p) > 1 {
			p = strings.ToUpper(p[:1]) + p[1:]
		}

		if ss := nonIdentChars.Split(p, -1); len(ss) > 1 {
			p = export(ss...)
		}

		out = out + p
	}

	return out
}

func unexport(pp ...string) (out string) {
	out = export(pp...)
	return strings.ToLower(out[:1]) + out[1:]
}

func removePtr(name string) string {
	return strings.TrimLeft(name, "*")
}

func hasPtr(name string) bool {
	return len(name) > 0 && name[0:1] == "*"
}

func toggleExport(e bool, pp ...string) (out string) {
	if e {
		return export(pp...)
	}

	return unexport(pp...)
}

// convets to underscore
func cc2underscore(cc string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	u := matchFirstCap.ReplaceAllString(cc, "${1}_${2}")
	u = matchAllCap.ReplaceAllString(u, "${1}_${2}")
	return strings.ToLower(u)
}

// Handle list of imports, adds quotes around each import
//
// If import string contains a space, assume import alias and
// quotes only the 2nd part
func normalizeImport(i string) string {
	if strings.Contains(i, " ") {
		p := strings.SplitN(i, " ", 2)
		return fmt.Sprintf(`%s "%s"`, p[0], strings.Trim(p[1], `"`))
	} else {
		return fmt.Sprintf(`"%s"`, strings.Trim(i, `"`))
	}
}
