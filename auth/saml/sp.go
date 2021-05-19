package saml

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"strings"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/pkg/errors"
)

const defaultNameIdentifier = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"

type (
	SamlSPService struct {
		IdpURL      url.URL
		Host        url.URL
		IDPUserMeta *IdpIdentityPayload
		IDPMeta     *saml.EntityDescriptor

		sp      saml.ServiceProvider
		handler *samlsp.Middleware
	}

	SamlSPArgs struct {
		AcsURL  string
		MetaURL string
		SloURL  string

		// user meta from idp
		IdentityPayload IdpIdentityPayload

		IdpURL  url.URL
		Host    url.URL
		Cert    tls.Certificate
		IdpMeta *saml.EntityDescriptor
	}
)

// NewSamlSPService loads the certificates and registers the
// already fetched IDP metadata into the SAML middleware
func NewSamlSPService(args SamlSPArgs) (s *SamlSPService, err error) {
	var (
		keyPair = args.Cert
	)

	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])

	if err != nil {
		return
	}

	metadataURL, _ := url.Parse(args.MetaURL)
	acsURL, _ := url.Parse(args.AcsURL)
	logoutURL, _ := url.Parse(args.SloURL)

	sp := saml.ServiceProvider{
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: args.IdpMeta,

		MetadataURL: *args.Host.ResolveReference(metadataURL),
		AcsURL:      *args.Host.ResolveReference(acsURL),
		SloURL:      *args.Host.ResolveReference(logoutURL),
	}

	opts := samlsp.Options{
		URL:         args.Host,
		Key:         sp.Key,
		Certificate: sp.Certificate,
		IDPMetadata: args.IdpMeta,
	}

	// internal samlsp service
	handler, err := samlsp.New(opts)

	if err != nil {
		err = errors.Wrap(err, "could not init SAML SP handler")
		return
	}

	handler.RequestTracker = samlsp.DefaultRequestTracker(opts, &handler.ServiceProvider)
	handler.ServiceProvider = sp

	s = &SamlSPService{
		sp:      sp,
		handler: handler,

		IdpURL:      args.IdpURL,
		Host:        args.Host,
		IDPUserMeta: &args.IdentityPayload,
	}

	return
}

func (ssp *SamlSPService) NameIdentifier() string {
	return strings.TrimPrefix(defaultNameIdentifier, "urn:oasis:names:tc:SAML:1.1:nameid-format:")
}

// GuessIdentifier tries to guess the necessary (email) key
// for external authentication
func (ssp *SamlSPService) GuessIdentifier(payload map[string][]string) string {
	tryValues := []string{
		ssp.IDPUserMeta.Identifier,
		ssp.NameIdentifier(),
		defaultNameIdentifier,
		"urn:oasis:names:tc:SAML:attribute:subject-id",
		"email",
		"mail",
	}

	for _, v := range tryValues {
		if _, ok := payload[v]; ok {
			return payload[v][0]
		}
	}

	return ""
}

func (ssp *SamlSPService) Handler() *samlsp.Middleware {
	return ssp.handler
}

// ServeHTTP enables us to use the service directly
// in the router
func (ssp SamlSPService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ssp.handler.ServeHTTP(w, r)
}
