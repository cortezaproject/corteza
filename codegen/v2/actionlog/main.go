package main

import (
	"flag"
	"fmt"
	"github.com/cortezaproject/corteza-server/codegen/v2/internal"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/handle"

	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	eventsTemplateFile = "codegen/v2/actionlog/*.go.tpl"
)

type (
	// List of event/log properties that can/will be captured
	// and injected into log or message string
	propsDef struct {
		Name         string
		Type         string
		Fields       []string
		DefaultField string
		Builtin      bool
	}

	actionDef struct {
		// Action name
		Action string `yaml:"action"`

		// String to log when action is successful
		Log string `yaml:"log"`

		// String to log when error was yield
		//ErrorLog string `yaml:"errorLog"`

		// Action severity
		Severity string `yaml:"severity"`
	}

	// Event definition
	errorDef struct {
		// Error key
		// message can contain {variables} from meta data
		Error string `yaml:"error"`

		// Error key
		// message can contain {variables} from meta data
		Message string `yaml:"message"`

		// Formatted and readable audit log message
		// message can contain {variables} from meta data
		Log string `yaml:"log"`

		// Reference to "safe" error
		// safe error should hide any information that might cause
		// personal data leakage or expose system internals
		Safe string `yaml:"safe"`

		// Error severity
		Severity string `yaml:"severity"`

		// HTTP Status code for this error
		HttpStatus string `yaml:"httpStatus"`
	}
)

const (
	// list of actinos and errors
	defSuffix = "_actions.yaml"
)

var (
	// Cut off this binary
	defs = os.Args[1:]

	overwrite bool
	preview   bool

	tpl *template.Template

	placeholderMatcher = regexp.MustCompile(`{(.+?)}`)
)

func main() {
	tpl = template.New("").Funcs(map[string]interface{}{
		"camelCase": internal.CamelCase,
	})

	tpl = template.Must(tpl.ParseGlob(eventsTemplateFile))

	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite all files")
	flag.BoolVar(&preview, "preview", false, "Output to stdout instead of outputPath")
	flag.Parse()

	for _, path := range defs {
		defs, err := filepath.Glob(path + "/*" + defSuffix)
		if err != nil {
			cli.HandleError(err)
		}

		for _, def := range defs {
			base := filepath.Base(def)
			procDef(def, filepath.Join(filepath.Dir(def), base[0:len(base)-len(defSuffix)]+"_actions.gen.go"))
		}
	}
}

func procDef(path, output string) {
	println(path, output)
	var (
		decoder *yaml.Decoder
		tplData = struct {
			Command string
			YAML    string

			Package string

			// List of imports
			// Used only by generated file and not pre-generated-user-file
			Import []string `yaml:"import"`

			Service  string `yaml:"service"`
			Resource string `yaml:"resource"`

			// Default severity for actions
			DefaultActionSeverity string `yaml:"defaultActionSeverity"`

			// Default severity for errors
			DefaultErrorSeverity string `yaml:"defaultErrorSeverity"`

			// If at least one of the errors has HTTP status defined,
			// add support for http errors
			SupportHttpErrors bool

			Props   []*propsDef
			Actions []*actionDef
			Errors  []*errorDef
		}{
			Package: "service",
			YAML:    path,

			DefaultActionSeverity: "info",
			DefaultErrorSeverity:  "error",
		}
	)

	if f, err := os.Open(path); err != nil {
		cli.HandleError(err)
	} else {
		decoder = yaml.NewDecoder(f)
	}

	cli.HandleError(decoder.Decode(&tplData))

	// Prepend generic error
	tplData.Errors = append([]*errorDef{{
		Error:    "generic",
		Message:  "failed to complete request due to internal error",
		Log:      "{err}",
		Severity: "error",
	}}, tplData.Errors...)

	for i := range tplData.Import {
		// Handle list of imports, adds quotes around each import
		//
		// If import string contains a space, assume import alias and
		// quotes only the 2nd part
		if strings.Contains(tplData.Import[i], " ") {
			p := strings.SplitN(tplData.Import[i], " ", 2)
			tplData.Import[i] = fmt.Sprintf(`%s "%s"`, p[0], p[1])
		} else {
			tplData.Import[i] = fmt.Sprintf(`"%s"`, tplData.Import[i])
		}
	}

	// index known meta fields and sanitize types (no type => string type)
	knownProps := map[string]bool{
		"err": true,
	}

	for _, m := range tplData.Props {
		if m.Type == "" {
			m.Type = "string"
		}

		// very optimistic check if referenced type is builtin or not
		m.Builtin = !strings.Contains(m.Type, ".")

		knownProps[m.Name] = true

		if len(m.Fields) > 0 {
			m.DefaultField = m.Fields[0]

		}

		for _, f := range m.Fields {
			knownProps[fmt.Sprintf("%s.%s", m.Name, f)] = true

		}
	}

	for _, a := range tplData.Actions {
		if a.Severity == "" {
			a.Severity = tplData.DefaultActionSeverity
		}
	}

	for _, e := range tplData.Errors {
		if e.Severity == "" {
			e.Severity = tplData.DefaultErrorSeverity
		}

		if e.HttpStatus != "" {
			tplData.SupportHttpErrors = true
		}
	}

	checkHandle := func(s string) {
		if !handle.IsValid(s) {
			cli.HandleError(fmt.Errorf(
				"%s: %s handle empty", path))

		}

		if !handle.IsValid(s) {
			cli.HandleError(fmt.Errorf(
				"%s: invalid handle format: %q", path, s))

		}
	}

	checkPlaceholders := func(def string, kind, s string) {
		for _, match := range placeholderMatcher.FindAllStringSubmatch(s, 1) {
			placeholder := match[1]
			if !knownProps[placeholder] {
				cli.HandleError(fmt.Errorf(
					"%s: unknown placeholder %q used in %s for %s", path, placeholder, def, kind))
			}
		}
	}

	for _, a := range tplData.Actions {
		checkHandle(a.Action)
		if a.Log == "" {
			// If no log is defined, use action handle
			a.Log = a.Action
		}

		checkPlaceholders(a.Action, "log", a.Log)
	}

	for _, e := range tplData.Errors {
		checkHandle(e.Error)

		if e.Message == "" {
			// If no error message is defined, use error handle
			e.Message = e.Error
		}

		if e.Log == "" {
			// If no error log is defined, use error message
			e.Log = e.Message
		}

		checkPlaceholders(e.Error, "message", e.Message)
		checkPlaceholders(e.Error, "log", e.Log)
	}

	internal.WriteTo(tpl, tplData, "actions.gen.go.tpl", output)
}

func (a actionDef) SeverityConstName() string {
	return severityConstName(a.Severity)
}

func (e errorDef) SeverityConstName() string {
	return severityConstName(e.Severity)
}

func severityConstName(s string) string {
	switch strings.ToLower(s) {
	case "emergency":
		return "actionlog.Emergency"
	case "alert":
		return "actionlog.Alert"
	case "crit", "critical":
		return "actionlog.Critical"
	case "warn", "warning":
		return "actionlog.Warning"
	case "notice":
		return "actionlog.Notice"
	case "info", "informational":
		return "actionlog.Info"
	case "debug":
		return "actionlog.Debug"
	default:
		return "actionlog.Error"
	}
}
