import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { cloneDeep, isEqual } from 'lodash-es'
import { ref } from 'vue'

// An unset connection ('0') means "whichever connection is primary". Resolving
// it into the draft would dirty an editor the user never touched — and worse,
// the next save would pin the module to one connection for good.

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('./DalFieldStoreEncoding.vue', () => ({ default: { template: '<div />' } }))

import { compose } from '@planetcrust/human-js'
import DalSettings from './DalSettings.vue'

const PRIMARY = { connectionID: 'C-primary', type: 'corteza::system:primary-dal-connection' }
const OTHER = { connectionID: 'C-other', type: 'corteza::system:dal-connection' }

function mountSettings(module) {
  const draft = ref(module)
  const wrapper = mount(DalSettings, {
    global: {
      provide: {
        moduleDraft: draft,
        $toast: { toastErrorHandler: () => vi.fn() },
        $SystemAPI: { dalConnectionList: () => Promise.resolve({ set: [PRIMARY, OTHER] }) },
      },
      mocks: { $t: k => k },
      renderStubDefaultSlot: false,
    },
    shallow: true,
  })
  return { wrapper, draft }
}

describe('DalSettings and the module draft', () => {
  it('leaves an unset connection unset once the connections load', async () => {
    const module = new compose.Module({ name: 'probe', fields: [] })
    expect(module.config.dal.connectionID).toBe('0')
    const before = cloneDeep(module)

    const { draft } = mountSettings(module)
    await flushPromises()

    expect(draft.value.config.dal.connectionID).toBe('0')
    expect(isEqual(draft.value, before)).toBe(true)
  })

  it('does not touch the draft at all while resolving the primary', async () => {
    const module = new compose.Module({
      name: 'probe',
      fields: [],
      config: { dal: { systemFieldEncoding: { id: { omit: true } } } },
    })
    const before = cloneDeep(module)

    const { draft } = mountSettings(module)
    await flushPromises()

    expect(isEqual(draft.value, before)).toBe(true)
  })

  it('keeps an explicitly chosen connection', async () => {
    const module = new compose.Module({
      name: 'probe',
      fields: [],
      config: { dal: { connectionID: 'C-other' } },
    })

    const { draft } = mountSettings(module)
    await flushPromises()

    expect(draft.value.config.dal.connectionID).toBe('C-other')
  })
})
