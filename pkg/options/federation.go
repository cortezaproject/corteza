package options

type (
	FederationOpt struct {
		Enabled bool `env:"FEDERATION_ENABLED"`
	}
)

func Federation() (o *FederationOpt) {

	o = &FederationOpt{
		Enabled: false,
	}

	fill(o)

	return
}
