// eslint-disable-next-line
import { default as component } from './RecordListFilter.vue'
import { compose } from '@cortezaproject/corteza-js'

const props = {
  selectedField: { name: 'CaseNumber' },
  namespace: new compose.Namespace(),
  module: new compose.Module({
    fields: [
      {
        isMulti: false,
        label: 'Account Name',
        name: 'AccountId',
      },
    ],
  }),
  recordListFilter: [{
    groupCondition: '',
    filter: [{
      name: '',
    }],
  }],
}

export default {
  name: 'Record list filter',
  group: ['Common'],
  component,
  props,
  controls: [],
}
