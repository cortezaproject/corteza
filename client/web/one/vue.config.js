const buildVueConfig = require('./vue.config-builder')

module.exports = buildVueConfig({
  appFlavour: 'One',
  appName: 'one',
  appLabel: 'Corteza One',
  theme: 'corteza-base',
  packageAlias: 'corteza-webapp-one',
})
