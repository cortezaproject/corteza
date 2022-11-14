import component from './CTranslatorForm.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  disabled: false,

  languages: [
    { tag: 'sl' },
    { tag: 'en' },
    { tag: 'de' },
  ],

  translations: [
    { resource: 'compose:module/34082935092', lang: 'sl', key: 'name', msg: 'slovensko ime' },
    { resource: 'compose:module/34082935092', lang: 'en', key: 'name', msg: 'english name' },
    { resource: 'compose:module/34082935092', lang: 'en', key: 'name', msg: 'de name' },
    { resource: 'compose:module/34082935092', lang: 'sl', key: 'description', msg: 'Kontakt' },
    { resource: 'compose:module/34082935092', lang: 'en', key: 'description', msg: 'Kontakt' },
    { resource: 'compose:module-field/582375902375', lang: 'en', key: 'name', msg: 'Contact' },
    { resource: 'compose:module-field/582375902375', lang: 'sl', key: 'name', msg: 'Contact' },
    { resource: 'compose:module-field/582375902375', lang: 'en', key: 'description', msg: 'Contact' },
    { resource: 'compose:module-field/582375902375', lang: 'sl', key: 'description', msg: 'Contact' },
    { resource: 'compose:module-field/582375902373', lang: 'en', key: 'description', msg: 'Contact' },
    { resource: 'compose:module-field/582375902373', lang: 'sl', key: 'description', msg: 'Contact' },
  ],
}

export default {
  name: 'Form',
  group: ['Translator'],
  component,
  props,
  controls: [
    checkbox('Disabled', 'disabled'),
  ],
}
