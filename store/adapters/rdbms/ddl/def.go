package ddl

//import (
//	"fmt"
//	"strings"
//)
//
//type (
//	tableManipulator  func(*Table)
//	indexManipulator  func(*Index)
//	columnManipulator func(*Column)
//
//	//Table struct {
//	//	Name       string
//	//	Columns    []*Column
//	//	Indexes    []*Index
//	//	PrimaryKey *Index
//	//}
//	//
//	//Index struct {
//	//	Name      string
//	//	Table     string
//	//	Fields    []*IField
//	//	Unique    bool
//	//	Condition string
//	//
//	//	// someday maybe index type
//	//}
//	//
//	//IField struct {
//	//	// Expression or a single column
//	//	Field string
//	//
//	//	Length int
//	//
//	//	// Wrap part in parentheses
//	//	Expr bool
//	//
//	//	// Ascending or descending
//	//	Desc bool
//	//}
//	//
//	//columnType uint
//	//ColumnType struct {
//	//	Type   columnType
//	//	Length int
//	//
//	//	// implementation variations
//	//	Flags map[string]interface{}
//	//}
//	//
//	//Column struct {
//	//	Name         string
//	//	Type         string
//	//	IsNull       bool
//	//	DefaultValue string
//	//	Comment      string
//	//}
//	//
//	//Columns []*Column
//)
//
////const (
////	// Subset of SQL types
////	//
////	// Just the ones we use in Corteza
////	// Each type should have its own formatter inside
////
////	ColumnTypeIdentifier columnType = iota
////	ColumnTypeVarchar
////	ColumnTypeText
////	ColumnTypeBinary
////	ColumnTypeTimestamp
////	ColumnTypeInteger
////	ColumnTypeJson
////	ColumnTypeBoolean
////)
////
//func TableDef(name string, mm ...tableManipulator) *Table {
//	var t = &Table{Ident: name}
//	return t.Apply(mm...)
//}
//
//func (t *Table) Apply(mm ...tableManipulator) *Table {
//	for _, m := range mm {
//		m(t)
//	}
//
//	return t
//}
//
//// Adds ID column and sets primary key
//func ID(t *Table) {
//	panic("obsolete")
//	t.Columns = append(t.Columns, &Column{
//		Name: "id",
//		// Type:    ColumnType{Type: ColumnTypeIdentifier},
//		Comment: fmt.Sprintf("Unique ID for %s", t.Name),
//	})
//
//	PrimaryKey(IColumn("id"))(t)
//}
//
//// Adds created_at/updated_at/deleted-at columns
//func CUDTimestamps(t *Table) {
//	panic("obsolete")
//	t.Apply(
//		ColumnDef("created_at", ColumnTypeTimestamp),
//		ColumnDef("updated_at", ColumnTypeTimestamp, Null),
//		ColumnDef("deleted_at", ColumnTypeTimestamp, Null),
//	)
//}
//
//// Adds created_at/updated_at/deleted-by columns
//func CUDUsers(t *Table) {
//	panic("obsolete")
//	t.Apply(
//		ColumnDef("created_by", ColumnTypeIdentifier),
//		ColumnDef("updated_by", ColumnTypeIdentifier, DefaultValue("0")),
//		ColumnDef("deleted_by", ColumnTypeIdentifier, DefaultValue("0")),
//	)
//}
//
//func ColumnDef(name string, cType columnType, mm ...columnManipulator) tableManipulator {
//	return func(t *Table) {
//		typ := ColumnType{
//			Type:  cType,
//			Flags: make(map[string]interface{}),
//
//			// Setting length to -1
//			//
//			// This way we'll know if it was actually set to 0
//			Length: -1,
//		}
//		col := &Column{Name: name, Type: typ}
//
//		for _, m := range mm {
//			m(col)
//		}
//
//		t.Columns = append(t.Columns, col)
//	}
//}
//
//// AddIndex adds index to table
//// If name has unique_ prefix it unique flag to true
//func AddIndex(name string, mm ...indexManipulator) tableManipulator {
//	return func(t *Table) {
//		i := &Index{
//			Name:   name,
//			Table:  t.Name,
//			Unique: strings.HasPrefix(name, "unique_"),
//		}
//
//		for _, m := range mm {
//			m(i)
//		}
//
//		t.Indexes = append(t.Indexes, i)
//	}
//}
//
//func Null(c *Column) {
//	c.IsNull = true
//}
//
//func ColumnTypeLength(l int) columnManipulator {
//	panic("obsolete")
//
//	return func(c *Column) {
//		//c.Type.Length = l
//	}
//}
//
//func ColumnTypeFlag(n string, v interface{}) columnManipulator {
//	panic("obsolete")
//	return func(c *Column) {
//		//c.Type.Flags[n] = v
//	}
//}
//
//// DefaultValue sets default value on column
////
//// Note: avoid using NOW(), CURRENT_TIMESTAMP or other (volatile) function
////       service layer should take care of this
//func DefaultValue(v string) columnManipulator {
//	return func(c *Column) {
//		c.DefaultValue = v
//	}
//}
//
//func PrimaryKey(mm ...indexManipulator) tableManipulator {
//	return func(t *Table) {
//		if t.PrimaryKey == nil {
//			t.PrimaryKey = &Index{}
//		}
//
//		for _, m := range mm {
//			m(t.PrimaryKey)
//		}
//	}
//}
//
//func Unique(i *Index) {
//	i.Unique = true
//}
//
//// IColumn adds one or more keys as columns
//func IColumn(cc ...string) indexManipulator {
//	return func(i *Index) {
//		for _, c := range cc {
//			i.Fields = append(i.Fields, &IField{Field: c})
//		}
//	}
//}
//
//// IColumn adds one or more keys as columns
//func IFieldFull(ff ...*IField) indexManipulator {
//	return func(i *Index) {
//		for _, f := range ff {
//			i.Fields = append(i.Fields, f)
//		}
//	}
//}
//
//// IExpr adds one or more keys as expressions
//func IExpr(ee ...string) indexManipulator {
//	return func(i *Index) {
//		for _, e := range ee {
//			i.Fields = append(i.Fields, &IField{Field: e, Expr: true})
//		}
//	}
//}
//
//func IWhere(cnd string) indexManipulator {
//	return func(i *Index) {
//		i.Condition = cnd
//	}
//}
//
//func (cc Columns) Get(name string) *Column {
//	for c := range cc {
//		if cc[c].Name == name {
//			return cc[c]
//		}
//	}
//
//	return nil
//}
