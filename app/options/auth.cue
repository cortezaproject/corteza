package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
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
			type:    "bool"
			default: "true"
			description: """
				Password security allows you to disable constraints to which passwords must conform to.

				[CAUTION]
				====
				Disabling password security can be useful for development environments as it removes the need for complex passwords.
				Password security *should be enabled* on production environments to avoid security incidents
				====
				"""
		}
		secret: {
			default: "getSecretFromEnv(\"jwt secret\")"
			description: """
				Secret used for signing JWT tokens.

				[IMPORTANT]
				====
				If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
				Generated secret will change if you change any of these variables.
				====
				"""
			env: "AUTH_JWT_SECRET"
		}
		access_token_lifetime: {
			type:        "time.Duration"
			default:     "time.Hour * 2"
			description: "Access token lifetime"
			env:         "AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME"
		}
		refresh_token_lifetime: {
			type:        "time.Duration"
			default:     "time.Hour * 24 * 3"
			description: "Refresh token lifetime"
			env:         "AUTH_OAUTH2_REFRESH_TOKEN_LIFETIME"
		}
		expiry: {
			type:        "time.Duration"
			default:     "time.Hour * 24 * 30"
			description: "Experation time for the auth JWT tokens."
			env:         "AUTH_JWT_EXPIRY"
		}
		external_redirect_URL: {
			default: "fullURL(\"/auth/external/{provider}/callback\")"
			description: """
				Redirect URL to be sent with OAuth2 authentication request to provider

				`provider` placeholder is replaced with the actual value when used.
				"""
		}
		external_cookie_secret: {
			default: "getSecretFromEnv(\"external cookie secret\")"
			description: """
				Secret used for securing cookies

				[IMPORTANT]
				====
				If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
				Generated secret will change if you change any of these variables.
				====
				"""
		}
		base_URL: {
			default: "fullURL(\"/auth\")"
			description: """
				Frontend base URL. Must be an absolute URL, with the domain.
				This is used for some redirects and links in auth emails.
				"""
		}
		session_cookie_name: {
			default:     "\"session\""
			description: "Session cookie name"
		}
		session_cookie_path: {
			default:     "pathPrefix(\"/auth\")"
			description: "Session cookie path"
		}
		session_cookie_domain: {
			default:     "guessHostname()"
			description: "Session cookie domain"
		}
		session_cookie_secure: {
			type:        "bool"
			default:     "isSecure()"
			description: "Defaults to true when HTTPS is used. Corteza will try to guess the this setting by"
		}
		session_lifetime: {
			type:        "time.Duration"
			default:     "24 * time.Hour"
			description: "How long do we keep the temporary session"
		}
		session_perm_lifetime: {
			type:        "time.Duration"
			default:     "360 * 24 * time.Hour"
			description: "How long do we keep the permanent session"
		}
		garbage_collector_interval: {
			type:        "time.Duration"
			default:     "15 * time.Minute"
			description: "How often are expired sessions and tokens purged from the database"
		}
		request_rate_limit: {
			type:    "int"
			default: "60"
			description: """
				How many requests from a cerain IP address are allowed in a time window.
				Set to zero to disable
				"""
		}
		request_rate_window_length: {
			type:        "time.Duration"
			default:     "time.Minute"
			description: "How many requests from a cerain IP address are allowed in a time window"
		}
		csrf_secret: {
			default: "getSecretFromEnv(\"csrf secret\")"
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
			type:        "bool"
			default:     "true"
			description: "Enable CSRF protection"
		}
		csrf_field_name: {
			default:     "\"same-site-authenticity-token\""
			description: "Form field name used for CSRF protection"
		}
		csrf_cookie_name: {
			default:     "\"same-site-authenticity-token\""
			description: "Cookie name used for CSRF protection"
		}
		default_client: {
			default: "\"corteza-webapp\""
			description: """
				Handle for OAuth2 client used for automatic redirect from /auth/oauth2/go endpoint.

				This simplifies configuration for OAuth2 flow for Corteza Web applications as it removes
				the need to suply redirection URL and client ID (oauth2/go endpoint does that internally)

				"""
		}
		assets_path: {
			default: ""
			description: """
				Path to js, css, images and template source files

				When corteza starts, if path exists it tries to load template files from it.
				If not it uses statically embedded files.

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
