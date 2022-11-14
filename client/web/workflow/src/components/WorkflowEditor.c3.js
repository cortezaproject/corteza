import { default as component } from './WorkflowEditor.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  workflowObject: {
    workflowID: '',
    handle: 'handle',
    meta: {
      name: 'Name',
      description: 'desc',
    },
    enabled: true,
  },
  workflowTriggers: [],
  changeDetected: false,
  canCreate: true,
}

export default {
  name: 'Workflow editor',
  group: ['Root components'],
  component,
  props,

  controls: [
    checkbox('CanCreate', 'canCreate'),
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
        workflowObject: {
          ...props.workflowObject,
          enabled: false,
          name: '',
        },
        canCreate: false,
      },
    },
  ],
}
