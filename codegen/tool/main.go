package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/cortezaproject/corteza-server/pkg/cli"
)

type (
	task struct {
		Template string      `json:"template"`
		Output   string      `json:"output"`
		Syntax   string      `json:"syntax"`
		Payload  interface{} `json:"payload"`
	}
)

var (
	verbose     bool
	showHelp    bool
	tplRootPath string
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.BoolVar(&verbose, "v", false, "be verbose")
	flag.StringVar(&tplRootPath, "p", "codegen/assets/templates", "location of the template files")
	flag.Parse()
}

// Takes JSON input with codegen tasks and definitions and generates files
func main() {
	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	var (
		input = json.NewDecoder(os.Stdin)
		tasks = make([]*task, 0)
		tpl   *template.Template
		err   error
	)

	print("Waiting for stdin ...\n")
	if err = input.Decode(&tasks); err != nil {
		cli.HandleError(fmt.Errorf("failed to decode input from standard input: %v", err))
	}

	if tpl, err = LoadTemplates(BaseTemplate(), tplRootPath); err != nil {
		cli.HandleError(fmt.Errorf("failed to load templates: %v", err))
	}

	for _, j := range tasks {
		switch j.Syntax {
		case "go":
			print(fmt.Sprintf("generating %s (from %s) ...", j.Output, j.Template))
			if err = GoTemplate(j.Output, tpl.Lookup(j.Template), j.Payload); err != nil {
				cli.HandleError(fmt.Errorf("failed to write template: %v", err))
			}
			print("done\n")
		}
	}
}

func print(msg string) {
	if verbose {
		_, _ = fmt.Fprint(os.Stderr, msg)
	}
}
