package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

var (
	DefaultOption options.DiscoveryOpt
)

// Initialize discovery service
func Initialize(_ context.Context, opt options.DiscoveryOpt) (err error) {
	// @todo maybe move pkg/discovery her or other way around
	DefaultOption = opt
	return
}
