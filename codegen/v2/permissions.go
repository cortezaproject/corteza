package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"encoding/json"
	"go/format"
	"io/ioutil"

	"github.com/crusttech/crust/internal/rules"
)

func main() {
	var (
		pkg        = flag.String("package", "main", "Package name")
		input      = flag.String("input", "", "Input .json filename")
		output     = flag.String("output", "", "Output .go filename")
		objectName = flag.String("object-name", "c Permissions() []rules.OperationGroup", "Default function declaration")
	)
	flag.Parse()

	var (
		fnPermissions       = fmt.Sprintf("func (*%s) Permissions() []rules.OperationGroup", *objectName)
		fnPermissionDefault = fmt.Sprintf("func (*%s) PermissionDefault(key string) rules.Access", *objectName)
	)

	export := func(s string, values string) []byte {
		formatCode := func(s string) string {
			s = strings.Replace(s, "true,", "true,\n", -1)
			s = strings.Replace(s, "false,", "false,\n", -1)
			s = strings.Replace(s, "{", "{\n", -1)
			s = strings.Replace(s, "}", ",\n}", -1)
			s = strings.Replace(s, "\", ", "\",\n", -1)

			s = strings.Replace(s, "Default:2,", "Default: rules.Allow,", -1)
			s = strings.Replace(s, "Default:1,", "Default: rules.Deny,", -1)
			s = strings.Replace(s, "Default:0,", "Default: rules.Inherit,", -1)
			return s
		}

		formatDefaults := func(s string) string {
			s = formatCode(s)
			s = strings.Replace(s, ", ", ",\n", -1)
			s = strings.Replace(s, ":2,", ": rules.Allow,", -1)
			s = strings.Replace(s, ":1,", ": rules.Deny,", -1)
			s = strings.Replace(s, ":0,", ": rules.Inherit,", -1)
			return s
		}

		s = formatCode(s)
		values = formatDefaults(values)

		var w bytes.Buffer

		fmt.Fprintln(&w, "package", *pkg)
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, "import \"github.com/crusttech/crust/internal/rules\"")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, "/* File is generated from", *input, "with permissions.go */")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, fnPermissions, "{")
		fmt.Fprintln(&w, "\treturn", s)
		fmt.Fprintln(&w, "}")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, fnPermissionDefault, "{")
		fmt.Fprintln(&w, "\tvalues := ", values)
		fmt.Fprintln(&w, "\tif value, ok := values[key]; ok {")
		fmt.Fprintln(&w, "\t\treturn value")
		fmt.Fprintln(&w, "\t}")
		fmt.Fprintln(&w, "\treturn rules.Inherit")
		fmt.Fprintln(&w, "}")

		fmtsrc, err := format.Source(w.Bytes())
		if err != nil {
			log.Printf("fmt warn: %v\n", err)
			fmtsrc = w.Bytes()
		}

		return fmtsrc
	}

	var result []rules.OperationGroup

	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		log.Fatal(err)
	}

	// fill default values from groups
	values := make(map[string]rules.Access)
	for _, group := range result {
		for _, row := range group.Operations {
			values[row.Key] = row.Default
		}
	}

	source := export(fmt.Sprintf("%#v", result), fmt.Sprintf("%#v", values))
	if err := ioutil.WriteFile(*output, source, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println(*output)
}
