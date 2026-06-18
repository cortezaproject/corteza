package dal

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/crusttech/human/server/pkg/dal"
	"github.com/crusttech/human/server/pkg/logger"
	"github.com/crusttech/human/server/store/adapters/rdbms/ddl"
)

// tableToModel converts a ddl.Table (from information_schema) to a dal.Model.
// connectionID is set by the caller (SearchExternalModels in pkg/dal/service.go).
func tableToModel(connectionID uint64, t *ddl.Table) (*dal.Model, error) {
	pk := primaryKeyColumns(t)
	m := &dal.Model{
		ConnectionID: connectionID,
		Ident:        t.Ident,
		Label:        t.Ident,
	}

	for _, col := range t.Columns {
		typ := sqlTypeToDalType(col.Type)
		m.Attributes = append(m.Attributes, &dal.Attribute{
			Ident:      col.Ident,
			Label:      col.Ident,
			PrimaryKey: pk[col.Ident],
			Sortable:   pk[col.Ident],
			Filterable: pk[col.Ident],
			Type:       typ,
			Store:      &dal.CodecPlain{},
		})
	}

	return m, nil
}

// primaryKeyColumns returns a set of column idents that form the PK.
// Falls back to any column named "id" if no index info is available.
func primaryKeyColumns(t *ddl.Table) map[string]bool {
	pk := map[string]bool{}

	for _, idx := range t.Indexes {
		if strings.EqualFold(idx.Ident, "PRIMARY") || idx.Unique {
			for _, f := range idx.Fields {
				if f.Column != "" {
					pk[f.Column] = true
				}
			}
			if len(pk) > 0 {
				return pk
			}
		}
	}

	// fallback: column named "id"
	for _, col := range t.Columns {
		if strings.EqualFold(col.Ident, "id") {
			pk[col.Ident] = true
			return pk
		}
	}

	return pk
}

// sqlTypeToDalType maps a SQL column type (as reported by information_schema)
// to a dal.Type. Unknown types fall back to TypeText with a warning.
func sqlTypeToDalType(ct *ddl.ColumnType) dal.Type {
	if ct == nil {
		return &dal.TypeText{}
	}

	// strip (precision, scale) or (length) suffix, lowercase
	raw := strings.ToLower(ct.Name)
	base := raw
	if idx := strings.IndexByte(raw, '('); idx >= 0 {
		base = raw[:idx]
	}
	base = strings.TrimSpace(base)

	nullable := ct.Null

	// parse optional precision/scale or length from parens
	precision, scale, length := 0, 0, uint(0)
	if start := strings.IndexByte(raw, '('); start >= 0 {
		if end := strings.IndexByte(raw, ')'); end > start {
			inner := raw[start+1 : end]
			parts := strings.SplitN(inner, ",", 2)
			if len(parts) == 2 {
				fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &precision)
				fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &scale)
			} else {
				var l int
				fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &l)
				length = uint(l)
			}
		}
	}

	switch {
	// MySQL convention: tinyint(1) is the canonical boolean; any wider tinyint
	// (tinyint, tinyint(4), ...) is a real integer and must stay numeric.
	case base == "tinyint":
		if raw == "tinyint(1)" {
			return &dal.TypeBoolean{Nullable: nullable}
		}
		return &dal.TypeNumber{Precision: precision, Scale: scale, Nullable: nullable}

	case isNumericType(base):
		return &dal.TypeNumber{
			Precision: precision,
			Scale:     scale,
			Nullable:  nullable,
		}

	case isTextType(base):
		return &dal.TypeText{
			Length:   length,
			Nullable: nullable,
		}

	case base == "timestamptz" || base == "timestamp with time zone":
		return &dal.TypeTimestamp{Timezone: true, Nullable: nullable}

	case base == "timestamp" || base == "timestamp without time zone" || base == "datetime":
		return &dal.TypeTimestamp{Timezone: false, Nullable: nullable}

	case base == "timetz" || base == "time with time zone":
		return &dal.TypeTime{Timezone: true, Nullable: nullable}

	case base == "time" || base == "time without time zone":
		return &dal.TypeTime{Timezone: false, Nullable: nullable}

	case base == "date":
		return &dal.TypeDate{Nullable: nullable}

	case base == "bool" || base == "boolean":
		return &dal.TypeBoolean{Nullable: nullable}

	case base == "json" || base == "jsonb":
		return &dal.TypeJSON{Nullable: nullable}

	case base == "uuid":
		return &dal.TypeUUID{Nullable: nullable}

	case base == "bytea" || base == "blob" || base == "longblob" || base == "mediumblob" || base == "tinyblob":
		return &dal.TypeBlob{Nullable: nullable}

	default:
		logger.Default().Warn("unknown SQL type, falling back to TypeText",
			zap.String("type", ct.Name),
		)
		return &dal.TypeText{Nullable: nullable}
	}
}

func isNumericType(base string) bool {
	switch base {
	case "numeric", "decimal",
		"int", "int2", "int4", "int8",
		"integer", "bigint", "smallint", "mediumint",
		"real", "double", "double precision", "float",
		"serial", "bigserial", "smallserial":
		return true
	}
	return false
}

func isTextType(base string) bool {
	switch base {
	case "varchar", "text", "char", "character", "bpchar",
		"character varying",
		"tinytext", "mediumtext", "longtext",
		"nvarchar", "nchar", "ntext":
		return true
	}
	return false
}
