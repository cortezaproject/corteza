package types

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/settings"
)

type (
	// Settings structured representation of current system settings
	Settings struct {
		General struct {
			Mail struct {
				Logo   string
				Header string `kv:"header.en"`
				Footer string `kv:"footer.en"`
			}
		}

		Auth struct {
			Internal struct {
				// Is internal authentication (username + password) enabled
				Enabled bool

				Signup struct {
					// Can users register
					Enabled bool

					// Users must confirm their emails when signing-up
					EmailConfirmationRequired bool `kv:"email-confirmation-required"`
				}

				// Can users reset their passwords
				PasswordReset struct{ Enabled bool } `kv:"password-reset"`
			}

			External struct {
				// Is external authentication
				Enabled bool

				// Where to redirect (url used for registration)
				RedirectUrl string `kv:"redirect-url"`

				// session secret to use
				SessionStoreSecret string `kv:"session-store-secret"`

				// session store should be secure
				SessionStoreSecure bool `kv:"session-store-secure"`

				// all external providers we know
				Providers ExternalAuthProviderSet
			}

			Frontend struct {
				Url struct {
					// Password reset path (<frontend password reset url> "?token=" + <token>)
					PasswordReset string `kv:"password-reset"`

					// EmailAddress confirmation path (<frontend  email confirmation url> "?token=" + <token>)
					EmailConfirmation string `kv:"email-confirmation"`

					// Where to redirect user after external auth flow
					Redirect string

					// Webapp Base URL
					Base string
				}
			}

			Mail struct {
				FromAddress string `kv:"from-name"`
				FromName    string `kv:"from-address"`

				EmailConfirmation struct {
					Subject string `kv:"subject.en"`
					Body    string `kv:"body.en"`
				} `kv:"email-confirmation"`

				PasswordReset struct {
					Subject string `kv:"subject.en"`
					Body    string `kv:"body.en"`
				} `kv:"password-reset"`
			}
		}
	}

	ExternalAuthProviderSet []*ExternalAuthProvider

	ExternalAuthProvider struct {
		Enabled     bool   `json:"enabled"`
		Handle      string `json:"handle"`
		Label       string `json:"label"`
		Key         string `json:"-"`
		Secret      string `json:"-"`
		RedirectUrl string `json:",omitempty" kv:"redirect"`
		IssuerUrl   string `json:",omitempty" kv:"issuer"`
		Weight      int    `json:"-"`
	}
)

// DecodeKV translates settings' KV into internal system external auth settings
func (set *ExternalAuthProviderSet) DecodeKV(kv settings.KV, prefix string) (err error) {
	if *set == nil {
		*set = ExternalAuthProviderSet{}
	}

	// create standard provider set
	providers := map[string]bool{"github": true, "facebook": true, "google": true, "linkedin": true}

	// remove prefix
	kv = kv.CutPrefix(prefix + ".")

	// add all additional providers (prefixed with "openid-connect.")
	oidcPrefix := "openid-connect."
	for p := range kv {
		if !strings.HasPrefix(p, oidcPrefix) {
			continue
		}

		l := len(oidcPrefix)
		dotPos := strings.Index(p[l:], ".") + l
		if dotPos > 0 {
			providers[p[:dotPos]] = true
		}
	}

	// go over all added providers again add decode KV into each one
	for handle := range providers {
		p := (*set).FindByHandle(handle)
		if p == nil {
			p = &ExternalAuthProvider{Handle: handle}
			(*set) = append((*set), p)
		}

		err = settings.DecodeKV(kv.CutPrefix(handle+"."), p)
		if err != nil {
			return
		}

		if p.Label == "" {
			switch p.Handle {
			case "corteza-iam", "corteza", "corteza-one":
				p.Label = "Corteza One"
			case "crust-iam", "crust", "crust-unify":
				p.Label = "Crust Unify"
			default:
				strings.Title(p.Handle)
			}
		}
	}

	return
}

func (set ExternalAuthProviderSet) FindByHandle(handle string) *ExternalAuthProvider {
	for p := range set {
		if set[p].Handle == handle {
			return set[p]
		}
	}

	return nil
}

func (set ExternalAuthProviderSet) Len() int           { return len(set) }
func (set ExternalAuthProviderSet) Swap(i, j int)      { set[i], set[j] = set[j], set[i] }
func (set ExternalAuthProviderSet) Less(i, j int) bool { return set[i].Weight < set[j].Weight }

var _ settings.KVDecoder = &ExternalAuthProviderSet{}

func (p ExternalAuthProvider) EncodeKV() (vv settings.ValueSet, err error) {
	if p.Handle == "" {
		return nil, errors.New("can not encode external auth provider without handle")
	}
	var (
		prefix = "auth.external.providers." + p.Handle + "."
		pairs  = map[string]interface{}{
			"enabled":  p.Enabled,
			"label":    p.Label,
			"key":      p.Key,
			"secret":   p.Secret,
			"issuer":   p.IssuerUrl,
			"redirect": p.RedirectUrl,
			"weight":   p.Weight,
		}
	)

	for key, value := range pairs {
		v := &settings.Value{Name: prefix + key}

		if err = v.SetValue(value); err != nil {
			return
		}

		vv = append(vv, v)
	}

	return
}
