import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CTranslatorForm from './CTranslatorForm.vue'

// jsdom has no scrollIntoView at all, so mounting with a highlightKey throws
// without this.
const scrollIntoView = vi.fn()
Element.prototype.scrollIntoView = scrollIntoView

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: (k: string) => k }) }))

const LANGUAGES = [
  { tag: 'en', name: 'English', localizedName: 'English' },
  { tag: 'fr', name: 'français', localizedName: 'French' },
]

const RESOURCE = 'compose:module-field/1/2/3'

const rows = (keys: string[]) =>
  keys.map(key => ({ resource: RESOURCE, key, lang: 'en', message: '' }))

function form(props: Record<string, unknown> = {}) {
  return mount(CTranslatorForm, {
    props: {
      languages: LANGUAGES,
      primaryResource: RESOURCE,
      translations: rows(['label', 'meta.options.new.text']),
      ...props,
    },
    global: {
      mocks: { $t: (k: string) => k },
      stubs: { Select: true },
    },
  })
}

function keyColumn(wrapper: ReturnType<typeof form>) {
  return wrapper.findAll('tbody tr td:first-child').map(td => td.text())
}

describe('CTranslatorForm key names', () => {
  it('formats a key generically when no prettifier is given', () => {
    expect(keyColumn(form())).toEqual(['Label', 'Meta › Options › New › Text'])
  })

  it('uses the prettifier where it names a key', () => {
    const wrapper = form({
      keyPrettifier: (k: string) => (k === 'label' ? 'Field label' : ''),
    })
    // '' is not a name — that key falls back rather than rendering blank.
    expect(keyColumn(wrapper)).toEqual(['Field label', 'Meta › Options › New › Text'])
  })
})

describe('CTranslatorForm highlightKey', () => {
  it('marks the row a caller asked for', () => {
    const wrapper = form({ highlightKey: 'meta.options.new.text' })
    const marked = wrapper.findAll('tbody tr').filter(tr => tr.classes('bg-highlight'))
    expect(marked).toHaveLength(1)
    expect(marked[0].text()).toContain('Meta › Options › New › Text')
  })

  it('marks nothing when no key is given', () => {
    expect(form().findAll('tbody tr.bg-highlight')).toHaveLength(0)
  })

  it('scrolls the marked row into view', async () => {
    scrollIntoView.mockClear()
    form({ highlightKey: 'label' })
    await flushPromises()
    expect(scrollIntoView).toHaveBeenCalledWith({ block: 'center' })
  })

  it('scrolls nothing when no key is given', async () => {
    scrollIntoView.mockClear()
    form()
    await flushPromises()
    expect(scrollIntoView).not.toHaveBeenCalled()
  })
})
