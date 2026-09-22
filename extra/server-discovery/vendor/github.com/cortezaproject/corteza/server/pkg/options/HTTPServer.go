package options

import (
	"strings"
)

func (o *HttpServerOpt) Defaults() {
	if Environment().IsDevelopment() {
		// enable web console and remove username, password defaults
		// if this is explicitly via ENV, it will override these defaults
		o.WebConsoleEnabled = true
		o.WebConsoleUsername = ""
		o.WebConsolePassword = ""
	}
}

func (o *HttpServerOpt) Cleanup() {
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

// CorsAnyOrigin allows cross-origin requests from any origin
func CorsAnyOrigin() []string {
	return []string{"http://*", "https://*"}
}

// GetCorsAllowedOrigins parses comma separated origins, falling back to any origin when empty
func (o HttpServerOpt) GetCorsAllowedOrigins() []string {
	var origins []string

	for _, origin := range strings.Split(o.CorsAllowedOrigins, ",") {
		// browsers never send origins with trailing slash
		origin = strings.TrimRight(strings.TrimSpace(origin), "/")

		switch origin {
		case "":
			continue
		case "*":
			return CorsAnyOrigin()
		}

		origins = append(origins, origin)
	}

	if len(origins) == 0 {
		return CorsAnyOrigin()
	}

	return origins
}
