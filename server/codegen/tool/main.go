package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
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
	outputBase  string
)

func init() {
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.BoolVar(&verbose, "v", false, "be verbose")
	flag.StringVar(&tplRootPath, "p", "codegen/assets/templates", "location of the template files")
	flag.StringVar(&outputBase, "b", ".", "base dir for output")
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

	if tpl, err = loadTemplates(baseTemplate(), tplRootPath); err != nil {
		cli.HandleError(fmt.Errorf("failed to load templates: %v", err))
	}

	for _, j := range tasks {
		output := path.Join(outputBase, j.Output)
		print(fmt.Sprintf("generating %s (from %s) ...", output, j.Template))

		switch j.Syntax {
		case "go":
			err = writeFormattedGo(output, tpl.Lookup(j.Template), j.Payload)
		default:
			err = write(output, tpl.Lookup(j.Template), j.Payload)
		}

		if err != nil {
			cli.HandleError(fmt.Errorf("failed to write template: %v", err))
		} else {
			print("done\n")
		}

	}
}

func print(msg string) {
	if verbose {
		_, _ = fmt.Fprint(os.Stderr, msg)
	}
}
