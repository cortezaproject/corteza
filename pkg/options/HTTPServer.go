package options

func (o *HTTPServerOpt) Defaults() {

	if o.WebappEnabled && o.ApiEnabled && o.ApiBaseUrl == "" {
		// api base URL is still on root (empty string)
		// but webapps are enabled (that means, server also serves static files from WebappBaseDir)
		//
		// Let's be nice and move API to /api
		o.ApiBaseUrl = "/api"
	}
}
