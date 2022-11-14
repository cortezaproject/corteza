// eslint-disable-next-line
import { default as component } from './List.vue'

const props = {
  reminders: [
    {
      payload: {
        notes: 'Urgent!',
        remindAt: 1800000,
        title: 'Remind me',
      },
      remindAt: 'Thu Aug 12 2021 13:09:40 GMT+0300 (Eastern European Summer Time)',
      reminderID: '244162818700476418',
      snoozeCount: 0,
    },
  ],
}

export default {
  name: 'List reminder',
  group: ['Right Panel'],
  component,
  props,
  controls: [],
}
