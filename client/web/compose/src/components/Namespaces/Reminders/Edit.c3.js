// eslint-disable-next-line
import { default as component } from './Edit.vue'
import { components } from '@cortezaproject/corteza-vue'
const { input } = components.C3.controls

const props = {
  edit: {
    payload: {
      title: 'Reminder',
      notes: 'Urgent!',
    },
  },
  myID: '242313586825756675',
  users: [],
}

export default {
  name: 'Edit reminder',
  group: ['Right Panel'],
  component,
  props,
  controls: [
    input('Title', 'edit.payload.title'),
    input('Notes', 'edit.payload.notes'),
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
        edit: {
          payload: {
            title: '',
            notes: '',
          },
        },
      },
    },
  ],
}
