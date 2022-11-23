package options

import "strings"

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
