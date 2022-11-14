// eslint-disable-next-line
import { default as component } from './Grid.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  editable: true,
  blocks: [
    {
      description: 'Description',
      title: 'Title',
      x: 0,
      w: 20,
      y: 0,
      h: 15,
      i: 0,
      xywh: [0, 0, 20, 15],
    },
  ],
}

export default {
  name: 'Grid',
  group: ['Common'],
  component,
  props,
  controls: [
    checkbox('Editable', 'editable'),
  ],

  scenarios: [
    {
      label: 'Full form',
      props,
    },
    {
      label: 'Empty form',
      props: {
        ...props,
        editable: false,
      },
    },
  ],
}
