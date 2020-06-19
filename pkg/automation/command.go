package automation

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/titpetric/factory"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

//
//		  _____                                _           _
//		 |  __ \                              | |         | |
//		 | |  | | ___ _ __  _ __ ___  ___ __ _| |_ ___  __| |
//		 | |  | |/ _ \ '_ \| '__/ _ \/ __/ _` | __/ _ \/ _` |
//		 | |__| |  __/ |_) | | |  __/ (_| (_| | ||  __/ (_| |
//		 |_____/ \___| .__/|_|  \___|\___\__,_|\__\___|\__,_|
//					 | |
//					 |_|
//
//
//
// This automation package is kept only to aid in export of deprecated
// trigger & script format from database into files
//
// Package was refactored and replaced by pkg/corredor, pkg/eventbus and pkg/scheduler
//
// Scheduled for removal in 2020.6
//

type (
	module struct {
		ID     uint64
		Handle string
		Name   string
	}
)

var (
	sanitizer = regexp.MustCompile(`[^a-zA-Z0-9_\-.]+`)

	modules []*module

	// Language=GoTemplate
	scriptTemplateRaw string = `
export default {
  label: {{ quote $.Name }},
  desc: '...',
  triggers ({ on }) {
{{- range $t := $.Triggers }}
    // auto-migrated
    //   ID:       {{ $t.ID }}
    //   Created:  {{ $t.CreatedAt }}
    //   Updated:  {{ $t.UpdatedAt }}
{{ if not $t.Enabled }}/* disabled {{ end }}
	return [
    on({{ quote $t.Event }}){{- if $.RunAs }}
      .as({{ quote $t.RunAs }}){{ end }}
      .for({{ quote $t.Resource }})
      {{- makeConditionFn $t -}}
{{ if not $t.Enabled }}*/{{ end }}
	]
{{ end -}}
  },

  async handler ({ $namespace, $module, $record }, { log, ComposeUI, Compose }) {
    {{ indent .Source 4 }}
  }
}
`

	tpl *template.Template
)

func init() {
	var err error
	tpl = template.New("").Funcs(map[string]interface{}{
		//"camelCase":  camelCase,
		//"makeEvents": makeEvents,

		"dump": func(s ...interface{}) string {
			return spew.Sdump(s...)
		},

		"quote": func(s string) string {
			return `'` + s + `'`
		},

		"makeConditionFn": func(t *Trigger) string {
			isStdBeforeAfter := (strings.HasSuffix(t.Event, "Create") ||
				strings.HasSuffix(t.Event, "Update") ||
				strings.HasSuffix(t.Event, "Delete")) &&
				(strings.HasPrefix(t.Event, "before") ||
					strings.HasPrefix(t.Event, "after"))

			if t.Condition != "" {
				if t.Resource == "compose:record" && (isStdBeforeAfter || t.Event == "manual") {
					id, _ := strconv.ParseUint(t.Condition, 10, 64)
					if id > 0 {
						for _, m := range modules {
							if m.ID == id {
								cnd := m.Handle
								if cnd == "" {
									cnd = t.Condition
								}

								return fmt.Sprintf("\n      .where('module', '%s'), // module (%d) %s\n", cnd, m.ID, m.Name)
							}
						}

						return fmt.Sprintf("\n      .where('module', '%s'), // module not found, could not translate ID to handle\n", t.Condition)
					}
				} else if t.Event == "deferred" {
					return fmt.Sprintf("\n      .where('timestamp', '%s'),\n", t.Condition)
				} else if t.Event == "interval" {
					return fmt.Sprintf("\n      .where('interval', '%s'),\n", t.Condition)
				}
				return fmt.Sprintf(", // unresolvable condition - %s \n", t.Condition)
			}

			return ",\n"
		},

		"makeEventFn": func(ev string) string {
			tpl := ".%s('%s')"

			if strings.HasPrefix(ev, "before") {
				return fmt.Sprintf(tpl, "before", strings.ToLower(ev))
			}

			if strings.HasPrefix(ev, "after") {
				return fmt.Sprintf(tpl, "after", strings.ToLower(ev))
			}

			if ev == "deferred" {
				ev = "timestamp"
			}

			return fmt.Sprintf(tpl, "on", ev)
		},

		"indent": func(s string, spaces int) (o string) {
			for _, l := range strings.Split(s, "\n") {
				o = o + strings.Repeat(" ", spaces) + strings.TrimRightFunc(l, unicode.IsSpace) + "\n"
			}

			return
		},
	})

	tpl, err = tpl.Parse(scriptTemplateRaw)
	if err != nil {
		panic(err)
	}
}

func ScriptExporter(subsys string) *cobra.Command {
	var (
		tblPrefix = subsys
		isCompose = subsys == "compose"
		isSystem  = subsys == "system"
	)

	if isSystem {
		tblPrefix = "sys"
	}

	cmd := &cobra.Command{
		Use:   "script-migrator",
		Short: "Migrates automation scripts",
		Long:  "Scans system & compose automation tables for scripts & Triggers and creates script files",

		Run: func(cmd *cobra.Command, args []string) {
			var (
				dstPath = cmd.Flags().Lookup("dst").Value.String()

				f *os.File

				err error
				ss  = make([]*Script, 0)
				tt  = make([]*Trigger, 0)

				ctx = cli.Context()

				db = factory.Database.MustGet().With(ctx) //.Quiet()

				// Skip deleted, scripts, ones named test and those with empty source
				scriptQuery = squirrel.
						Select("*").
						From(tblPrefix+"_automation_script").
						Where("name <> ?", "test").
						Where("source <> ?", "").
						Where("deleted_at IS NULL").
						OrderBy("RAND()")

				// Skip deleted triggers
				triggerQuery = squirrel.
						Select("*").
						From(tblPrefix + "_automation_trigger").
						Where("deleted_at IS NULL")

				// (for compose)
				// preload modules - we want to refer to them (if possible) by handle
				moduleQuery = squirrel.
						Select("id", "handle", "name").
						From(tblPrefix + "_module")
			)

			err = rh.FetchAll(db, scriptQuery, &ss)
			cli.HandleError(err)

			err = rh.FetchAll(db, triggerQuery, &tt)
			cli.HandleError(err)

			if isCompose {
				err = rh.FetchAll(db, moduleQuery, &modules)
				cli.HandleError(err)
			}

			cmd.Printf("Found %d scripts and %d triggers\n", len(ss), len(tt))

			if len(dstPath) == 0 {
				// No destination, just output list of scripts
				for _, s := range ss {
					cmd.Printf("%s\n", s.Name)
				}
			} else {
				done := make(map[string]bool)

				for _, s := range ss {
					// sanitize script name into safe file name
					sname := sanitizer.ReplaceAllString(strings.ReplaceAll(s.Name, " ", "_"), "")
					fname := ""
					names := []string{
						fmt.Sprintf("%s.js", sname),
						fmt.Sprintf("%s_%d.js", sname, s.ID),
					}

					for _, fname = range names {
						if !done[fname] {
							break
						}
					}

					for _, t := range tt {
						if t.ScriptID == s.ID {
							s.Triggers = append(s.Triggers, t)
						}
					}

					fullpath := path.Join(dstPath, fname)

					cmd.Printf(" exporting to %s\n", fullpath)
					f, err = os.Create(fullpath)
					cli.HandleError(err)
					cli.HandleError(tpl.Execute(f, s))
					done[fname] = true
				}
			}
		},
	}

	cmd.Flags().String("dst", "", "Where to export the scripts")

	return cmd
}
