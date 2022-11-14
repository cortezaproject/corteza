// eslint-disable-next-line
import { default as component } from './Reminder.vue'

const props = {
  feed: {
    options: { color: '#9d5c5c' },
  },
  modules: [],
}

export default {
  name: 'Reminder',
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
          options: { color: '#ffff' },
        },
      },
    },
  ],
}
