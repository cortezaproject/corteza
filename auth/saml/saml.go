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
	saml := source.Auth.External.Saml

	dest.Saml.Enabled = saml.Enabled
	dest.Saml.Name = saml.Name
	dest.Saml.Cert = saml.Cert
	dest.Saml.Key = saml.Key
	dest.Saml.SignRequests = saml.SignRequests
	dest.Saml.SignMethod = saml.SignMethod
	dest.Saml.Binding = saml.Binding

	dest.Saml.IDP.URL = saml.IDP.URL
	dest.Saml.IDP.IdentName = saml.IDP.IdentName
	dest.Saml.IDP.IdentHandle = saml.IDP.IdentHandle
	dest.Saml.IDP.IdentIdentifier = saml.IDP.IdentIdentifier
}
