package types

import (
	"errors"
	"strings"
)

const (
	oidcProviderPrefix = "openid-connect." // must match const in "github.com/cortezaproject/corteza-server/auth/external" external.go
)

type (
	// AppSettings type is structured representation of all application settings
	//
	// Raw settings keys are hypen (kebab) case, separated with a dot (.) that indicates sub-level
	// JSON properties for settings are NOT converted (lower-cased, etc...)
	// Use `json:"-"` tag to hide settings on REST endpoint
	AppSettings struct {
		Privacy struct {
			Mask struct {
				// Enable masking of user's email (value replaced with ######)
				Email bool

				// Enable masking of user's name (value replaced with ######)
				Name bool
			}
		} `json:"-"`

		General struct {
			Mail struct {
				Logo   string
				Header string `kv:"header.en"`
				Footer string `kv:"footer.en"`
			}
		} `json:"-"`

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

				// Splits credentials check into 2 parts
				// If user has password credentials it offers him to enter the password
				// Otherwise we offer the user to choose among the enabled external providers
				// If only one ext. provider is enabled, user is automatically redirected there
				SplitCredentialsCheck bool `kv:"split-credentials-check"`
			}

			External struct {
				// Is external authentication
				Enabled bool

				// Saml
				Saml struct {
					Enabled bool

					// SAML certificate
					Cert string `kv:"cert"`

					// SAML certificate private key
					Key string `kv:"key"`

					// Identity provider settings
					IDP struct {
						URL  string `kv:"url"`
						Name string

						// identifier payload from idp
						IdentName       string `kv:"ident-name"`
						IdentHandle     string `kv:"ident-handle"`
						IdentIdentifier string `kv:"ident-identifier"`
					} `kv:"idp"`
				}

				// all external providers we know
				Providers ExternalAuthProviderSet
			}

			MultiFactor struct {
				EmailOTP struct {
					// Can users use email for MFA
					Enabled bool

					// Is MFA with email enforced?
					Enforced bool

					// Require fresh Email OTP on every client authorization
					//Strict bool

					Expires uint
				} `kv:"email-otp"`

				TOTP struct {
					// Can users use TOTP for MFA
					Enabled bool

					// Is MFA with TOTP enforced?
					Enforced bool

					// Require fresh TOTP on every client authorization
					//Strict bool

					// TOTP issuer, defaults to "Corteza"
					Issuer string
				} `kv:"totp"`
			} `kv:"multi-factor"`

			Mail struct {
				FromAddress string `kv:"from-address"`
				FromName    string `kv:"from-name"`
			} `json:"-"`
		} `json:"-"`

		Compose struct {
			// UI related settings
			// (placeholder)
			UI struct{} `kv:"ui"`

			// Record related settings
			Record struct {
				// @todo implementation
				Attachments struct {
					// What is max size (in MB, so: MaxSize x 2^20)
					MaxSize uint `kv:"max-size"`

					// List of mime-types we support,
					Mimetypes []string
				}
			}

			// Page related settings
			Page struct {
				// @todo implementation
				Attachments struct {
					// What is max size (in MB, so: MaxSize x 2^20)
					MaxSize uint `kv:"max-size"`

					// List of mime-types we support,
					Mimetypes []string
				}
			}
		} `kv:"compose" json:"compose"`

		// UserInterface settings
		UI struct {
			MainLogo string `kv:"main-logo" json:"mainLogo"`
			IconLogo string `kv:"icon-logo" json:"iconLogo"`
		} `kv:"ui" json:"ui"`
	}

	ExternalAuthProviderSet []*ExternalAuthProvider

	ExternalAuthProvider struct {
		Enabled     bool   `json:"enabled"`
		Handle      string `json:"handle"`
		Label       string `json:"label"`
		Key         string `json:"-"`
		Secret      string `json:"-"`
		RedirectUrl string `json:"-" kv:"redirect"`
		IssuerUrl   string `json:"-" kv:"issuer"`
		Weight      int    `json:"-"`
	}
)

func (set *ExternalAuthProvider) ValidConfiguration() bool {
	if !set.Enabled || set.Handle == "" || set.Key == "" || set.Secret == "" {
		return false
	}

	if strings.HasPrefix(set.Handle, oidcProviderPrefix) && set.IssuerUrl == "" {
		// OIDC IdPs need to have issuer URL
		return false
	}

	return true
}

// DecodeKV translates settings' KV into internal system external auth settings
func (set *ExternalAuthProviderSet) DecodeKV(kv SettingsKV, prefix string) (err error) {
	if *set == nil {
		*set = ExternalAuthProviderSet{}
	}

	// create standard provider set
	providers := map[string]bool{"github": true, "facebook": true, "google": true, "linkedin": true}

	// remove prefix
	kv = kv.CutPrefix(prefix + ".")

	// add all additional providers (prefixed with "openid-connect.")
	for p := range kv {
		if !strings.HasPrefix(p, oidcProviderPrefix) {
			continue
		}

		l := len(oidcProviderPrefix)
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

		err = DecodeKV(kv.CutPrefix(handle+"."), p)
		if err != nil {
			return
		}

		if p.Label == "" {
			switch p.Handle {
			case "github":
				p.Label = "GitHub"
			case "linkedin":
				p.Label = "LinkedIn"
			case "corteza-iam", "corteza", "corteza-one":
				p.Label = "Corteza IAM"
			case "crust-iam", "crust", "crust-unify":
				p.Label = "Crust IAM"
			default:
				if strings.HasPrefix(p.Handle, oidcProviderPrefix) {
					p.Label = strings.Title(p.Handle[len(oidcProviderPrefix):])
				} else {
					p.Label = strings.Title(p.Handle)
				}
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

func (set ExternalAuthProviderSet) Len() int      { return len(set) }
func (set ExternalAuthProviderSet) Swap(i, j int) { set[i], set[j] = set[j], set[i] }
func (set ExternalAuthProviderSet) Less(i, j int) bool {
	if set[i].Weight != set[j].Weight {
		// Sort by weight
		return set[i].Weight < set[j].Weight
	}

	if set[i].Label+set[j].Label != "" {
		// If at least one of the
		return set[i].Label < set[j].Label
	}

	return set[i].Handle < set[j].Handle
}

// Returns enabled providers, sorted with their redirect-URLs set...
func (set ExternalAuthProviderSet) Valid() (out ExternalAuthProviderSet) {
	for _, eap := range set {
		if !eap.Enabled {
			continue
		}

		out = append(out, eap)
	}

	return
}

var _ KVDecoder = &ExternalAuthProviderSet{}

func (p ExternalAuthProvider) EncodeKV() (vv SettingValueSet, err error) {
	if p.Handle == "" {
		return nil, errors.New("cannot encode external auth provider without handle")
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
		v := &SettingValue{Name: prefix + key}

		if err = v.SetValue(value); err != nil {
			return
		}

		vv = append(vv, v)
	}

	return
}
