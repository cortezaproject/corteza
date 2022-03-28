package schema

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type (
	CommonDialect struct {
		tpl *template.Template
	}
)

const (
	// table creation
	commonCreateTable = `
	CREATE TABLE {{template "if-not-exists-clause" .}} {{ .Name }} (
	{{ range $n, $c := .Columns -}}
	{{ if $n }}, {{ else }}  {{ end }}{{ template "create-table-column" . }}
	{{ end -}}
	{{- if .PrimaryKey }}
	, PRIMARY KEY {{ template "index-fields" .PrimaryKey.Fields }}
	{{- end }}
	) {{- template "create-table-suffix" . -}}
	`

	commonCreateTableSuffix = ``

	commonCreateTableColumn = `
	{{- .Name }} {{ columnType .Type }}
	{{- if not .IsNull }} NOT NULL{{end}}
	{{- if .DefaultValue }} DEFAULT {{ .DefaultValue }} {{end}}
	`

	commonAddColumn     = `ALTER TABLE {{ .Table }} ADD {{ template "create-table-column" .Column }}`
	commonAddPrimaryKey = `ALTER TABLE {{ .Table }} ADD CONSTRAINT PRIMARY KEY {{ template "index-fields" .PrimaryKey.Fields }}`
	commonDropColumn    = `ALTER TABLE {{ .Table }} DROP {{ .Column }}`
	commonRenameColumn  = `ALTER TABLE {{ .Table }} RENAME COLUMN {{ .OldName }} TO {{ .NewName }}`

	// index creation
	commonCreateIndex = `CREATE {{ if .Unique }}UNIQUE {{ end }}INDEX {{ template "if-not-exists-clause" . }} {{ template "index-name" . }} ON {{ .Table }} {{ template "index-fields" .Fields }}{{ template "index-condition" . }}`

	commonIndexName      = `{{ .Table }}_{{ .Name }}`
	commonIndexCondition = `{{- if .Condition }} WHERE ({{ .Condition }}){{ end }}`

	// index creation
	commonIndexFields = `
	({{ range $n, $f := . -}}
		{{ if $n }}, {{ end }}
		{{- if .Expr}}({{ end }}
		{{- .Field }}
		{{- if .Expr}}){{ end }}
		{{- if .Desc }} DESC{{ end }}
	{{- end }})
	`
	// table/index exist or not clause
	commonIfNotExistsClause = `IF NOT EXISTS`
)

func NewCommonDialect() *CommonDialect {
	g := &CommonDialect{tpl: template.New("")}
	g.AddTemplateFunc("columnType", GenColumnType)
	g.AddTemplate("create-table", commonCreateTable)
	g.AddTemplate("create-table-suffix", commonCreateTableSuffix)
	g.AddTemplate("create-table-column", commonCreateTableColumn)
	g.AddTemplate("add-column", commonAddColumn)
	g.AddTemplate("add-primary-key", commonAddPrimaryKey)
	g.AddTemplate("drop-column", commonDropColumn)
	g.AddTemplate("rename-column", commonRenameColumn)
	g.AddTemplate("create-index", commonCreateIndex)
	g.AddTemplate("index-condition", commonIndexCondition)
	g.AddTemplate("index-name", commonIndexName)
	g.AddTemplate("index-fields", commonIndexFields)
	g.AddTemplate("if-not-exists-clause", commonIfNotExistsClause)

	return g
}

func (g *CommonDialect) executeTemplate(name string, data interface{}) string {
	buf := &bytes.Buffer{}
	if err := g.tpl.ExecuteTemplate(buf, name, data); err != nil {
		panic(err)
	}

	return buf.String()
}

func (g *CommonDialect) AddTemplateFunc(name string, fn interface{}) {
	g.tpl.Funcs(template.FuncMap{name: fn})
}

func (g *CommonDialect) AddTemplate(name, tpl string) {
	funcMap := template.FuncMap{
		"trimExpression": func(s string) string {
			re := regexp.MustCompile(`^.*\((\w+)\).*$`)
			if str := re.FindAllStringSubmatch(s, 1); len(str) > 0 && len(str[0]) > 0 && len(str[0][1]) > 0 {
				return str[0][1]
			}
			return s
		},
	}
	template.Must(g.tpl.Funcs(funcMap).New(name).Parse(strings.TrimSpace(tpl)))
}

func (g *CommonDialect) CreateTable(t *Table) string {
	return g.executeTemplate("create-table", t)
}

func (g *CommonDialect) CreateIndex(i *Index) string {
	return g.executeTemplate("create-index", i)
}

func (g *CommonDialect) AddColumn(table string, c *Column) string {
	return g.executeTemplate("add-column", map[string]interface{}{
		"Table":  table,
		"Column": c,
	})
}

func (g *CommonDialect) DropColumn(table, column string) string {
	return g.executeTemplate("drop-column", map[string]interface{}{
		"Table":  table,
		"Column": column,
	})
}

func (g *CommonDialect) RenameColumn(table, oldName, newName string) string {
	return g.executeTemplate("rename-column", map[string]interface{}{
		"Table":   table,
		"OldName": oldName,
		"NewName": newName,
	})
}

func (g *CommonDialect) AddPrimaryKey(table string, pk *Index) string {
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
