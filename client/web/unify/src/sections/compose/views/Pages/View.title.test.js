import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { reactive } from 'vue'

// A page layout may override the page title with its own (`config.useTitle`),
// interpolated against the signed-in user. Record pages have always had this;
// these pin it for plain pages, where there is no record to read.

const route = reactive({
  name: 'page',
  params: { slug: 'ns', pageID: 'P1' },
  query: {},
})

const router = { push: vi.fn(), replace: vi.fn() }
const goBack = vi.fn()

vi.mock('vue-router', () => ({
  useRoute: () => route,
  useRouter: () => router,
}))

let page
let layouts

const pageStore = { getByID: () => page }
const pageLayoutStore = { getByPageID: () => layouts }

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('@planetcrust/human-vue', () => ({
  usePageStore: () => pageStore,
  usePageLayoutStore: () => pageLayoutStore,
  useHistoryBack: () => goBack,
}))

vi.mock('@planetcrust/human-js', () => ({
  NoID: '0',
  compose: {
    PageBlockMaker: b => b,
    // The real implementation, so the test exercises real template semantics
    // (including throwing on a malformed template) rather than a stand-in.
    interpolateTemplate(template, { record, user, recordID, ownerID, userID }) {
      const evaluate = new Function(
        'record',
        'user',
        'recordID',
        'ownerID',
        'userID',
        'return `' + template + '`',
      )
      return evaluate(record, user, recordID, ownerID, userID)
    },
  },
}))

vi.mock('@/sections/compose/composables/usePageVisibility', () => ({
  fetchBlockID: b => b.blockID,
  refuseOnce: () => true,
  clearRefusal: () => {},
  usePageVisibility: () => ({
    buildExpressionVariables: () => ({}),
    determineLayout: () => Promise.resolve(layouts[0]),
    evaluateBlocks: () => Promise.resolve(new Set()),
  }),
}))

vi.mock('@/sections/compose/composables/useResourceTranslations', () => ({
  useResourceTranslations: () => ({ showTranslatorButton: false }),
}))

vi.mock('@/sections/compose/components/PageBlocks/Grid.vue', () => ({
  default: { template: '<div />' },
}))

vi.mock('@/sections/compose/components/Admin/Page/PageTranslator.vue', () => ({
  default: { template: '<div />' },
}))

import View from './View.vue'

let wrapper

function layoutWith(config, meta = {}) {
  return [{ pageLayoutID: 'L1', blocks: [], config, meta }]
}

async function mountView() {
  wrapper = mount(View, {
    props: { namespace: { namespaceID: 'N1' } },
    global: {
      stubs: { teleport: true, Teleport: true, Message: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $SystemAPI: {},
        $Auth: { user: { userID: 'U1', name: 'Ada' } },
      },
      renderStubDefaultSlot: false,
    },
    shallow: true,
  })

  await flushPromises()
  return wrapper
}

beforeEach(() => {
  page = {
    pageID: 'P1',
    namespaceID: 'N1',
    moduleID: '0',
    title: 'Plain page title',
    blocks: [],
  }
  layouts = layoutWith({})
})

afterEach(() => {
  wrapper?.unmount()
  wrapper = null
  vi.clearAllMocks()
})

describe('View title', () => {
  it('uses the page title when the layout does not override it', async () => {
    layouts = layoutWith({ useTitle: false }, { title: 'Layout name' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
    expect(wrapper.text()).not.toContain('Layout name')
  })

  it('uses the layout title when useTitle is on', async () => {
    layouts = layoutWith({ useTitle: true }, { title: 'Custom title' })
    await mountView()

    expect(wrapper.text()).toContain('Custom title')
    expect(wrapper.text()).not.toContain('Plain page title')
  })

  it('interpolates the signed-in user into the layout title', async () => {
    layouts = layoutWith({ useTitle: true }, { title: 'Welcome, ${user.name} (${userID})' })
    await mountView()

    expect(wrapper.text()).toContain('Welcome, Ada (U1)')
  })

  it('falls back to the page title for a template that needs a record', async () => {
    layouts = layoutWith({ useTitle: true }, { title: '${record.values.name}' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
  })

  it('falls back to the page title when useTitle is on but no title is set', async () => {
    layouts = layoutWith({ useTitle: true }, { title: '' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
  })

  it('falls back to the page title when the template cannot be evaluated', async () => {
    layouts = layoutWith({ useTitle: true }, { title: 'Welcome, ${user.name' })
    await mountView()

    expect(wrapper.text()).toContain('Plain page title')
  })
})
