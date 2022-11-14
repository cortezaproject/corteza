import { default as component } from './CTemplateEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  template: {
    templateID: '',
    meta: {
      short: 'Name',
      description: 'Desc',
    },
    handle: 'Handle',
  },
  processing: false,
  success: false,
  canCreate: true,
}

export default {
  name: 'Editor info',
  group: ['Templates'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('CanCreate', 'canCreate'),
    checkbox('Enable delete', {
      value (p) { return p.template.templateID.length > 0 },
      update (p, val) { p.template.templateID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        template: {
          ...props.template,
          meta: {
            short: '',
            description: '',
          },
          handle: '',
        },
        canCreate: false,
      },
    },
  ],
}
