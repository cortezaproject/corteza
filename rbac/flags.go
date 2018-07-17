package rbac

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	configuration struct {
		auth    string
		tenant  string
		baseURL string
		timeout int
	}
)

var config configuration

func (c configuration) validate() error {
	if c.auth == "" {
		return errors.New("No authentication provided for RBAC")
	}
	if c.tenant == "" {
		return errors.New("No tenant provided for RBAC")
	}
	if c.baseURL == "" {
		return errors.New("No Base URL provided for RBAC")
	}
	return nil
}

// Flags should be called from main to register flags
func Flags() {
	flag.StringVar(&config.auth, "rbac-auth", "username:password", "Credentials to use for RBAC queries")
	flag.StringVar(&config.tenant, "rbac-tenant", "", "Tenant ID")
	flag.StringVar(&config.baseURL, "rbac-base-url", "", "RBAC Base URL")
	flag.IntVar(&config.timeout, "rbac-timeout", 30, "RBAC request timeout (seconds)")
}
