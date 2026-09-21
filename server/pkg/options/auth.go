package options

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"
)

func getSecretFromEnv(salt string) string {
	gen := salt
	// generate default secrets from virtualhost/hostname and DB_DSN value.
	// this will keep the secret the same through restarts
	gen += EnvString("DB_DSN", "memory")
	// pick one of the env that holds hostname
	gen += EnvString("HOSTNAME", "localhost")

	return fmt.Sprintf("%x", md5.Sum([]byte(gen)))
}

// GetDefaultRedirectURIs returns redirect URIs allowed for auth clients without their own
//
// Origins of auth base URL and webapp are always on the list
func (o AuthOpt) GetDefaultRedirectURIs() (out []string) {
	var (
		seen = make(map[string]bool)
		add  = func(uri string) {
			if uri != "" && !seen[uri] {
				seen[uri] = true
				out = append(out, uri)
			}
		}
	)

	add(urlOrigin(o.BaseURL))
	add(urlOrigin(FullWebappURL()))

	for _, uri := range strings.Fields(o.DefaultRedirectURIs) {
		add(uri)
	}

	return
}

// urlOrigin returns scheme and host of the absolute URL
func urlOrigin(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || !u.IsAbs() || u.Host == "" {
		return ""
	}

	return u.Scheme + "://" + u.Host
}
