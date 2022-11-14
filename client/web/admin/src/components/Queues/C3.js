import { default as component } from './CQueueEditorInfo.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  queue: {
    queueID: '',
    queue: 'Name',
    meta: {
      poll_delay: '1h',
      dispatch_events: 'test',
    },
    consumer: null,
  },
  processing: false,
  success: false,
  canCreate: true,
  consumers: [
    { value: 'dummy', text: 'Dummy' },
  ],
}

export default {
  name: 'Editor info',
  group: ['Queues'],

  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('CanCreate', 'canCreate'),
    checkbox('Enable delete', {
      value (p) { return p.queue.queueID.length > 0 },
      update (p, val) { p.queue.queueID = val ? '123456789' : '' },
    }),
  ],
}
