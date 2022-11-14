// eslint-disable-next-line
import { default as component } from './ReportItem.vue'

const props = {
  report: {},
  fixed: true,
}

export default {
  name: 'Report item',
  group: ['Chart'],
  component,
  props,
  controls: [
  ],
}
