// eslint-disable-next-line
import { default as component } from './ReportEdit.vue'

const props = {
  report: {
    moduleID: '',
    filter: '',
    colorScheme: '#FFFFFF',
    metrics: [
      {
        aggregate: 'AVG',
        field: 'count',
        moduleID: '',
      },
    ],
    dimensions: [{
      field: 'Rating',
      modifier: '(no grouping / buckets)',
      skipMissing: true,
      default: '',
    }],
  },
  modules: [
    {
      moduleID: '',
      handle: 'Pool',
      name: 'Pool',
      fields: [{
        kind: 'String',
        name: 'Sample',
        options: {
          multiLine: false,
          useRichTextEditor: false,
        },
      }],
    },
    {
      moduleID: '',
      handle: 'Party',
      name: 'Party',
      fields: [{
        kind: 'String',
        name: 'Party',
        options: {
          multiLine: false,
          useRichTextEditor: false,
        },
      }],
    },
    {
      moduleID: '',
      handle: 'Test',
      name: 'Test',
      fields: [{
        kind: 'String',
        name: 'Test',
        options: {
          multiLine: false,
          useRichTextEditor: false,
        },
      }],
    },
  ],
  supportedMetrics: -1,
  dimensionFieldKind: ['DateTime', 'Select', 'Number', 'Bool', 'String'],
  unSkippable: false,
}

export default {
  name: 'Edit',
  group: ['Chart', 'Report'],
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
        modules: [],
      },
    },
  ],
}
