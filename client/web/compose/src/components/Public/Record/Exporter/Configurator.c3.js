// eslint-disable-next-line
import { default as component } from './Configurator.vue'
import { compose } from '@cortezaproject/corteza-js'

// TypeError: this.module.systemFields is not a function"
const props = {
  allowJSON: true,
  allowCSV: true,
  module: new compose.Module(),
  preselectedFields: [],
  recordCount: 0,
  query: undefined,
  selection: [],
  filterRangeType: 'all',
  filterRangeBy: 'created_at',
  dateRange: 'lastMonth',
  startDate: null,
  endDate: null,
  systemFields: ['ownedBy', 'createdAt', 'createdBy', 'updatedAt', 'updatedBy'],
  disabledTypes: ['User', 'Record', 'File'],
}

export default {
  name: 'Configurator',
  group: ['Exporter'],
  component,
  props,
  controls: [],
}
