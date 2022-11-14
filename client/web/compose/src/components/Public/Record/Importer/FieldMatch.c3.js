// eslint-disable-next-line
import { default as component } from './FieldMatch.vue'

const props = {
  session: {
    fields: {
      Test1: '',
      Test2: '',
      Test3: '',
    },
  },
  module: {
    fields: [
      {
        label: 'Address Street',
        name: 'Street',
      },
      {
        label: 'Address City',
        name: 'Street',
      },
      {
        label: 'Party',
        name: 'Street',
      },
    ],
  },
}

export default {
  name: 'Field match',
  group: ['Public/Record'],
  component,
  props,
  controls: [],
}
