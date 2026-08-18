import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { ref } from 'vue'

// A tab label is interpolated at render time, so it is authored in the shared
// expression input like every other templated string — which only works if the
// row holding it survives being typed into.

const SCOPE = [{ name: 'recordID' }]

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CResourceTable: {
      name: 'CResourceTable',
      props: ['items', 'fields', 'primaryKey', 'emptyMessage'],
      // Keyed exactly the way PrimeVue's DataTable keys rows from `data-key`,
      // so a key that changes remounts the row here too.
      template: `<div>
        <div v-for="row in items" :key="row[primaryKey]" class="row">
          <slot name="body-title" :data="row" />
        </div>
      </div>`,
    },
  },
}))

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@/sections/compose/composables/useExpressionScope', () => ({
  useExpressionScope: () => ({ scope: SCOPE }),
}))

import TabsConfigurator from './TabsConfigurator.vue'

const stub = (name, props = []) => ({ name, props, template: '<div><slot /></div>' })

const ExpressionStub = {
  name: 'CInputExpression',
  props: ['modelValue', 'dialect', 'scope', 'minLines', 'placeholder'],
  methods: {
    insert(text) {
      this.$emit('update:modelValue', (this.modelValue || '') + text)
    },
  },
  template: '<div class="expr" />',
}

const HintStub = { name: 'CExpressionHint', props: ['scope'], template: '<div class="hint" />' }

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
        CInputExpression: ExpressionStub,
        CExpressionHint: HintStub,
      },
    },
  })

  return { blockDraft, wrapper }
}

const titleInputs = w => w.findAllComponents(ExpressionStub)
const draftTabs = d => d.value.options.tabs

describe('TabsConfigurator tab labels', () => {
  it('authors the label in an interpolation expression input', () => {
    const { wrapper } = mountConfigurator()
    const input = titleInputs(wrapper)[0]

    expect(input.props('dialect')).toBe('interpolation')
    expect(input.props('scope')).toEqual(SCOPE)
    expect(input.props('minLines')).toBe(1)
  })

  it('writes the label through to the block draft', async () => {
    const { blockDraft, wrapper } = mountConfigurator()

    await titleInputs(wrapper)[0].vm.$emit('update:modelValue', '${record.values.name}')

    expect(draftTabs(blockDraft)[0].title).toBe('${record.values.name}')
  })

  // The row key used to carry the title, so every keystroke remounted the cell
  // and threw away the editor's cursor and its open suggestion list.
  it('keeps the same editor across an edit', async () => {
    const { wrapper } = mountConfigurator()
    const before = titleInputs(wrapper)[0].vm

    await titleInputs(wrapper)[0].vm.$emit('update:modelValue', 'Deta')
    await titleInputs(wrapper)[0].vm.$emit('update:modelValue', 'Details')

    expect(titleInputs(wrapper)[0].vm).toBe(before)
  })

  it('inserts a hint chip into the row that was last focused', async () => {
    const { blockDraft, wrapper } = mountConfigurator([
      { title: 'one', blockID: 'B2', lazy: true },
      { title: 'two', blockID: '', lazy: true },
    ])

    await wrapper.findAll('.expr')[1].trigger('focusin')
    await wrapper.findComponent(HintStub).vm.$emit('insert', '${recordID}')

    expect(draftTabs(blockDraft)[1].title).toBe('two${recordID}')
    expect(draftTabs(blockDraft)[0].title).toBe('one')
  })

  it('inserts into the first row before anything has been focused', async () => {
    const { blockDraft, wrapper } = mountConfigurator([
      { title: 'one', blockID: 'B2', lazy: true },
      { title: 'two', blockID: '', lazy: true },
    ])

    await wrapper.findComponent(HintStub).vm.$emit('insert', '${recordID}')

    expect(draftTabs(blockDraft)[0].title).toBe('one${recordID}')
  })
})
