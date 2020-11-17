package options

import "time"

type (
	FederationOpt struct {
		Enabled                  bool          `env:"FEDERATION_ENABLED"`
		Label                    string        `env:"FEDERATION_LABEL"`
		Host                     string        `env:"FEDERATION_HOST"`
		StructureMonitorInterval time.Duration `env:"FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL"`
		StructurePageSize        int           `env:"FEDERATION_SYNC_STRUCTURE_PAGE_SIZE"`
		DataMonitorInterval      time.Duration `env:"FEDERATION_SYNC_DATA_MONITOR_INTERVAL"`
		DataPageSize             int           `env:"FEDERATION_SYNC_DATA_PAGE_SIZE"`
	}
)

func Federation() (o *FederationOpt) {

	o = &FederationOpt{
		Enabled:                  false,
		Label:                    `Example host`,
		Host:                     `example.tld`,
		StructureMonitorInterval: 3600,
		StructurePageSize:        100,
		DataMonitorInterval:      3600,
		DataPageSize:             100,
	}

	fill(o)

	return
}
