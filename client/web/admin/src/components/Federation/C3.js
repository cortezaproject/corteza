import { default as component } from './CFederationEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  node: {
    nodeID: '',
    name: 'Name',
    baseURL: 'https://example.com/federation',
    contact: 'email@example.com',
  },
  processing: false,
  success: false,
}

export default {
  name: 'Editor info',
  group: ['Federation'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('Enable delete', {
      value (p) { return p.node.nodeID.length > 0 },
      update (p, val) { p.node.nodeID = val ? '123456789' : '' },
    }),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        node: {
          ...props.node,
          name: '',
          baseURL: '',
          contact: '',
        },
      },
    },
  ],
}
