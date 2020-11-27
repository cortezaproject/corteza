package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/DB.yaml

type (
	DBOpt struct {
		DSN string `env:"DB_DSN"`
	}
)

// DB initializes and returns a DBOpt with default values
func DB() (o *DBOpt) {
	o = &DBOpt{
		DSN: "sqlite3://file::memory:?cache=shared&mode=memory",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *DB) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
