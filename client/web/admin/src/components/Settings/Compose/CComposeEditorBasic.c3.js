import { default as component } from './CComposeEditorBasic.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  basic: {
    'compose.page.attachments.max-size': 22,
    'compose.record.attachments.max-size': 22,
  },
  processing: false,
  success: false,
  canManage: true,
}

export default {
  name: 'Compose editor basic',
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
        basic: {
          ...props.basic,
          'compose.page.attachments.max-size': '',
          'compose.record.attachments.max-size': '',
        },
        canCreate: false,
      },
    },
  ],
}
