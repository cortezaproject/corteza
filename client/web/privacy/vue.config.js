const buildVueConfig = require('./vue.config-builder')

module.exports = buildVueConfig({
  appFlavour: 'Privacy',
  appName: 'privacy',
  appLabel: 'Corteza Privacy',
  theme: 'corteza-base',
  packageAlias: 'corteza-webapp-privacy',
})
