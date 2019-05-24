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
	value := make({{ $set }}Set, 3)

	// check walk with no errors
	{
		err := value.Walk(func(*{{ $set }}) error {
			return nil
		})
		test.NoError(t, err, "Expected no returned error from Walk, got %+v", err)
	}

	// check walk with error
	test.Error(t, value.Walk(func(*{{ $set }}) error { return errors.New("Walk error") }), "Expected error from walk, got nil")
}

func Test{{ $set }}SetFilter(t *testing.T) {
	value := make({{ $set }}Set, 3)

	// filter nothing
	{
		set, err := value.Filter(func(*{{ $set }}) (bool, error) {
			return true, nil
		})
		test.NoError(t, err, "Didn't expect error when filtering set: %+v", err)
		test.Assert(t, len(set) == len(value), "Expected equal length filter: %d != %d", len(value), len(set))
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
		test.NoError(t, err, "Didn't expect error when filtering set: %+v", err)
		test.Assert(t, len(set) == 1, "Expected single item, got %d", len(value))
	}

	// filter error
	{
		_, err := value.Filter(func(*{{ $set }}) (bool, error) {
			return false, errors.New("Filter error")
		})
		test.Error(t, err, "Expected error, got %#v", err)
	}
}

{{ if $.WithPrimaryKey }}
func Test{{ $set }}SetIDs(t *testing.T) {
	value := make({{ $set }}Set, 3)
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
		test.Assert(t, val.ID == 2, "Expected ID 2, got %d", val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		test.Assert(t, val == nil, "Expected no value, got %#v", val)
	}

	// List IDs from set
	{
		val := value.IDs()
		test.Assert(t, len(val) == len(value), "Expected ID count mismatch, %d != %d", len(val), len(value))
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
		"",
		"github.com/cortezaproject/corteza-server/internal/test",
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
