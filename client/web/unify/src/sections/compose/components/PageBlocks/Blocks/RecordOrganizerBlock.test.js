import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// One organizer block is one column of a board, and options.group is the
// groupField value it holds. Corteza builds that condition with the shared
// filter helper (RecordOrganizerBase.vue) and this has to match: interpolating
// the value straight into the query breaks on an apostrophe, and spells an
// empty group `= ''` — nothing on a text column, a postgres error on a numeric
// one — where the helper spells it IS NULL, the ungrouped column.

const list = vi.fn(() => Promise.resolve({ set: [] }))

const moduleStore = {
  findByID: vi.fn(() => Promise.resolve({})),
  getByID: () => ({
    moduleID: 'M1',
    namespaceID: 'N1',
    fields: [
      { name: 'status', kind: 'Select' },
      { name: 'rank', kind: 'Number' },
      { name: 'title', kind: 'String' },
    ],
  }),
}

vi.mock('@planetcrust/human-vue', () => ({
  useRecordStore: () => ({ list }),
  useModuleStore: () => moduleStore,
  usePageStore: () => ({ set: [] }),
  components: { CInputConfirm: { template: '<div />' } },
}))

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn(), resolve: () => ({ href: '' }) }),
  useRoute: () => ({ query: {}, params: {} }),
}))

vi.mock('./PageBlock.vue', () => ({ default: { template: '<div><slot /></div>' } }))

import RecordOrganizerBlock from './RecordOrganizerBlock.vue'

async function mountWith(options) {
  const w = mount(RecordOrganizerBlock, {
    props: {
      block: { options: { moduleID: 'M1', labelField: 'title', ...options } },
      namespace: { namespaceID: 'N1' },
    },
    global: {
      stubs: { Button: true, ProgressSpinner: true, draggable: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: { $ComposeAPI: {}, $Auth: { user: { userID: 'U1' } }, $eventBus: null },
    },
  })
  await flushPromises()
  return w
}

const lastQuery = () => list.mock.calls.at(-1)[0].query

beforeEach(() => list.mockClear())

describe('RecordOrganizerBlock group filter', () => {
  it('builds the group condition through the shared filter helper', async () => {
    await mountWith({ groupField: 'status', group: 'backlog' })

    expect(lastQuery()).toBe("((status = 'backlog'))")
  })

  it('escapes a group value that contains a quote', async () => {
    await mountWith({ groupField: 'status', group: "O'Brien" })

    // an unescaped value closes the string early and the query fails to parse
    expect(lastQuery()).not.toBe("((status = 'O'Brien'))")
    expect(lastQuery()).toContain('O')
  })

  it('reads an empty group as the ungrouped column, not an empty string', async () => {
    await mountWith({ groupField: 'status', group: '' })

    expect(lastQuery()).toBe('((status IS NULL))')
  })

  it('does not compare a numeric group field against an empty string', async () => {
    // `rank = ''` is `pq: invalid input syntax for type numeric: ""`
    await mountWith({ groupField: 'rank', group: '' })

    expect(lastQuery()).not.toContain("= ''")
    expect(lastQuery()).toBe('((rank IS NULL))')
  })

  it('quotes a numeric group value as the helper does', async () => {
    await mountWith({ groupField: 'rank', group: '7' })

    expect(lastQuery()).toBe("((rank = '7'))")
  })
})
