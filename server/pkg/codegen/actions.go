package codegen

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/cortezaproject/corteza/server/pkg/handle"
	"gopkg.in/yaml.v3"
)

type (
	// definitions are in multiple files and each definition
	// should produce one output
	actionsDef struct {
		Component string
		Source    string
		outputDir string

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
	}

	// List of event/log properties that can/will be captured
	// and injected into log or message string
	propsDef struct {
		Name    string
		Type    string
		Fields  []string
		Builtin bool
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

		// Longer message or error description that can help resolving the error
		Details string `yaml:"details"`

		// Relative link to content in the documentation
		Documentation string `yaml:"documentation"`

		// Reference to "safe" error
		// safe error should hide any information that might cause
		// personal data leakage or expose system internals
		MaskedWith string `yaml:"maskedWith"`

		// Error severity
		Severity string `yaml:"severity"`

		// HTTP Status code for this error
		HttpStatus string `yaml:"httpStatus"`
	}
)

// Processes multiple action definitions
func procActions(mm ...string) (dd []*actionsDef, err error) {
	var (
		f io.ReadCloser
		d *actionsDef
	)

	dd = make([]*actionsDef, 0)
	for _, m := range mm {
		err = func() error {

			if f, err = os.Open(m); err != nil {
				return err
			}

			defer f.Close()

			d = &actionsDef{Component: strings.SplitN(m, string(filepath.Separator), 2)[0]}

			if err := yaml.NewDecoder(f).Decode(d); err != nil {
				return err
			}

			if err = actionNormalize(d); err != nil {
				return err
			}

			d.Source = m
			d.outputDir = path.Dir(m)

			dd = append(dd, d)
			return nil
		}()

		if err != nil {
			return nil, fmt.Errorf("could not process %s: %w", m, err)
		}
	}

	return dd, nil
}

func actionNormalize(d *actionsDef) error {
	// Prepend generic error
	d.Errors = append([]*errorDef{{
		Error:    "generic",
		Message:  "failed to complete request due to internal error",
		Log:      "{err}",
		Severity: "error",
	}}, d.Errors...)

	// index known meta fields and sanitize types (no type => string type)
	knownProps := map[string]bool{
		"err": true,
	}

	for _, m := range d.Props {
		if m.Type == "" {
			m.Type = "string"
		}

		// very optimistic check if referenced type is builtin or not
		m.Builtin = !strings.Contains(m.Type, ".")

		knownProps[m.Name] = true

		for _, f := range m.Fields {
			knownProps[fmt.Sprintf("%s.%s", m.Name, f)] = true

		}
	}

	for _, a := range d.Actions {
		if a.Severity == "" {
			a.Severity = d.DefaultActionSeverity
		}
	}

	for _, e := range d.Errors {
		if e.Severity == "" {
			e.Severity = d.DefaultErrorSeverity
		}

		if e.HttpStatus != "" {
			d.SupportHttpErrors = true
		}
	}

	checkHandle := func(s string) error {
		if !handle.IsValid(s) {
			return fmt.Errorf("handle empty")

		}

		if !handle.IsValid(s) {
			return fmt.Errorf("invalid handle format: %q", s)

		}

		return nil
	}

	placeholderMatcher := regexp.MustCompile(`\{\{(.+?)\}\}`)
	checkPlaceholders := func(def string, kind, s string) error {
		for _, match := range placeholderMatcher.FindAllStringSubmatch(s, 1) {
			placeholder := match[1]
			if !knownProps[placeholder] {
				return fmt.Errorf("unknown placeholder %q used in %s for %s", placeholder, def, kind)
			}
		}

		return nil
	}

	for _, a := range d.Actions {
		checkHandle(a.Action)
		if a.Log == "" {
			// If no log is defined, use action handle
			a.Log = a.Action
		}

		if err := checkPlaceholders(a.Action, "log", a.Log); err != nil {
			return err
		}
	}

	for _, e := range d.Errors {
		if err := checkHandle(e.Error); err != nil {
			return err
		}

		if err := checkPlaceholders(e.Error, "message", e.Message); err != nil {
			return err
		}
		if err := checkPlaceholders(e.Error, "log", e.Log); err != nil {
			return err
		}
	}

	return nil
}

func (a actionsDef) Package() string {
	return path.Base(path.Dir(a.Source))
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
		return "actionlog.Err"
	}
}

func genActions(tpl *template.Template, dd ...*actionsDef) (err error) {
	var (
		// Will only be generated if file does not exist previously
		tplActionsGen = tpl.Lookup("actions.gen.go.tpl")

		dst string
	)

	for _, d := range dd {
		// Generic code, actions for every resource goes to a separated file
		dst = path.Join(d.outputDir, path.Base(d.Source)[:strings.LastIndex(path.Base(d.Source), ".")]+".gen.go")
		err = goTemplate(dst, tplActionsGen, d)
		if err != nil {
			return
		}

	}

	return nil
}
