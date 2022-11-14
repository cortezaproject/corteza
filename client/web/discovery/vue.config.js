const buildVueConfig = require('./vue.config-builder')

module.exports = buildVueConfig({
  appFlavour: 'Discovery',
  appName: 'discovery',
  appLabel: 'Corteza Discovery Discovery',
  theme: 'corteza-base',
  packageAlias: 'corteza-webapp-discovery',
})
