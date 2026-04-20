package options

import (
	"fmt"
	"github.com/crusttech/human/server/pkg/options"
)

type (
	HumanOpt struct {
		BaseUrl      string
		AuthUrl      string
		DiscoveryUrl string
	}
)

const (
	envKeyBaseUrl      = "HUMAN_SERVER_BASE_URL"
	envKeyAuthUrl      = "HUMAN_SERVER_AUTH_URL"
	envKeyDiscoveryUrl = "HUMAN_SERVER_DISCOVERY_URL"
)

func Human() (o *HumanOpt, err error) {
	o = &HumanOpt{}

	return o, func() error {
		baseUrl := options.EnvString(envKeyBaseUrl, "http://server:80")
		o.BaseUrl = baseUrl

		o.AuthUrl = options.EnvString(envKeyAuthUrl, baseUrl+"/auth")
		if o.AuthUrl == "" {
			return fmt.Errorf("Human Auth endpoint value empty, set it directly with %s or indirectly with %s", envKeyAuthUrl, envKeyBaseUrl)
		}

		o.DiscoveryUrl = options.EnvString(envKeyDiscoveryUrl, baseUrl+"/api/discovery")
		if o.DiscoveryUrl == "" {
			return fmt.Errorf("Human Discovery API endpoint value empty, set it directly with %s or indirectly with %s", envKeyDiscoveryUrl, envKeyBaseUrl)
		}

		return nil
	}()
}
