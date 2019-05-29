package options

type (
	DBOpt struct {
		DSN      string `env:"DB_DSN"`
		Profiler string `env:"DB_PROFILER"`
	}
)

func DB(pfix string) (o *DBOpt) {
	o = &DBOpt{
		DSN:      "corteza:corteza@tcp(db:3306)/corteza?collation=utf8mb4_general_ci",
		Profiler: "none",
	}

	fill(o, pfix)

	return
}
