package request

import (
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/gorilla/sessions"
	"net/url"
)

const (
	keyPermanent              = "permanent"
	keyOriginalSession        = "originalSession"
	keyAuthUser               = "authUser"
	keyRoles                  = "roles"
	keyOAuth2AuthParams       = "oauth2AuthParams"
	keyOAuth2Client           = "oauth2ClientID"
	keyOAuth2ClientAuthorized = "oauth2ClientAuthorized"
)

// GetUser is wrapper to get value from session
func GetAuthUser(ses *sessions.Session) *authUser {
	val, has := ses.Values[keyAuthUser]
	if !has {
		return nil
	}

	return val.(*authUser)
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
