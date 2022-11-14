// eslint-disable-next-line
import { default as component } from './Grid.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  editable: true,
  blocks: [
    {
      description: 'Description',
      kind: 'Chart',
      options: { chartID: '242313186186166275' },
      style: {
        variants: {
          bodyBg: 'white',
          border: 'primary',
          headerBg: 'white',
          headerText: 'primary',
        },
        title: 'Title',
        wrap: {
          kind: 'card',
        },
      },
      title: 'Leads by type',
      xywh: [9, 32, 3, 8],
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
