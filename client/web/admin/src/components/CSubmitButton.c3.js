import { default as component } from './CSubmitButton.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox, input } = components.C3.controls

const props = {
  processing: false,
  success: false,
  disabled: true,
  buttonClass: 'mt-2 pt-1',
  variant: 'outline-primary',
  iconVariant: '',
}

export default {
  name: 'Submit button',
  group: ['Root components'],
  component,
  props,

  controls: [
    checkbox('Processing', 'processing'),
    checkbox('Success', 'success'),
    checkbox('Disabled', 'disabled'),
    input('Button class', 'buttonClass'),
  ],

  scenarios: [
    { label: 'Full form',
      props,
    },
    { label: 'Empty form',
      props: {
        ...props,
        disabled: false,
        buttonClass: '',
        variant: '',
        iconVariant: '',
      },
    },
  ],
}
