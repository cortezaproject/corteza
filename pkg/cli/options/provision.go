package options

type (
	ProvisionOpt struct {
		MigrateDatabase bool `env:"PROVISION_MIGRATE_DATABASE"`
		AutoSetup       bool `env:"PROVISION_AUTO_SETUP"`
	}
)

func Provision(pfix string) (o *ProvisionOpt) {
	o = &ProvisionOpt{
		MigrateDatabase: true,
		AutoSetup:       true,
	}

	fill(o, pfix)

	return
}
