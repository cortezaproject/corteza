// Human API location
window.HumanAPI = 'https://api.your-domain.tld'

// HumanAuth can be autoconfigured by replacing /api with /auth in HumanAPI
// or by appending /auth to the end of HumanAPI string
// When this is not possible and your configuration is more exotic you can set it
// explicitly:
// window.HumanAuth = 'https://api.your-domain.tld/auth';

// Configure HumanWebapp when your web applications are not placed on the root.
// This is autoconfigured from the value of <base> tag href attribute in most cases.
// window.HumanWebapp = 'https://your-domain.tld';

// Set to true to enable i18next-pseudo
// Used to test translation string in a development environment
// Even if set to true, it will be disabled in production
window.i18nPseudoModeEnabled = false
