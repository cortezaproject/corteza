package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

	export := func(s string) string {
		s = strings.Replace(s, "{", "{\n", -1)
		s = strings.Replace(s, "}", ",\n}", -1)
		s = strings.Replace(s, "\", ", "\",\n", -1)

		var w bytes.Buffer

		fmt.Fprintln(&w, "package", *pkg)
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, "import \"github.com/crusttech/crust/internal/rbac\"")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, "/* File is generated from", *input, "& main.go */")
		fmt.Fprintln(&w)
		fmt.Fprintln(&w, *fname, "{")
		fmt.Fprintln(&w, "\treturn", s)
		fmt.Fprintln(&w, "}")

		return w.String()
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
	if err := ioutil.WriteFile(*output, []byte(source), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println(*output)
}
