// eslint-disable-next-line
import { default as component } from './MStyle.vue'
import { components } from '@cortezaproject/corteza-vue'
const { input } = components.C3.controls

const props = {
  options: {
    backgroundColor: '#f6efef',
    color: '#1d3a44',
    fontSize: '18',
  },
}

export default {
  name: 'MStyle',
  group: ['PageBlocks', 'MetricConfigurator'],
  component,
  props,
  controls: [
    input('Background color', 'options.backgroundColor'),
    input('Color', 'options.color'),
    input('Font size', 'options.fontSize'),
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
        options: {
          backgroundColor: '',
          color: '',
          fontSize: '',
        },
      },
    },
  ],
}
