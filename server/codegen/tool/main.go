package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/cli"
)

type (
	inOut struct {
		Template string `json:"template"`
		Output   string `json:"output"`
		Syntax   string `json:"syntax"`
	}

	task struct {
		inOut

		Bulk []inOut `json:"bulk"`

		Payload interface{} `json:"payload"`
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

	started := time.Now()
	print("Waiting for stdin ...")
	if err = input.Decode(&tasks); err != nil {
		cli.HandleError(fmt.Errorf("failed to decode input from standard input: %v", err))
	}

	println(time.Now().Sub(started).Round(time.Second)/time.Second, "sec")

	if tpl, err = loadTemplates(baseTemplate(), tplRootPath); err != nil {
		cli.HandleError(fmt.Errorf("failed to load templates: %v", err))
	}

	for _, j := range tasks {
		if len(j.Bulk) == 0 {
			j.Bulk = append(j.Bulk, j.inOut)
		}

		for _, o := range j.Bulk {
			output := path.Join(outputBase, o.Output)
			print(fmt.Sprintf("generating %s (from %s) ...", output, o.Template))

			switch o.Syntax {
			case "go":
				err = writeFormattedGo(output, tpl.Lookup(o.Template), j.Payload)
			default:
				err = write(output, tpl.Lookup(o.Template), j.Payload)
			}

			if err != nil {
				cli.HandleError(fmt.Errorf("failed to write template: %v", err))
			} else {
				print("done\n")
			}
		}
	}
}

func print(msg string) {
	if verbose {
		_, _ = fmt.Fprint(os.Stderr, msg)
	}
}
