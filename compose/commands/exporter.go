package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	sqlTypes "github.com/jmoiron/sqlx/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func Exporter(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Specify one ("modules", "pages", "charts", "permissions") or more resources to export`,

		Run: func(cmd *cobra.Command, args []string) {

			c.InitServices(ctx, c)

			ctx = auth.SetSuperUserContext(ctx)

			var (
				nsFlag = cmd.Flags().Lookup("namespace").Value.String()
				ns     *types.Namespace
				err    error

				out = Compose{
					Namespaces: map[string]Namespace{},
				}
				nsOut = Namespace{}
			)

			if namespaceID, _ := strconv.ParseUint(nsFlag, 10, 64); namespaceID > 0 {
				ns, err = service.DefaultNamespace.FindByID(namespaceID)
				if err != repository.ErrNamespaceNotFound {
					cli.HandleError(err)
				}
			} else if ns, err = service.DefaultNamespace.FindByHandle(nsFlag); err != nil {
				if err != repository.ErrNamespaceNotFound {
					cli.HandleError(err)
				}
			}

			roles, err = service.DefaultSystemRole.Find(ctx)
			cli.HandleError(err)

			modules, _, err := service.DefaultModule.Find(types.ModuleFilter{NamespaceID: ns.ID})
			cli.HandleError(err)

			pages, _, err := service.DefaultPage.Find(types.PageFilter{NamespaceID: ns.ID})
			cli.HandleError(err)

			charts, _, err := service.DefaultChart.Find(types.ChartFilter{NamespaceID: ns.ID})
			cli.HandleError(err)

			y := yaml.NewEncoder(cmd.OutOrStdout())

			nsOut.Name = ns.Name
			nsOut.Handle = ns.Slug
			nsOut.Enabled = ns.Enabled
			nsOut.Meta = ns.Meta

			nsOut.Allow = expResourcePermissions(permissions.Allow, ns.PermissionResource())
			nsOut.Deny = expResourcePermissions(permissions.Deny, ns.PermissionResource())

			for _, arg := range args {
				switch arg {
				case "module", "modules":
					nsOut.Modules = expModules(modules)
				case "chart", "charts":
					nsOut.Charts = expCharts(charts, modules)
				case "page", "pages":
					nsOut.Pages = expPages(0, pages, modules, charts)
				case "allow", "deny", "permission", "permissions":
					out.Allow = expServicePermissions(permissions.Allow)
					out.Deny = expServicePermissions(permissions.Deny)
				}
			}

			out.Namespaces[ns.Slug] = nsOut

			_, _ = y, out
			cli.HandleError(y.Encode(out))
		},
	}

	cmd.Flags().String("namespace", "crm", "Export namespace resources (by ID or string)")

	return cmd
}

// This is PoC for exporting compose resources
//

type (
	Compose struct {
		Namespaces map[string]Namespace

		Allow map[string]map[string][]string `yaml:",omitempty"`
		Deny  map[string]map[string][]string `yaml:",omitempty"`
	}

	Namespace struct {
		Name    string              `yaml:",omitempty"`
		Handle  string              `yaml:",omitempty"`
		Enabled bool                `yaml:",omitempty"`
		Meta    types.NamespaceMeta `yaml:",omitempty"`

		Modules map[string]Module `yaml:",omitempty"`
		Pages   map[string]Page   `yaml:",omitempty"`
		Charts  map[string]Chart  `yaml:",omitempty"`

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	Module struct {
		Name   string
		Meta   string `yaml:"meta,omitempty"`
		Fields map[string]Field

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	Field struct {
		Label string `yaml:",omitempty"`
		Kind  string

		Options types.ModuleFieldOptions `yaml:",omitempty"`

		Private  bool `yaml:",omitempty"`
		Required bool `yaml:",omitempty"`
		Visible  bool `yaml:",omitempty"`
		Multi    bool `yaml:",omitempty"`
		// DefaultValue string

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	Page struct {
		Module      string `yaml:",omitempty"`
		Title       string `yaml:",omitempty"`
		Description string `yaml:",omitempty"`

		Blocks types.PageBlocks `yaml:",omitempty"`

		Pages map[string]Page `yaml:",omitempty"`

		Visible bool

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	Chart struct {
		Name string `yaml:",omitempty"`

		Config ChartConfig `yaml:",omitempty"`

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	ChartConfig struct {
		Reports []*ChartConfigReport
	}
	ChartConfigReport struct {
		types.ChartConfigReport
		Module string `json:"module"`
	}
)

var (
	roles sysTypes.RoleSet
)

func expModules(mm types.ModuleSet) (o map[string]Module) {
	o = map[string]Module{}

	for _, m := range mm {
		module := Module{
			Name:   m.Name,
			Fields: expModuleFields(m.Fields, mm),

			Allow: expResourcePermissions(permissions.Allow, types.ModulePermissionResource),
			Deny:  expResourcePermissions(permissions.Deny, types.ModulePermissionResource),
		}

		if meta := expModuleMetaCleanup(m.Meta); len(meta) > 0 {
			module.Meta = meta
		}

		handle := makeHandleFromName(m.Name, m.Handle, "module-%d", m.ID)
		o[handle] = module
	}

	return
}

func expModuleMetaCleanup(meta sqlTypes.JSONText) string {
	var aux interface{}
	err := meta.Unmarshal(&aux)
	cli.HandleError(err)

	if kv, ok := aux.(map[string]interface{}); !ok {
		return ""
	} else if _, has := kv["admin"]; has {
		delete(kv, "admin")
		if len(kv) == 0 {
			return ""
		}
		meta, err = json.Marshal(kv)
		cli.HandleError(err)
	}

	return meta.String()
}

func expModuleFields(ff types.ModuleFieldSet, modules types.ModuleSet) (o map[string]Field) {
	o = make(map[string]Field)

	for _, f := range ff {
		o[f.Name] = Field{
			Label:    f.Label,
			Kind:     f.Kind,
			Options:  expModuleFieldOptions(f, modules),
			Private:  f.Private,
			Required: f.Required,
			Visible:  f.Visible,
			Multi:    f.Multi,

			Allow: expResourcePermissions(permissions.Allow, types.ModuleFieldPermissionResource),
			Deny:  expResourcePermissions(permissions.Deny, types.ModuleFieldPermissionResource),
		}
	}

	return
}

func expModuleFieldOptions(f *types.ModuleField, modules types.ModuleSet) types.ModuleFieldOptions {
	out := f.Options

	if moduleIDstr, has := out["moduleID"].(string); has {
		delete(out, "moduleID")
		out["module"] = "Error: module with ID " + moduleIDstr + " does not exist"
		if moduleID, _ := strconv.ParseUint(moduleIDstr, 10, 64); moduleID > 0 {
			if module := modules.FindByID(moduleID); module != nil {
				out["module"] = makeHandleFromName(module.Name, module.Handle, "module-%d", module.ID)
			}
		}
	}

	// Remove extra options to keep the output tidy

	rmDefault := func(k string, def interface{}) {
		if v, ok := out[k]; ok && v == def {
			delete(out, k)
		}
	}

	rmFalse := func(f string) {
		rmDefault(f, false)
	}

	rmDefault("multiDelimiter", "\n")

	switch f.Kind {
	case "Number":
		rmDefault("format", 0)
		rmDefault("precision", 0)
		rmDefault("prefix", "")
		rmDefault("suffix", "")
	case "DateTime":
		rmFalse("onlyDate")
		rmFalse("onlyFutureValues")
		rmFalse("onlyPastValues")
		rmFalse("onlyTime")
		rmFalse("outputRelative")
		rmDefault("format", "")
	case "Url":
		rmFalse("onlySecure")
		rmFalse("outputPlain")
		rmFalse("trimFragment")
		rmFalse("trimPath")
		rmFalse("trimQuery")
	case "Email":
		rmFalse("outputPlain")
	case "String":
		rmFalse("multiLine")
		rmFalse("useRichTextEditor")
	}

	return out
}

func expPages(parentID uint64, pages types.PageSet, modules types.ModuleSet, charts types.ChartSet) (o map[string]Page) {
	var (
		children = pages.FindByParent(parentID)
		handle   string
	)
	o = map[string]Page{}

	for _, child := range children {
		page := Page{
			Title:       child.Title,
			Description: child.Description,
			Blocks:      expPageBlocks(child.Blocks, pages, modules, charts),
			Pages:       expPages(child.ID, pages, modules, charts),
			Visible:     child.Visible,

			Allow: expResourcePermissions(permissions.Allow, types.PagePermissionResource),
			Deny:  expResourcePermissions(permissions.Deny, types.PagePermissionResource),
		}

		if child.ModuleID > 0 {
			m := modules.FindByID(child.ModuleID)
			if m == nil {
				page.Module = fmt.Sprintf("Error: module with ID %d does not exist", child.ModuleID)
			} else {
				page.Module = makeHandleFromName(m.Name, m.Handle, "module-%d", child.ModuleID)

				if child.Handle == "" {
					// Reuse module's handle for page
					handle = makeHandleFromName(m.Name, m.Handle, "record-page-%d", child.ModuleID)
				}
			}
		} else {
			handle = makeHandleFromName(child.Title, child.Handle, "page-%d", child.ID)
		}

		o[handle] = page
	}

	return
}

func expPageBlocks(in types.PageBlocks, pages types.PageSet, modules types.ModuleSet, charts types.ChartSet) types.PageBlocks {
	out := types.PageBlocks(in)

	// Remove extra options to keep the output tidy

	rmDefault := func(kv map[string]interface{}, k string, def interface{}) {
		if v, ok := kv[k]; ok && v == def {
			delete(kv, k)
		}
	}

	rmFalse := func(kv map[string]interface{}, f string) {
		rmDefault(kv, f, false)
	}

	for i := range out {
		if ff, has := out[i].Options["fields"].([]interface{}); has {
			// Trim out obsolete field info
			for fi := range ff {
				if f, ok := ff[fi].(map[string]interface{}); ok {
					ff[fi] = map[string]string{
						"name": f["name"].(string),
					}
				}
			}
		}

		if moduleIDstr, has := out[i].Options["moduleID"].(string); has {
			delete(out[i].Options, "moduleID")
			out[i].Options["module"] = "Error: module with ID " + moduleIDstr + " does not exist"
			if moduleID, _ := strconv.ParseUint(moduleIDstr, 10, 64); moduleID > 0 {
				if module := modules.FindByID(moduleID); module != nil {
					out[i].Options["module"] = makeHandleFromName(module.Name, module.Handle, "module-%d", module.ID)
				}
			}
		}

		if pageIDstr, has := out[i].Options["pageID"].(string); has {
			delete(out[i].Options, "pageID")
			out[i].Options["page"] = "Error: page with ID " + pageIDstr + " does not exist"
			if pageID, _ := strconv.ParseUint(pageIDstr, 10, 64); pageID > 0 {
				if page := pages.FindByID(pageID); page != nil {
					out[i].Options["page"] = makeHandleFromName(page.Title, page.Handle, "page-%d", page.ID)
				}
			}
		}

		if chartIDstr, has := out[i].Options["chartID"].(string); has {
			delete(out[i].Options, "chartID")
			out[i].Options["chart"] = "Error: chart with ID " + chartIDstr + " does not exist"
			if chartID, _ := strconv.ParseUint(chartIDstr, 10, 64); chartID > 0 {
				if chart := charts.FindByID(chartID); chart != nil {
					out[i].Options["chart"] = makeHandleFromName(chart.Name, chart.Handle, "chart-%d", chart.ID)
				}
			}
		}

		rmFalse(out[i].Options, "hideAddButton")
		rmFalse(out[i].Options, "hideHeader")
		rmFalse(out[i].Options, "hidePaging")
		rmFalse(out[i].Options, "hideSearch")
		rmFalse(out[i].Options, "hideSorting")
		rmFalse(out[i].Options, "allowExport")
	}

	return out
}

func expCharts(charts types.ChartSet, modules types.ModuleSet) (o map[string]Chart) {
	o = map[string]Chart{}

	for _, c := range charts {
		chart := Chart{
			Name:   c.Name,
			Config: ChartConfig{Reports: make([]*ChartConfigReport, len(c.Config.Reports))},

			Allow: expResourcePermissions(permissions.Allow, types.ChartPermissionResource),
			Deny:  expResourcePermissions(permissions.Deny, types.ChartPermissionResource),
		}

		for i, r := range c.Config.Reports {
			chart.Config.Reports[i] = &ChartConfigReport{
				ChartConfigReport: *r,
			}

			if r.ModuleID > 0 {
				module := modules.FindByID(r.ModuleID)
				chart.Config.Reports[i].ModuleID = 0
				chart.Config.Reports[i].Module =
					makeHandleFromName(module.Name, module.Handle, "module-%d", module.ID)
			}

		}

		handle := makeHandleFromName(c.Name, c.Handle, "chart-%d", c.ID)

		o[handle] = chart
	}

	return
}

func expServicePermissions(access permissions.Access) map[string]map[string][]string {
	var (
		has   bool
		res   string
		rules permissions.RuleSet
		sp    = make(map[string]map[string][]string)
	)

	for _, r := range roles {
		rules = service.DefaultPermissions.FindRulesByRoleID(r.ID)

		if len(rules) == 0 {
			continue
		}

		for _, rule := range rules {
			if rule.Resource.GetService() != rule.Resource && !rule.Resource.HasWildcard() {
				continue
			}

			res = strings.TrimRight(rule.Resource.String(), ":*")

			if _, has = sp[r.Handle]; !has {
				sp[r.Handle] = map[string][]string{}
			}

			if _, has = sp[r.Handle][res]; !has {
				sp[r.Handle][res] = make([]string, 0)
			}

			sp[r.Handle][res] = append(sp[r.Handle][res], rule.Operation.String())
		}
	}

	return sp
}

func expResourcePermissions(access permissions.Access, resource permissions.Resource) map[string][]string {
	var (
		has   bool
		rules permissions.RuleSet
		sp    = make(map[string][]string)
	)

	for _, r := range roles {
		rules = service.DefaultPermissions.FindRulesByRoleID(r.ID)

		if len(rules) == 0 {
			continue
		}

		for _, rule := range rules {
			if rule.Resource != resource {
				continue
			}

			if _, has = sp[r.Handle]; !has {
				sp[r.Handle] = make([]string, 0)

			}

			sp[r.Handle] = append(sp[r.Handle], rule.Operation.String())
		}
	}

	return sp
}

func makeHandleFromName(name, currentHandle, def string, id uint64) string {
	if currentHandle != "" {
		return currentHandle
	}

	newHandle := strings.ReplaceAll(name, " ", "_")
	newHandle = regexp.MustCompile(`[^0-9A-Za-z_\-.]+`).ReplaceAllString(newHandle, "")
	if handle.IsValid(newHandle) {
		return newHandle
	}

	return fmt.Sprintf(def, id)
}
