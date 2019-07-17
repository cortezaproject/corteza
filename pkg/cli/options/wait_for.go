package options

import (
	"strings"
	"time"
)

type (
	WaitForOpt struct {
		Delay                 time.Duration `env:"WAIT_FOR"`
		StatusPage            bool          `env:"WAIT_FOR_STATUS_PAGE"`
		Services              string        `env:"WAIT_FOR_SERVICES"`
		ServicesTimeout       time.Duration `env:"WAIT_FOR_SERVICES_TIMEOUT"`
		ServicesProbeTimeout  time.Duration `env:"WAIT_FOR_SERVICES_PROBE_TIMEOUT"`
		ServicesProbeInterval time.Duration `env:"WAIT_FOR_SERVICES_PROBE_INTERVAL"`
	}
)

func WaitFor(pfix string) (o *WaitForOpt) {
	o = &WaitForOpt{
		Delay:                 0,
		StatusPage:            true,
		Services:              "satosa.iam.crust.tech:443 https://satosa.iam.crust.tech/.well-known/openid-configuration",
		ServicesTimeout:       time.Minute,
		ServicesProbeTimeout:  time.Second * 30,
		ServicesProbeInterval: time.Second * 5,
	}

	fill(o, pfix)

	return
}

// Parses hosts and return slice of strings, one per host
func (o WaitForOpt) GetServices() []string {
	if len(o.Services) == 0 {
		return []string{}
	}

	return strings.Split(o.Services, " ")
}
