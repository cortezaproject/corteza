package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

auth: schema.#optionsGroup & {
	handle: "auth"

	imports: [
		"\"time\"",
	]

	options: {
		log_enabled: {
			type:        "bool"
			description: "Enable extra logging for authentication flows"
		}
		password_security: {
			type:          "bool"
			defaultGoExpr: "true"
			description: """
				Password security allows you to disable constraints to which passwords must conform to.

				[CAUTION]
				====
				Disabling password security can be useful for development environments as it removes the need for complex passwords.
				Password security *should be enabled* on production environments to avoid security incidents
				====
				"""
		}
		jwt_algorithm: {
			defaultGoExpr: "\"HS512\""
			defaultValue:  "HS512"
			description: """
				Algoritm to be use for JWT signature.

				Supported valus:
				 - HS256, HS384, HS512
				 - PS256, PS384, PS512,
				 - RS256, RS384, RS512

				Provide shared secret string for HS256, HS384, HS512 and full private key or path to the file PS* and RS* algorithms.
				"""
		}
		secret: {
			defaultGoExpr: "getSecretFromEnv(\"jwt secret\")"
			description: """
				Secret used for signing JWT tokens.
				Value is used only when HS256, HS384 or HS512 algorithm is used.

				[IMPORTANT]
				====
				If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
				Generated secret will change if you change any of these variables.
				====
				"""
			env: "AUTH_JWT_SECRET"
		}
		jwt_key: {
			description: """
				Raw private key or absolute or relative path to the file containing one.
				"""
		}
		access_token_lifetime: {
			type: "time.Duration"
			description: """
				Lifetime of the access token. Should be shorter than lifetime of the refresh token.
				"""
			env: "AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME"

			defaultGoExpr: "time.Hour * 2"
			defaultValue:  "2h"
		}
		refresh_token_lifetime: {
			type: "time.Duration"
			description: """
				Lifetime of the refresh token. Should be much longer than lifetime of the access token.

				Refresh tokens are used to exchange expired access tokens with new ones.
				"""
			env: "AUTH_OAUTH2_REFRESH_TOKEN_LIFETIME"

			defaultGoExpr: "time.Hour * 24 * 3"
			defaultValue:  "72h"
		}
		external_redirect_URL: {
			description: """
				Redirect URL to be sent with OAuth2 authentication request to provider

				`provider` placeholder is replaced with the actual value when used.
				"""
			defaultGoExpr: "fullURL(\"/auth/external/{provider}/callback\")"
		}
		external_cookie_secret: {
			description: """
				Secret used for securing cookies

				[IMPORTANT]
				====
				If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
				Generated secret will change if you change any of these variables.
				====
				"""

			defaultGoExpr: "getSecretFromEnv(\"external cookie secret\")"
		}
		base_URL: {
			description: """
				Frontend base URL. Must be an absolute URL, with the domain.
				This is used for some redirects and links in auth emails.
				"""

			defaultGoExpr: "fullURL(\"/auth\")"
		}
		session_cookie_name: {
			description:  "Session cookie name"
			defaultValue: "session"
		}
		session_cookie_path: {
			description:   "Session cookie path"
			defaultGoExpr: "pathPrefix(\"/auth\")"
		}
		session_cookie_domain: {
			defaultGoExpr: "guessHostname()"
			description:   "Session cookie domain"
		}
		session_cookie_secure: {
			type:          "bool"
			defaultGoExpr: "isSecure()"
			description:   "Defaults to true when HTTPS is used. Corteza will try to guess the this setting by"
		}
		session_lifetime: {
			type: "time.Duration"
			description: """
				Maximum time user is allowed to stay idle when logged in without \"remember-me\" option and before session is expired.

				Recomended value is between an hour and a day.

				[IMPORTANT]
				====
				This affects only profile (/auth) pages. Using applications (admin, compose, ...) does not prolong the session.
				====

				"""
			defaultGoExpr: "24 * time.Hour"
			defaultValue:  "24h"
		}
		session_perm_lifetime: {
			type: "time.Duration"
			description: """
				Duration of the session in /auth lasts when user logs-in with \"remember-me\" option.

				If set to 0, \"remember-me\" option is removed.
				"""

			defaultGoExpr: "360 * 24 * time.Hour"
			defaultValue:  "8640h"
		}
		garbage_collector_interval: {
			type:          "time.Duration"
			description:   "How often are expired sessions and tokens purged from the database"
			defaultGoExpr: "15 * time.Minute"
			defaultValue:  "15min"
		}
		request_rate_limit: {
			type: "int"
			description: """
				How many requests from a cerain IP address are allowed in a time window.
				Set to zero to disable
				"""
			defaultGoExpr: "60"
			defaultValue:  "60"
		}
		request_rate_window_length: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute"
			defaultValue:  "1m"
			description:   "How many requests from a cerain IP address are allowed in a time window"
		}
		csrf_secret: {
			defaultGoExpr: "getSecretFromEnv(\"csrf secret\")"
			description: """
				Secret used for securing CSRF protection

				[IMPORTANT]
				====
				If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
				Generated secret will change if you change any of these variables.
				====
				"""
		}
		csrf_enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable CSRF protection"
		}
		csrf_field_name: {
			defaultValue: "same-site-authenticity-token"
			description:  "Form field name used for CSRF protection"
		}
		csrf_cookie_name: {
			defaultValue: "same-site-authenticity-token"
			description:  "Cookie name used for CSRF protection"
		}
		default_client: {
			defaultValue: "corteza-webapp"
			description: """
				Handle for OAuth2 client used for automatic redirect from /auth/oauth2/go endpoint.

				This simplifies configuration for OAuth2 flow for Corteza Web applications as it removes
				the need to suply redirection URL and client ID (oauth2/go endpoint does that internally)

				"""
		}
		assets_path: {
			description: """
				Path to js, css, images and template source files

				When corteza starts, if path exists it tries to load template files from it.

				When empty path is set (default value), embedded files are used.
				"""
		}
		development_mode: {
			type: "bool"
			description: """
				When enabled, corteza reloads template before every execution.
				Enable this for debugging or when developing auth templates.

				Should be disabled in production where templates do not change between server restarts.
				"""
		}
	}
	title: "Authentication"
}
