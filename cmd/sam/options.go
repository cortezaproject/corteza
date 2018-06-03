package main

type RouteOptions struct {
	enableLogging bool
}

func (RouteOptions) New() (*RouteOptions, error) {
	opts := &RouteOptions{}
	opts.enableLogging = true
	return opts, nil
}

func (o *RouteOptions) EnableLogging(enable bool) *RouteOptions {
	o.enableLogging = enable
	return o
}
