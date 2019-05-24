package flags

import (
	"github.com/spf13/cobra"
)

type (
	ProvisionOpt struct {
		MigrateDatabase bool

		AutoSetup bool
	}
)

func Provision(cmd *cobra.Command, pfix string) (o *ProvisionOpt) {
	o = &ProvisionOpt{}

	BindBool(cmd, &o.MigrateDatabase,
		pFlag(pfix, "provision-migrate-database"), true,
		"Run database migration")

	BindBool(cmd, &o.AutoSetup,
		pFlag(pfix, "provision-auto-setup"), true,
		"Run auto-setup procedures on service")

	return
}
