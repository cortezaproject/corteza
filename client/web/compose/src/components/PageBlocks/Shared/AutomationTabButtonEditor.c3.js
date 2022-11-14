// eslint-disable-next-line
import { default as component } from './AutomationTabButtonEditor.vue'
import { components } from '@cortezaproject/corteza-vue'
const { input } = components.C3.controls

const props = {
  button: {
    label: 'Dummy',
    variant: 'danger',
  },
  script: {},
  trigger: {},
}

export default {
  name: 'Automation editor',
  group: ['PageBlocks', 'Shared'],
  component,
  props,
  controls: [
    input('Label', 'button.label'),
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
        button: {
          label: '',
          variant: '',
        },
      },
    },
  ],
}
