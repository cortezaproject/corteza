// eslint-disable-next-line
import { default as component } from './Record.vue'

// TypeError: this.module.systemFields is not a function or its return value is not iterable
const props = {
  feed: {
    options: {
      color: '#9d5c5c',
      moduleID: '5555555',
      prefilter: '',
    },
    allDay: true,
    endField: 'LastViewedDate',
    startField: 'LastViewedDate',
    titleField: 'BillingCity',
  },
  modules: [
    {
      moduleID: '242313186336833539',
      handle: 'Pool',
      name: 'Pool',
      fields: [
        { label: 'Pool', name: 'Pool', kind: 'String' },
      ],
    },
    {
      moduleID: '242313186335633539',
      handle: 'Party',
      name: 'Party',
      fields: [
        { label: 'Party', name: 'Party', kind: 'String' },
      ],
    },
    {
      moduleID: '242356786336833539',
      handle: 'Cool',
      name: 'Cool',
      fields: [
        { label: 'Cool', name: 'Cool', kind: 'String' },
      ],
    },
  ],
}

export default {
  name: 'Record',
  group: ['PageBlocks', 'Feed Source'],
  component,
  props,
  controls: [],

  scenarios: [
    {
      label: 'Full form',
      props,
    },
    {
      label: 'Empty form',
      props: {
        ...props,
        feed: {
          options: {
            color: '#FFFFFF',
            moduleID: '5555555',
            prefilter: '',
          },
          allDay: false,
          endField: '',
          startField: '',
          titleField: '',
        },
      },
    },
  ],
}
