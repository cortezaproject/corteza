package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/federation.yaml

import (
	"time"
)

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

// Federation initializes and returns a FederationOpt with default values
func Federation() (o *FederationOpt) {
	o = &FederationOpt{
		Enabled:                  false,
		Label:                    "Federation label",
		Host:                     "local.cortezaproject.org",
		StructureMonitorInterval: time.Minute * 2,
		StructurePageSize:        1,
		DataMonitorInterval:      time.Second * 60,
		DataPageSize:             100,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *Federation) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
