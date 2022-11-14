import { default as component } from './CSystemEmailServer.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  value: {},
  processing: false,
  success: false,
  disabled: false,
}

export default {
  name: 'Email server',
  group: ['Settings'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('Disabled', 'disabled'),
  ],
}
