package options

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/options"
)

type (
	CortezaOpt struct {
		BaseUrl      string
		AuthUrl      string
		DiscoveryUrl string
	}
)

const (
	envKeyBaseUrl      = "CORTEZA_SERVER_BASE_URL"
	envKeyAuthUrl      = "CORTEZA_SERVER_AUTH_URL"
	envKeyDiscoveryUrl = "CORTEZA_SERVER_DISCOVERY_URL"
)

func Corteza() (o *CortezaOpt, err error) {
	o = &CortezaOpt{}

	return o, func() error {
		baseUrl := options.EnvString(envKeyBaseUrl, "http://server:80")
		o.BaseUrl = baseUrl

		o.AuthUrl = options.EnvString(envKeyAuthUrl, baseUrl+"/auth")
		if o.AuthUrl == "" {
			return fmt.Errorf("corteza Auth endpoint value empty, set it directly with %s or indirectly with %s", envKeyAuthUrl, envKeyBaseUrl)
		}

		o.DiscoveryUrl = options.EnvString(envKeyDiscoveryUrl, baseUrl+"/api/discovery")
		if o.DiscoveryUrl == "" {
			return fmt.Errorf("corteza Discovery API endpoint value empty, set it directly with %s or indirectly with %s", envKeyDiscoveryUrl, envKeyBaseUrl)
		}

		return nil
	}()
}
