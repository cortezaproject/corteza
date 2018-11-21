package main

// Makeshift *Set type generator for functions like Walk(), Filter(), FindByID() & IDs()

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

const tmplTypeSet = `package types

type (
{{ range .Sets }}
	// {{ . }}Set slice of {{ . }}
	{{ . }}Set []*{{ . }}
{{- end }}
)

{{ range .Sets }}
// Walk iterates through every slice item and calls w({{ . }}) err
func (set {{ . }}Set) Walk(w func(*{{ . }}) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f({{ . }}) (bool, err) and return filtered slice
func (set {{ . }}Set) Filter(f func(*{{ . }}) (bool, error)) (out {{ . }}Set, err error) {
	var ok bool
	out = {{ . }}Set{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

// Finds slice item by its ID property
func (set {{ . }}Set) FindByID(ID uint64) *{{ . }} {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// Returns a slice of uint64s from all items in the set
func (set {{ . }}Set) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
{{- end }}
`

func main() {
	t := template.New("")

	tpl := template.Must(template.Must(t.Clone()).Parse(tmplTypeSet))

	if len(os.Args) < 3 {
		exit(errors.New("type-set.go <output-file> <type type2 type3 ...>"))
	}

	l := len(os.Args) - 1
	outputFile := os.Args[l]
	payload := struct{ Sets []string }{Sets: os.Args[1:l]}

	buf := bytes.Buffer{}

	if err := tpl.Execute(&buf, payload); err != nil {
		exit(err)
	} else {
		fmtsrc, err := format.Source(buf.Bytes())
		if err != nil {
			stderr("fmt warn: %v", err)
			fmtsrc = buf.Bytes()
		}

		output, err := os.Create(outputFile)
		if err != nil {
			exit(err)
		}

		defer output.Close()

		if _, err := output.Write(fmtsrc); err != nil {
			exit(err)
		}
	}
}

func stderr(format string, a ...interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf(format, a...))
}

func exit(err error) {
	stderr("error: %v", err)
	os.Exit(1)
}
