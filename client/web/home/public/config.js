// Corteza API location
window.CortezaAPI = 'http://localhost:1043/api'
// window.CortezaAPI = 'https://20239-qc.cortezaproject.org/api'
// window.CortezaAPI = 'https://internal.crust.tech/api'
// window.CortezaAPI = 'https://pyd.staging.crust.tech/api'
// window.CortezaAPI = 'https://nocode-api-qc.cloud.planetcrust.net/api'
// window.CortezaAPI = 'https://20249-qc.cortezaproject.org/api'

// window.CortezaAPI = 'https://20223-qc.cortezaproject.org/api'
// window.CortezaAPI = 'https://portal.signalhillproducts.com/api'

// CortezaAuth can be autoconfigured by replacing /api with /auth in CortezaAPI
// or by appending /auth to the end of CortezaAPI string
// When this is not possible and your configuration is more exotic you can set it
// explicitly:
// window.CortezaAuth = 'https://api.cortezaproject.your-domain.tld/auth';

// Configure CortezaWebapp when your web applications are not placed on the root.
// This is autoconfigured from the value of <base> tag href attribute in most cases.
// window.CortezaWebapp = 'https://cortezaproject.your-domain.tld';

// Set to true to enable i18next-pseudo
// Used to test translation string in a development environment
// Even if set to true, it will be disabled in production
window.i18nPseudoModeEnabled = false
