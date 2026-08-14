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

vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))

vi.mock('./PageBlock.vue', () => ({ default: { template: '<div><slot /></div>' } }))

import RecordOrganizerBlock from './RecordOrganizerBlock.vue'

const recordExec = vi.fn(() => Promise.resolve({}))
const emit = vi.fn()

async function mountWith(options, records = []) {
  list.mockResolvedValueOnce({ set: records })
  const w = mount(RecordOrganizerBlock, {
    props: {
      block: { options: { moduleID: 'M1', labelField: 'title', ...options } },
      namespace: { namespaceID: 'N1' },
    },
    global: {
      stubs: { Button: true, ProgressSpinner: true, draggable: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: { recordExec },
        $Auth: { user: { userID: 'U1' } },
        $eventBus: { on: () => () => {}, emit },
      },
    },
  })
  await flushPromises()
  return w
}

const lastQuery = () => list.mock.calls.at(-1)[0].query

// A DataTransfer stand-in: jsdom has no drag support, and `types` is the only
// thing a dragover handler may read, which is why the moduleID rides in the
// MIME type rather than the payload.
function transfer(mime, payload) {
  const store = payload === undefined ? {} : { [mime]: JSON.stringify(payload) }
  return {
    // a getter, not a snapshot: setData during dragstart has to show up here
    get types() {
      return Object.keys(store)
    },
    setData: (t, v) => {
      store[t] = v
    },
    getData: t => store[t] || '',
    effectAllowed: '',
    dropEffect: '',
  }
}

const MIME = 'application/x-human-record-M1'

beforeEach(() => {
  list.mockClear()
  recordExec.mockClear()
  emit.mockClear()
})

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

// A board is several organizer blocks over one module, and dropping a card into
// a column is what sets its key field to that column's key value — the promise
// the configurator already makes ("Value that will be set to the key field").
describe('RecordOrganizerBlock as a board', () => {
  const cards = [
    { recordID: 'R1', values: { title: 'One', pos: '1' } },
    { recordID: 'R2', values: { title: 'Two', pos: '2' } },
  ]

  const zone = w =>
    w.find('[ref="dropZone"]').exists()
      ? w.find('[ref="dropZone"]')
      : w
          .findAll('div')
          .find(d => d.attributes('draggable') === undefined && d.classes('overflow-auto'))

  it('marks cards draggable when a position or group field is configured', async () => {
    const w = await mountWith({ groupField: 'status', group: 'doing', positionField: 'pos' }, cards)

    expect(w.findAll('[data-organizer-card]')[0].attributes('draggable')).toBe('true')
  })

  it('leaves cards undraggable when neither field is configured', async () => {
    const w = await mountWith({}, cards)

    expect(w.findAll('[data-organizer-card]')[0].attributes('draggable')).toBe('false')
  })

  it('puts the moduleID in the drag type so a foreign board refuses the card', async () => {
    const w = await mountWith({ groupField: 'status', group: 'doing', positionField: 'pos' }, cards)
    const dataTransfer = transfer(MIME)

    await w.findAll('[data-organizer-card]')[0].trigger('dragstart', { dataTransfer })

    expect(dataTransfer.types).toContain(MIME)
    expect(JSON.parse(dataTransfer.getData(MIME))).toMatchObject({
      recordID: 'R1',
      moduleID: 'M1',
    })
  })

  it('organizes the record into this column on drop, setting the key field', async () => {
    const w = await mountWith({ groupField: 'status', group: 'done', positionField: 'pos' }, cards)

    await zone(w).trigger('drop', {
      clientY: 0,
      dataTransfer: transfer(MIME, { recordID: 'R9', moduleID: 'M1' }),
    })
    await flushPromises()

    expect(recordExec).toHaveBeenCalledTimes(1)
    const call = recordExec.mock.calls[0][0]
    expect(call.procedure).toBe('organize')
    expect(call.moduleID).toBe('M1')
    expect(Object.fromEntries(call.args.map(a => [a.name, a.value]))).toMatchObject({
      recordID: 'R9',
      groupField: 'status',
      group: 'done',
      positionField: 'pos',
    })
    // every column reloads, the one the card left included
    expect(emit).toHaveBeenCalledWith('refetch-records')
  })

  it('sends an empty group for the ungrouped column rather than omitting it', async () => {
    const w = await mountWith({ groupField: 'status', group: '', positionField: 'pos' }, cards)

    await zone(w).trigger('drop', {
      clientY: 0,
      dataTransfer: transfer(MIME, { recordID: 'R9', moduleID: 'M1' }),
    })
    await flushPromises()

    const args = Object.fromEntries(recordExec.mock.calls[0][0].args.map(a => [a.name, a.value]))
    expect(args.group).toBe('')
  })

  it('ignores a card dropped from another module', async () => {
    const w = await mountWith({ groupField: 'status', group: 'done', positionField: 'pos' }, cards)

    await zone(w).trigger('drop', {
      clientY: 0,
      dataTransfer: transfer('application/x-human-record-M2', { recordID: 'R9', moduleID: 'M2' }),
    })
    await flushPromises()

    expect(recordExec).not.toHaveBeenCalled()
  })
})
