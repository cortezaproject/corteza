package types

import (
	"errors"
	"strings"
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
			// Not a setting, only a holder for potential
			// warning for from subscription
			SignupWarning string `json:",omitempty" kv:"-"`

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
				RedirectUrl string `json:"-" kv:"redirect-url"`

				// session secret to use
				SessionStoreSecret string `json:"-" kv:"session-store-secret"`

				// session store should be secure
				SessionStoreSecure bool `json:"-" kv:"session-store-secure"`

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
				FromAddress string `kv:"from-address"`
				FromName    string `kv:"from-name"`

				EmailConfirmation struct {
					Subject string `kv:"subject.en"`
					Body    string `kv:"body.en"`
				} `kv:"email-confirmation"`

				PasswordReset struct {
					Subject string `kv:"subject.en"`
					Body    string `kv:"body.en"`
				} `kv:"password-reset"`
			} `json:"-"`
		} `json:"-"`

		UI struct {
			// Corteza One configuration settings
			One struct {
				Logo  string `kv:"logo" json:"logo"`
				Title string `kv:"title" json:"title"`

				PanelsEnabled bool `kv:"panels-enabled" json:"panelsEnabled"`
				Panels        []struct {
					Visible        bool `json:"visible"`
					Sticky         bool `json:"sticky"`
					Width          uint `json:"width"`
					Height         uint `json:"height"`
					ActiveTabIndex uint `json:"activeTabIndex"`
					Tabs           []struct {
						Title  string `json:"title"`
						Url    string `json:"url"`
						Logo   string `json:"logo"`
						Icon   string `json:"icon"`
						Sticky bool   `json:"sticky"`
					} `json:"tabs"`
				} `kv:"panels,final" json:"panels"`
			} `kv:"one" json:"One"`

			// Admin struct {} `kv:"admin"`
		} `kv:"ui" json:"ui"`

		Compose struct {
			// UI related settings
			UI struct {
				// Emoji
				// @todo implementation
				NamespaceSwitcher struct {
					Enabled     bool
					DefaultOpen bool
				} `kv:"namespace-switcher"`
			} `kv:"ui"`

			// Message related settings
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

		Messaging struct {
			// UI related settings
			UI struct {
				// Emoji
				// @todo implementation
				Emoji struct {
					Enabled bool
				}

				// In-browser notifications
				// @todo implementation
				BrowserNotifications struct {
					Enabled     bool
					Header      string
					MessageTrim uint `kv:"message-trim"`
				} `kv:"browser-notifications"`
			} `kv:"ui"`

			// Message related settings
			Message struct {
				// @todo implementation
				Attachments struct {
					// Completely disable attachments
					Enabled bool

					// What is max size (in MB, so: MaxSize x 2^20)
					MaxSize uint `kv:"max-size"`

					// List of mime-types we support,
					Mimetypes []string

					// Enable/disable individual attachment sources (mobile)
					Source struct {
						Gallery struct{ Enabled bool }
						Camera  struct{ Enabled bool }
					}
				}
			}
		} `kv:"messaging" json:"messaging"`
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
				if strings.HasPrefix(p.Handle, oidcPrefix) {
					p.Label = strings.Title(p.Handle[len(oidcPrefix):])
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
		v := &SettingValue{Name: prefix + key}

		if err = v.SetValue(value); err != nil {
			return
		}

		vv = append(vv, v)
	}

	return
}
