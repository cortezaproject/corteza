// eslint-disable-next-line
import { default as component } from './Card.vue'
import { compose } from '@cortezaproject/corteza-js'

const props = {
  block: new compose.PageBlock({
    description: '',
    kind: 'Chart',
    title: 'Monthly sales',
    xywh: [4, 15, 4, 7],
  }),
  scrollableBody: true,
}

export default {
  name: 'Card',
  group: ['PageBlocks', 'Wrap'],
  component,
  props,
  controls: [],
}
