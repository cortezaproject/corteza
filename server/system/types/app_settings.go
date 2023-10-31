package types

import (
	"fmt"
	"strings"
	"time"
)

const (
	oidcProviderPrefix = "openid-connect." // must match const in "github.com/cortezaproject/corteza/server/auth/external" external.go

	ExternalProviderUsageIdentity = "identity"
	ExternalProviderUsageAPI      = "api"
)

type (
	// AppSettings type is structured representation of all application settings
	//
	// Raw settings keys are hyphen (kebab) case, separated with a dot (.) that indicates sub-level
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

		SMTP struct {
			Servers []SmtpServers `json:"-" kv:"servers,final"`
		} `json:"-" kv:"smtp"`

		Auth struct {
			Internal struct {
				// Is internal authentication (username + password) enabled
				Enabled bool `json:"-"`

				Signup struct {
					// Can users register
					Enabled bool

					// Users must confirm their emails when signing-up
					EmailConfirmationRequired bool `kv:"email-confirmation-required"`
				} `json:"-"`

				// Can users reset their passwords
				PasswordReset struct{ Enabled bool } `json:"-" kv:"password-reset"`

				// When enabled, users added via CLI will receive an email with a link to reset their password.
				SendUserInviteEmail struct {
					Enabled bool
					Expires uint
				} `kv:"send-user-invite-email" json:"sendUserInviteEmail"`

				// PasswordCreate setting for create password for user via generated link with token
				// If user has no password then link redirects to create password page
				// Otherwise it redirects to profile page of that user
				// link can be generated through useradd cli command with `make-password-link` flag
				PasswordCreate struct {
					Enabled bool
					Expires uint
				} `json:"-" kv:"password-create"`

				// Splits credentials check into 2 parts
				// If user has password credentials it offers him to enter the password
				// Otherwise we offer the user to choose among the enabled external providers
				// If only one ext. provider is enabled, user is automatically redirected there
				SplitCredentialsCheck bool `json:"-" kv:"split-credentials-check"`

				PasswordConstraints PasswordConstraints `kv:"password-constraints" json:"passwordConstraints"`

				ProfileAvatar struct{ Enabled bool } `kv:"profile-avatar" json:"profile-avatar"`
			} `json:"internal"`

			External struct {
				// Is external authentication
				Enabled bool `json:"enabled"`

				// Saml
				Saml struct {
					Enabled bool

					// IdP name used on the login form
					Name string `kv:"name"`

					// SAML certificate
					Cert string `kv:"cert"`

					// SAML certificate private key
					Key string `kv:"key"`

					// Sign AuthNRequest and assertion
					SignRequests bool `kv:"sign-requests"`

					// Signature method for signing
					SignMethod string `kv:"sign-method"`

					// Post or redirect binding
					Binding string `kv:"binding"`

					// Identity provider settings
					IDP struct {
						URL string `kv:"url"`

						// identifier payload from idp
						IdentName       string `kv:"ident-name"`
						IdentHandle     string `kv:"ident-handle"`
						IdentIdentifier string `kv:"ident-identifier"`
					} `kv:"idp"`

					Security ExternalAuthProviderSecurity `json:"-" kv:"security,final"`
				} `json:"-"`

				// all external providers we know
				Providers ExternalAuthProviderSet `json:"providers"`
			} `json:"external"`

			MultiFactor struct {
				EmailOTP struct {
					// Can users use email for MFA
					Enabled bool

					// Is MFA with email enforced?
					Enforced bool

					// Require fresh Email OTP on every client authorization
					// Strict bool

					Expires uint
				} `kv:"email-otp"`

				TOTP struct {
					// Can users use TOTP for MFA
					Enabled bool

					// Is MFA with TOTP enforced?
					Enforced bool

					// Require fresh TOTP on every client authorization
					// Strict bool

					// TOTP issuer, defaults to "Corteza"
					Issuer string
				} `kv:"totp"`
			} `json:"-" kv:"multi-factor"`

			Mail struct {
				FromAddress string `kv:"from-address"`
				FromName    string `kv:"from-name"`
			} `json:"-"`

			// Auth Background Image settings
			UI struct {
				BackgroundImageSrc string `kv:"background-image-src" json:"backgroundImageSrc"`
				Styles             string `kv:"styles" json:"styles"`
			} `kv:"ui" json:"ui"`
		} `json:"auth"`

		Compose struct {
			// Compose UI settings
			UI struct {
				// Sidebar specific settings
				Sidebar struct {
					// Hide namespace list
					HideNamespaceList bool `json:"hideNamespaceList"`

					// Hide namespace link at the end of the list
					HideNamespaceListLink bool `json:"hideNamespaceListLink"`
				} `kv:"sidebar,final" json:"sidebar"`

				// Record toolbar specific settings
				RecordToolbar struct {
					HideNew    bool `json:"hideNew"`
					HideEdit   bool `json:"hideEdit"`
					HideSubmit bool `json:"hideSubmit"`
					HideDelete bool `json:"hideDelete"`
					HideClone  bool `json:"hideClone"`
					HideBack   bool `json:"hideBack"`
				} `kv:"record-toolbar,final" json:"record-toolbar"`
			} `kv:"ui" json:"ui"`

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

			// Icon related settings
			Icon struct {
				// @todo implementation
				Attachments struct {
					// What is max size (in MB, so: MaxSize x 2^20)
					MaxSize uint `kv:"max-size"`

					// List of mime-types we support,
					Mimetypes []string
				}
			}
		} `kv:"compose" json:"compose"`

		// Federation settings
		Federation struct {
			// This only holds the value of FEDERATION_ENABLED for now
			//
			Enabled bool `kv:"-" json:"enabled"`
		} `kv:"federation" json:"federation"`

		// Integration gateway settings
		Apigw ApigwSettings `kv:"apigw" json:"apigw"`

		// UserInterface settings
		UI struct {
			MainLogo string `kv:"main-logo" json:"mainLogo"`
			IconLogo string `kv:"icon-logo" json:"iconLogo"`

			Sidebar struct {
				// General sidebar settings
				Disabled bool `json:"disabled"`
			} `kv:"sidebar,final" json:"sidebar"`

			CustomCSS string `kv:"custom-css" json:"customCSS"`
			Studio struct {
				BrandingSASS  string `kv:"branding-sass" json:"branding-sass"`
				SassInstalled bool   `kv:"sass-installed" json:"sass-installed"`
			} `kv:"studio" json:"studio"`

			Topbar struct {
				HideAppSelector        bool `json:"hideAppSelector"`
				HideHelp               bool `json:"hideHelp"`
				HideForumLink          bool `json:"hideForumLink"`
				HideDocumentationLink  bool `json:"hideDocumentationLink"`
				HideFeedbackLink       bool `json:"hideFeedbackLink"`
				HideProfile            bool `json:"hideProfile"`
				HideChangePasswordLink bool `json:"hideChangePasswordLink"`
				HideProfileLink        bool `json:"hideProfileLink"`

				HelpLinks []struct {
					Handle string `json:"handle"`
					URL    string `json:"url"`
					NewTab bool   `json:"newTab"`
				} `json:"helpLinks"`

				ProfileLinks []struct {
					Handle string `json:"handle"`
					URL    string `json:"url"`
					NewTab bool   `json:"newTab"`
				} `json:"profileLinks"`
			} `kv:"topbar,final" json:"topbar"`

			Charts struct {
				ColorSchemes []struct {
					ID     string   `json:"id"`
					Name   string   `json:"name"`
					Colors []string `json:"colors"`
				} `kv:"colorSchemes" json:"colorSchemes"`
			} `kv:"charts" json:"charts"`
		} `kv:"ui" json:"ui"`

		ResourceTranslations struct {
			// List of all languages (resource translations) enabled and
			// available for resource translations (these are module names,
			// field labels, descriptions, ...)

			// This is always a subset of all languages available
			// in Corteza instance (LOCALE_LANGUAGES)
			//
			// Note: later, we will enable this to contain languages
			//       that are not part of LOCALE_LANGUAGES
			//
			// 1st language in the set is also a default one
			//
			// Empty slice defaults to LOCALE_LANGUAGES
			Languages []string `kv:"languages" json:"languages"`
		} `kv:"resource-translations" json:"resourceTranslations"`

		Discovery struct {
			// Enable indexing
			Enabled bool `kv:"enabled" json:"enabled"`

			SystemUsers struct {
				// Enable indexing of users
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"system-users" json:"system-users"`

			SystemApplications struct {
				// Enable indexing of applications
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"system-applications" json:"system-applications"`

			SystemRoles struct {
				// Enable indexing of roles
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"system-roles" json:"system-roles"`

			SystemTemplates struct {
				// Enable indexing of templates
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"system-templates" json:"system-templates"`

			AutomationWorkflows struct {
				// Enable indexing of workflows
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"automation-workflows" json:"automation-workflows"`

			ComposeNamespaces struct {
				// Enable indexing of compose namespaces
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"compose-namespaces" json:"compose-namespaces"`

			ComposeCharts struct {
				// Enable indexing of compose charts
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"compose-charts" json:"compose-charts"`

			ComposePages struct {
				// Enable indexing of compose pages
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"compose-pages" json:"compose-pages"`

			ComposeModules struct {
				// Enable indexing of compose modules
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"compose-modules" json:"compose-modules"`

			ComposeRecords struct {
				// Enable indexing of compose records
				Enabled bool `kv:"enabled" json:"enabled"`
			} `kv:"compose-records" json:"compose-records"`
		} `kv:"discovery" json:"discovery"`
	}

	ExternalAuthProviderSet []*ExternalAuthProvider

	ExternalAuthProvider struct {
		Enabled     bool   `json:"enabled"`
		Handle      string `json:"handle"`
		Label       string `json:"label"`
		Key         string `json:"-"`
		Secret      string `json:"-"`
		Scope       string `json:"scope"`
		RedirectUrl string `json:"-" kv:"redirect"`
		IssuerUrl   string `json:"-" kv:"issuer"`
		Weight      int    `json:"-"`

		// potentially just use home-brew scopes?
		Usage []string `json:"usage" kv:"usage"`

		Security ExternalAuthProviderSecurity `json:"-" kv:"security,final"`
	}

	ExternalAuthProviderSecurity struct {
		// Subset of roles, permitted to be used with this client
		// when authorizing via this auth provider.
		//
		// IDs are intentionally stored as strings to support JS (int64 only)
		//
		PermittedRoles []string `json:"permittedRoles,omitempty"`

		// Subset of roles, prohibited to be used with this client
		// when authorizing via this auth provider.
		//
		// IDs are intentionally stored as strings to support JS (int64 only)
		//
		ProhibitedRoles []string `json:"prohibitedRoles,omitempty"`

		// Set of additional roles that are forced on this user
		// when authorizing via this auth provider.
		//
		// IDs are intentionally stored as strings to support JS (int64 only)
		ForcedRoles []string `json:"forcedRoles,omitempty"`

		// Map external roles or groups to internal
		//
		// If IdP provides a list of roles (groups) along side authenticated user
		// these roles can be mapped to the valid local roles
		//
		// @todo implement mapped roles
		// MappedRoles map[string]string `json:"mappedRoles,omitempty"`
	}

	SmtpServers struct {
		Host          string `json:"host"`
		Port          int    `json:"port,string"`
		User          string `json:"user"`
		Pass          string `json:"pass"`
		From          string `json:"from"`
		TlsInsecure   bool   `json:"tlsInsecure"`
		TlsServerName string `json:"tlsServerName"`
	}

	ApigwSettings struct {
		Enabled bool `kv:"enabled" json:"enabled"`

		Profiler struct {
			Enabled bool `kv:"enabled" json:"enabled"`
			Global  bool `kv:"global" json:"global"`
		} `kv:"profiler" json:"profiler"`

		Proxy struct {
			FollowRedirects bool          `kv:"follow-redirects" json:"follow-redirects"`
			OutboundTimeout time.Duration `kv:"outbound-timeout" json:"outbound-timeout"`
		} `kv:"proxy" json:"proxy"`
	}

	PasswordConstraints struct {
		// Should the environment not enforce the constraints
		PasswordSecurity bool `kv:"-" json:"passwordSecurity"`

		// The min password length
		MinLength uint `kv:"min-length" json:"minLength"`

		// Minimum number of uppercase letters in password
		MinUpperCase uint `kv:"min-upper-case" json:"minUpperCase"`

		// Minimum number of lowercase letters in password
		MinLowerCase uint `kv:"min-lower-case" json:"minLowerCase"`

		// The min number of numeric characters
		MinNumCount uint `kv:"min-num-count" json:"minNumCount"`

		// The min number of special characters
		MinSpecialCount uint `kv:"min-special-count" json:"minSpecialCount"`
	}
)

// WithDefaults sets defaults on copy (!!) of settings
// to avoid any unintended corruption or leaks
func (cs AppSettings) WithDefaults() *AppSettings {
	if len(strings.TrimSpace(cs.UI.IconLogo)) == 0 {
		cs.UI.IconLogo = "/assets/favicon32x32.png"
	}
	if len(strings.TrimSpace(cs.UI.MainLogo)) == 0 {
		cs.UI.MainLogo = "/assets/logo.png"
	}

	return &cs
}

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
	permanent := map[string]bool{"github": true, "facebook": true, "google": true, "linkedin": true, "nylas": true}
	// and make a working copy
	providers := make(map[string]bool)
	for k, v := range permanent {
		providers[k] = v
	}

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
			*set = append(*set, p)
		}

		err = DecodeKV(kv.CutPrefix(handle+"."), p)
		if err != nil {
			return
		}
	}

	// Cleanup
	var clean = ExternalAuthProviderSet{}
	for _, p := range *set {
		if p.empty() && !permanent[p.Handle] {
			continue
		}

		if p.Label == "" {
			switch p.Handle {
			case "github":
				p.Label = "GitHub"
			case "nylas":
				p.Label = "Nylas"
			case "linkedin":
				p.Label = "LinkedIn"
			case "corteza-iam", "corteza", "corteza-one":
				// Some legacy provider naming
				p.Label = "Corteza IAM"
			case "crust-iam", "crust", "crust-unify":
				// Some legacy provider naming
				p.Label = "Crust IAM"
			default:
				if strings.HasPrefix(p.Handle, oidcProviderPrefix) {
					p.Label = strings.Title(p.Handle[len(oidcProviderPrefix):])
				} else {
					p.Label = strings.Title(p.Handle)
				}
			}
		}

		clean = append(clean, p)
	}

	*set = clean

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
		if eap.empty() || !eap.Enabled {
			continue
		}

		out = append(out, eap)
	}

	return
}

var _ KVDecoder = &ExternalAuthProviderSet{}

func (eap ExternalAuthProvider) EncodeKV() (vv SettingValueSet, err error) {
	if eap.Handle == "" {
		return nil, fmt.Errorf("cannot encode external auth provider without handle")
	}
	var (
		prefix = "auth.external.providers." + eap.Handle + "."
		pairs  = map[string]interface{}{
			"enabled":  eap.Enabled,
			"label":    eap.Label,
			"key":      eap.Key,
			"secret":   eap.Secret,
			"scope":    eap.Scope,
			"issuer":   eap.IssuerUrl,
			"redirect": eap.RedirectUrl,
			"weight":   eap.Weight,
		}
	)

	for key, value := range pairs {
		v := &SettingValue{Name: prefix + key}

		if err = v.SetSetting(value); err != nil {
			return
		}

		vv = append(vv, v)
	}

	return
}

// returns true if all relevant props are empty
func (eap ExternalAuthProvider) empty() bool {
	return len(
		eap.RedirectUrl+
			eap.Secret+
			eap.Key,
	) == 0
}
