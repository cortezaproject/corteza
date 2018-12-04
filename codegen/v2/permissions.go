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

	"github.com/crusttech/crust/internal/rbac"
)

func main() {
	var (
		pkg    = flag.String("package", "main", "Package name")
		input  = flag.String("input", "", "Input .json filename")
		output = flag.String("output", "", "Output .go filename")
		fname  = flag.String("function", "func Permissions() []rbac.OperationGroup", "Default function declaration")
	)
	flag.Parse()

	export := func(s string) []byte {
		s = strings.Replace(s, "true,", "true,\n", -1)
		s = strings.Replace(s, "false,", "false,\n", -1)
		s = strings.Replace(s, "{", "{\n", -1)
		s = strings.Replace(s, "}", ",\n}", -1)
		s = strings.Replace(s, "\", ", "\",\n", -1)

		var w bytes.Buffer

		fmt.Fprintln(&w, "package", *pkg)
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, "import \"github.com/crusttech/crust/internal/rbac\"")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, "/* File is generated from", *input, "with permissions.go */")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, *fname, "{")
		fmt.Fprintln(&w, "\treturn", s)
		fmt.Fprintln(&w, "}")

		fmtsrc, err := format.Source(w.Bytes())
		if err != nil {
			log.Printf("fmt warn: %v\n", err)
			fmtsrc = w.Bytes()
		}

		return fmtsrc
	}

	var result []rbac.OperationGroup
	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		log.Fatal(err)
	}
	source := export(fmt.Sprintf("%#v", result))
	if err := ioutil.WriteFile(*output, source, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println(*output)
}
