package flags

import (
	"github.com/spf13/cobra"
)

type (
	ProvisionOpt struct {
		Database bool
	}
)

func Provision(cmd *cobra.Command, pfix string) (o *ProvisionOpt) {
	o = &ProvisionOpt{}

	bindBool(cmd, &o.Database,
		pFlag(pfix, "provision-database"), true,
		"Run database migration scripts")

	return
}
