package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/extra/server-discovery/pkg/options"
)

type (
	credentials struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		expiresAt   time.Time

		authBaseUri string
		key         string
		secret      string
	}

	client struct {
		baseUri          string
		discoveryBaseUrl string
		credentials      *credentials
	}

	ClientService interface {
		HttpClient() *http.Client
		Mappings() (*http.Request, error)
		Feed(url.Values) (*http.Request, error)
		Resources(string, url.Values) (*http.Request, error)
		Namespaces() (*http.Request, error)
		Modules(uint64) (*http.Request, error)
		Request(string) (*http.Request, error)
		Authenticate() error
	}
)

func Client(opt options.CortezaOpt, key, secret string) (c *client, err error) {
	c = &client{baseUri: opt.BaseUrl, discoveryBaseUrl: opt.DiscoveryUrl}
	c.credentials = &credentials{
		authBaseUri: opt.AuthUrl,
		key:         key,
		secret:      secret,
	}
	return c, err
}

func (*client) HttpClient() *http.Client {
	return http.DefaultClient
}

func (c *client) Mappings() (*http.Request, error) {
	return c.Request(fmt.Sprintf("%s/mappings/", c.discoveryBaseUrl))
}

func (c *client) Feed(qs url.Values) (*http.Request, error) {
	query := ""
	if len(qs.Encode()) > 0 {
		query = fmt.Sprintf("from=%s", qs.Get("from"))
	}
	return c.Request(fmt.Sprintf("%s/feed/?%s", c.discoveryBaseUrl, query))
}

func (c *client) Resources(endpoint string, qs url.Values) (*http.Request, error) {
	return c.Request(fmt.Sprintf("%s/resources/%s?%s", c.discoveryBaseUrl, strings.TrimLeft(endpoint, "/"), qs.Encode()))
}

func (c *client) Namespaces() (*http.Request, error) {
	return c.Request(fmt.Sprintf("%s/api/compose/namespace/", c.baseUri))
}

func (c *client) Modules(namespaceID uint64) (*http.Request, error) {
	return c.Request(fmt.Sprintf("%s/api/compose/namespace/%d/module/?sort=name+ASC", c.baseUri, namespaceID))
}

func (c *client) Request(endpoint string) (req *http.Request, err error) {
	if err = c.Authenticate(); err != nil {
		return
	}

	if req, err = http.NewRequest(http.MethodGet, endpoint, nil); err != nil {
		return
	}

	req.Header.Set("User-Agent", "corteza-server-discovery/0.1")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.credentials.AccessToken))
	return
}

func (c *client) Authenticate() (err error) {
	if c.credentials == nil {
		return fmt.Errorf("missing credentials")
	}

	tokenNeedsRefresh := c.credentials.expiresAt.IsZero() || time.Now().After(c.credentials.expiresAt.Add(-60*time.Second))

	if tokenNeedsRefresh {
		newCredentials, err := c.authToken()
		if err != nil {
			return fmt.Errorf("failed to refresh token: %w", err)
		}

		c.credentials = newCredentials
	}

	return nil
}

func (c *client) authToken() (crd *credentials, err error) {
	var (
		req         *http.Request
		rsp         *http.Response
		form        = url.Values{}
		authBaseUri = c.credentials.authBaseUri
		key         = c.credentials.key
		secret      = c.credentials.secret
	)

	form.Set("grant_type", "client_credentials")
	form.Set("scope", "profile api discovery")

	req, err = http.NewRequest(http.MethodPost, authBaseUri+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	req.SetBasicAuth(key, secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//d, _ := httputil.DumpRequest(req, true)
	//println(string(d))

	rsp, err = c.HttpClient().Do(req)
	if err != nil {
		return
	}

	defer rsp.Body.Close()
	crd = &credentials{
		authBaseUri: authBaseUri,
		key:         key,
		secret:      secret,
	}

	if rsp.StatusCode != http.StatusOK {
		aux := struct{ Error string }{}
		if err = json.NewDecoder(rsp.Body).Decode(&aux); err != nil {
			return
		} else if aux.Error != "" {
			return nil, fmt.Errorf(aux.Error)
		} else {
			return nil, fmt.Errorf("can not authenticate, unexpected error")
		}

	}

	//d, _ := httputil.DumpResponse(rsp, true)
	//println(string(d))

	err = json.NewDecoder(rsp.Body).Decode(crd)
	if err != nil {
		return
	}

	crd.expiresAt, err = parseJWTExpiry(crd.AccessToken)
	if err != nil {
		crd.expiresAt = time.Now().Add(time.Second * time.Duration(crd.ExpiresIn))
	}

	return
}

func parseJWTExpiry(tokenString string) (time.Time, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	// Parse exp claim
	var claims struct {
		Exp int64 `json:"exp"`
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, fmt.Errorf("failed to unmarshal JWT claims: %w", err)
	}

	if claims.Exp == 0 {
		return time.Time{}, fmt.Errorf("JWT missing exp claim")
	}

	return time.Unix(claims.Exp, 0), nil
}
