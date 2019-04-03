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

const tmplTypeSet = `package {{ .Package }}

{{ if .Imports }}
import (
	{{ range $i, $import := .Imports }}
	"{{ $import }}"
	{{ end }}
)
{{ end }}

// 	Hello! This file is auto-generated.

type (
{{ range $i, $set := .Sets }}
	// {{ $set }}Set slice of {{ $set }}
	//
	// This type is auto-generated.
	{{ $set }}Set []*{{ $set }}
{{ end }}
)

{{ range $i, $set := .Sets }}
// Walk iterates through every slice item and calls w({{ $set }}) err
//
// This function is auto-generated.
func (set {{ $set }}Set) Walk(w func(*{{ $set }}) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f({{ $set }}) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set {{ $set }}Set) Filter(f func(*{{ $set }}) (bool, error)) (out {{ $set }}Set, err error) {
	var ok bool
	out = {{ $set }}Set{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

{{ if $.WithPrimaryKey }}
// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set {{ $set }}Set) FindByID(ID uint64) *{{ $set }} {
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
func (set {{ $set }}Set) IDs() (IDs []uint64) {
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

	var (
		importsStr,
		stdTypesStr,
		outputFile string
		output io.Writer
	)

	payload := struct {
		Package string
		Sets    []string
		Imports []string

		ResourceType string

		WithPrimaryKey bool
		WithResources  bool
	}{}

	flag.StringVar(&stdTypesStr, "types", "", "Comma separated list of types")
	flag.StringVar(&importsStr, "imports", "", "Comma separated list of imports")
	flag.StringVar(&payload.Package, "package", "types", "Package name")

	flag.BoolVar(&payload.WithPrimaryKey, "with-primary-key", true, "Generate types with ID field")

	flag.BoolVar(&payload.WithResources, "with-resources", false, "Generate Resources() functions")
	flag.StringVar(&payload.ResourceType, "resource-type", "Resource", "Resource type name to use (with imports)")

	flag.StringVar(&outputFile, "output", "", "JWT Expiration in minutes")

	flag.Parse()
	for _, name := range strings.Split(importsStr, ",") {
		if name = strings.TrimSpace(name); len(name) > 0 {
			payload.Imports = append(payload.Imports, name)
		}
	}

	for _, name := range strings.Split(stdTypesStr, ",") {
		if name = strings.TrimSpace(name); len(name) > 0 {
			payload.Sets = append(payload.Sets, name)
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
			fmt.Println(outputFile)
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
