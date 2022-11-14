// eslint-disable-next-line
import { default as component } from './FieldPicker.vue'
import { compose } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const f1 = new compose.ModuleField({ name: 'f1', label: 'Field One Label' })
const f2 = new compose.ModuleField({ name: 'f2' })
const f3 = new compose.ModuleField({ name: 'f3', isRequired: true })
const f4 = new compose.ModuleField({ name: 'f4', label: 'required with label', isRequired: true })

const module = new compose.Module({
  name: 'modName',
  fields: [f1, f2, f3, f4],
})

const props = {
  module,
  fields: [f2, f4],
  disableSystemFields: false,
}

export default {
  name: 'FieldPicker',
  group: ['Picker'],
  component,
  props,
  controls: [
    checkbox('Disabled system fields', 'disableSystemFields'),
  ],
}
