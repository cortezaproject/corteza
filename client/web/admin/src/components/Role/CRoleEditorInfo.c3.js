import { default as component } from './CRoleEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  role: {
    roleID: '',
    name: 'Federation',
    handle: 'federation',
  },
  processing: false,
  success: false,
  canCreate: true,
}

export default {
  name: 'Editor info',
  group: ['Role'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('CanCreate', 'canCreate'),
    checkbox('Enable delete and archive', {
      value (p) { return p.role.roleID.length > 0 },
      update (p, val) { p.role.roleID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        role: {
          ...props.role,
          name: '',
          handle: '',
        },
        canCreate: false,
      },
    },
  ],
}
