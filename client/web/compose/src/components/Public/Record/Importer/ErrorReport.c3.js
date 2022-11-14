// eslint-disable-next-line
import { default as component } from './ErrorReport.vue'

const props = {
  session: {
    fields: {},
    progress: {
      completed: 0,
      entryCount: 37080,
      failLog: {
        errors: {
          'empty field Company': 1,
          'empty field LastName': 1,
          'invalidValue field Country': 1,
          'invalidValue value Level 1': 1,
        },
      },
      failReason: 'failed to complete transaction: store encoder encode corteza::compose:record [242313185917403139]: 3 issue(s) found',
      failed: 1,
      finishedAt: '2021-08-12T13:13:20.1322898Z',
      startedAt: '2021-08-12T13:13:18.6814197Z',
    },
  },
  noPool: false,
}

export default {
  name: 'Error report',
  group: ['Public/Record'],
  component,
  props,
  controls: [],
}
