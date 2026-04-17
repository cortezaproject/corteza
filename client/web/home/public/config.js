// Human API location
window.HumanAPI = 'http://localhost:1043/api'
// window.HumanAPI = 'https://20219-qc.planetcrust.io/api'
// window.HumanAPI = 'https://internal.crust.tech/api'
// window.HumanAPI = 'https://pyd.staging.crust.tech/api'
// window.HumanAPI = 'https://nocode-api-qc.cloud.planetcrust.net/api'
// window.HumanAPI = 'https://20249-qc.planetcrust.io/api'

// window.HumanAPI = 'https://20223-qc.planetcrust.io/api'
// window.HumanAPI = 'https://portal.signalhillproducts.com/api'

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
