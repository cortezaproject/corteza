package options

func (o *HTTPServerOpt) Defaults() {
	o.BaseUrl = CleanBase(o.BaseUrl)
	o.ApiBaseUrl = CleanBase(o.ApiBaseUrl)
	o.WebappBaseUrl = CleanBase(o.WebappBaseUrl)

	if o.WebappEnabled && o.ApiEnabled && (o.ApiBaseUrl == "/" || o.ApiBaseUrl == "") {
		// api base URL is still on root (empty string)
		// but webapps are enabled (that means, server also serves static files from WebappBaseDir)
		//
		// Let's be nice and move API to /api
		o.ApiBaseUrl = CleanBase("api")
	}
}
