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
	{{ if $import }}
	"{{ $import }}"
	{{ else }}

	{{ end }}
	{{ end }}
)
{{ end }}

// 	Hello! This file is auto-generated.

{{ range $i, $set := .Sets }}

func Test{{ $set }}SetWalk(t *testing.T) {
	var (
		value = make({{ $set }}Set, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*{{ $set }}) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*{{ $set }}) error { return errors.New("walk error") }))

}

func Test{{ $set }}SetFilter(t *testing.T) {
	var (
		value = make({{ $set }}Set, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*{{ $set }}) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*{{ $set }}) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*{{ $set }}) (bool, error) {
			return false, errors.New("filter error")
		})
		req.Error(err)
	}
}

{{ if $.WithPrimaryKey }}
func Test{{ $set }}SetIDs(t *testing.T) {
	var (
		value = make({{ $set }}Set, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new({{ $set }})
	value[1] = new({{ $set }})
	value[2] = new({{ $set }})
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
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

	payload.Imports = []string{
		"testing",
		"errors",
		"github.com/stretchr/testify/require",
	}

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
