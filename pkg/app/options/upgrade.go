package options

type (
	UpgradeOpt struct {
		Always bool `env:"UPGRADE_ALWAYS"`
	}
)

func Upgrade(pfix string) (o *UpgradeOpt) {
	o = &UpgradeOpt{
		Always: true,
	}

	fill(o, pfix)

	return
}
