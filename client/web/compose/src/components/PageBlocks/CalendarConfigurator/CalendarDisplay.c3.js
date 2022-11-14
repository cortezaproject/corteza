// eslint-disable-next-line
import { default as component } from './CalendarDisplay.vue'
import { compose } from '@cortezaproject/corteza-js'

const props = {
  block: new compose.PageBlock({
    options: {
      header: {
        hide: false,
        hidePrevNext: 'true',
        hideTitle: true,
        hideToday: true,
      },
      defaultView: 'dayGridMonth',
    },
  }),
  blockIndex: 7,
  boundingRect: {},
  mode: 'configurator',
  module: null,
  namespace: new compose.Namespace(),
  page: new compose.Page(),
  record: null,
}

export default {
  name: 'Calendar display',
  group: ['PageBlocks', 'Calendar configurator'],
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
        block: new compose.PageBlock({
          options: {
            header: {
              hide: false,
              hidePrevNext: 'false',
              hideTitle: false,
              hideToday: false,
            },
            defaultView: '',
          },
        }),
      },
    },
  ],
}
