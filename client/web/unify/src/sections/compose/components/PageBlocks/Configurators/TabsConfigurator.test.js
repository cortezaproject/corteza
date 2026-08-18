import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'

// A tab row carried its own title in its key, so every keystroke gave the row a
// new identity: the cell remounted, the input lost focus, and a title could not
// be typed past its first character.

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CResourceTable: {
      name: 'CResourceTable',
      props: ['items', 'fields', 'primaryKey', 'emptyMessage'],
      // Keyed the way PrimeVue's DataTable keys rows from `data-key`, so a key
      // that changes remounts the row here too.
      template: `<div>
        <div v-for="row in items" :key="row[primaryKey]" class="row">
          <slot name="body-title" :data="row" />
        </div>
      </div>`,
    },
  },
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

import TabsConfigurator from './TabsConfigurator.vue'

const stub = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

const InputStub = {
  name: 'InputText',
  props: ['modelValue', 'placeholder'],
  template: '<input class="title" />',
}

function mountConfigurator(tabs = [{ title: '', blockID: 'B2', lazy: true }]) {
  const blockDraft = ref({ blockID: 'B1', kind: 'Tabs', options: { tabs } })

  const wrapper = mount(TabsConfigurator, {
    props: {
      namespace: { namespaceID: 'N1' },
      page: { pageID: 'P1' },
      blocks: [{ blockID: 'B2', kind: 'Content', title: 'Inner' }],
    },
    global: {
      provide: { blockDraft },
      mocks: { $t: k => k },
      directives: { tooltip: {} },
      stubs: {
        Select: stub('Select', ['modelValue', 'options', 'placeholder', 'disabled']),
        Fieldset: stub('Fieldset', ['legend']),
        Divider: stub('Divider'),
        Column: stub('Column', ['field']),
        Button: stub('Button', ['label', 'icon', 'disabled']),
        Checkbox: stub('Checkbox', ['modelValue']),
        CFormGroup: stub('CFormGroup', ['label']),
        InputText: InputStub,
      },
    },
  })

  return { blockDraft, wrapper }
}

const titleInputs = w => w.findAllComponents(InputStub)
const draftTabs = d => d.value.options.tabs

describe('TabsConfigurator tab labels', () => {
  it('writes the label through to the block draft', async () => {
    const { blockDraft, wrapper } = mountConfigurator()

    await titleInputs(wrapper)[0].vm.$emit('update:modelValue', 'Details')

    expect(draftTabs(blockDraft)[0].title).toBe('Details')
  })

  it('keeps the same row across an edit', async () => {
    const { wrapper } = mountConfigurator()
    const before = wrapper.findAll('.row')[0].element

    await titleInputs(wrapper)[0].vm.$emit('update:modelValue', 'Deta')
    await titleInputs(wrapper)[0].vm.$emit('update:modelValue', 'Details')

    expect(wrapper.findAll('.row')[0].element).toBe(before)
  })
})
