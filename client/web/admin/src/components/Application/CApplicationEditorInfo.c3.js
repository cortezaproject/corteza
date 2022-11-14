import { default as component } from './CApplicationEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  application: { applicationID: '' },
  processing: false,
  success: false,
  canCreate: true,
}

export default {
  name: 'Editor info',
  group: ['Applications'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('CanCreate', 'canCreate'),
    checkbox('Enable delete', {
      value (p) { return p.application.applicationID.length > 0 },
      update (p, val) { p.application.applicationID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        application: {
          ...props.application,
          enabled: false,
          name: '',
        },
        canCreate: false,
      },
    },
  ],
}
