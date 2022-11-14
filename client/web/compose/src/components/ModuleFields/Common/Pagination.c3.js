// eslint-disable-next-line
import { default as component } from './Pagination.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  hasPrevPage: true,
  hasNextPage: true,
}

export default {
  name: 'Pagination',
  group: ['ModuleFields'],
  component,
  props,
  controls: [
    checkbox('HasPrevPage', 'hasPrevPage'),
    checkbox('HasNextPage', 'hasNextPage'),
  ],

  scenarios: [
    {
      label: 'Full form',
      props,
    },
    {
      label: 'Empty form',
      props: {
        hasPrevPage: false,
        hasNextPage: false,
      },
    },
  ],
}
