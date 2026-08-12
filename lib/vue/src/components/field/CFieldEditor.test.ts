import { describe, it, expect, vi } from 'vitest'
import { defineComponent } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import CFieldEditor from './CFieldEditor.vue'

// Stub registry so sub-editor is a simple sync component — no async loading
vi.mock('./registry', () => ({
  resolveFieldEditor: () =>
    defineComponent({
      props: ['field', 'modelValue', 'disabled', 'namespace'],
      emits: ['update:modelValue'],
      template:
        '<input data-testid="sub-editor" :value="modelValue" :disabled="disabled || undefined" @input="$emit(\'update:modelValue\', $event.target.value)"/>',
    }),
}))

function field(overrides: Record<string, unknown> = {}) {
  return { kind: 'String', isMulti: false, options: {}, ...overrides }
}

async function mountEditor(props: Record<string, unknown>) {
  const wrapper = mount(CFieldEditor, { props })
  await flushPromises()
  return wrapper
}

describe('CFieldEditor', () => {
  describe('single value (non-multi)', () => {
    it('renders the sub-editor', async () => {
      const wrapper = await mountEditor({ field: field(), modelValue: 'hello' })
      expect(wrapper.find('[data-testid="sub-editor"]').exists()).toBe(true)
    })

    it('passes modelValue to sub-editor', async () => {
      const wrapper = await mountEditor({ field: field(), modelValue: 'test-value' })
      expect((wrapper.find('[data-testid="sub-editor"]').element as HTMLInputElement).value).toBe(
        'test-value',
      )
    })

    it('emits update:modelValue when sub-editor emits', async () => {
      const wrapper = await mountEditor({ field: field(), modelValue: '' })
      await wrapper.find('[data-testid="sub-editor"]').setValue('new-val')
      expect(wrapper.emitted('update:modelValue')![0]).toEqual(['new-val'])
    })

    it('passes disabled=true to sub-editor', async () => {
      const wrapper = await mountEditor({ field: field(), modelValue: '', disabled: true })
      expect(wrapper.find('[data-testid="sub-editor"]').attributes('disabled')).toBeDefined()
    })

    it('passes disabled=false (no attribute) when not disabled', async () => {
      const wrapper = await mountEditor({ field: field(), modelValue: '', disabled: false })
      expect(wrapper.find('[data-testid="sub-editor"]').attributes('disabled')).toBeUndefined()
    })
  })

  describe('multi value', () => {
    it('renders one sub-editor per value', async () => {
      const wrapper = await mountEditor({
        field: field({ isMulti: true }),
        modelValue: ['a', 'b', 'c'],
      })
      expect(wrapper.findAll('[data-testid="sub-editor"]')).toHaveLength(3)
    })

    it('emits updated array when value changes in one entry', async () => {
      const wrapper = await mountEditor({
        field: field({ isMulti: true }),
        modelValue: ['a', 'b'],
      })
      const editors = wrapper.findAll('[data-testid="sub-editor"]')
      await editors[0].setValue('x')
      const emitted = wrapper.emitted('update:modelValue')!
      expect(emitted[emitted.length - 1][0]).toEqual(['x', 'b'])
    })

    it('adds empty entry and emits when add button clicked', async () => {
      const wrapper = await mountEditor({
        field: field({ isMulti: true }),
        modelValue: ['a'],
      })
      const buttons = wrapper.findAll('button')
      await buttons[buttons.length - 1].trigger('click') // last button = add
      const emitted = wrapper.emitted('update:modelValue')!
      const last = emitted[emitted.length - 1][0] as string[]
      expect(last).toHaveLength(2)
      expect(last[1]).toBe('')
    })

    it('removes entry and emits remaining values', async () => {
      const wrapper = await mountEditor({
        field: field({ isMulti: true }),
        modelValue: ['a', 'b', 'c'],
        allowEmpty: true,
      })
      const buttons = wrapper.findAll('button')
      await buttons[0].trigger('click') // first button = remove first entry
      const emitted = wrapper.emitted('update:modelValue')!
      expect(emitted[emitted.length - 1][0]).toEqual(['b', 'c'])
    })

    it('keeps one empty entry when last item removed with allowEmpty=false', async () => {
      const wrapper = await mountEditor({
        field: field({ isMulti: true }),
        modelValue: ['only'],
        allowEmpty: false,
      })
      const buttons = wrapper.findAll('button')
      await buttons[0].trigger('click')
      const emitted = wrapper.emitted('update:modelValue')!
      const last = emitted[emitted.length - 1][0] as string[]
      expect(last).toEqual([''])
    })
  })

  describe('multiAbsorbing', () => {
    it('File field: renders single sub-editor receiving the full array', async () => {
      const wrapper = await mountEditor({
        field: field({ kind: 'File' }),
        modelValue: ['att1', 'att2'],
      })
      // Only one sub-editor, no per-entry loop
      expect(wrapper.findAll('[data-testid="sub-editor"]')).toHaveLength(1)
      // No add/remove buttons
      expect(wrapper.findAll('button')).toHaveLength(0)
    })

    it('Select+multiple is multi-absorbing', async () => {
      const wrapper = await mountEditor({
        field: field({ kind: 'Select', isMulti: true, options: { selectType: 'multiple' } }),
        modelValue: ['x', 'y'],
      })
      expect(wrapper.findAll('[data-testid="sub-editor"]')).toHaveLength(1)
      expect(wrapper.findAll('button')).toHaveLength(0)
    })

    it('multiAbsorbing editor emits update:modelValue directly', async () => {
      const wrapper = await mountEditor({
        field: field({ kind: 'File' }),
        modelValue: ['a'],
      })
      await wrapper.find('[data-testid="sub-editor"]').setValue('new')
      expect(wrapper.emitted('update:modelValue')![0]).toEqual(['new'])
    })
  })
})
