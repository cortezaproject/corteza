import { default as component } from './CApplicationEditorUnify.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  application: {},

  unify: {
    config: '"data"',
    listed: true,
    logo: '/applications/low-code-platform.png',
    pinned: false,
    name: 'Low Code',
    url: '/compose',
  },

  canPin: true,
  processing: false,
  success: false,
  canCreate: true,
}

export default {
  name: 'Editor unify',
  group: ['Applications'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        unify: {
          ...props.unify,
          listed: false,
          name: '',
          config: '',
        },
        canPin: false,
        canCreate: false,
      },
    },
  ],
}
