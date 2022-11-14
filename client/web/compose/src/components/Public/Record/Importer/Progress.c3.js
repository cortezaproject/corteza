// eslint-disable-next-line
import { default as component } from './Progress.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

// when !noPool => TypeError: Cannot read property 'recordImportProgress' of undefined
const props = {
  session: {
    progress: {
      completed: 0,
      entryCount: 1,
      failed: 0,
      finishedAt: null,
    },
  },
  noPool: false,
}

export default {
  name: 'Progress WIP',
  group: ['Public', 'Record'],
  component,
  props,
  controls: [
    checkbox('No pool', 'noPool'),
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
        noPool: false,
      },
    },
  ],
}
