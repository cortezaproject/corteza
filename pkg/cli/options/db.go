package options

type (
	DBOpt struct {
		DSN      string
		Profiler string
	}
)

func DB(pfix string) (o *DBOpt) {
	o = &DBOpt{
		DSN:      EnvString(pfix, "DB_DSN", "corteza:corteza@tcp(db:3306)/corteza?collation=utf8mb4_general_ci"),
		Profiler: EnvString(pfix, "DB_PROFILER", "none"),
	}

	return
}
