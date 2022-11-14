import { default as component } from './CAuthclientEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  authclient: {
    authClientID: '',
    handle: 'corteza-webapp',
    meta: {
      description: '',
      name: 'Corteza Web Applications',
    },
  },
  roles: [],
  processing: false,
  success: false,
}

export default {
  name: 'Editor info',
  group: ['AuthClient'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('Enable delete', {
      value (p) { return p.authclient.authClientID.length > 0 },
      update (p, val) { p.authclient.authClientID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        authclient: {
          ...props.authclient,
          handle: '',
          meta: {
            name: '',
          },
        },
      },
    },
  ],
}
