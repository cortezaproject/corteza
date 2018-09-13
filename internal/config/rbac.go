package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	RBAC struct {
		Auth    string
		Tenant  string
		BaseURL string
		Timeout int
	}
)

var rbac *RBAC

func (c *RBAC) Validate() error {
	if c.Auth == "" {
		return errors.New("No authentication provided for RBAC")
	}
	if c.Tenant == "" {
		return errors.New("No tenant provided for RBAC")
	}
	if c.BaseURL == "" {
		return errors.New("No Base URL provided for RBAC")
	}
	return nil
}

func (*RBAC) Init(prefix ...string) *RBAC {
	if rbac != nil {
		return rbac
	}
	rbac = new(RBAC)
	flag.StringVar(&rbac.Auth, "rbac-auth", "username:password", "Credentials to use for RBAC queries")
	flag.StringVar(&rbac.Tenant, "rbac-tenant", "", "Tenant ID")
	flag.StringVar(&rbac.BaseURL, "rbac-base-url", "", "RBAC Base URL")
	flag.IntVar(&rbac.Timeout, "rbac-timeout", 30, "RBAC request timeout (seconds)")
	return rbac
}
