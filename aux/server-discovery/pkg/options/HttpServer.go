package options

type (
	HttpServerOpt struct {
		Addr                   string `env:"HTTP_ADDR"`
		EnableHealthcheckRoute bool   `env:"HTTP_ENABLE_HEALTHCHECK_ROUTE"`
		EnableVersionRoute     bool   `env:"HTTP_ENABLE_VERSION_ROUTE"`
		BaseUrl                string `env:"HTTP_BASE_URL"`
		ApiBaseUrl             string `env:"HTTP_API_BASE_URL"`
	}
)

// HttpServer initializes and returns a HTTPServerOpt with default values
func HttpServer() (o *HttpServerOpt) {
	o = &HttpServerOpt{
		Addr:                   ":80",
		EnableHealthcheckRoute: true,
		EnableVersionRoute:     true,
		BaseUrl:                "/",
		ApiBaseUrl:             "/",
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *HTTPServer) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}
