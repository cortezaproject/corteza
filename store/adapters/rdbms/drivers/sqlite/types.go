package sqlite

import "github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"

func columnTypeTranslator(ct ddl.ColumnType) string {
	switch ct.Type {
	case ddl.ColumnTypeTimestamp:
		return "TIMESTAMP"
	case ddl.ColumnTypeBinary:
		return "BLOB"
	default:
		return ddl.ColumnTypeTranslator(ct)
	}

}
