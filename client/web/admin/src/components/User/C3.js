import { default as component } from './CUserEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  user: {
    userID: '',
    email: 'email@mail.bg',
    name: 'Stefan',
    handle: 'SS',
  },
  uiPage: 'user/editor',
  resourceType: 'system:user',
  uiSlot: 'infoFooter',
  processing: false,
  success: false,
}

export default {
  name: 'Editor info',
  group: ['User'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('Enable delete and suspend', {
      value (p) { return p.user.userID.length > 0 },
      update (p, val) { p.user.userID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        user: {
          ...props.user,
          email: '',
          name: '',
          handle: '',
        },
      },
    },
  ],
}
