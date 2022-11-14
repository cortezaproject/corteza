import { default as component } from './ConfirmationToggle.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  disabled: false,
}

export default {
  name: 'Confirmation toggle',
  group: ['Root components'],
  component,
  props,
  controls: [
    checkbox('Disabled', 'disabled'),
  ],
}
