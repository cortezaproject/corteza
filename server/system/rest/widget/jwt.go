package widget

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// sessionClaims is the minimal set of claims signed into the widget session JWT.
// HS256 over a secret derived from sha256(widgetKey + serverSecret) so key
// regeneration invalidates all live sessions for that agent.
type sessionClaims struct {
	Sid string `json:"sid"`         // session ID (random)
	Cid uint64 `json:"cid,string"`  // conversation ID
	Dbsid uint64 `json:"dbsid,string"` // database session ID
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

func deriveSessionSecret(widgetKey, serverSecret string) []byte {
	s := sha256.Sum256([]byte(widgetKey + "|" + serverSecret))
	return s[:]
}

func signSession(c sessionClaims, secret []byte) (string, error) {
	header := []byte(`{"alg":"HS256","typ":"JWT"}`)
	payload, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	he := base64.RawURLEncoding.EncodeToString(header)
	pe := base64.RawURLEncoding.EncodeToString(payload)
	signingInput := he + "." + pe

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig, nil
}

func verifySession(token string, secret []byte) (*sessionClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("widget: malformed token")
	}
	signingInput := parts[0] + "." + parts[1]

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, errors.New("widget: invalid signature")
	}

	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("widget: bad payload")
	}
	var c sessionClaims
	if err = json.Unmarshal(raw, &c); err != nil {
		return nil, errors.New("widget: bad payload")
	}
	if time.Now().Unix() >= c.Exp {
		return nil, errors.New("widget: expired")
	}
	return &c, nil
}
