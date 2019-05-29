package options

type (
	ProvisionOpt struct {
		MigrateDatabase bool

		AutoSetup bool
	}
)

func Provision(pfix string) (o *ProvisionOpt) {
	o = &ProvisionOpt{
		MigrateDatabase: EnvBool(pfix, "PROVISION_MIGRATE_DATABASE", true),
		AutoSetup:       EnvBool(pfix, "PROVISION_AUTO_SETUP", true),
	}

	return
}
