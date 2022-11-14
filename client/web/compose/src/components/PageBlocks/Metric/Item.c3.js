// eslint-disable-next-line
import { default as component } from './Item.vue'
import { components } from '@cortezaproject/corteza-vue'
const { input } = components.C3.controls

const props = {
  metric: {
    valueStyle: {
      backgroundColor: '#f6efef',
      color: '#1d3a50',
      fontSize: '18',
    },
    prefix: 'prefix',
    suffix: 'suffix',
  },
  value: { value: 'value' },
}

export default {
  name: 'Item',
  group: ['PageBlocks', 'Metric'],
  component,
  props,
  controls: [
    input('prefix', 'metric.prefix'),
    input('suffix', 'metric.suffix'),
    input('value', 'value.value'),
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
        metric: {
          valueStyle: {
            backgroundColor: '',
            color: '',
            fontSize: '',
          },
          prefix: '',
          suffix: '',
        },
        value: { value: '' },
      },
    },
  ],
}
