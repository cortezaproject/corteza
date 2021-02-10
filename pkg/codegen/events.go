package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/cortezaproject/corteza-server/pkg/slice"
	"gopkg.in/yaml.v3"
)

type (
	// definitions are in one file
	eventsDef struct {
		Package   string
		App       string
		Source    string
		outputDir string

		// List of imports
		// Used only by generated file and not pre-generated-user-file
		Imports []string

		Resources evResourceDefMap
	}

	evResourceDefMap map[string]*evResourceDef

	evResourceDef struct {
		// used as string
		ResourceString string

		// used as (go) ident
		ResourceIdent string

		// used for filename
		ResourceFile string

		On          []string      `yaml:"on"`
		BeforeAfter []string      `yaml:"ba"`
		Properties  []*eventProps `yaml:"props"`
		Result      string        `yaml:"result"`
	}

	eventProps struct {
		Name     string
		Type     string
		ExprType string

		// Import path for prop type, use package's type by default (see importTypePathTpl)
		Import string

		// Set property internally only, not via constructor
		Internal bool

		// Do not allow change of the variable through
		Immutable bool
	}
)

func procEvents(mm ...string) (dd []*eventsDef, err error) {
	// <app>/service/event/events.yaml
	const (
		importTypePathTpl = "github.com/cortezaproject/corteza-server/%s/types"
		importAutoPathTpl = "github.com/cortezaproject/corteza-server/%s/automation"
		importAuthPath    = "github.com/cortezaproject/corteza-server/pkg/auth"
	)

	dd = make([]*eventsDef, 0)
	for _, m := range mm {
		f, err := os.Open(m)
		if err != nil {
			return nil, fmt.Errorf("%s read failed: %w", m, err)
		}

		defer f.Close()

		var (
			e = evResourceDefMap{}
			d = &eventsDef{
				Package:   "event",
				Source:    m,
				App:       m[:strings.Index(m, "/")],
				outputDir: path.Dir(m),
				Resources: map[string]*evResourceDef{},
			}
		)

		if err := yaml.NewDecoder(f).Decode(e); err != nil {
			return nil, err
		}

		for resName, evDef := range e {

			d.Imports = []string{
				fmt.Sprintf(importTypePathTpl, d.App),
			}

			if d.App != "messaging" {
				d.Imports = append(d.Imports, fmt.Sprintf(importAutoPathTpl, d.App))
			}

			evDef.ResourceString = resName

			if l := strings.Index(resName, ":"); l > 0 {
				evDef.ResourceIdent = resName[l+1:]
			} else {
				evDef.ResourceIdent = resName
			}

			// make filename
			evDef.ResourceFile = strings.ReplaceAll(evDef.ResourceIdent, ":", "_")
			evDef.ResourceFile = strings.ReplaceAll(evDef.ResourceFile, "-", "_")

			// make identifier (string that will be used for struct name)
			evDef.ResourceIdent = camelCase(strings.Split(evDef.ResourceFile, "_")...)

			// Prepare the data

			// no default ("result") result set, use first one from properties
			if evDef.Result == "" && len(evDef.Properties) > 0 {
				evDef.Result = evDef.Properties[0].Name
			}

			// Invoker - user that invoked (triggered) the event
			evDef.Properties = append(evDef.Properties, &eventProps{
				Name:      "invoker",
				Type:      "auth.Identifiable",
				Import:    importAuthPath,
				Immutable: false,
				Internal:  true,
			})

			// Ensure all imports are checked
			for _, p := range evDef.Properties {
				if p.Import == "" {
					if strings.HasPrefix(p.Type, "*types.") || strings.HasPrefix(p.Type, "types.") {
						p.Import = fmt.Sprintf(importTypePathTpl, d.App)
					}
				}

				if p.Import != "" && !slice.HasString(d.Imports, p.Import) {
					d.Imports = append(d.Imports, p.Import)
				}

				p.Import = ""
			}

			d.Resources[resName] = evDef
		}

		dd = append(dd, d)
	}

	return dd, nil
}

func expandEventTypes(ee []*eventsDef, tt []*exprTypesDef) {
	// index of all known types
	expTypes := make(map[string]*exprTypeDef)
	goTypes := make(map[string]string)

	for _, t := range tt {
		for typ, d := range t.Types {
			expTypes[typ] = d
			goTypes[d.As] = typ
		}
	}

	for _, e := range ee {
		for _, r := range e.Resources {
			for _, p := range r.Properties {
				if p.ExprType != "" && expTypes[p.ExprType] == nil {
					fmt.Printf("unknown type %q used for param %q for events on resource %s\n", p.ExprType, p.Name, r.ResourceString)
				}

				if p.ExprType == "" && goTypes[p.Type] != "" {
					p.ExprType = goTypes[p.Type]
				}
			}
		}
	}
}

func genEvents(tpl *template.Template, dd ...*eventsDef) (err error) {
	var (
		// Will only be generated if file does not exist previously
		tplEvents = tpl.Lookup("events.go.tpl")

		// Always regenerated
		tplEventsGen = tpl.Lookup("events.gen.go.tpl")

		// List of event-type definitions for automation REST endpoint
		tplAutomationRestDefGen = tpl.Lookup("events_rest_def.gen.go.tpl")

		dst string
	)

	for _, d := range dd {
		// Generic code, all events go into one file (per app)
		err = goTemplate(path.Join(d.outputDir, "events.gen.go"), tplEventsGen, d)
		if err != nil {
			return
		}

		for _, r := range d.Resources {
			dst = path.Join(d.outputDir, r.ResourceFile+".go")
			_, err = os.Stat(dst)
			if os.IsNotExist(err) {
				err = goTemplate(dst, tplEvents, map[string]interface{}{
					"Package":       d.Package,
					"ResourceIdent": r.ResourceIdent,
				})
			}

			if err != nil {
				return
			}
		}
	}

	// Remove messaging
	var msgIndex = -1
	for i := range dd {
		if dd[i].App == "messaging" {
			msgIndex = i
		}
	}

	if msgIndex > -1 {
		dd = append(dd[:msgIndex], dd[msgIndex+1:]...)
	}

	err = goTemplate(
		path.Join("automation", "rest", "eventTypes.gen.go"),
		tplAutomationRestDefGen,
		map[string]interface{}{
			"Definitions": dd,
			"Imports":     collectEventDefImports("", dd...),
		})
	if err != nil {
		return
	}

	return nil
}

// Merge on/before/after events
func (def evResourceDef) Events() []string {
	return append(
		makeEventGroup("on", def.On),
		append(
			makeEventGroup("before", def.BeforeAfter),
			makeEventGroup("after", def.BeforeAfter)...,
		)...,
	)
}

func makeEventGroup(pfix string, ee []string) (out []string) {
	for _, e := range ee {
		out = append(out, pfix+strings.ToUpper(e[:1])+e[1:])
	}

	return
}

func genEventsDocs(tpl *template.Template, docsPath string, dd ...*eventsDef) (err error) {
	var (
		tplEventsAdoc = tpl.Lookup("events.gen.adoc.tpl")

		dst string
	)

	dst = path.Join(docsPath, "resource-events.gen.adoc")
	return plainTemplate(dst, tplEventsAdoc, map[string]interface{}{
		"Definitions": dd,
	})
}

func collectEventDefImports(basePkg string, dd ...*eventsDef) []string {
	ii := make([]string, 0, len(dd))
	for _, d := range dd {
		for _, i := range d.Imports {
			if !slice.HasString(ii, i) && (basePkg == "" || !strings.HasSuffix(i, basePkg)) {
				ii = append(ii, i)
			}
		}
	}

	return ii
}
