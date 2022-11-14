// eslint-disable-next-line
import { compose } from '@cortezaproject/corteza-js'
import * as fieldTypes from './loader'

console.error({ ...fieldTypes })
const props = {
  namespace: new compose.Namespace(),
  module: new compose.Module(),
  field: new compose.ModuleFieldBool({
    options: {
      trueLabel: 'true',
      falseLabel: 'false',
    },
  }),
}

export default {
  name: { ...fieldTypes },
  group: ['ModuleFields', 'Configurator'],
  ...fieldTypes,
  props,
  controls: [],
}
