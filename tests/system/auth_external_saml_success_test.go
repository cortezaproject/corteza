package system

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/saml"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	s "github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/steinfletcher/apitest"
)

func loadSAMLService(ctx context.Context) (srvc *saml.SamlSPService, err error) {
	links := handlers.GetLinks()

	certManager := saml.NewCertManager(&saml.CertStoreLoader{})

	cert, err := certManager.Parse(
		readStaticFile("static/spCert.cert"),
		readStaticFile("static/spCert.key"))

	if err != nil {
		return
	}

	idpUrl, err := url.Parse("")
	if err != nil {
		return
	}

	// idp metadata needs to be loaded before
	// the internal samlsp package
	md, err := samlsp.ParseMetadata(readStaticFile("static/idp_metadata.xml"))
	ru, err := url.Parse("http://localhost:8084")

	rootURL := &url.URL{
		Scheme: ru.Scheme,
		User:   ru.User,
		Host:   ru.Host,
	}

	if err != nil {
		return
	}

	srvc, err = saml.NewSamlSPService(saml.SamlSPArgs{
		AcsURL:  links.SamlCallback,
		MetaURL: links.SamlMetadata,
		SloURL:  links.SamlLogout,

		IdpURL: *idpUrl,
		Host:   *rootURL,

		Cert:    cert,
		IdpMeta: md,
	})

	srvc.Handler().ServiceProvider.AllowIDPInitiated = true

	return
}

func TestAuthExternalSAMLSuccess(t *testing.T) {
	var (
		h = newHelper(t)

		cookieSessionIDPtoSP = apitest.
					NewCookie("saml_tCu5PV6EgxcvUAa9e57uJ2g-bTkqnNkyyHHaOu15yEfZjgWKt02AtXGe").
					Value(string(readStaticFile("static/idp_to_sp.cookie")))

		cookieTokenIDPtoSPAfterLogin = apitest.
						NewCookie("token").
						Value(string(readStaticFile("static/idp_to_sp_token.cookie")))
	)

	s.MaxClockSkew = time.Hour
	s.MaxIssueDelay = time.Hour

	jwt.TimeFunc = func() time.Time {
		tm, _ := time.Parse("2006-01-2 15:04:05", "2021-05-17 09:17:10")
		return tm
	}

	s.TimeNow = func() time.Time {
		tm, _ := time.Parse("2006-01-2 15:04:05", "2021-05-17 09:17:10")
		return tm
	}

	// first step, there is no session cookie, redirect to idp
	// in this case, host from parsed metadata
	t.Log("start login process")

	h.apiInit().
		Get(handlers.GetLinks().SamlInit).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			loc, _ := res.Location()

			h.assertBody("SSO Error: saml: session not present", res.Body)
			h.a.NotEmpty(loc.Query().Get("RelayState"))
			h.a.NotEmpty(loc.Query().Get("SAMLRequest"))
			h.a.Equal(http.StatusFound, res.StatusCode)
			return nil
		}).
		End()

	cookies := []*apitest.Cookie{}

	// coming back from idp, posting to sp-related endpoint
	// mocking session cookie and saml response
	// if everything is ok, redirect back to SAML init
	t.Log("post from idp to sp")

	h.apiInit().
		Post(handlers.GetLinks().SamlCallback).
		Header("Content-Type", "application/x-www-form-urlencoded").
		FormData("RelayState", "tCu5PV6EgxcvUAa9e57uJ2g-bTkqnNkyyHHaOu15yEfZjgWKt02AtXGe").
		Cookies(cookieSessionIDPtoSP).
		FormData("SAMLResponse", string(readStaticFile("static/idp_to_sp.post"))).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			loc, _ := res.Location()

			h.a.NotNil(getSessionCookie("token", res.Cookies()...))
			h.a.Equal(http.StatusFound, res.StatusCode)
			h.a.Equal(handlers.GetLinks().SamlInit, loc.String())

			cookies = append(cookies, cookieSessionIDPtoSP)
			return nil
		}).
		End()

	// idp sends a token session cookie also
	cookies = append(cookies, cookieTokenIDPtoSPAfterLogin)

	// once everything is set and the external authentication via
	// internal Corteza services is done, redirect to default path (profile)
	t.Log("redirect to profile after session is created")

	h.apiInit().
		Get(handlers.GetLinks().SamlInit).
		Cookies(cookies...).
		Expect(t).
		Assert(func(res *http.Response, req *http.Request) error {
			loc, _ := res.Location()

			ss, _, err := service.DefaultStore.SearchAuthSessions(context.Background(), types.AuthSessionFilter{})
			h.a.NoError(err)
			h.a.Len(ss, 1)

			h.a.NotNil(getSessionCookie("session", res.Cookies()...))
			h.a.Equal(http.StatusSeeOther, res.StatusCode)
			h.a.Equal(handlers.GetLinks().Profile, loc.String())
			return nil
		}).
		End()
}

func getSessionCookie(name string, cc ...*http.Cookie) (found *http.Cookie) {
	for _, c := range cc {
		if c.Name == name {
			found = c
			return
		}
	}

	return
}
