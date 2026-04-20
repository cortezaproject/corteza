package options

import (
	"strings"
)

// Parses hosts and return slice of strings, one per host
func (o WaitForOpt) GetServices() []string {
	if len(o.Services) == 0 {
		return []string{}
	}

	return strings.Split(o.Services, " ")
}
