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
	storeDef struct {
		Package  string
		App      string
		Source   string
		Filename string

		Import []string `yaml:"import"`

		// List of all locations where we should export the store interface to
		Interface []string `yaml:"interface"`

		// Tries to autogenerate type by changing it to singular and prefixing it with *types.
		Types storeTypeDef `yaml:"types"`

		// All known fields that we need to store on a particular type
		//
		// For now, this set does not variate between different implementation
		// To support that, a (sub)set will need to be defined under each implementation (rdbms, mysql, mongo...)
		//
		Fields         storeTypeFieldSetDef      `yaml:"fields"`
		Lookups        []*storeTypeLookups       `yaml:"lookups"`
		PartialUpdates []*storeTypePartialUpdate `yaml:"partialUpdates"`
		RDBMS          *storeTypeRdbmsDef        `yaml:"rdbms"`

		Search storeTypeSearchDef `yaml:"search"`
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

		CustomRowScanner      bool `yaml:"customRowScanner"`
		CustomFilterConverter bool `yaml:"customFilterConverter"`
		CustomEncoder         bool `yaml:"customEncoder"`
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

		// If field is flagged as PK it is used in update & remove conditions
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
		Suffix      string            `yaml:"suffix"`
		Description string            `yaml:"description"`
		Fields      []string          `yaml:"fields"`
		Filter      map[string]string `yaml:"filter"`
		fields      storeTypeFieldSetDef
	}

	storeTypePartialUpdate struct {
		Name        string            `yaml:"name"`
		Description string            `yaml:"description"`
		Set         map[string]string `yaml:"set"`
		XX_Args     []string          `yaml:"args"`
		fields      storeTypeFieldSetDef
	}

	storeTypeSearchDef struct {
		Disable              bool `yaml:"disable"`
		DisablePaging        bool `yaml:"disablePaging"`
		DisableSorting       bool `yaml:"disableSorting"`
		DisableFilterCheckFn bool `yaml:"disableFilterCheckFunction"`
	}
)

var (
	outputDir string = "store"
)

func procStore() ([]*storeDef, error) {
	procDef := func(m string) (*storeDef, error) {
		def := &storeDef{Source: m}
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

		if def.Search.Disable {
			// No use for any of that if search is disabled...
			def.Search.DisablePaging = true
			def.Search.DisableSorting = true
			def.Search.DisableFilterCheckFn = true
		}

		// Always generate interface in store/tests and store/bulk
		def.Interface = append(def.Interface, "store/tests", "store/bulk")

		if def.Types.Base == "" {
			def.Types.Base = pubIdent(strings.Split(def.Filename, "_")...)
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
			def.Types.GoType = def.Types.Package + "." + pubIdent(def.Types.Singular)
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

		var hasPrimaryKey = false
		for _, f := range def.Fields {
			if f.IsPrimaryKey {
				hasPrimaryKey = true
				break
			}
		}

		for _, f := range def.Fields {
			if !hasPrimaryKey && f.Field == "ID" {
				f.IsPrimaryKey = true
				f.IsSortable = true
			}

			// copy alias from global spec so we can
			// generate aliased columsn
			f.alias = def.RDBMS.Alias

			if f.Column == "" {
				switch {
				case f.Field != "ID" && strings.HasSuffix(f.Field, "ID"):
					f.Column = "rel_" + cc2underscore(f.Field[:len(f.Field)-2])
				default:
					f.Column = cc2underscore(f.Field)
				}
			}

			switch {
			case f.Type != "":
				// type set
			case strings.HasSuffix(f.Field, "ID") || strings.HasSuffix(f.Field, "By"):
				f.Type = "uint64"
			case f.Field == "CreatedAt":
				f.Type = "time.Time"
			case strings.HasSuffix(f.Field, "At"):
				f.Type = "uint64"
			default:
				f.Type = "string"
			}
		}

		if len(def.PartialUpdates) > 0 && def.Fields.Find("ID") == nil {
			return nil, fmt.Errorf("partial updates without ID field are not supported")
		}

		// Checking if filters exist in the fields
		for i, p := range def.PartialUpdates {
			// Check and normalize set
			for f, v := range p.Set {
				if def.Fields.Find(f) == nil {
					return nil, fmt.Errorf("undefined field %q used in partialUpdate #%d set", f, i)
				}

				if v == "" {
					// Set empty strings to nil
					p.Set[f] = "nil"
				}
			}

			for _, a := range p.XX_Args {
				if def.Fields.Find(a) == nil {
					return nil, fmt.Errorf("undefined field %q used in partialUpdate #%d arguments", a, i)
				}
			}

			p.fields = def.Fields

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

	mm, err := filepath.Glob(filepath.Join(outputDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob: %w", err)
	}

	dd := []*storeDef{}
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
func genStore(tpl *template.Template, dd []*storeDef) (err error) {
	var (
		// general interfaces
		tplInterfacesJoined = tpl.Lookup("store_interfaces_joined.gen.go.tpl")
		tplInterfaces       = tpl.Lookup("store_interfaces.gen.go.tpl")

		// general tests
		tplTestAll = tpl.Lookup("store_test_all.gen.go.tpl")

		// @todo in-memory

		// rdbms specific
		tplRdbms = tpl.Lookup("store_rdbms.gen.go.tpl")

		// @todo redis
		// @todo mongodb
		// @todo elasticsearch

		// bulk specific
		tplBulk = tpl.Lookup("store_bulk.gen.go.tpl")

		dst             string
		joinedInterface = make(map[string][]*storeDef)
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

		dst = path.Join(outputDir, "bulk", d.Filename+".gen.go")
		if err = goTemplate(dst, tplBulk, d); err != nil {
			return
		}

		// Collect and map all interface output locations
		// and their corresponding definitions
		for _, dst = range d.Interface {
			if err = genStoreInterfaces(tplInterfaces, path.Join(dst, "store_interface_"+d.Filename+".gen.go"), path.Base(dst), d); err != nil {
				return
			}

			if joinedInterface[dst] == nil {
				joinedInterface[dst] = make([]*storeDef, 0, len(dd))
			}

			joinedInterface[dst] = append(joinedInterface[dst], d)
		}
	}

	// Add joined interfaces for each interface destination
	for dst, dd := range joinedInterface {
		if err = genStoreInterfacesJoined(tplInterfacesJoined, path.Join(dst, "store_interface.gen.go"), path.Base(dst), dd); err != nil {
			return
		}
	}

	//for _, d := range dd {
	//	d.Package = "tests"
	//	dst = path.Join("store/tests", "store_interface_"+d.Filename+".gen.go")
	//	if err = goTemplate(dst, tplInterfaces, d); err != nil {
	//		return
	//	}
	//}

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

func (s storeTypeFieldSetDef) Find(name string) *storeTypeFieldDef {
	for _, f := range s {
		if f.Field == name {
			return f
		}
	}

	return nil
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

func (p storeTypePartialUpdate) Args() []*storeTypeFieldDef {
	ff := make([]*storeTypeFieldDef, len(p.XX_Args))
	for a := range p.XX_Args {
		ff[a] = p.fields.Find(p.XX_Args[a])
	}

	return ff
}
