// eslint-disable-next-line
import { default as component } from './CircleStep.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox, input } = components.C3.controls

const props = {
  done: true,
  disabled: false,
  optional: false,
  small: false,
  stepNumber: '111',
}

export default {
  name: 'Circle step',
  group: ['Common'],
  component,
  props,

  controls: [
    checkbox('Done', 'done'),
    checkbox('Disabled', 'disabled'),
    checkbox('Small', 'small'),
    input('Number of step', 'stepNumber'),
  ],
}
