package options

type (
	UpgradeOpt struct {
		Debug  bool `env:"UPGRADE_DEBUG"`
		Always bool `env:"UPGRADE_ALWAYS"`
	}
)

func Upgrade(pfix string) (o *UpgradeOpt) {
	o = &UpgradeOpt{
		Debug:  false,
		Always: true,
	}

	fill(o)

	return
}
