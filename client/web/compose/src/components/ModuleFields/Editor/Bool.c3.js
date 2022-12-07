// eslint-disable-next-line
import { default as component } from './Bool.vue'
import { compose } from '@cortezaproject/corteza-js'

// const namespace = ({
//   canCreateChart: true,
//   canCreateModule: true,
//   canCreatePage: true,
//   canDeleteNamespace: true,
//   canGrant: true,
//   canManageNamespace: true,
//   canUpdateNamespace: true,
//   createdAt: 'Fri Jul 30 2021 18:25:13 GMT+0300 (Eastern European Summer Time)',
//   deletedAt: undefined,
//   enabled: true,
//   // labels: Object (empty),
//   meta: {
//     iconID: '0',
//     logoID: '0',
//   },
//   name: 'CRM',
//   namespaceID: '242313184189546499',
//   slug: 'crm',
// })

const props = {
  // compose.Namespace
  namespace: new compose.Namespace(),
  // compose.Module
  module: new compose.Module(),
  // compose.ModuleField
  field: new compose.ModuleFieldBool({
    options: { trueLabel: 'foo' },
  }),
  // options
}

export default {
  name: 'Bool',
  group: ['ModuleFields', '/Configurator'],
  component,
  props,
  controls: [],
}
