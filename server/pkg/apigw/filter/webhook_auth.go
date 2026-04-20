package filter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/crusttech/human/server/pkg/apigw/types"
	pe "github.com/crusttech/human/server/pkg/errors"
)

type (
	webhookAuth struct {
		types.FilterMeta

		cfg types.Config

		params struct {
			// hmac | bearer | query | noop
			Type string `json:"type"`

			HMAC struct {
				Secret string `json:"secret"`
				Header string `json:"header"`
				// sha256 (default) | sha1
				Algorithm string `json:"algorithm"`
				Prefix    string `json:"prefix"`
			} `json:"hmac,omitempty"`
			Bearer struct {
				Secret string `json:"secret"`
				Header string `json:"header"`
				Prefix string `json:"prefix"`
			} `json:"bearer,omitempty"`
			Query struct {
				Secret string `json:"secret"`
				Param  string `json:"param"`
			} `json:"query,omitempty"`
		}
	}
)

func NewWebhookAuth(cfg types.Config) *webhookAuth {
	f := &webhookAuth{cfg: cfg}

	f.Name = "webhookAuth"
	f.Label = "Webhook auth"
	f.Kind = types.PreFilter

	f.Args = []*types.FilterMetaArg{
		{Type: "type", Label: "type", Options: map[string]interface{}{}},
	}

	return f
}

func (h webhookAuth) New(cfg types.Config) types.Handler {
	return NewWebhookAuth(cfg)
}

func (h webhookAuth) Enabled() bool { return true }

func (h webhookAuth) String() string {
	return fmt.Sprintf("apigw filter %s (%s)", h.Name, h.Label)
}

func (h webhookAuth) Meta() types.FilterMeta {
	return h.FilterMeta
}

func (h *webhookAuth) Merge(params []byte, cfg types.Config) (types.Handler, error) {
	h.cfg = cfg

	if err := json.NewDecoder(bytes.NewBuffer(params)).Decode(&h.params); err != nil {
		return nil, err
	}

	switch h.params.Type {
	case "hmac":
		if h.params.HMAC.Secret == "" {
			return nil, fmt.Errorf("webhookAuth: hmac.secret is required")
		}

		if h.params.HMAC.Header == "" {
			return nil, fmt.Errorf("webhookAuth: hmac.header is required")
		}

		if h.params.HMAC.Algorithm == "" {
			h.params.HMAC.Algorithm = "sha256"
		}

		if h.params.HMAC.Algorithm != "sha256" && h.params.HMAC.Algorithm != "sha1" {
			return nil, fmt.Errorf("webhookAuth: hmac.algorithm must be sha256 or sha1")
		}

	case "bearer":
		if h.params.Bearer.Secret == "" {
			return nil, fmt.Errorf("webhookAuth: bearer.secret is required")
		}

		if h.params.Bearer.Header == "" {
			return nil, fmt.Errorf("webhookAuth: bearer.header is required")
		}

	case "query":
		if h.params.Query.Secret == "" {
			return nil, fmt.Errorf("webhookAuth: query.secret is required")
		}

		if h.params.Query.Param == "" {
			return nil, fmt.Errorf("webhookAuth: query.param is required")
		}

	case "noop":
		// explicit opt-out — always allowed

	default:
		return nil, fmt.Errorf("webhookAuth: unknown type %q", h.params.Type)
	}

	return h, nil
}

func (h *webhookAuth) Handler() types.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) error {
		switch h.params.Type {
		case "noop":
			return nil

		case "hmac":
			return h.handleHMAC(r)

		case "bearer":
			return h.handleBearer(r)

		case "query":
			return h.handleQuery(r)
		}

		// impossible state as validation happens before
		return pe.Unauthorized("unauthorized")
	}
}

func (h *webhookAuth) handleHMAC(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return pe.Unauthorized("unauthorized")
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body))

	var mac hash.Hash
	switch h.params.HMAC.Algorithm {
	case "sha1":
		mac = hmac.New(sha1.New, []byte(h.params.HMAC.Secret))
	default:
		mac = hmac.New(sha256.New, []byte(h.params.HMAC.Secret))
	}
	mac.Write(body)
	computed := []byte(hex.EncodeToString(mac.Sum(nil)))

	provided := []byte(stripPrefix(r.Header.Get(h.params.HMAC.Header), h.params.HMAC.Prefix))

	if subtle.ConstantTimeCompare(computed, provided) != 1 {
		return pe.Unauthorized("unauthorized")
	}

	return nil
}

func (h *webhookAuth) handleBearer(r *http.Request) error {
	provided := []byte(stripPrefix(r.Header.Get(h.params.Bearer.Header), h.params.Bearer.Prefix))
	secret := []byte(h.params.Bearer.Secret)

	if subtle.ConstantTimeCompare(provided, secret) != 1 {
		return pe.Unauthorized("unauthorized")
	}

	return nil
}

func (h *webhookAuth) handleQuery(r *http.Request) error {
	provided := []byte(r.URL.Query().Get(h.params.Query.Param))
	secret := []byte(h.params.Query.Secret)

	if subtle.ConstantTimeCompare(provided, secret) != 1 {
		return pe.Unauthorized("unauthorized")
	}

	return nil
}

func stripPrefix(s, prefix string) string {
	if prefix == "" || len(s) < len(prefix) {
		return s
	}
	if s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
