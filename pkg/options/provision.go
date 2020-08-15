package options

type (
	ProvisionOpt struct {
		Always bool `env:"PROVISION_ALWAYS"`
	}
)

func Provision(pfix string) (o *ProvisionOpt) {
	o = &ProvisionOpt{
		Always: true,
	}

	fill(o)

	return
}
