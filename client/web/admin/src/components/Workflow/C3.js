import { default as component } from './CWorkflowEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  workflow: {
    workflowID: '',
    meta: {
      name: 'CRM-Test',
    },
    handle: 'CRMT',
    enabled: true,
  },
  processing: false,
  success: false,
  canCreate: true,
}

export default {
  name: 'Editor info',
  group: ['Workflow'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('CanCreate', 'canCreate'),
    checkbox('Enable delete', {
      value (p) { return p.workflow.workflowID.length > 0 },
      update (p, val) { p.workflow.workflowID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        workflow: {
          ...props.workflow,
          meta: {
            name: '',
          },
          handle: '',
          enabled: false,
        },
        canCreate: false,
      },
    },
  ],
}
