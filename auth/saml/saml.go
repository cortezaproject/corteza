package saml

import (
	"context"
	"net/http"
	"net/url"

	"github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

type (
	templateProvider struct {
		Label, Handle, Icon string
	}
)

// FetchIDPMetadata loads the idp metadata, usually the url
// is configured in settings
func FetchIDPMetadata(ctx context.Context, u url.URL) (*saml.EntityDescriptor, error) {
	return samlsp.FetchMetadata(ctx, http.DefaultClient, u)
}

// TemplateProvider adds a wrapper to the button
// data that is displayed on the login form
func TemplateProvider(url, name string) templateProvider {
	if name == "" {
		name = url
	}

	return templateProvider{
		Label:  name,
		Handle: "saml/init",
		Icon:   "key",
	}
}

// UpdateSettings applies the app settings to the
// auth specific settings
func UpdateSettings(source *types.AppSettings, dest *settings.Settings) {
	cas := source.Auth

	dest.Saml.Enabled = cas.External.Saml.Enabled
	dest.Saml.Cert = cas.External.Saml.Cert
	dest.Saml.Key = cas.External.Saml.Key

	dest.Saml.IDP.URL = cas.External.Saml.IDP.URL
	dest.Saml.IDP.Name = cas.External.Saml.IDP.Name
	dest.Saml.IDP.IdentName = cas.External.Saml.IDP.IdentName
	dest.Saml.IDP.IdentHandle = cas.External.Saml.IDP.IdentHandle
	dest.Saml.IDP.IdentIdentifier = cas.External.Saml.IDP.IdentIdentifier
}
