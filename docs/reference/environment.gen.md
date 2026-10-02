---
title: Environment variables
description: Every server option, with its type and default.
outline: [2, 2]
---

<!-- This file is auto-generated from server/app/options. -->

# Environment variables

The server reads these options from the environment, or from a `.env` file
next to the binary or passed with `--env-file`. Environment variables always
win.

## Connection to data store backend

### `DB_DSN` {#DB_DSN}

<div class="env-meta">

**Type** `string` · **Default** `sqlite3://file::memory:?cache=shared&mode=memory`

</div>

Database connection string.

Connection pool and retry settings can be added to the query string of the DSN.
They are read by the server and removed before the DSN is passed to the database driver,
so they can be combined with the driver's own parameters (for example `postgres://user:pass@host/human?sslmode=disable&*connMaxOpen=50`).
The same parameters work in the DSN of a DAL connection.

* `*connMaxOpen` (default `256`): maximum number of open connections to the database. Keep it below the server's connection limit (PostgreSQL defaults to 100).
* `*connMaxIdle` (default `32`): maximum number of idle connections kept in the pool.
* `*connMaxLifetime` (default `10m`): maximum time a connection may be reused.
* `*connTryTimeout` (default `30s`): timeout for a single connection attempt on startup.
* `*connTryBackoffDelay` (default `10s`): delay between failed connection attempts.
* `*connMaxTries` (default `99`): number of connection attempts before giving up.
* `*connTryPatience` (default `0`): time window in which failed connection attempts are not logged as errors.

## HTTP Client

### `HTTP_CLIENT_RESTRICTED_NETWORKS` {#HTTP_CLIENT_RESTRICTED_NETWORKS}

<div class="env-meta">

**Type** `string` · **Default** `link-local`

</div>

Networks that can not be reached by HTTP requests to destinations controlled by users
(workflow HTTP request and OAuth2 functions, API gateway proxy).

Supported values:
 - link-local: blocks link-local addresses and well known cloud metadata endpoints (default)
 - private: additionally blocks loopback, private and other non-public addresses
 - none: requests can reach any address

::: info
Restrictions are checked against the resolved IP address when connecting.
When outbound HTTP proxy is used, restrictions need to be enforced on the proxy.
:::

### `HTTP_CLIENT_TIMEOUT` {#HTTP_CLIENT_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `30s`

</div>

Default timeout for clients.

### `HTTP_CLIENT_TLS_INSECURE` {#HTTP_CLIENT_TLS_INSECURE}

<div class="env-meta">

**Type** `bool`

</div>

Allow insecure (invalid, expired TLS/SSL certificates) connections.

::: warning
We strongly recommend keeping this value set to false except for local development or demos.
:::

## HTTP Server

### `HTTP_ADDR` {#HTTP_ADDR}

<div class="env-meta">

**Type** `string` · **Default** `:80`

</div>

IP and port for the HTTP server.

### `HTTP_API_BASE_URL` {#HTTP_API_BASE_URL}

<div class="env-meta">

**Type** `string` · **Default** `/`

</div>

When webapps are enabled (HTTP_WEBAPP_ENABLED) this is moved to '/api' if not explicitly set otherwise.
API base URL is internally prefixed with baseUrl

### `HTTP_API_ENABLED` {#HTTP_API_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

### `HTTP_SERVER_ASSETS_PATH` {#HTTP_SERVER_ASSETS_PATH}

<div class="env-meta">

**Type** `string`

</div>

Human will directly serve these assets (static files).
When empty path is set (default value), embedded files are used.

### `HTTP_BASE_URL` {#HTTP_BASE_URL}

<div class="env-meta">

**Type** `string` · **Default** `/`

</div>

Base URL (prefix) for all routes (&lt;baseUrl>/auth, &lt;baseUrl>/api, ...)

### `HTTP_CORS_ALLOWED_ORIGINS` {#HTTP_CORS_ALLOWED_ORIGINS}

<div class="env-meta">

**Type** `string` · **Default** `http://*,https://*`

</div>

Comma separated list of origins allowed to make cross-origin requests (CORS).
Each origin can contain one wildcard (https://*.example.com).
Human web applications served by this server do not need to be listed.
When empty, requests from any origin are allowed.

### `DOMAIN` {#DOMAIN}

<div class="env-meta">

**Type** `string` · **Default** `localhost`

</div>

Domain for the HTTP server.

### `DOMAIN_WEBAPP` {#DOMAIN_WEBAPP}

<div class="env-meta">

**Type** `string` · **Default** `localhost`

</div>

Domain for the HTTP webapp.

### `HTTP_ENABLE_DEBUG_ROUTE` {#HTTP_ENABLE_DEBUG_ROUTE}

<div class="env-meta">

**Type** `bool`

</div>

Enable `/debug` route.

### `HTTP_ENABLE_HEALTHCHECK_ROUTE` {#HTTP_ENABLE_HEALTHCHECK_ROUTE}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

### `HTTP_METRICS` {#HTTP_METRICS}

<div class="env-meta">

**Type** `bool`

</div>

Enable (prometheus) metrics.

### `HTTP_REPORT_PANIC` {#HTTP_REPORT_PANIC}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Report HTTP panic to Sentry.

### `HTTP_ENABLE_VERSION_ROUTE` {#HTTP_ENABLE_VERSION_ROUTE}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enable `/version` route.

### `HTTP_LOG_REQUEST` {#HTTP_LOG_REQUEST}

<div class="env-meta">

**Type** `bool`

</div>

Log HTTP requests.

### `HTTP_LOG_RESPONSE` {#HTTP_LOG_RESPONSE}

<div class="env-meta">

**Type** `bool`

</div>

Log HTTP responses.

### `HTTP_METRICS_PASSWORD` {#HTTP_METRICS_PASSWORD}

<div class="env-meta">

**Type** `string` · **Default** 5 random alphanumeric characters, drawn again on every start

</div>

Password for the metrics endpoint.

### `HTTP_METRICS_NAME` {#HTTP_METRICS_NAME}

<div class="env-meta">

**Type** `string` · **Default** `human`

</div>

Name for metrics endpoint.

### `HTTP_METRICS_USERNAME` {#HTTP_METRICS_USERNAME}

<div class="env-meta">

**Type** `string` · **Default** `metrics`

</div>

Username for the metrics endpoint.

### `HTTP_SSL_TERMINATED` {#HTTP_SSL_TERMINATED}

<div class="env-meta">

**Type** `bool` · **Default** true when LETSENCRYPT_HOST is set, false otherwise

</div>

Is SSL termination enabled in ingres, proxy or load balancer that is in front of Human?
By default, Human checks for presence of LETSENCRYPT_HOST environmental variable.
This DOES NOT enable SSL termination in Corteza!

### `HTTP_ERROR_TRACING` {#HTTP_ERROR_TRACING}

<div class="env-meta">

**Type** `bool`

</div>

### `HTTP_SERVER_WEB_CONSOLE_ENABLED` {#HTTP_SERVER_WEB_CONSOLE_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `false` · **Default** false, and true when ENVIRONMENT names a development environment

</div>

Enable web console. When running in dev environment, web console is enabled by default.

### `HTTP_SERVER_WEB_CONSOLE_PASSWORD` {#HTTP_SERVER_WEB_CONSOLE_PASSWORD}

<div class="env-meta">

**Type** `string` · **Default** 32 random alphanumeric characters drawn on every start, and empty in a development environment

</div>

Password for the web console endpoint. When running in dev environment, password is not required.

Human intentionally sets default password to random chars to prevent security incidents.

### `HTTP_SERVER_WEB_CONSOLE_USERNAME` {#HTTP_SERVER_WEB_CONSOLE_USERNAME}

<div class="env-meta">

**Type** `string` · **Default** `admin` · **Default** admin, and empty in a development environment

</div>

Username for the web console endpoint.

### `HTTP_WEBAPP_BASE_DIR` {#HTTP_WEBAPP_BASE_DIR}

<div class="env-meta">

**Type** `string` · **Default** `./webapp/public`

</div>

### `HTTP_WEBAPP_BASE_URL` {#HTTP_WEBAPP_BASE_URL}

<div class="env-meta">

**Type** `string` · **Default** `/`

</div>

Webapp base URL is internally prefixed with baseUrl

### `HTTP_WEBAPP_ENABLED` {#HTTP_WEBAPP_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

## RBAC options

### `RBAC_ANONYMOUS_ROLES` {#RBAC_ANONYMOUS_ROLES}

<div class="env-meta">

**Type** `string` · **Default** `anonymous`

</div>

Space delimited list of role handles.
These roles are automatically assigned to anonymous user.
Memberships can not be managed for these roles.

### `RBAC_AUTHENTICATED_ROLES` {#RBAC_AUTHENTICATED_ROLES}

<div class="env-meta">

**Type** `string` · **Default** `authenticated`

</div>

Space delimited list of role handles.
These roles are automatically assigned to authenticated user.
Memberships can not be managed for these roles.
System will refuse to start if roles listed here are also listed under anonymous roles

### `RBAC_BYPASS_ROLES` {#RBAC_BYPASS_ROLES}

<div class="env-meta">

**Type** `string` · **Default** `super-admin`

</div>

Space delimited list of role handles.
These roles causes short-circuiting access control check and allowing all operations.
System will refuse to start if check-bypassing roles are also listed as authenticated or anonymous auto-assigned roles.

### `RBAC_LOG` {#RBAC_LOG}

<div class="env-meta">

**Type** `bool`

</div>

Log RBAC related events and actions

### `RBAC_SERVICE_USER` {#RBAC_SERVICE_USER}

<div class="env-meta">

**Type** `string`

</div>

## SCIM Server

### `SCIM_BASE_URL` {#SCIM_BASE_URL}

<div class="env-meta">

**Type** `string` · **Default** `/scim`

</div>

Prefix for SCIM API endpoints

### `SCIM_ENABLED` {#SCIM_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Enable SCIM subsystem

### `SCIM_EXTERNAL_ID_AS_PRIMARY` {#SCIM_EXTERNAL_ID_AS_PRIMARY}

<div class="env-meta">

**Type** `bool`

</div>

Use external IDs in SCIM API endpoints

### `SCIM_EXTERNAL_ID_VALIDATION` {#SCIM_EXTERNAL_ID_VALIDATION}

<div class="env-meta">

**Type** `string` · **Default** `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`

</div>

Validates format of external IDs. Defaults to UUID

### `SCIM_SECRET` {#SCIM_SECRET}

<div class="env-meta">

**Type** `string`

</div>

Secret to use to validate requests on SCIM API endpoints

## Email sending

Configure your local SMTP server or use one of the available providers.

These values are copied to settings when the server starts and can be managed from the administration console.
We recommend you remove these values after they are copied to settings.
If server detects difference between these options and settings, it shows a warning in the log on server start.

### `SMTP_FROM` {#SMTP_FROM}

<div class="env-meta">

**Type** `string`

</div>

The SMTP `from` email parameter

### `SMTP_HOST` {#SMTP_HOST}

<div class="env-meta">

**Type** `string`

</div>

The SMTP server hostname.

### `SMTP_PASS` {#SMTP_PASS}

<div class="env-meta">

**Type** `string`

</div>

The SMTP password.

### `SMTP_PORT` {#SMTP_PORT}

<div class="env-meta">

**Type** `int`

</div>

The SMTP post.

### `SMTP_TLS_INSECURE` {#SMTP_TLS_INSECURE}

<div class="env-meta">

**Type** `bool`

</div>

Allow insecure (invalid, expired TLS certificates) connections.

### `SMTP_TLS_SERVER_NAME` {#SMTP_TLS_SERVER_NAME}

<div class="env-meta">

**Type** `string`

</div>

### `SMTP_USER` {#SMTP_USER}

<div class="env-meta">

**Type** `string`

</div>

The SMTP username.

## Actionlog

### `ACTIONLOG_DB_DSN` {#ACTIONLOG_DB_DSN}

<div class="env-meta">

**Type** `string`

</div>

Database connection string for the action log; when empty the primary DB connection is used.

### `ACTIONLOG_DEBUG` {#ACTIONLOG_DEBUG}

<div class="env-meta">

**Type** `bool`

</div>

### `ACTIONLOG_ENABLED` {#ACTIONLOG_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

### `ACTIONLOG_WORKFLOW_FUNCTIONS_ENABLED` {#ACTIONLOG_WORKFLOW_FUNCTIONS_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

## API Gateway

### `APIGW_DEBUG` {#APIGW_DEBUG}

<div class="env-meta">

**Type** `bool`

</div>

Enable API Gateway debugging info

### `APIGW_ENABLED` {#APIGW_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enable API Gateway

### `APIGW_LOG_ENABLED` {#APIGW_LOG_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Enable extra logging

### `APIGW_LOG_REQUEST_BODY` {#APIGW_LOG_REQUEST_BODY}

<div class="env-meta">

**Type** `bool`

</div>

Enable incoming request body output in logs

### `APIGW_PROFILER_ENABLED` {#APIGW_PROFILER_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enable profiler

### `APIGW_PROFILER_GLOBAL` {#APIGW_PROFILER_GLOBAL}

<div class="env-meta">

**Type** `bool` · **Default** `false`

</div>

Profiler enabled for all routes

### `APIGW_PROXY_ENABLE_DEBUG_LOG` {#APIGW_PROXY_ENABLE_DEBUG_LOG}

<div class="env-meta">

**Type** `bool`

</div>

Enable full debug log on requests / responses - warning, includes sensitive data

### `APIGW_PROXY_FOLLOW_REDIRECTS` {#APIGW_PROXY_FOLLOW_REDIRECTS}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Follow redirects on proxy requests

### `APIGW_PROXY_OUTBOUND_TIMEOUT` {#APIGW_PROXY_OUTBOUND_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `30s`

</div>

Outbound request timeout

## Authentication

### `AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME` {#AUTH_OAUTH2_ACCESS_TOKEN_LIFETIME}

<div class="env-meta">

**Type** `time.Duration` · **Default** `2h`

</div>

Lifetime of the access token. Should be shorter than lifetime of the refresh token.

### `AUTH_ASSETS_PATH` {#AUTH_ASSETS_PATH}

<div class="env-meta">

**Type** `string`

</div>

Path to js, css, images and template source files

When human starts, if path exists it tries to load template files from it.

When empty path is set (default value), embedded files are used.

### `AUTH_BASE_URL` {#AUTH_BASE_URL}

<div class="env-meta">

**Type** `string` · **Default** DOMAIN and HTTP_BASE_URL joined, e.g. http://localhost/auth

</div>

Frontend base URL. Must be an absolute URL, with the domain.
This is used for some redirects and links in auth emails.

### `AUTH_CSRF_COOKIE_NAME` {#AUTH_CSRF_COOKIE_NAME}

<div class="env-meta">

**Type** `string` · **Default** `same-site-authenticity-token`

</div>

Cookie name used for CSRF protection

### `AUTH_CSRF_ENABLED` {#AUTH_CSRF_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enable CSRF protection

### `AUTH_CSRF_FIELD_NAME` {#AUTH_CSRF_FIELD_NAME}

<div class="env-meta">

**Type** `string` · **Default** `same-site-authenticity-token`

</div>

Form field name used for CSRF protection

### `AUTH_CSRF_SECRET` {#AUTH_CSRF_SECRET}

<div class="env-meta">

**Type** `string` · **Default** md5 of DB_DSN and HOSTNAME; stable across restarts, changes when either does

</div>

Secret used for securing CSRF protection

::: warning
If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
Generated secret will change if you change any of these variables.
:::

### `AUTH_DEFAULT_CLIENT` {#AUTH_DEFAULT_CLIENT}

<div class="env-meta">

**Type** `string` · **Default** `human-webapp`

</div>

Handle for OAuth2 client used for automatic redirect from /auth/oauth2/go endpoint.

This simplifies configuration for OAuth2 flow for Human Web applications as it removes
the need to supply redirection URL and client ID (oauth2/go endpoint does that internally)

### `AUTH_DEFAULT_REDIRECT_URIS` {#AUTH_DEFAULT_REDIRECT_URIS}

<div class="env-meta">

**Type** `string`

</div>

Space separated list of redirect URIs allowed for auth clients that do not have their own redirect URIs set.
Origins of AUTH_BASE_URL and webapp domain (DOMAIN_WEBAPP) are always allowed and do not need to be listed.

Scheme, host and port must match exactly, path of the redirect URI must be under the listed path.
Use * to allow redirection to any URI (not recommended).

### `AUTH_DEFAULT_SUB_USER_GROUP` {#AUTH_DEFAULT_SUB_USER_GROUP}

<div class="env-meta">

**Type** `string` · **Default** `default-sub-root`

</div>

Default sub-user group.
If default sub-user group is set, all users are assigned here.

Each user must belong to a user group if it wishes to use the system.

### `AUTH_DEFAULT_USER_GROUP` {#AUTH_DEFAULT_USER_GROUP}

<div class="env-meta">

**Type** `string` · **Default** `default-root`

</div>

Default user group.
If default sub user group is not set, all users are assigned here.

Each user must belong to a user group if it wishes to use the system.

### `AUTH_DEVELOPMENT_MODE` {#AUTH_DEVELOPMENT_MODE}

<div class="env-meta">

**Type** `bool`

</div>

When enabled, human reloads template before every execution.
Enable this for debugging or when developing auth templates.

Should be disabled in production where templates do not change between server restarts.

### `AUTH_EXTERNAL_COOKIE_SECRET` {#AUTH_EXTERNAL_COOKIE_SECRET}

<div class="env-meta">

**Type** `string` · **Default** md5 of DB_DSN and HOSTNAME; stable across restarts, changes when either does

</div>

Secret used for securing cookies

::: warning
If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
Generated secret will change if you change any of these variables.
:::

### `AUTH_EXTERNAL_REDIRECT_URL` {#AUTH_EXTERNAL_REDIRECT_URL}

<div class="env-meta">

**Type** `string` · **Default** DOMAIN and HTTP_BASE_URL joined, e.g. http://localhost/auth/external/{provider}/callback

</div>

Redirect URL to be sent with OAuth2 authentication request to provider

`provider` placeholder is replaced with the actual value when used.

### `AUTH_GARBAGE_COLLECTOR_INTERVAL` {#AUTH_GARBAGE_COLLECTOR_INTERVAL}

<div class="env-meta">

**Type** `time.Duration` · **Default** `15m`

</div>

How often are expired sessions and tokens purged from the database

### `AUTH_JWT_ALGORITHM` {#AUTH_JWT_ALGORITHM}

<div class="env-meta">

**Type** `string` · **Default** `HS512`

</div>

Algorithm to be use for JWT signature.

Supported values:
 - HS256, HS384, HS512
 - PS256, PS384, PS512,
 - RS256, RS384, RS512

Provide shared secret string for HS256, HS384, HS512 and full private key or path to the file PS* and RS* algorithms.

### `AUTH_JWT_KEY` {#AUTH_JWT_KEY}

<div class="env-meta">

**Type** `string`

</div>

Raw private key or absolute or relative path to the file containing one.

### `AUTH_LOG_ENABLED` {#AUTH_LOG_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Enable extra logging for authentication flows

### `AUTH_PASSWORD_SECURITY` {#AUTH_PASSWORD_SECURITY}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Password security allows you to disable constraints to which passwords must conform to.

::: danger
Disabling password security can be useful for development environments as it removes the need for complex passwords.
Password security *should be enabled* on production environments to avoid security incidents
:::

### `AUTH_PROVISION_SUPER_USER` {#AUTH_PROVISION_SUPER_USER}

<div class="env-meta">

**Type** `string`

</div>

When set, Human creates one or more users with the configured values using provided email as a password.
It skips existing (email, handle). All new users are assigned to all bypass roles.

When set in production, Human stops and reports an error

### `AUTH_OAUTH2_REFRESH_TOKEN_LIFETIME` {#AUTH_OAUTH2_REFRESH_TOKEN_LIFETIME}

<div class="env-meta">

**Type** `time.Duration` · **Default** `72h`

</div>

Lifetime of the refresh token. Should be much longer than lifetime of the access token.

Refresh tokens are used to exchange expired access tokens with new ones.

### `AUTH_REQUEST_RATE_LIMIT` {#AUTH_REQUEST_RATE_LIMIT}

<div class="env-meta">

**Type** `int` · **Default** `60`

</div>

How many requests from a certain IP address are allowed in a time window.
Set to zero to disable

### `AUTH_REQUEST_RATE_WINDOW_LENGTH` {#AUTH_REQUEST_RATE_WINDOW_LENGTH}

<div class="env-meta">

**Type** `time.Duration` · **Default** `1m`

</div>

How many requests from a certain IP address are allowed in a time window

### `AUTH_JWT_SECRET` {#AUTH_JWT_SECRET}

<div class="env-meta">

**Type** `string` · **Default** md5 of DB_DSN and HOSTNAME; stable across restarts, changes when either does

</div>

Secret used for signing JWT tokens.
Value is used only when HS256, HS384 or HS512 algorithm is used.

::: warning
If secret is not set, system auto-generates one from DB_DSN and HOSTNAME environment variables.
Generated secret will change if you change any of these variables.
:::

### `AUTH_SESSION_COOKIE_DOMAIN` {#AUTH_SESSION_COOKIE_DOMAIN}

<div class="env-meta">

**Type** `string` · **Default** DOMAIN, falling back to LETSENCRYPT_HOST, VIRTUAL_HOST, then localhost

</div>

Session cookie domain

### `AUTH_SESSION_COOKIE_NAME` {#AUTH_SESSION_COOKIE_NAME}

<div class="env-meta">

**Type** `string` · **Default** `session`

</div>

Session cookie name

### `AUTH_SESSION_COOKIE_PATH` {#AUTH_SESSION_COOKIE_PATH}

<div class="env-meta">

**Type** `string` · **Default** HTTP_BASE_URL joined with /auth, so /auth unless HTTP_BASE_URL is set

</div>

Session cookie path

### `AUTH_SESSION_COOKIE_SECURE` {#AUTH_SESSION_COOKIE_SECURE}

<div class="env-meta">

**Type** `bool` · **Default** follows HTTP_SSL_TERMINATED, which is true when LETSENCRYPT_HOST is set

</div>

Defaults to true when HTTPS is used. Human will try to guess the this setting by

### `AUTH_SESSION_LIFETIME` {#AUTH_SESSION_LIFETIME}

<div class="env-meta">

**Type** `time.Duration` · **Default** `24h`

</div>

Maximum time user is allowed to stay idle when logged in without "remember-me" option and before session is expired.

Recommended value is between an hour and a day.

::: warning
This affects only profile (/auth) pages. Using applications (admin, compose, ...) does not prolong the session.
:::

### `AUTH_SESSION_PERM_LIFETIME` {#AUTH_SESSION_PERM_LIFETIME}

<div class="env-meta">

**Type** `time.Duration` · **Default** `8640h`

</div>

Duration of the session in /auth lasts when user logs-in with "remember-me" option.

If set to 0, "remember-me" option is removed.

## Connection to Corredor

### `CORREDOR_ADDR` {#CORREDOR_ADDR}

<div class="env-meta">

**Type** `string` · **Default** `localhost:50051`

</div>

Hostname and port of the Corredor gRPC server.

### `CORREDOR_DEFAULT_EXEC_TIMEOUT` {#CORREDOR_DEFAULT_EXEC_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `1m`

</div>

### `CORREDOR_ENABLED` {#CORREDOR_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Enable/disable Corredor integration

### `CORREDOR_LIST_REFRESH` {#CORREDOR_LIST_REFRESH}

<div class="env-meta">

**Type** `time.Duration` · **Default** `5s`

</div>

### `CORREDOR_LIST_TIMEOUT` {#CORREDOR_LIST_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `2s`

</div>

### `CORREDOR_MAX_BACKOFF_DELAY` {#CORREDOR_MAX_BACKOFF_DELAY}

<div class="env-meta">

**Type** `time.Duration` · **Default** `1m`

</div>

Max delay for backoff on connection.

### `CORREDOR_MAX_RECEIVE_MESSAGE_SIZE` {#CORREDOR_MAX_RECEIVE_MESSAGE_SIZE}

<div class="env-meta">

**Type** `int` · **Default** `16777216`

</div>

Max message size that can be received.

### `CORREDOR_RUN_AS_ENABLED` {#CORREDOR_RUN_AS_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

### `CORREDOR_CLIENT_CERTIFICATES_CA` {#CORREDOR_CLIENT_CERTIFICATES_CA}

<div class="env-meta">

**Type** `string` · **Default** `ca.crt` · **Default** ca.crt, resolved against CORREDOR_CLIENT_CERTIFICATES_PATH

</div>

### `CORREDOR_CLIENT_CERTIFICATES_ENABLED` {#CORREDOR_CLIENT_CERTIFICATES_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

### `CORREDOR_CLIENT_CERTIFICATES_PATH` {#CORREDOR_CLIENT_CERTIFICATES_PATH}

<div class="env-meta">

**Type** `string` · **Default** `/certs/corredor/client`

</div>

### `CORREDOR_CLIENT_CERTIFICATES_PRIVATE` {#CORREDOR_CLIENT_CERTIFICATES_PRIVATE}

<div class="env-meta">

**Type** `string` · **Default** `private.key` · **Default** private.key, resolved against CORREDOR_CLIENT_CERTIFICATES_PATH

</div>

### `CORREDOR_CLIENT_CERTIFICATES_PUBLIC` {#CORREDOR_CLIENT_CERTIFICATES_PUBLIC}

<div class="env-meta">

**Type** `string` · **Default** `public.crt` · **Default** public.crt, resolved against CORREDOR_CLIENT_CERTIFICATES_PATH

</div>

### `CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME` {#CORREDOR_CLIENT_CERTIFICATES_SERVER_NAME}

<div class="env-meta">

**Type** `string`

</div>

## Environment

### `ENVIRONMENT` {#ENVIRONMENT}

<div class="env-meta">

**Type** `string` · **Default** `production`

</div>

## Events and scheduler

### `EVENTBUS_SCHEDULER_ENABLED` {#EVENTBUS_SCHEDULER_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enable eventbus scheduler.

### `EVENTBUS_SCHEDULER_INTERVAL` {#EVENTBUS_SCHEDULER_INTERVAL}

<div class="env-meta">

**Type** `time.Duration` · **Default** `60s`

</div>

Set time interval for `eventbus` scheduler.

## federation

### `FEDERATION_SYNC_DATA_MONITOR_INTERVAL` {#FEDERATION_SYNC_DATA_MONITOR_INTERVAL}

<div class="env-meta">

**Type** `time.Duration` · **Default** `1m`

</div>

Delay in seconds for data sync

### `FEDERATION_SYNC_DATA_PAGE_SIZE` {#FEDERATION_SYNC_DATA_PAGE_SIZE}

<div class="env-meta">

**Type** `int` · **Default** `100`

</div>

Bulk size in fetching for data sync

### `FEDERATION_ENABLED` {#FEDERATION_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Federation enabled on system, it toggles rest API endpoints, possibility to map modules in Compose and sync itself

### `FEDERATION_HOST` {#FEDERATION_HOST}

<div class="env-meta">

**Type** `string` · **Default** `local.example.com`

</div>

Host that is used during node pairing, also included in invitation

### `FEDERATION_LABEL` {#FEDERATION_LABEL}

<div class="env-meta">

**Type** `string` · **Default** `federated`

</div>

Federation label

### `FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL` {#FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL}

<div class="env-meta">

**Type** `time.Duration` · **Default** `2m`

</div>

Delay in seconds for structure sync

### `FEDERATION_SYNC_STRUCTURE_PAGE_SIZE` {#FEDERATION_SYNC_STRUCTURE_PAGE_SIZE}

<div class="env-meta">

**Type** `int` · **Default** `1`

</div>

Bulk size in fetching for structure sync

## Limits

### `LIMIT_RECORD_COUNT_PER_MODULE` {#LIMIT_RECORD_COUNT_PER_MODULE}

<div class="env-meta">

**Type** `int`

</div>

Maximum number of records per module

### `LIMIT_SYSTEM_USERS` {#LIMIT_SYSTEM_USERS}

<div class="env-meta">

**Type** `int`

</div>

Maximum number of valid (not deleted, not suspended) users

## locale

### `LOCALE_DEVELOPMENT_MODE` {#LOCALE_DEVELOPMENT_MODE}

<div class="env-meta">

**Type** `bool`

</div>

When enabled, Human reloads language files on every request
Enable this for debugging or developing.

### `LOCALE_LANGUAGES` {#LOCALE_LANGUAGES}

<div class="env-meta">

**Type** `string` · **Default** `en`

</div>

List of comma delimited languages (language tags) to enable.
In case when an enabled language can not be loaded, error is logged.

When loading language configurations (config.xml) from the configured path(s).

### `LOCALE_LOG` {#LOCALE_LOG}

<div class="env-meta">

**Type** `bool`

</div>

Log locale related events and actions

### `LOCALE_PATH` {#LOCALE_PATH}

<div class="env-meta">

**Type** `string` · **Default** empty, and ../locale when ENVIRONMENT names a development environment or LOCALE_DEVELOPMENT_MODE is on

</div>

One or more paths to locale config and translation files, separated by colon

When with LOCALE_DEVELOPMENT_MODE=true, default value for path is ../locale

### `LOCALE_QUERY_STRING_PARAM` {#LOCALE_QUERY_STRING_PARAM}

<div class="env-meta">

**Type** `string` · **Default** `lng`

</div>

Name of the query string parameter used to pass the language tag (it overrides Accept-Language header).
Set it to empty string to disable detection from the query string.
This parameter is ignored if only one language is enabled

### `LOCALE_RESOURCE_TRANSLATIONS_ENABLED` {#LOCALE_RESOURCE_TRANSLATIONS_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

When enabled, an editor for resource translations is enabled in UI

## log

### `LOG_DEBUG` {#LOG_DEBUG}

<div class="env-meta">

**Type** `bool`

</div>

Disables json format for logging and enables more human-readable output with colors.

Disable for production.

### `LOG_FILTER` {#LOG_FILTER}

<div class="env-meta">

**Type** `string`

</div>

Log filtering rules by level and name (log-level:log-namespace).
Please note that level (LOG_LEVEL) is applied before filter and it affects the final output!

Leave unset for production.

Example:
`warn+:* *:auth,workflow.*`
Log warnings, errors, panic, fatals. Everything from auth and workflow is logged.


See more examples and documentation here: https://github.com/moul/zapfilter

### `LOG_INCLUDE_CALLER` {#LOG_INCLUDE_CALLER}

<div class="env-meta">

**Type** `bool`

</div>

Set to true to see where the logging was called from.

Disable for production.

### `LOG_LEVEL` {#LOG_LEVEL}

<div class="env-meta">

**Type** `string` · **Default** `warn`

</div>

Minimum logging level. If set to "warn",
Levels warn, error, dpanic panic and fatal will be logged.

Recommended value for production: warn

Possible values: debug, info, warn, error, dpanic, panic, fatal

### `LOG_STACKTRACE_LEVEL` {#LOG_STACKTRACE_LEVEL}

<div class="env-meta">

**Type** `string` · **Default** `dpanic`

</div>

Include stack-trace when logging at a specified level or below.
Disable for production.

Possible values: debug, info, warn, error, dpanic, panic, fatal

## Messaging queue

### `MESSAGEBUS_ENABLED` {#MESSAGEBUS_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enable messagebus

### `MESSAGEBUS_LOG_ENABLED` {#MESSAGEBUS_LOG_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Enable extra logging for messagebus watchers

## Monitoring

### `MONITOR_INTERVAL` {#MONITOR_INTERVAL}

<div class="env-meta">

**Type** `time.Duration` · **Default** `5m`

</div>

Output (log) interval for monitoring.

## Object (file) storage

The MinIO integration allows you to replace local storage with cloud storage. When configured, `STORAGE_PATH` is not needed.

### `MINIO_ACCESS_KEY` {#MINIO_ACCESS_KEY}

<div class="env-meta">

**Type** `string`

</div>

### `MINIO_BUCKET` {#MINIO_BUCKET}

<div class="env-meta">

**Type** `string` · **Default** `{component}`

</div>

`component` placeholder is replaced with service name (e.g system).

### `MINIO_ENDPOINT` {#MINIO_ENDPOINT}

<div class="env-meta">

**Type** `string`

</div>

### `MINIO_PATH_PREFIX` {#MINIO_PATH_PREFIX}

<div class="env-meta">

**Type** `string`

</div>

`component` placeholder is replaced with service name (e.g system).

### `MINIO_SSEC_KEY` {#MINIO_SSEC_KEY}

<div class="env-meta">

**Type** `string`

</div>

### `MINIO_SECRET_KEY` {#MINIO_SECRET_KEY}

<div class="env-meta">

**Type** `string`

</div>

### `MINIO_SECURE` {#MINIO_SECURE}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

### `MINIO_STRICT` {#MINIO_STRICT}

<div class="env-meta">

**Type** `bool`

</div>

### `STORAGE_PATH` {#STORAGE_PATH}

<div class="env-meta">

**Type** `string` · **Default** `var/store`

</div>

Location where uploaded files are stored.

## Provisioning

Provisioning allows you to configure a Human instance when deployed.
It occurs automatically after the Human server starts.

::: warning
We recommend you to keep provisioning enabled as it simplifies version updates by updating the database and updating settings.

If you're doing local development or some debugging, you can disable this.
:::

### `PROVISION_ALWAYS` {#PROVISION_ALWAYS}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Controls if provision should run when the server starts.

### `PROVISION_PATH` {#PROVISION_PATH}

<div class="env-meta">

**Type** `string` · **Default** `provision/*`

</div>

Colon seperated paths to config files for provisioning.

## Sentry monitoring

::: info
These parameters help in the development and testing process.
When you are deploying to production, these should be disabled to improve performance and reduce storage usage.

You should configure external services such as Sentry or ELK to keep track of logs and error reports.
:::

### `SENTRY_DSN` {#SENTRY_DSN}

<div class="env-meta">

**Type** `string`

</div>

Set to enable Sentry client.

### `SENTRY_ATTACH_STACKTRACE` {#SENTRY_ATTACH_STACKTRACE}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Attach stacktraces

### `SENTRY_DEBUG` {#SENTRY_DEBUG}

<div class="env-meta">

**Type** `bool`

</div>

Print out debugging information.

### `SENTRY_DIST` {#SENTRY_DIST}

<div class="env-meta">

**Type** `string`

</div>

Set reported distribution.

### `SENTRY_ENVIRONMENT` {#SENTRY_ENVIRONMENT}

<div class="env-meta">

**Type** `string`

</div>

Set reported environment.

### `SENTRY_MAX_BREADCRUMBS` {#SENTRY_MAX_BREADCRUMBS}

<div class="env-meta">

**Type** `int` · **Default** `0`

</div>

Maximum number of breadcrumbs.

### `SENTRY_RELEASE` {#SENTRY_RELEASE}

<div class="env-meta">

**Type** `string` · **Default** the version this server was built with, or "development" when it carries no build stamp

</div>

Set reported Release.

### `SENTRY_SAMPLE_RATE` {#SENTRY_SAMPLE_RATE}

<div class="env-meta">

**Type** `float64` · **Default** `1.0`

</div>

Sample rate for event submission (0.0 - 1.0. defaults to 1.0)

### `SENTRY_SERVERNAME` {#SENTRY_SERVERNAME}

<div class="env-meta">

**Type** `string`

</div>

Set reported Server name.

### `SENTRY_WEBAPP_DSN` {#SENTRY_WEBAPP_DSN}

<div class="env-meta">

**Type** `string`

</div>

Set to enable Sentry client for webapp.

## Rendering engine

### `TEMPLATE_RENDERER_GOTENBERG_ADDRESS` {#TEMPLATE_RENDERER_GOTENBERG_ADDRESS}

<div class="env-meta">

**Type** `string`

</div>

Gotenberg rendering container address.

### `TEMPLATE_RENDERER_GOTENBERG_ENABLED` {#TEMPLATE_RENDERER_GOTENBERG_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Is Gotenberg rendering container enabled.

## Data store (database) upgrade

### `UPGRADE_ALWAYS` {#UPGRADE_ALWAYS}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Controls if the upgradable systems should be upgraded when the server starts.

### `UPGRADE_DEBUG` {#UPGRADE_DEBUG}

<div class="env-meta">

**Type** `bool`

</div>

Enable/disable debug logging.
    To enable debug logging set `UPGRADE_DEBUG=true`.

## Delay system startup

You can configure these options to defer API execution until another external (HTTP) service is up and running.

::: tip
Delaying API execution can come in handy in complex setups where execution order is important.
:::

### `WAIT_FOR` {#WAIT_FOR}

<div class="env-meta">

**Type** `time.Duration`

</div>

Delays API startup for the amount of time specified (10s, 2m...).
    This delay happens before service (`WAIT_FOR_SERVICES`) probing.

### `WAIT_FOR_SERVICES` {#WAIT_FOR_SERVICES}

<div class="env-meta">

**Type** `string`

</div>

Space delimited list of hosts and/or URLs to probe.
    Host format: `host` or `host:443` (port will default to 80).

::: info
Services are probed in parallel.
:::

### `WAIT_FOR_SERVICES_PROBE_INTERVAL` {#WAIT_FOR_SERVICES_PROBE_INTERVAL}

<div class="env-meta">

**Type** `time.Duration` · **Default** `5s`

</div>

Interval between service probes.

### `WAIT_FOR_SERVICES_PROBE_TIMEOUT` {#WAIT_FOR_SERVICES_PROBE_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `30s`

</div>

Timeout for each service probe.

### `WAIT_FOR_SERVICES_TIMEOUT` {#WAIT_FOR_SERVICES_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `1m`

</div>

Max time for each service probe.

### `WAIT_FOR_STATUS_PAGE` {#WAIT_FOR_STATUS_PAGE}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Show temporary status web page.

## Websocket server

### `WEBSOCKET_LOG_ENABLED` {#WEBSOCKET_LOG_ENABLED}

<div class="env-meta">

**Type** `bool`

</div>

Enable extra logging for authentication flows

### `WEBSOCKET_PING_PERIOD` {#WEBSOCKET_PING_PERIOD}

<div class="env-meta">

**Type** `time.Duration` · **Default** `108s`

</div>

### `WEBSOCKET_PING_TIMEOUT` {#WEBSOCKET_PING_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `120s`

</div>

### `WEBSOCKET_TIMEOUT` {#WEBSOCKET_TIMEOUT}

<div class="env-meta">

**Type** `time.Duration` · **Default** `15s`

</div>

Time before `WsServer` gets timed out.

## Workflow

### `WORKFLOW_CALL_STACK_SIZE` {#WORKFLOW_CALL_STACK_SIZE}

<div class="env-meta">

**Type** `int` · **Default** `16`

</div>

Defines the maximum call stack size between workflows

### `WORKFLOW_EXEC_DEBUG` {#WORKFLOW_EXEC_DEBUG}

<div class="env-meta">

**Type** `bool`

</div>

Enables verbose logging for workflow execution

### `WORKFLOW_REGISTER` {#WORKFLOW_REGISTER}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Registers enabled and valid workflows and executes them when triggered

### `WORKFLOW_STACK_TRACE_ENABLED` {#WORKFLOW_STACK_TRACE_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `true`

</div>

Enables execution stack trace construction

### `WORKFLOW_STACK_TRACE_FULL` {#WORKFLOW_STACK_TRACE_FULL}

<div class="env-meta">

**Type** `bool` · **Default** `false`

</div>

Forces the stack trace to record all steps; intended for testing, memory use grows with every executed step

## Discovery

### `DISCOVERY_BASE_URL` {#DISCOVERY_BASE_URL}

<div class="env-meta">

**Type** `string`

</div>

Indicates host of human discovery server

### `DISCOVERY_DEBUG` {#DISCOVERY_DEBUG}

<div class="env-meta">

**Type** `bool` · **Default** `false`

</div>

Enable discovery related activity info

### `DISCOVERY_EMBEDDINGS_DIMENSION` {#DISCOVERY_EMBEDDINGS_DIMENSION}

<div class="env-meta">

**Type** `int` · **Default** `384`

</div>

Embeddings dimension

### `DISCOVERY_EMBEDDINGS_ENABLED` {#DISCOVERY_EMBEDDINGS_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `false`

</div>

Enable discovery embeddings generation

### `DISCOVERY_ENABLED` {#DISCOVERY_ENABLED}

<div class="env-meta">

**Type** `bool` · **Default** `false`

</div>

Enable discovery endpoints

### `DISCOVERY_HNSW_EF_CONSTRUCTION` {#DISCOVERY_HNSW_EF_CONSTRUCTION}

<div class="env-meta">

**Type** `int` · **Default** `128`

</div>

HNSW ef construction parameter

### `DISCOVERY_HNSW_M` {#DISCOVERY_HNSW_M}

<div class="env-meta">

**Type** `int` · **Default** `16`

</div>

HNSW m parameter

### `DISCOVERY_HUMAN_DOMAIN` {#DISCOVERY_HUMAN_DOMAIN}

<div class="env-meta">

**Type** `string`

</div>

Indicates host of human compose webapp

### `DISCOVERY_JWT_SECRET` {#DISCOVERY_JWT_SECRET}

<div class="env-meta">

**Type** `string`

</div>

JWT secret shared with the discovery searcher for signing search tokens

## attachment

### `AVATAR_INITIALS_BACKGROUND_COLOR` {#AVATAR_INITIALS_BACKGROUND_COLOR}

<div class="env-meta">

**Type** `string` · **Default** `#F3F3F3`

</div>

Avatar initials background color

### `AVATAR_INITIALS_COLOR` {#AVATAR_INITIALS_COLOR}

<div class="env-meta">

**Type** `string` · **Default** `#0B344E`

</div>

Avatar initials text color

### `AVATAR_INITIALS_FONT_PATH` {#AVATAR_INITIALS_FONT_PATH}

<div class="env-meta">

**Type** `string` · **Default** `fonts/Poppins-Regular.ttf`

</div>

Avatar initials font file path

### `ATTACHMENT_AVATAR_MAX_FILE_SIZE` {#ATTACHMENT_AVATAR_MAX_FILE_SIZE}

<div class="env-meta">

**Type** `int64` · **Default** `1000000`

</div>

Avatar image maximum upload size, default value is 1MB

## webapp

### `WEBAPP_SCSS_DIR_PATH` {#WEBAPP_SCSS_DIR_PATH}

<div class="env-meta">

**Type** `string`

</div>

Path to custom SCSS source files directory

## Agent observability

### `OBSERVABILITY_LANGFUSE_HOST` {#OBSERVABILITY_LANGFUSE_HOST}

<div class="env-meta">

**Type** `string`

</div>

Langfuse server host. When set, agent traces are exported to Langfuse.

### `OBSERVABILITY_LANGFUSE_PUBLIC_KEY` {#OBSERVABILITY_LANGFUSE_PUBLIC_KEY}

<div class="env-meta">

**Type** `string`

</div>

Langfuse public key.

### `OBSERVABILITY_LANGFUSE_SECRET_KEY` {#OBSERVABILITY_LANGFUSE_SECRET_KEY}

<div class="env-meta">

**Type** `string`

</div>

Langfuse secret key.

### `OBSERVABILITY_OTEL_EXPORTER_OTLP_ENDPOINT` {#OBSERVABILITY_OTEL_EXPORTER_OTLP_ENDPOINT}

<div class="env-meta">

**Type** `string`

</div>

OpenTelemetry OTLP endpoint. When set, agent traces are exported via OTLP HTTP.

## Agentic

### `AGENTIC_ANTHROPIC_API_VERSION` {#AGENTIC_ANTHROPIC_API_VERSION}

<div class="env-meta">

**Type** `string` · **Default** `2023-06-01`

</div>

Anthropic API version header sent with every request.

### `AGENTIC_MCP_SERVER_NAME` {#AGENTIC_MCP_SERVER_NAME}

<div class="env-meta">

**Type** `string` · **Default** `Human MCP`

</div>

MCP server name reported to clients.

### `AGENTIC_MCP_SERVER_VERSION` {#AGENTIC_MCP_SERVER_VERSION}

<div class="env-meta">

**Type** `string` · **Default** `v1`

</div>

MCP server version reported to clients.

