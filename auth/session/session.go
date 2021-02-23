package session

import (
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/gorilla/sessions"
	"net/url"
	"time"
)

const (
	keyPermanent              = "permanent"
	keyOriginalSession        = "originalSession"
	keySessionKind            = "sessionKind"
	keyUser                   = "user"
	keyRoles                  = "roles"
	keyOAuth2AuthParams       = "oauth2AuthParams"
	keyOAuth2Client           = "oauth2ClientID"
	keyOAuth2ClientAuthorized = "oauth2ClientAuthorized"
)

// GetUser is wrapper to get value from session
func GetUser(ses *sessions.Session) *types.User {
	val, has := ses.Values[keyUser]
	if !has {
		return nil
	}

	return val.(*types.User)
}

// SetUser is a session value setting wrapper for User
func SetUser(ses *sessions.Session, val *types.User) {
	if val != nil {
		ses.Values[keyUser] = val
	} else {
		delete(ses.Values, keyUser)
	}
}

// GetRoleMemberships is wrapper to get value from session
func GetRoleMemberships(ses *sessions.Session) []uint64 {
	val, has := ses.Values[keyRoles]
	if !has {
		return nil
	}

	return val.([]uint64)
}

// SetRoleMemberships is a session value setting wrapper for RoleMemberships
func SetRoleMemberships(ses *sessions.Session, val []uint64) {
	if val != nil {
		ses.Values[keyRoles] = val
	} else {
		delete(ses.Values, keyRoles)
	}
}

// GetOAuth2AuthParams is wrapper to get value from session
func GetOAuth2AuthParams(ses *sessions.Session) url.Values {
	val, has := ses.Values[keyOAuth2AuthParams]
	if !has {
		return nil
	}

	return val.(url.Values)
}

// SetOauth2AuthParams is a session value setting wrapper for Oauth2AuthParams
func SetOauth2AuthParams(ses *sessions.Session, val url.Values) {
	if val != nil {
		ses.Values[keyOAuth2AuthParams] = val
	} else {
		delete(ses.Values, keyOAuth2AuthParams)
	}
}

// GetOauth2Client is wrapper to get value from session
func GetOauth2Client(ses *sessions.Session) *types.AuthClient {
	val, has := ses.Values[keyOAuth2Client]
	if !has {
		return nil
	}

	return val.(*types.AuthClient)
}

// SetOauth2Client is a session value setting wrapper for Oauth2Client
func SetOauth2Client(ses *sessions.Session, val *types.AuthClient) {
	if val != nil {
		ses.Values[keyOAuth2Client] = val
	} else {
		delete(ses.Values, keyOAuth2Client)
	}
}

// IsOauth2ClientAuthorized is wrapper to get value from session
func IsOauth2ClientAuthorized(ses *sessions.Session) bool {
	val, has := ses.Values[keyOAuth2ClientAuthorized]
	if !has {
		return false
	}

	return val.(bool)
}

// SetOauth2ClientAuthorized is a session value setting wrapper for Oauth2ClientAuthorized
func SetOauth2ClientAuthorized(ses *sessions.Session, val bool) {
	if val {
		ses.Values[keyOAuth2ClientAuthorized] = true
	} else {
		delete(ses.Values, keyOAuth2ClientAuthorized)
	}
}

// IsOauth2ClientAuthorized is wrapper to get value from session
func IsPerm(ses *sessions.Session) bool {
	_, has := ses.Values[keyPermanent]
	return has
}

// SetOauth2ClientAuthorized is a session value setting wrapper for Oauth2ClientAuthorized
func SetPerm(ses *sessions.Session, ttl time.Duration) {
	if ttl > 0 {
		ses.Options.MaxAge = int(ttl / time.Second)
		ses.Values[keyPermanent] = true
	} else {
		ses.Options.MaxAge = 0
		delete(ses.Values, keyPermanent)
	}
}
