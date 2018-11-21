package main

// *Set type generator for functions like Walk(), Filter(), FindByID() & IDs()

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
	"text/template"
)

const tmplTypeSet = `package types

// 	Hello! This file is auto-generated.

type (
{{ range $i, $set := .Sets }}
	// {{ $set.Name }}Set slice of {{ $set.Name }}
	//
	// This type is auto-generated.
	{{ $set.Name }}Set []*{{ $set.Name }}
{{ end }}
)

{{ range $i, $set := .Sets }}
// Walk iterates through every slice item and calls w({{ $set.Name }}) err
//
// This function is auto-generated.
func (set {{ $set.Name }}Set) Walk(w func(*{{ $set.Name }}) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f({{ $set.Name }}) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set {{ $set.Name }}Set) Filter(f func(*{{ $set.Name }}) (bool, error)) (out {{ $set.Name }}Set, err error) {
	var ok bool
	out = {{ $set.Name }}Set{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

{{ if $set.PK }}
// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set {{ $set.Name }}Set) FindByID(ID uint64) *{{ $set.Name }} {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set {{ $set.Name }}Set) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
{{ end }}
{{ end }}
`

func main() {
	t := template.New("")

	tpl := template.Must(template.Must(t.Clone()).Parse(tmplTypeSet))

	type (
		set struct {
			Name string
			PK   bool
		}
	)
	var (
		stdTypesStr, nopkTypesStr, outputFile string
		output                                io.Writer
	)

	flag.StringVar(&stdTypesStr, "types", "", "Comma separated list of types")
	flag.StringVar(&nopkTypesStr, "no-pk-types", "", "Comma separated list of types without ID field")
	flag.StringVar(&outputFile, "output", "", "JWT Expiration in minutes")

	flag.Parse()

	payload := struct {
		Sets []set
	}{}

	for _, name := range strings.Split(stdTypesStr, ",") {
		if name = strings.TrimSpace(name); len(name) > 0 {
			payload.Sets = append(payload.Sets, set{Name: name, PK: true})
		}
	}

	for _, name := range strings.Split(nopkTypesStr, ",") {
		if name = strings.TrimSpace(name); len(name) > 0 {
			payload.Sets = append(payload.Sets, set{Name: name, PK: false})
		}
	}

	buf := bytes.Buffer{}

	if err := tpl.Execute(&buf, payload); err != nil {
		exit(err)
	} else {
		fmtsrc, err := format.Source(buf.Bytes())
		if err != nil {
			stderr("fmt warn: %v", err)
			fmtsrc = buf.Bytes()
		}

		if outputFile == "" || outputFile == "-" {
			output = os.Stdout
		} else {
			if output, err = os.Create(outputFile); err != nil {
				exit(err)
			}
			defer output.(io.Closer).Close()
		}

		if _, err := output.Write(fmtsrc); err != nil {
			exit(err)
		}
	}
}

func stderr(format string, a ...interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf(format+"\n", a...))
}

func exit(err error) {
	stderr("error: %v", err)
	os.Exit(1)
}
