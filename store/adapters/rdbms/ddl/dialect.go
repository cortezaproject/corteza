package ddl

//import (
//	"bytes"
//	"fmt"
//	"regexp"
//	"strings"
//	"text/template"
//)
//
//type (
//	CommonDialect struct {
//		tpl *template.Template
//
//		CreateTable       string
//		CreateTableSuffix string
//		CreateTableColumn string
//		AddColumn         string
//		AddPrimaryKey     string
//		DropColumn        string
//		RenameColumn      string
//		CreateIndex       string
//		IndexName         string
//		IndexCondition    string
//		IndexFields       string
//		IfNotExistsClause string
//	}
//)
//
//const (
//	CreateTable       = "create-table"
//	CreateTableSuffix = "create-table-suffix"
//	CreateTableColumn = "create-table-column"
//	AddColumn         = "add-column"
//	AddPrimaryKey     = "add-primary-key"
//	DropColumn        = "drop-column"
//	RenameColumn      = "rename-column"
//	CreateIndex       = "create-index"
//	IndexCondition    = "index-condition"
//	IndexName         = "index-name"
//	IndexFields       = "index-fields"
//	IfNotExistsClause = "if-not-exists-clause"
//	TableExists       = "table-exists"
//
//	CreateTableDDL = `
//	CREATE TABLE {{template "if-not-exists-clause" .}} {{ .Name }} (
//	{{ range $n, $c := .Columns -}}
//	{{ if $n }}, {{ else }}  {{ end }}{{ template "create-table-column" . }}
//	{{ end -}}
//	{{- if .PrimaryKey }}
//	, PRIMARY KEY {{ template "index-fields" .PrimaryKey.Fields }}
//	{{- end }}
//	) {{- template "create-table-suffix" . -}}
//	`
//	IfNotExistsClauseDDL = `IF NOT EXISTS`
//
//	CreateTableSuffixDDL = ``
//
//	CreateTableColumnDDL = `
//	{{- .Name }} {{ columnType .Type }}
//	{{- if not .IsNull }} NOT NULL{{end}}
//	{{- if .DefaultValue }} DEFAULT {{ .DefaultValue }} {{end}}
//	`
//
//	//AddColumnDDL     = `ALTER TABLE {{ .Table }} ADD {{ template "create-table-column" .Column }}`
//	//AddPrimaryKeyDDL = `ALTER TABLE {{ .Table }} ADD CONSTRAINT PRIMARY KEY {{ template "index-fields" .PrimaryKey.Fields }}`
//	//DropColumnDDL    = `ALTER TABLE {{ .Table }} DROP {{ .Column }}`
//	//RenameColumnDDL  = `ALTER TABLE {{ .Table }} RENAME COLUMN {{ .OldName }} TO {{ .NewName }}`
//
//	//CreateIndexDDL = `CREATE {{ if .Unique }}UNIQUE {{ end }}INDEX {{ template "if-not-exists-clause" . }} {{ template "index-name" . }} ON {{ .Table }} {{ template "index-fields" .Fields }}{{ template "index-condition" . }}`
//
//	//IndexNameDDL      = `{{ .Table }}_{{ .Name }}`
//	//IndexConditionDDL = `{{- if .Condition }} WHERE ({{ .Condition }}){{ end }}`
//
//	IndexFieldsDDL = `
//	({{ range $n, $f := . -}}
//		{{ if $n }}, {{ end }}
//		{{- if .Expr}}({{ end }}
//		{{- .Field }}
//		{{- if .Expr}}){{ end }}
//		{{- if .Desc }} DESC{{ end }}
//	{{- end }})
//	`
//
//	TableExistDDL = `
//	SELECT COUNT(*) > 0"
//	  FROM INFORMATION_SCHEMA.TABLES
//     WHERE TABLE_NAME = '{{ .Table }}'
//       AND TABLE_SCHEMA = '{{ .Schema }}'
//	`
//)
//
//func NewCommonDialect() *CommonDialect {
//	g := &CommonDialect{tpl: template.New("")}
//	g.AddTemplateFunc("columnType", GenColumnType)
//	g.AddTemplate(CreateTable, CreateTableDDL)
//	g.AddTemplate(CreateTableSuffix, CreateTableSuffixDDL)
//	g.AddTemplate(CreateTableColumn, CreateTableColumnDDL)
//	//g.AddTemplate(AddColumn, AddColumnDDL)
//	//g.AddTemplate(AddPrimaryKey, AddPrimaryKeyDDL)
//	//g.AddTemplate(DropColumn, DropColumnDDL)
//	//g.AddTemplate(RenameColumn, RenameColumnDDL)
//	//g.AddTemplate(CreateIndex, CreateIndexDDL)
//	//g.AddTemplate(IndexCondition, IndexConditionDDL)
//	//g.AddTemplate(IndexName, IndexNameDDL)
//	g.AddTemplate(IndexFields, IndexFieldsDDL)
//	g.AddTemplate(IfNotExistsClause, IfNotExistsClauseDDL)
//	g.AddTemplate(TableExists, TableExistDDL)
//
//	return g
//}
//
//func (g *CommonDialect) GenerateSQL(name string, data interface{}) string {
//	buf := &bytes.Buffer{}
//	if err := g.tpl.ExecuteTemplate(buf, name, data); err != nil {
//		panic(err)
//	}
//
//	return buf.String()
//}
//
//func (g *CommonDialect) AddTemplateFunc(name string, fn interface{}) {
//	g.tpl.Funcs(template.FuncMap{name: fn})
//}
//
//func (g *CommonDialect) AddTemplate(name, tpl string) {
//	funcMap := template.FuncMap{
//		"trimExpression": func(s string) string {
//			re := regexp.MustCompile(`^.*\((\w+)\).*$`)
//			if str := re.FindAllStringSubmatch(s, 1); len(str) > 0 && len(str[0]) > 0 && len(str[0][1]) > 0 {
//				return str[0][1]
//			}
//			return s
//		},
//	}
//	template.Must(g.tpl.Funcs(funcMap).New(name).Parse(strings.TrimSpace(tpl)))
//}
//
////func (g *CommonDialect) CreateTable(t *Table) string {
////	return g.GenerateSQL("create-table", t)
////}
////
////func (g *CommonDialect) CreateIndex(i *Index) string {
////	return g.GenerateSQL("create-index", i)
////}
//
////func (g *CommonDialect) AddColumn(table string, c *Column) string {
////	return g.GenerateSQL("add-column", map[string]interface{}{
////		"Table":  table,
////		"Column": c,
////	})
////}
////
////func (g *CommonDialect) DropColumn(table, column string) string {
////	return g.GenerateSQL("drop-column", map[string]interface{}{
////		"Table":  table,
////		"Column": column,
////	})
////}
////
////func (g *CommonDialect) RenameColumn(table, oldName, newName string) string {
////	return g.GenerateSQL("rename-column", map[string]interface{}{
////		"Table":   table,
////		"OldName": oldName,
////		"NewName": newName,
////	})
////}
////
////func (g *CommonDialect) AddPrimaryKey(table string, pk *Index) string {
////	return g.GenerateSQL("add-primary-key", map[string]interface{}{
////		"Table":      table,
////		"PrimaryKey": pk,
////	})
////}
//
//func GenColumnType(ct *ColumnType) string {
//	switch ct.Type {
//	case ColumnTypeIdentifier:
//		return "BIGINT"
//	case ColumnTypeVarchar:
//		if ct.Length > 0 {
//			// VARCHAR(0) is useless
//			return fmt.Sprintf("VARCHAR(%d)", ct.Length)
//		}
//		return "VARCHAR"
//	case ColumnTypeText:
//		return "TEXT"
//	case ColumnTypeJson:
//		return "JSON"
//	case ColumnTypeBinary:
//		return "BYTEA"
//	case ColumnTypeTimestamp:
//		if ct.Length > -1 {
//			// TIMESTAMPTZ(0) strips out milliseconds
//			return fmt.Sprintf("TIMESTAMPTZ(%d)", ct.Length)
//		}
//
//		return "TIMESTAMPTZ"
//	case ColumnTypeInteger:
//		return "INTEGER"
//	case ColumnTypeBoolean:
//		return "BOOLEAN"
//	default:
//		panic(fmt.Sprintf("unhandled column type: %d ", ct.Type))
//	}
//}
