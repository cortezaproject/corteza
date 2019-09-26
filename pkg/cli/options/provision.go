package options

type (
	ProvisionOpt struct {
		MigrateDatabase bool `env:"PROVISION_MIGRATE_DATABASE"`
		Configuration   bool `env:"PROVISION_CONFIGURATION"`
	}
)

func Provision(pfix string) (o *ProvisionOpt) {
	o = &ProvisionOpt{
		MigrateDatabase: true,
		Configuration:   true,
	}

	fill(o, pfix)

	return
}
