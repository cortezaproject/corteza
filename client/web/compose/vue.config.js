const buildVueConfig = require('./vue.config-builder')

module.exports = buildVueConfig({
  appFlavour: 'Namespaces',
  appName: 'compose',
  appLabel: 'Corteza Compose',
  theme: 'corteza-base',
  packageAlias: 'corteza-webapp-compose',
})
