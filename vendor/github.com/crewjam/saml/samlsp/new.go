// Package samlsp provides helpers that can be used to protect web services using SAML.
package samlsp

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"net/http"
	"net/url"
	"time"

	dsig "github.com/russellhaering/goxmldsig"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"
)

// Options represents the parameters for creating a new middleware
type Options struct {
	EntityID          string
	URL               url.URL
	Key               *rsa.PrivateKey
	Certificate       *x509.Certificate
	Intermediates     []*x509.Certificate
	AllowIDPInitiated bool
	IDPMetadata       *saml.EntityDescriptor
	SignRequest       bool
	ForceAuthn        bool // TODO(ross): this should be *bool
	CookieSameSite    http.SameSite

	// The following fields exist <= 0.3.0, but are superceded by the new
	// SessionProvider and RequestTracker interfaces.
	Logger         logger.Interface // DEPRECATED: this field will be removed, instead provide a custom OnError function to handle errors
	IDPMetadataURL *url.URL         // DEPRECATED: this field will be removed, instead use FetchMetadata
	HTTPClient     *http.Client     // DEPRECATED: this field will be removed, instead pass httpClient to FetchMetadata
	CookieMaxAge   time.Duration    // DEPRECATED: this field will be removed. Instead, assign a custom CookieRequestTracker or CookieSessionProvider
	CookieName     string           // DEPRECATED: this field will be removed. Instead, assign a custom CookieRequestTracker or CookieSessionProvider
	CookieDomain   string           // DEPRECATED: this field will be removed. Instead, assign a custom CookieRequestTracker or CookieSessionProvider
	CookieSecure   bool             // DEPRECATED: this field will be removed, the Secure flag is set on cookies when the root URL uses the https scheme
}

// DefaultSessionCodec returns the default SessionCodec for the provided options,
// a JWTSessionCodec configured to issue signed tokens.
func DefaultSessionCodec(opts Options) JWTSessionCodec {
	// for backwards compatibility, support CookieMaxAge
	maxAge := defaultSessionMaxAge
	if opts.CookieMaxAge > 0 {
		maxAge = opts.CookieMaxAge
	}

	return JWTSessionCodec{
		SigningMethod: defaultJWTSigningMethod,
		Audience:      opts.URL.String(),
		Issuer:        opts.URL.String(),
		MaxAge:        maxAge,
		Key:           opts.Key,
	}
}

// DefaultSessionProvider returns the default SessionProvider for the provided options,
// a CookieSessionProvider configured to store sessions in a cookie.
func DefaultSessionProvider(opts Options) CookieSessionProvider {
	// for backwards compatibility, support CookieMaxAge
	maxAge := defaultSessionMaxAge
	if opts.CookieMaxAge > 0 {
		maxAge = opts.CookieMaxAge
	}

	// for backwards compatibility, support CookieName
	cookieName := defaultSessionCookieName
	if opts.CookieName != "" {
		cookieName = opts.CookieName
	}

	// for backwards compatibility, support CookieDomain
	cookieDomain := opts.URL.Host
	if opts.CookieDomain != "" {
		cookieDomain = opts.CookieDomain
	}

	// for backwards compatibility, support CookieDomain
	cookieSecure := opts.URL.Scheme == "https"
	if opts.CookieSecure {
		cookieSecure = true
	}

	return CookieSessionProvider{
		Name:     cookieName,
		Domain:   cookieDomain,
		MaxAge:   maxAge,
		HTTPOnly: true,
		Secure:   cookieSecure,
		SameSite: opts.CookieSameSite,
		Codec:    DefaultSessionCodec(opts),
	}
}

// DefaultTrackedRequestCodec returns a new TrackedRequestCodec for the provided
// options, a JWTTrackedRequestCodec that uses a JWT to encode TrackedRequests.
func DefaultTrackedRequestCodec(opts Options) JWTTrackedRequestCodec {
	return JWTTrackedRequestCodec{
		SigningMethod: defaultJWTSigningMethod,
		Audience:      opts.URL.String(),
		Issuer:        opts.URL.String(),
		MaxAge:        saml.MaxIssueDelay,
		Key:           opts.Key,
	}
}

// DefaultRequestTracker returns a new RequestTracker for the provided options,
// a CookieRequestTracker which uses cookies to track pending requests.
func DefaultRequestTracker(opts Options, serviceProvider *saml.ServiceProvider) CookieRequestTracker {
	return CookieRequestTracker{
		ServiceProvider: serviceProvider,
		NamePrefix:      "saml_",
		Codec:           DefaultTrackedRequestCodec(opts),
		MaxAge:          saml.MaxIssueDelay,
		SameSite:        opts.CookieSameSite,
	}
}

// DefaultServiceProvider returns the default saml.ServiceProvider for the provided
// options.
func DefaultServiceProvider(opts Options) saml.ServiceProvider {
	metadataURL := opts.URL.ResolveReference(&url.URL{Path: "saml/metadata"})
	acsURL := opts.URL.ResolveReference(&url.URL{Path: "saml/acs"})
	sloURL := opts.URL.ResolveReference(&url.URL{Path: "saml/slo"})

	var forceAuthn *bool
	if opts.ForceAuthn {
		forceAuthn = &opts.ForceAuthn
	}
	signatureMethod := dsig.RSASHA1SignatureMethod
	if !opts.SignRequest {
		signatureMethod = ""
	}

	return saml.ServiceProvider{
		EntityID:          opts.EntityID,
		Key:               opts.Key,
		Certificate:       opts.Certificate,
		Intermediates:     opts.Intermediates,
		MetadataURL:       *metadataURL,
		AcsURL:            *acsURL,
		SloURL:            *sloURL,
		IDPMetadata:       opts.IDPMetadata,
		ForceAuthn:        forceAuthn,
		SignatureMethod:   signatureMethod,
		AllowIDPInitiated: opts.AllowIDPInitiated,
	}
}

// New creates a new Middleware with the default providers for the
// given options.
//
// You can customize the behavior of the middleware in more detail by
// replacing and/or changing Session, RequestTracker, and ServiceProvider
// in the returned Middleware.
func New(opts Options) (*Middleware, error) {
	// for backwards compatibility, support Logger
	onError := DefaultOnError
	if opts.Logger != nil {
		onError = defaultOnErrorWithLogger(opts.Logger)
	}

	// for backwards compatibility, support IDPMetadataURL
	if opts.IDPMetadataURL != nil && opts.IDPMetadata == nil {
		httpClient := opts.HTTPClient
		if httpClient == nil {
			httpClient = http.DefaultClient
		}
		metadata, err := FetchMetadata(context.TODO(), httpClient, *opts.IDPMetadataURL)
		if err != nil {
			return nil, err
		}
		opts.IDPMetadata = metadata
	}

	m := &Middleware{
		ServiceProvider: DefaultServiceProvider(opts),
		Binding:         "",
		OnError:         onError,
		Session:         DefaultSessionProvider(opts),
	}
	m.RequestTracker = DefaultRequestTracker(opts, &m.ServiceProvider)

	return m, nil
}
