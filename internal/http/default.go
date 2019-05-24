package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// SetupDefaults Reconfigures defaults for  HTTP client & transport
func SetupDefaults(timeout time.Duration, tslInsecure bool) {
	if tslInsecure {
		// This will allow HTTPS requests to insecure hosts (expired, wrong host, self signed, untrusted root...)
		// With this enabled, features like OIDC auto-discovery should work on any of examples found on badssl.com.
		//
		// With SYSTEM_HTTP_CLIENT_TSL_INSECURE=0 (default) next command returns 404 error (expected)
		// > ./system external-auth auto-discovery foo-tsl-1 https://expired.badssl.com/
		//
		// Without SYSTEM_HTTP_CLIENT_TSL_INSECURE=1 next command returns "x509: certificate has expired or is not yet valid"
		// > ./system external-auth auto-discovery foo-tsl-1 https://expired.badssl.com/
		//
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		http.DefaultTransport.(*http.Transport).DialContext = (&net.Dialer{Timeout: timeout}).DialContext
		http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout = timeout
	}

	if timeout > 0 {
		http.DefaultClient.Timeout = timeout
	}

	http.DefaultClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
}
