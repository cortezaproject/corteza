package codegen

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/slice"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

type (
	// definitions are in one file
	storeDef struct {
		Package  string
		App      string
		Source   string
		Filename string

		Import []string `yaml:"import"`

		// Tries to autogenerate type by changing it to singular and prefixing it with *types.
		Types storeTypeDef `yaml:"types"`

		// All known fields that we need to store on a particular type
		//
		// For now, this set does not variate between different implementation
		// To support that, a (sub)set will need to be defined under each implementation (rdbms, mysql, mongo...)
		//
		Fields    storeTypeFieldSetDef         `yaml:"fields"`
		RDBMS     storeTypeRdbmsDef            `yaml:"rdbms"`
		Functions []*storeTypeFunctionsDef     `yaml:"functions"`
		Arguments []*storeTypeExtraArgumentDef `yaml:"arguments"`

		Search   storeTypeSearchDef   `yaml:"search"`
		Lookups  []*storeTypeLookups  `yaml:"lookups"`
		Create   storeTypeCreateDef   `yaml:"create"`
		Update   storeTypeUpdateDef   `yaml:"update"`
		Upsert   storeTypeUpsertDef   `yaml:"upsert"`
		Delete   storeTypeDeleteDef   `yaml:"Delete"`
		Truncate storeTypeTruncateDef `yaml:"truncate"`

		// Make interfaces and store functions
		Publish bool `yaml:"publish"`
	}

	storeTypeDef struct {
		// Name of the package where type can be found
		// (defaults to types)
		Package string `yaml:"package"`

		// Name of the base type
		// (defaults to base name of the yaml file)
		Base string `yaml:"base"`

		// Singular variation of name
		// (defaults to <Base> (s trimmed))
		Singular string `yaml:"singular"`

		// Plural variantion of name
		// (defaults to <Singular> (s appended))
		Plural string `yaml:"plural"`

		// Name of the set go type
		// (defaults to <Package>.<Singular>)
		GoType string `yaml:"type"`

		// Name of the set go type
		// (defaults to <GoType>Set)
		GoSetType string `yaml:"setType"`

		// Name of the filter go type
		// (defaults to <GoType>Filter)
		GoFilterType string `yaml:"filterType"`
	}

	storeTypeRdbmsDef struct {
		// Alias used in SQL queries
		Alias string `yaml:"alias,omitempty"`
		Table string `yaml:"table,omitempty"`

		CustomRowScanner        bool `yaml:"customRowScanner"`
		CustomFilterConverter   bool `yaml:"customFilterConverter"`
		CustomSortConverter     bool `yaml:"customSortConverter"`
		CustomCursorCollector   bool `yaml:"customCursorCollector"`
		CustomPostLoadProcessor bool `yaml:"customPostLoadProcessor"`
		CustomEncoder           bool `yaml:"customEncoder"`
	}

	storeTypeFunctionsDef struct {
		Name      string                      `yaml:"name"`
		Arguments []storeTypeExtraArgumentDef `yaml:"arguments"`
		Return    []string                    `yaml:"return"`
	}

	storeTypeExtraArgumentDef struct {
		Name string
		Type string
	}

	storeTypeFieldSetDef []*storeTypeFieldDef

	storeTypeFieldDef struct {
		Field string `yaml:"field"`

		// Autodiscovery logic (when not explicitly set)
		//   uint64: 		when field has "ID" suffix
		//   time.Time: 	when field equals with "created_at"
		//   *time.Time: 	when field ends with "_at"
		//   string: 		default
		Type string `yaml:"type"`

		// When not explicitly set, defaults to snake-cased value from field
		//
		// Exceptions:
		//  If field name ends with ID (<base>ID), it converts that to rel_<snake-cased-base>
		Column string `yaml:"column"`

		// If field is flagged as PK it is used in update & Delete conditions
		// Note: if no other field is set as primary and field with ID name
		//       exists, that field is auto-set as primary.
		IsPrimaryKey bool `yaml:"isPrimaryKey"`

		// FilterPreprocess sets preprocessing function used on
		// conditions for lookup functions
		//
		// See specific implementation for details
		LookupFilterPreprocess string `yaml:"lookupFilterPreprocessor"`

		// Is field sortable?
		IsSortable bool `yaml:"sortable"`

		// When sorting is disabled and paging enabled and we need to have fixed default
		// sorting set and sometimes default sorting needs to be in descending order
		SortDescending bool `yaml:"sortDescending"`

		// Is field unique
		IsUnique bool `yaml:"unique"`

		// @todo implementation
		FullTextSearch bool `yaml:"fts"`

		alias string
	}

	storeTypeLookups struct {
		// LookupBy<suffix>
		// When not explicitly defined, it names of all fields
		Export                bool              `yaml:"export"`
		Suffix                string            `yaml:"suffix"`
		Description           string            `yaml:"description"`
		UniqueConstraintCheck bool              `yaml:"uniqueConstraintCheck"`
		Fields                []string          `yaml:"fields"`
		Filter                map[string]string `yaml:"filter"`
		fields                storeTypeFieldSetDef
	}

	storeTypeSearchDef struct {
		Enable              bool `yaml:"enable"`
		Export              bool `yaml:"export"`
		EnablePaging        bool `yaml:"enablePaging"`
		EnableSorting       bool `yaml:"enableSorting"`
		EnableFilterCheckFn bool `yaml:"enableFilterCheckFunction"`
	}

	storeTypeCreateDef struct {
		Enable bool `yaml:"enable"`
		Export bool `yaml:"export"`
	}

	storeTypeUpdateDef struct {
		Enable bool `yaml:"enable"`
		Export bool `yaml:"export"`
	}

	storeTypeUpsertDef struct {
		Enable bool `yaml:"enable"`
		Export bool `yaml:"export"`
	}

	storeTypeDeleteDef struct {
		Enable bool `yaml:"enable"`
		Export bool `yaml:"export"`
	}

	storeTypeTruncateDef struct {
		Export bool `yaml:"export"`
	}
)

var (
	outputDir  string = "store"
	spaceSplit        = regexp.MustCompile(`\s+`)
)

func procStore(mm ...string) ([]*storeDef, error) {
	procDef := func(m string) (*storeDef, error) {
		// initialize & set default
		def := &storeDef{
			Source: m,

			RDBMS: storeTypeRdbmsDef{
				CustomRowScanner:      false,
				CustomFilterConverter: false,
				CustomEncoder:         false,
			},

			Search: storeTypeSearchDef{
				Enable:              true,
				Export:              true,
				EnablePaging:        true,
				EnableSorting:       true,
				EnableFilterCheckFn: true,
			},

			Create: storeTypeCreateDef{
				Enable: true,
				Export: true,
			},

			Update: storeTypeUpdateDef{
				Enable: true,
				Export: true,
			},

			Upsert: storeTypeUpsertDef{
				Enable: true,
				Export: true,
			},

			Delete: storeTypeDeleteDef{
				Enable: true,
				Export: true,
			},

			Truncate: storeTypeTruncateDef{
				Export: true,
			},

			Publish: true,
		}
		f, err := os.Open(m)
		if err != nil {
			return nil, fmt.Errorf("%s read failed: %w", m, err)
		}

		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(&def); err != nil {
			return nil, err
		}

		def.Filename = path.Base(m)
		def.Filename = def.Filename[:len(def.Filename)-5]

		if !def.Search.Enable {
			// No use for any of that if search is disabled...
			def.Search.EnablePaging = false
			def.Search.EnableSorting = false
			def.Search.EnableFilterCheckFn = false
		}

		if !def.Create.Enable {
			// No use for any of that if operation is disabled...
			def.Create.Export = false
		}

		if !def.Update.Enable {
			// No use for any of that if operation is disabled...
			def.Update.Export = false
		}

		if !def.Upsert.Enable {
			// No use for any of that if operation is disabled...
			def.Upsert.Export = false
		}

		if !def.Delete.Enable {
			// No use for any of that if operation is disabled...
			def.Delete.Export = false
		}

		if def.Types.Base == "" {
			def.Types.Base = export(strings.Split(def.Filename, "_")...)
		}

		if def.Types.Singular == "" {
			def.Types.Singular = strings.TrimRight(def.Types.Base, "s")
		}
		if def.Types.Plural == "" {
			def.Types.Plural = def.Types.Singular + "s"
		}

		if def.Types.Package == "" {
			def.Types.Package = "types"
		}

		if def.Types.GoType == "" {
			def.Types.GoType = def.Types.Package + "." + export(def.Types.Singular)
		}

		if def.Types.GoSetType == "" {
			def.Types.GoSetType = def.Types.GoType + "Set"
		}

		if def.Types.GoFilterType == "" {
			def.Types.GoFilterType = def.Types.GoType + "Filter"
		}

		if def.RDBMS.Alias == "" {
			def.RDBMS.Alias = def.Types.Base[0:1]
		}

		var hasPrimaryKey = def.Fields.HasPrimaryKey()
		for _, field := range def.Fields {
			if !hasPrimaryKey && field.Field == "ID" {
				field.IsPrimaryKey = true
				field.IsSortable = true
			}

			// copy alias from global spec so we can
			// generate aliased columsn
			field.alias = def.RDBMS.Alias

			if field.Column == "" {
				switch {
				case field.Field != "ID" && strings.HasSuffix(field.Field, "ID"):
					field.Column = "rel_" + cc2underscore(field.Field[:len(field.Field)-2])
				default:
					field.Column = cc2underscore(field.Field)
				}
			}

			switch {
			case field.Type != "":
				// type set
			case strings.HasSuffix(field.Field, "ID") || strings.HasSuffix(field.Field, "By"):
				field.Type = "uint64"
			case field.Field == "CreatedAt":
				field.Type = "time.Time"
			case strings.HasSuffix(field.Field, "At"):
				field.Type = "uint64"
			default:
				field.Type = "string"
			}
		}

		for i, l := range def.Lookups {
			if len(l.Fields) == 0 {
				return nil, fmt.Errorf("define at least one lookup field in lookup #%d", i)
			}

			// Checking if fields exist in the fields
			for _, f := range l.Fields {
				if def.Fields.Find(f) == nil {
					return nil, fmt.Errorf("undefined lookup field %q used", f)

				}
			}

			// Checking if filters exist in the fields
			for f, v := range l.Filter {
				if def.Fields.Find(f) == nil {
					return nil, fmt.Errorf("undefined lookup filter %q used", f)
				}

				if v == "" {
					// Set empty strings to nil
					l.Filter[f] = "nil"
				}

			}

			if l.Suffix == "" {
				l.Suffix = strings.Join(l.Fields, "")
			}

			l.fields = def.Fields
		}

		return def, nil
	}

	dd := make([]*storeDef, 0, len(mm))
	for _, m := range mm {
		def, err := procDef(m)
		if err != nil {
			return nil, fmt.Errorf("failed to process %s: %w", m, err)
		}

		dd = append(dd, def)
	}

	return dd, nil
}

// genStore generates all store related code, functions, interfaces...
//
// Templates can be found under assets/store*.tpl
func genStore(tpl *template.Template, dd ...*storeDef) (err error) {
	var (
		// general interfaces
		tplInterfacesJoined = tpl.Lookup("store_interfaces_joined.gen.go.tpl")
		tplBase             = tpl.Lookup("store_base.gen.go.tpl")

		// general tests
		tplTestAll = tpl.Lookup("store_test_all.gen.go.tpl")

		// @todo in-memory

		// rdbms specific
		tplRdbms = tpl.Lookup("store_rdbms.gen.go.tpl")

		// @todo redis
		// @todo mongodb
		// @todo elasticsearch

		dst string
	)

	// Output all test setup into a single file
	dst = path.Join(outputDir, "tests", "gen_test.go")
	if err = goTemplate(dst, tplTestAll, dd); err != nil {
		return
	}

	// Multi-file output
	for _, d := range dd {
		dst = path.Join(outputDir, "rdbms", d.Filename+".gen.go")
		if err = goTemplate(dst, tplRdbms, d); err != nil {
			return
		}

		dst = path.Join(outputDir, d.Filename+".gen.go")
		if err = goTemplate(dst, tplBase, d); err != nil {
			return
		}
	}

	if err = genStoreInterfacesJoined(tplInterfacesJoined, path.Join("store", "interfaces.gen.go"), path.Base(dst), dd); err != nil {
		return
	}

	return nil
}

func genStoreInterfaces(tpl *template.Template, dst, pkg string, d *storeDef) error {
	d.Package = pkg
	return goTemplate(dst, tpl, d)
}

func genStoreInterfacesJoined(tpl *template.Template, dst, pkg string, dd []*storeDef) error {
	payload := map[string]interface{}{
		"Package":     pkg,
		"Definitions": dd,
		"Import":      collectStoreDefImports("", dd...),
	}

	return goTemplate(dst, tpl, payload)
}

func collectStoreDefImports(basePkg string, dd ...*storeDef) []string {
	ii := make([]string, 0, len(dd))
	for _, d := range dd {
		for _, i := range d.Import {
			if !slice.HasString(ii, i) && (basePkg == "" || !strings.HasSuffix(i, basePkg)) {
				ii = append(ii, i)
			}
		}
	}

	return ii
}

// Exported returns true if at least one of the functions is exported
func (d storeDef) Exported() bool {
	return d.Search.Export ||
		d.Create.Export ||
		d.Update.Export ||
		d.Upsert.Export ||
		d.Delete.Export
}

func (ff storeTypeFieldSetDef) Find(name string) *storeTypeFieldDef {
	for _, f := range ff {
		if f.Field == name {
			return f
		}
	}

	return nil
}

func (ff storeTypeFieldSetDef) HasPrimaryKey() bool {
	for _, f := range ff {
		if f.IsPrimaryKey {
			return true
		}
	}

	return false
}

func (ff storeTypeFieldSetDef) PrimaryKeyFields() storeTypeFieldSetDef {
	pkSet := storeTypeFieldSetDef{}
	for _, f := range ff {
		if f.IsPrimaryKey {
			pkSet = append(pkSet, f)
		}
	}

	return pkSet
}

func (f storeTypeFieldDef) Arg() string {
	if f.Field == "ID" {
		return f.Field
	}

	return strings.ToLower(f.Field[:1]) + f.Field[1:]
}

func (f storeTypeFieldDef) AliasedColumn() string {
	return fmt.Sprintf("%s.%s", f.alias, f.Column)
}

// UnmarshalYAML makes sure that export flag is set to true when not explicity disabled
func (d *storeTypeLookups) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type dAux storeTypeLookups
	var aux = (*dAux)(d)
	aux.Export = true
	return unmarshal(aux)
}
