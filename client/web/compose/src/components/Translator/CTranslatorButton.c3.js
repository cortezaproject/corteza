import component from './CTranslatorButton.vue'
import { components } from '@cortezaproject/corteza-vue'
const { input } = components.C3.controls

const props = {
  buttonClass: 'p-0 m-0',
  buttonVariant: 'link',
}

export default {
  name: 'Button',
  group: ['Translator'],
  component,
  props,
  controls: [
    input('Button class', 'buttonClass'),
    input('Button variant', 'buttonVariant'),
  ],
}
