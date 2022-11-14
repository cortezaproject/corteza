import { default as component } from './CSystemEditorAuth.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  settings: {
    'auth.internal.enabled': true,
    'auth.internal.signup.enabled': true,
  },
  processing: false,
  success: false,
  canManage: true,
}

export default {
  name: 'System editor auth',
  group: ['Settings'],
  component,
  props,
  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('CanManage', 'canManage'),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        settings: {
          'auth.internal.enabled': false,
          'auth.internal.signup.enabled': false,
        },
        canManage: false,
      },
    },
  ],
}
