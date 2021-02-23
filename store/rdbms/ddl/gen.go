package ddl

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"text/template"
)

type (
	Generator struct {
		log *zap.Logger
		tpl *template.Template
	}
)

// @todo all DDL operations (altering, droping, adding...) from all *sql* implementations should be moved here
//       logic is more or less the same with a differet approach on how to read schema specs (tables, columns)

const (
	// table creation
	genericCreateTable = `
CREATE TABLE {{ .Name }} (
{{ range $n, $c := .Columns -}}
{{ if $n }}, {{ else }}  {{ end }}{{ template "create-table-column" . }}
{{ end -}}
{{- if .PrimaryKey }}
, PRIMARY KEY {{ template "index-fields" .PrimaryKey.Fields }}
{{- end }}
) {{- template "create-table-suffix" . -}}
`
	genericCreateTableSuffix = ``

	genericCreateTableColumn = `
	{{- .Name }} {{ columnType .Type }}
	{{- if not .IsNull }} NOT NULL{{end}}
	{{- if .DefaultValue }} DEFAULT {{ .DefaultValue }} {{end}}
`

	genericAddColumn     = `ALTER TABLE {{ .Table }} ADD {{ template "create-table-column" .Column }}`
	genericAddPrimaryKey = `ALTER TABLE {{ .Table }} ADD CONSTRAINT PRIMARY KEY {{ template "index-fields" .PrimaryKey.Fields }}`
	genericDropColumn    = `ALTER TABLE {{ .Table }} DROP {{ .Column }}`
	genericRenameColumn  = `ALTER TABLE {{ .Table }} RENAME COLUMN {{ .OldName }} TO {{ .NewName }}`

	// index creation
	genericCreateIndex = `CREATE {{ if .Unique }}UNIQUE {{ end }}INDEX {{ template "index-name" . }} ON {{ .Table }} {{ template "index-fields" .Fields }}{{ template "index-condition" . }}`

	genericIndexName      = `{{ .Table }}_{{ .Name }}`
	genericIndexCondition = `{{- if .Condition }} WHERE ({{ .Condition }}){{ end }}`

	// index creation
	genericIndexFields = `
({{ range $n, $f := . -}}
	{{ if $n }}, {{ end }}
	{{- if .Expr}}({{ end }}
	{{- .Field }}
	{{- if .Expr}}){{ end }}
	{{- if .Desc }} DESC{{ end }}
{{- end }})
`
)

func NewGenerator(log *zap.Logger) *Generator {
	g := &Generator{log: log, tpl: template.New("")}
	g.AddTemplateFunc("columnType", GenColumnType)
	g.AddTemplate("create-table", genericCreateTable)
	g.AddTemplate("create-table-suffix", genericCreateTableSuffix)
	g.AddTemplate("create-table-column", genericCreateTableColumn)
	g.AddTemplate("add-column", genericAddColumn)
	g.AddTemplate("add-primary-key", genericAddPrimaryKey)
	g.AddTemplate("drop-column", genericDropColumn)
	g.AddTemplate("rename-column", genericRenameColumn)
	g.AddTemplate("create-index", genericCreateIndex)
	g.AddTemplate("index-condition", genericIndexCondition)
	g.AddTemplate("index-name", genericIndexName)
	g.AddTemplate("index-fields", genericIndexFields)

	return g
}

func (g *Generator) executeTemplate(name string, data interface{}) string {
	buf := &bytes.Buffer{}
	if err := g.tpl.ExecuteTemplate(buf, name, data); err != nil {
		panic(err)
	}

	return buf.String()
}

func (g *Generator) AddTemplateFunc(name string, fn interface{}) {
	g.tpl.Funcs(template.FuncMap{name: fn})
}

func (g *Generator) AddTemplate(name, tpl string) {
	template.Must(g.tpl.New(name).Parse(strings.TrimSpace(tpl)))
}

func (g *Generator) CreateTable(t *Table) string {
	return g.executeTemplate("create-table", t)
}

func (g *Generator) CreateIndex(i *Index) string {
	return g.executeTemplate("create-index", i)
}

func (g *Generator) AddColumn(table string, c *Column) string {
	return g.executeTemplate("add-column", map[string]interface{}{
		"Table":  table,
		"Column": c,
	})
}

func (g *Generator) DropColumn(table, column string) string {
	return g.executeTemplate("drop-column", map[string]interface{}{
		"Table":  table,
		"Column": column,
	})
}

func (g *Generator) RenameColumn(table, oldName, newName string) string {
	return g.executeTemplate("rename-column", map[string]interface{}{
		"Table":   table,
		"OldName": oldName,
		"NewName": newName,
	})
}

func (g *Generator) AddPrimaryKey(table string, pk *Index) string {
	return g.executeTemplate("add-primary-key", map[string]interface{}{
		"Table":      table,
		"PrimaryKey": pk,
	})
}

func GenColumnType(ct *ColumnType) string {
	switch ct.Type {
	case ColumnTypeIdentifier:
		return "BIGINT"
	case ColumnTypeVarchar:
		if ct.Length > 0 {
			// VARCHAR(0) is useless
			return fmt.Sprintf("VARCHAR(%d)", ct.Length)
		}
		return "VARCHAR"
	case ColumnTypeText:
		return "TEXT"
	case ColumnTypeJson:
		return "JSON"
	case ColumnTypeBinary:
		return "BYTEA"
	case ColumnTypeTimestamp:
		if ct.Length > -1 {
			// TIMESTAMPTZ(0) strips out milliseconds
			return fmt.Sprintf("TIMESTAMPTZ(%d)", ct.Length)
		}

		return "TIMESTAMPTZ"
	case ColumnTypeInteger:
		return "INTEGER"
	case ColumnTypeBoolean:
		return "BOOLEAN"
	default:
		panic(fmt.Sprintf("unhandled column type: %d ", ct.Type))
	}
}
