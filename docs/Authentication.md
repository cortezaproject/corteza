# Authentication

Crust support a fixed set of standard OAuth 2 authentication providers 
(facebook, gplus, github and linkedin) and a arbitrary number of custom
issuers (over OpenID Connect).

# Available settings 

Settings for external providers are stored under keys 
`auth.external.providers.<provider>.<prop>` and 
`auth.external.providers.openid-connect.<provider>.<prop>`. 

Prop is one of: `key`, `secret`, `enabled`. OIDC settings also have `issuer` prop.

Example settings (`system-cli settings list --prefix=auth.external`):


```
auth.external.callback-endpoint	"http://system.api.local.crust.tech:3002/auth/external/%s/callback"
auth.external.enabled	true
auth.external.providers.facebook.enabled	true
auth.external.providers.facebook.key	"24226007270326"
auth.external.providers.facebook.secret	"7vtfXx213cfc125804a226afcae777fe47"
auth.external.providers.gplus	true
auth.external.providers.gplus.enabled	true
auth.external.providers.gplus.key	"10629818561-7a8vr0avs47dqic43h2lkrurhr.apps.googleusercontent.com"
auth.external.providers.gplus.secret	"bkHmIFdk2YvtfXx"
auth.external.providers.github.enabled	false
auth.external.providers.github.key	null
auth.external.providers.github.secret	null
auth.external.providers.linkedin.enabled	false
auth.external.providers.linkedin.key	null
auth.external.providers.linkedin.secret	null
auth.external.providers.openid-connect.didmos2.enabled	true
auth.external.providers.openid-connect.didmos2.key	"tXM2ouiovowzGabk"
auth.external.providers.openid-connect.didmos2.issuer "https://satosa.didmos.latest.crust.tech"
auth.external.providers.openid-connect.didmos2.secret	"e1d68bfd7718468ba8fd36131f5176b1"
auth.external.redirect-url	"http://system.api.local.crust.tech:3002/auth/external/%s/callback"
auth.external.session-store-secret	"fCVFSRWjVEcoYuhXSf3f6zVWO1p38XEWz2yS8WH7wKDbvpxFrZq7zlEuiUTvk4QF"
```

# Changing settings

Authentication settings can be changed in the administration (via the API) and with cli 
command (`system-cli settings set <key> <value>`). Please bare in mind that values passed 
to CLI tool must always be in raw JSON format.

Changing values requires system service restart.

On startup, you should see log entries similar to these:
```
initializing external authentication providers (3)
external authentication provider "facebook" added
external authentication provider "gplus" added
external authentication provider "openid-connect.didmos2" added
```


# OIDC Auto discovery/configuration

Crust CLI comes with auto-discovery tool:
```bash
system-cli external-auth auto-discovery name url
```

```bash
system-cli external-auth auto-discovery didmos2 https://satosa.didmos.crust.example.tld
```

This will autodiscover and autoconfigure new OIDC provider. 
If entry with this name already exists it will override it.

Please note that this provider is disabled by default.

To enable it, run:
```bash
system-cli settings key auth.external.providers.openid-connect.didmos2.enabled true
```
