package discovery

import "context"

type (
	Bundler interface {
		Register(context.Context, []string) error
		Deregister(context.Context, []string) error

		Validate(context.Context, []string) (byte, []string)
		Type() byte
	}
)
