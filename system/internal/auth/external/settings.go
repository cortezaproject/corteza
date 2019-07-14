package external

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	intset "github.com/cortezaproject/corteza-server/internal/settings"
)

const (
	settingsKeyBase = "auth.external."

	settingsKeyProviders          = settingsKeyBase + "providers."
	settingsKeyRedirectUrl        = settingsKeyBase + "redirect-url"
	settingsKeySessionStoreSecret = settingsKeyBase + "session-store-secret"
	settingsKeySessionStoreSecure = settingsKeyBase + "session-store-secure"
)

type (
	externalAuthSettings struct {
		enabled            bool
		redirectUrl        string
		sessionStoreSecret string
		sessionStoreSecure bool
		providers          map[string]externalAuthProvider
	}

	externalAuthProvider struct {
		enabled     bool
		key         string
		secret      string
		redirectUrl string
		issuerUrl   string
	}
)

func ExternalAuthProvider(kv intset.KV) (eap externalAuthProvider, err error) {
	for k, v := range kv {
		ld := strings.LastIndex(k, ".")

		switch k[ld+1:] {
		case "enabled":
			err = v.Unmarshal(&eap.enabled)
		case "key":
			err = v.Unmarshal(&eap.key)
		case "secret":
			err = v.Unmarshal(&eap.secret)
		case "issuer":
			err = v.Unmarshal(&eap.issuerUrl)
		}

		if err != nil {
			return
		}
	}

	return
}

func (eap externalAuthProvider) MakeValueSet(name string) (vv intset.ValueSet, err error) {
	set := func(name string, value interface{}) error {
		v := &intset.Value{Name: name}
		if v.Value, err = json.Marshal(value); err != nil {
			return err
		}

		vv = append(vv, v)
		return nil
	}

	prefix := settingsKeyProviders + name

	if err = set(prefix+".enabled", eap.enabled); err != nil {
		return nil, err
	}

	if err = set(prefix+".key", eap.key); err != nil {
		return nil, err
	}

	if err = set(prefix+".secret", eap.secret); err != nil {
		return nil, err
	}

	if err = set(prefix+".issuer", eap.issuerUrl); err != nil {
		return nil, err
	}

	return vv, err
}

// ExternalAuthSettings maps from plain values to externalAuthSettings struct
//
// see settings.Initialize() func
func ExternalAuthSettings(s intset.Service) (eas *externalAuthSettings, err error) {
	// Read all settings and populate struct
	settings, err := s.FindByPrefix(settingsKeyBase)
	if err != nil {
		return nil, errors.Wrap(err, "could not load settings for external auth provider")
	}

	kv := settings.KV()

	eas = &externalAuthSettings{
		enabled:            kv.Bool(settingsKeyBase + "enabled"),
		redirectUrl:        kv.String(settingsKeyRedirectUrl),
		sessionStoreSecret: kv.String(settingsKeySessionStoreSecret),
	}

	if !kv.Has(settingsKeySessionStoreSecure) {
		// If auth.external.session-store-secure is not explicitly set;
		// check if redirectUrl uses HTTPS schema and assume we want secure session store
		eas.sessionStoreSecure = strings.Index(eas.redirectUrl, "https://") == 0
	} else {
		eas.sessionStoreSecure = kv.Bool(settingsKeySessionStoreSecure)
	}

	if eas.providers, err = extractProviders(eas.redirectUrl, kv); err != nil {
		return nil, err
	}

	return
}

func extractProviders(redirectUrl string, kv intset.KV) (providers map[string]externalAuthProvider, err error) {
	// Standard providers:
	providers = map[string]externalAuthProvider{
		"github":   {},
		"facebook": {},
		"gplus":    {},
		"linkedin": {},
	}

	oidcKeyBase := settingsKeyProviders + "openid-connect."
	for k := range kv.Filter(oidcKeyBase) {
		if len(k) < len(oidcKeyBase)+2 {
			// skip invalid keys
			continue
		}

		// find next dot:
		name := k[len(oidcKeyBase):]
		dotPos := strings.Index(name, ".")
		if dotPos > 0 {
			name = name[:dotPos]
		}

		providers["openid-connect."+name] = externalAuthProvider{}
	}

	for provider := range providers {
		if eap, err := ExternalAuthProvider(kv.Filter(settingsKeyProviders + provider)); err != nil {
			return nil, err
		} else {
			if eap.enabled {
				eap.redirectUrl = fmt.Sprintf(redirectUrl, provider)
			}

			providers[provider] = eap
		}
	}

	return
}
