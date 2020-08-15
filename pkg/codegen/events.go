package codegen

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
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

	evResourceDefMap map[string]evResourceDef

	evResourceDef struct {
		// used as string
		ResourceString string

		// used as (go) ident
		ResourceIdent string

		// used for filename
		ResourceFile string

		On          []string     `yaml:"on"`
		BeforeAfter []string     `yaml:"ba"`
		Properties  []eventProps `yaml:"props"`
		Result      string       `yaml:"result"`
	}

	eventProps struct {
		Name string
		Type string

		// Import path for prop type, use package's type by default (see importTypePathTpl)
		Import string

		// Set property internally only, not via constructor
		Internal bool

		// Do not allow change of the variable through
		Immutable bool
	}
)

func procEvents() ([]*eventsDef, error) {
	// <app>/service/event/events.yaml
	const (
		importTypePathTpl = "github.com/cortezaproject/corteza-server/%s/types"
		importAuthPath    = "github.com/cortezaproject/corteza-server/pkg/auth"
	)

	var (
		dd = make([]*eventsDef, 0)
	)

	mm, err := filepath.Glob(filepath.Join("*", "service", "event", "events.yaml"))
	if err != nil {
		return nil, fmt.Errorf("glob failed: %w", err)
	}

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
				Resources: map[string]evResourceDef{},
			}
		)

		if err := yaml.NewDecoder(f).Decode(e); err != nil {
			return nil, err
		}

		for resName, evDef := range e {

			d.Imports = []string{fmt.Sprintf(importTypePathTpl, d.App)}
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
			evDef.Properties = append(evDef.Properties, eventProps{
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

func genEvents(tpl *template.Template, dd []*eventsDef) (err error) {
	var (
		// Will only be generated if file does not exist previously
		tplEvents = tpl.Lookup("events.go.tpl")

		// Always regenerated
		tplEventsGen = tpl.Lookup("events.gen.go.tpl")

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
