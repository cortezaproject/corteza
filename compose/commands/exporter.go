package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	sysExporter "github.com/cortezaproject/corteza-server/system/exporter"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	sqlTypes "github.com/jmoiron/sqlx/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Exporter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Specify one ("modules", "pages", "charts", "permissions") or more resources to export`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx    = auth.SetSuperUserContext(cli.Context())
				nsFlag = cmd.Flags().Lookup("namespace").Value.String()
				sFlag  = cmd.Flags().Lookup("settings").Changed
				pFlag  = cmd.Flags().Lookup("permissions").Changed

				out = &Compose{
					Namespaces: map[string]Namespace{},
					Settings:   yaml.MapSlice{},
				}
			)
			if nsFlag == "" && !sFlag && !pFlag {
				cli.HandleError(fmt.Errorf("Specify namespace or setting or permissions flag"))
			}

			if nsFlag != "" {
				nsExporter(ctx, out, nsFlag, args)
			}

			if sFlag {
				settingExporter(ctx, out)
			}

			if pFlag {
				permissionExporter(ctx, out)
			}

			y := yaml.NewEncoder(cmd.OutOrStdout())
			cli.HandleError(y.Encode(out))
		},
	}

	cmd.Flags().String("namespace", "", "Export namespace resources (by ID or string)")
	cmd.Flags().BoolP("settings", "s", false, "Export settings")
	cmd.Flags().BoolP("permissions", "p", false, "Export system permissions")

	return cmd
}

func nsExporter(ctx context.Context, out *Compose, nsFlag string, args []string) {
	var (
		ns  *types.Namespace
		err error

		nsOut = Namespace{}
	)

	if namespaceID, _ := strconv.ParseUint(nsFlag, 10, 64); namespaceID > 0 {
		ns, err = service.DefaultNamespace.FindByID(ctx, namespaceID)
		if errors.IsNotFound(err) {
			cli.HandleError(err)
		}
	} else if ns, err = service.DefaultNamespace.FindByHandle(ctx, nsFlag); err != nil {
		if errors.IsNotFound(err) {
			cli.HandleError(err)
		}
	}

	// roles, err = service.DefaultSystemRole.Find(ctx)
	// cli.HandleError(err)
	// At the moment, we can not load roles from system service
	// so we'll just use static set of known roles
	//
	// Roles are use for resolving access control
	roles = sysTypes.RoleSet{
		&sysTypes.Role{ID: rbac.EveryoneRoleID, Handle: "everyone"},
		&sysTypes.Role{ID: rbac.AdminsRoleID, Handle: "admins"},
	}

	modules, _, err := service.DefaultModule.Find(ctx, types.ModuleFilter{NamespaceID: ns.ID})
	cli.HandleError(err)

	pages, _, err := service.DefaultPage.Find(ctx, types.PageFilter{NamespaceID: ns.ID})
	cli.HandleError(err)

	charts, _, err := service.DefaultChart.Find(ctx, types.ChartFilter{NamespaceID: ns.ID})
	cli.HandleError(err)

	// nsOut.Name = ns.Name
	// nsOut.Handle = ns.Slug
	// nsOut.Always = ns.Always
	// nsOut.Meta = ns.Meta
	//
	// nsOut.Allow = sysExporter.ExportableResourcePermissions(roles, service.DefaultRBAC, permissions.Allow, ns.RBACResource())
	// nsOut.Deny = sysExporter.ExportableResourcePermissions(roles, service.DefaultRBAC, permissions.Deny, ns.RBACResource())

	for _, arg := range args {
		switch arg {
		case "module", "modules":
			nsOut.Modules = expModules(modules)
		case "chart", "charts":
			nsOut.Charts = expCharts(charts, modules)
		case "page", "pages":
			nsOut.Pages = expPages(0, pages, modules, charts)
		}
	}

	nsOut.Namespace = ns.Slug
	out.Namespaces[ns.Slug] = nsOut
}

func settingExporter(ctx context.Context, out *Compose) {
	var (
		err error
	)

	ss, err := sysService.DefaultSettings.FindByPrefix(ctx)
	cli.HandleError(err)

	out.Settings = settings.Export(ss)
}

func permissionExporter(ctx context.Context, out *Compose) {
	roles := sysTypes.RoleSet{
		&sysTypes.Role{ID: rbac.EveryoneRoleID, Handle: "everyone"},
		&sysTypes.Role{ID: rbac.AdminsRoleID, Handle: "admins"},
	}

	out.Allow = sysExporter.ExportableServicePermissions(roles, rbac.Global(), rbac.Allow)
	out.Deny = sysExporter.ExportableServicePermissions(roles, rbac.Global(), rbac.Deny)
}

// This is PoC for exporting compose resources
//

type (
	Compose struct {
		Namespaces map[string]Namespace `yaml:",omitempty"`
		Settings   yaml.MapSlice        `yaml:",omitempty"`

		Allow map[string]map[string][]string `yaml:",omitempty"`
		Deny  map[string]map[string][]string `yaml:",omitempty"`
	}

	Namespace struct {
		// This is used when exporting one single namespace
		Namespace string `yaml:",omitempty"`

		Name    string              `yaml:",omitempty"`
		Handle  string              `yaml:",omitempty"`
		Enabled bool                `yaml:",omitempty"`
		Meta    types.NamespaceMeta `yaml:",omitempty"`

		Modules map[string]Module `yaml:",omitempty"`
		Pages   yaml.MapSlice     `yaml:",omitempty"`
		Charts  map[string]Chart  `yaml:",omitempty"`
		Scripts map[string]Script `yaml:",omitempty"`

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	Module struct {
		Name   string
		Meta   string `yaml:"meta,omitempty"`
		Fields yaml.MapSlice

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

		Pages yaml.MapSlice `yaml:",omitempty"`

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

	Script struct {
		Source   string `yaml:"source"`
		Async    bool   `yaml:"async"`
		RunInUA  bool   `yaml:"runInUA"`
		Critical bool   `yaml:"critical"`
		Enabled  bool   `yaml:"enabled"`
		Timeout  uint   `yaml:"timeout"`

		Triggers []map[string]interface{} `yaml:"triggers,omitempty"`

		Allow map[string][]string `yaml:",omitempty"`
		Deny  map[string][]string `yaml:",omitempty"`
	}

	ChartConfig struct {
		Reports     []map[string]interface{}
		ColorScheme string `yaml:"colorScheme,omitempty"`
	}
)

var (
	// preloaded roles so we can
	//
	roles sysTypes.RoleSet

	// list of used page handles
	// we're exporting pages in a tree structure
	// so we need this to know if we've used a handle before
	//
	// non-autogenerated handles should not have this problem
	pagesHandles = make(map[string]bool)
)

func expSettings(ss sysTypes.SettingValueSet) (o yaml.MapSlice) {
	o = yaml.MapSlice{}
	for _, s := range ss {
		setting := yaml.MapItem{
			Key:   s.Name,
			Value: s.Value.String(),
		}
		o = append(o, setting)
	}

	return o
}

func expModules(mm types.ModuleSet) (o map[string]Module) {
	o = map[string]Module{}

	for _, m := range mm {
		module := Module{
			Name:   m.Name,
			Fields: expModuleFields(m.Fields, mm),

			Allow: sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Allow, types.ModuleRBACResource),
			Deny:  sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Deny, types.ModuleRBACResource),
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

func expModuleFields(ff types.ModuleFieldSet, modules types.ModuleSet) (o yaml.MapSlice) {
	o = make(yaml.MapSlice, len(ff))

	for i, f := range ff {
		o[i] = yaml.MapItem{
			Key: f.Name,
			Value: Field{
				Label:    f.Label,
				Kind:     f.Kind,
				Options:  expModuleFieldOptions(f, modules),
				Private:  f.Private,
				Required: f.Required,
				Visible:  f.Visible,
				Multi:    f.Multi,

				Allow: sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Allow, types.ModuleFieldRBACResource),
				Deny:  sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Deny, types.ModuleFieldRBACResource),
			},
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

func expPages(parentID uint64, pages types.PageSet, modules types.ModuleSet, charts types.ChartSet) (o yaml.MapSlice) {
	var (
		children = pages.FindByParent(parentID)
		handle   string
	)
	o = yaml.MapSlice{}

	for _, child := range children {
		page := Page{
			Title:       child.Title,
			Description: child.Description,
			Blocks:      expPageBlocks(child.Blocks, pages, modules, charts),
			Pages:       expPages(child.ID, pages, modules, charts),
			Visible:     child.Visible,

			Allow: sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Allow, types.PageRBACResource),
			Deny:  sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Deny, types.PageRBACResource),
		}

		if child.ModuleID > 0 {
			m := modules.FindByID(child.ModuleID)
			if m == nil {
				page.Module = fmt.Sprintf("Error: module with ID %d does not exist", child.ModuleID)
			} else {
				page.Module = makeHandleFromName(m.Name, m.Handle, "module-%d", child.ModuleID)
			}
		}

		handle = makeHandleFromName(child.Title, child.Handle, "page-%d", child.ID)
		if handle == "" || pagesHandles[handle] {
			// if handle exists, force simple handle with id
			handle = makeHandleFromName("", "", "page-%d", child.ID)
		}

		pagesHandles[handle] = true

		o = append(o, yaml.MapItem{
			Key:   handle,
			Value: page,
		})
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
		expEncodeRefFields(out[i].Options, "fields")
		expEncodeRefFields(out[i].Options, "editFields")

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

		switch out[i].Kind {
		case "Automation":
			expPageBlocks_Automation(&out[i])
		case "Calendar":
			expPageBlocks_Calendar(&out[i], modules)
		case "Metric":
			expPageBlocks_Metric(&out[i], modules)
		}
	}

	return out
}

func expPageBlocks_Automation(b *types.PageBlock) {
	bb := make([]interface{}, 0)
	_ = each(b.Options["buttons"], func(_ int, _ string, btn interface{}) error {
		button := map[string]interface{}{}

		_ = each(btn, func(_ int, k string, v interface{}) error {
			switch k {
			case "triggerID", "scriptID":
				// if s := scripts.FindByID(deinterfacer.ToUint64(v)); s != nil {
				// 	button["script"] = makeHandleFromName(s.Name, "", "automation-script-%d", s.ID)
				// }
			default:
				button[k] = v
			}

			return nil
		})

		bb = append(bb, button)
		return nil
	})

	b.Options["buttons"] = bb
}

func expPageBlocks_Calendar(b *types.PageBlock, modules types.ModuleSet) {
	ff := make([]interface{}, 0)
	_ = each(b.Options["feeds"], func(_ int, _ string, def interface{}) error {
		feed := map[string]interface{}{}

		_ = each(def, func(_ int, k string, v interface{}) error {
			switch k {
			case "options":
				expEncodeModuleInMap(v, modules)
				feed["options"] = v
			default:
				feed[k] = v
			}

			return nil
		})

		ff = append(ff, feed)
		return nil
	})

	b.Options["feeds"] = ff
}

func expPageBlocks_Metric(b *types.PageBlock, modules types.ModuleSet) {
	mm := make([]interface{}, 0)
	_ = each(b.Options["metrics"], func(_ int, _ string, def interface{}) error {
		metric := map[string]interface{}{}

		_ = each(def, func(_ int, k string, v interface{}) error {
			switch k {
			case "moduleID":
				metric["module"] = expEncodeModule(v, modules)
			default:
				metric[k] = v
			}

			return nil
		})

		mm = append(mm, metric)
		return nil
	})

	b.Options["metrics"] = mm
}

func expEncodeModule(id interface{}, modules types.ModuleSet) string {
	if moduleID := cast.ToUint64(id); moduleID > 0 {
		if module := modules.FindByID(moduleID); module != nil {
			return makeHandleFromName(module.Name, module.Handle, "module-%d", module.ID)
		}
	}

	return fmt.Sprintf("Error: module with ID  %v does not exist", id)
}

func expEncodeModuleInMap(in interface{}, modules types.ModuleSet) {
	m, ok := in.(map[string]interface{})
	if !ok {
		return
	}

	id, has := m["moduleID"]
	if !has {
		return
	}

	delete(m, "moduleID")
	m["module"] = expEncodeModule(id, modules)
}

// Trims out any extra info on referenced fields, keeps just name
func expEncodeRefFields(m map[string]interface{}, name string) {
	if ff, has := m[name].([]interface{}); has {
		// Trim out obsolete field info
		for fi := range ff {
			if f, ok := ff[fi].(map[string]interface{}); ok {
				ff[fi] = map[string]string{
					"name": f["name"].(string),
				}
			}
		}
	}
}

func expCharts(charts types.ChartSet, modules types.ModuleSet) (o map[string]Chart) {
	o = map[string]Chart{}

	for _, c := range charts {

		chart := Chart{
			Name: c.Name,
			Config: ChartConfig{
				Reports:     make([]map[string]interface{}, len(c.Config.Reports)),
				ColorScheme: c.Config.ColorScheme,
			},

			Allow: sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Allow, types.ChartRBACResource),
			Deny:  sysExporter.ExportableResourcePermissions(roles, rbac.Global(), rbac.Deny, types.ChartRBACResource),
		}

		for i, r := range c.Config.Reports {
			if len(r.Metrics) == 0 {
				continue
			}

			rOut := map[string]interface{}{
				"filter":     r.Filter,
				"metrics":    r.Metrics,
				"dimensions": r.Dimensions,
				"renderer":   r.Renderer,
			}

			if r.YAxis != nil {
				rOut["yAxis"] = r.YAxis
			}

			if r.ModuleID > 0 {
				module := modules.FindByID(r.ModuleID)
				rOut["module"] = makeHandleFromName(module.Name, module.Handle, "module-%d", module.ID)
			}

			chart.Config.Reports[i] = rOut
		}

		handle := makeHandleFromName(c.Name, c.Handle, "chart-%d", c.ID)

		o[handle] = chart
	}

	return
}

func makeHandleFromName(name, currentHandle, def string, id uint64) string {
	if currentHandle != "" {
		return currentHandle
	}

	newHandle := strings.ReplaceAll(strings.Title(name), " ", "")
	newHandle = regexp.MustCompile(`[^0-9A-Za-z_\-.]+`).ReplaceAllString(newHandle, "")
	if handle.IsValid(newHandle) {
		return newHandle
	}

	return fmt.Sprintf(def, id)
}

// each() tries to resolve interface{} into map[interface{}]interface{} or []interface{}
// and calls fn() for each entry
func each(i interface{}, fn func(int, string, interface{}) error) (err error) {
	if kv, ok := i.(map[interface{}]interface{}); ok {
		for k, v := range kv {
			if key, ok := k.(string); !ok {
				err = fmt.Errorf("unsupported key type: %T (%+v)", k, k)
			} else {
				err = fn(-1, key, v)
			}

			if err != nil {
				return
			}
		}
	} else if kv, ok := i.(map[string]interface{}); ok {
		for k, v := range kv {
			err = fn(-1, k, v)

			if err != nil {
				return
			}
		}
	} else if slice, ok := i.([]interface{}); ok {
		_ = slice
		for index, i := range slice {
			if err = fn(index, "", i); err != nil {
				return
			}
		}
	}

	return
}
