# Corteza Setup

If you are not already familiar with it, please read documentation about [Corteza Command Line Interface](CLI.md).

## First steps

Corteza will pre-initialize itself to allow you to access all of it's features as quickly as possible.
Internal authentication enabled, sign-up without email confirmation (in case you do not have your SMTP configured just 
yet)...

## Configuring system

### Configuring authentication

Review your current (auto-configure) settings with `settings list`:

```
auth.external.enabled                                      	false
auth.external.redirect-url                                 	"http://system.api.local.crust.tech/auth/external/%s/callback"
auth.external.session-store-secret                         	"PBVta4xKfQ0LIQEOtycxXqZZrGbZdTCuF4hw1cxrly1YA2AY5uO8a0SyY4Tbd1bk"
auth.external.session-store-secure                         	false
auth.internal.enabled                                      	true
auth.internal.password-reset.enabled                       	true
auth.internal.signup-email-confirmation-required           	false
auth.internal.signup.enabled                               	true
auth.mail.from-address                                     	"change-me@example.tld"
auth.mail.from-name                                        	"Corteza Team"
```

| Key | Description |
| ---- | ---- |
| auth.external.enabled                             | Enable external authentication, see [ExternalAuth.md](ExternalAuth.md) for details
| auth.external.redirect-url                        | Where to redirect after successful external authentication. This is the URL that you usually need to insert into your provider's auth app configuration page
| auth.external.session-store-secret                | Keep session values secret
| auth.external.session-store-secure                | Secure sessino store (set to false if not using TLS/HTTPS)
| auth.internal.enabled                             | Enable/disable internal authentication (will users be able to use Corteza username and password to login)
| auth.internal.password-reset.enabled              | Enable password reset
| auth.internal.signup-email-confirmation-required  | Is email confirmation required on sign-up.
| auth.internal.signup.enabled                      | Is sign-up enabled.
| auth.mail.from-address                            | Who (email) is sending auth emails (password reset, email confirmation)
| auth.mail.from-name                               | Who (name) is sending auth emails (password reset, email confirmation)

## Configuring messaging

_To be implemented_


## Configuring corteza

_To be implemented_
