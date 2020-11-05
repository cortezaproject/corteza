package options

import "strings"

type (
	EnvironmentOpt struct {
		Environment string `env:"ENVIRONMENT"`
	}
)

func Environment() (o *EnvironmentOpt) {
	o = &EnvironmentOpt{
		Environment: "production",
	}

	fill(o)

	return
}

func (e EnvironmentOpt) IsDevelopment() bool {
	return strings.HasPrefix(e.Environment, "dev")
}

func (e EnvironmentOpt) IsTest() bool {
	return strings.HasPrefix(e.Environment, "test")
}

func (e EnvironmentOpt) IsProduction() bool {
	return !e.IsDevelopment() &&
		!e.IsTest()
}
