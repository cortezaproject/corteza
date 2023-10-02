// eslint-disable-next-line
import { default as component } from './Plain.vue'
import { compose } from '@cortezaproject/corteza-js'

const props = {
  block: new compose.PageBlock(),
}

export default {
  name: 'Plain',
  group: ['PageBlocks', 'Wrap'],
  component,
  props,
  controls: [],
}
