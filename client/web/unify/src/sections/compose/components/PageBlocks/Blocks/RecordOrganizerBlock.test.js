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
  components: {
    CInputConfirm: { template: '<div />' },
    CFieldViewer: { name: 'CFieldViewer', props: ['field', 'record'], template: '<span />' },
    // Stands in for the shared draggable: renders its items and lets a test
    // emit the moves SortableJS would, which jsdom cannot produce.
    CDraggableList: {
      name: 'CDraggableList',
      props: ['modelValue', 'dragKey', 'group', 'disabled', 'handle', 'itemSelector', 'tag'],
      emits: ['add', 'update', 'remove', 'end'],
      template: '<div><slot /></div>',
    },
  },
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

// The block's own refetch-records handler, so a test can deliver an event the
// way another column on the board would.
let onRefetch = null

const BLOCK_ID = 3

async function mountWith(options, records = []) {
  list.mockResolvedValueOnce({ set: records })
  const w = mount(RecordOrganizerBlock, {
    props: {
      block: { blockID: BLOCK_ID, options: { moduleID: 'M1', labelField: 'title', ...options } },
      namespace: { namespaceID: 'N1' },
    },
    global: {
      stubs: { Button: true, ProgressSpinner: true },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: { recordExec },
        $Auth: { user: { userID: 'U1' } },
        $eventBus: {
          on: (name, fn) => {
            if (name === 'refetch-records') onRefetch = fn
            return () => {}
          },
          emit,
        },
        $toast: { toastErrorHandler: () => () => {} },
      },
    },
  })
  await flushPromises()
  return w
}

const dragList = w => w.findComponent({ name: 'CDraggableList' })
// The records the column is actually rendering, read off the card's viewer.
const cardRecords = w => w.findAllComponents({ name: 'CFieldViewer' }).map(c => c.props('record'))
const organizeArgs = call => Object.fromEntries(call.args.map(a => [a.name, a.value]))

const lastQuery = () => list.mock.calls.at(-1)[0].query

beforeEach(() => {
  list.mockClear()
  recordExec.mockClear()
  recordExec.mockResolvedValue({})
  emit.mockClear()
  onRefetch = null
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

  const moved = { recordID: 'R9', values: { title: 'Nine' } }
  const board = { groupField: 'status', group: 'done', positionField: 'pos' }

  it('drags when a position or group field is configured', async () => {
    const w = await mountWith({ groupField: 'status', group: 'doing' }, cards)

    expect(dragList(w).props('disabled')).toBe(false)
  })

  it('does not drag when neither field is configured', async () => {
    const w = await mountWith({}, cards)

    expect(dragList(w).props('disabled')).toBe(true)
  })

  it('names the drag group after the module so a foreign board refuses the card', async () => {
    const w = await mountWith(board, cards)

    expect(dragList(w).props('group')).toMatchObject({ name: 'record-organizer-M1' })
  })

  it('organizes the record into this column on drop, setting the key field', async () => {
    const w = await mountWith(board, cards)

    dragList(w).vm.$emit('add', { item: moved, index: 0, fromKey: '7' })
    await flushPromises()

    expect(recordExec).toHaveBeenCalledTimes(1)
    const call = recordExec.mock.calls[0][0]
    expect(call.procedure).toBe('organize')
    expect(call.moduleID).toBe('M1')
    expect(organizeArgs(call)).toMatchObject({
      recordID: 'R9',
      groupField: 'status',
      group: 'done',
      positionField: 'pos',
    })
  })

  it('positions the card one past the card it was dropped behind', async () => {
    const w = await mountWith(board, cards)

    // dropped at the end: the card now above it is R2, at position 2
    dragList(w).vm.$emit('add', { item: moved, index: 2, fromKey: '7' })
    await flushPromises()

    expect(organizeArgs(recordExec.mock.calls[0][0]).position).toBe('3')
  })

  it('positions a card dropped at the top at zero', async () => {
    const w = await mountWith(board, cards)

    dragList(w).vm.$emit('add', { item: moved, index: 0, fromKey: '7' })
    await flushPromises()

    expect(organizeArgs(recordExec.mock.calls[0][0]).position).toBe('0')
  })

  it('repositions a card moved within the column, naming itself both ends', async () => {
    const w = await mountWith(board, cards)

    dragList(w).vm.$emit('update', { item: cards[1], index: 0, oldIndex: 1 })
    await flushPromises()

    expect(recordExec).toHaveBeenCalledTimes(1)
    expect(emit).toHaveBeenCalledWith('refetch-records', {
      organized: { recordID: 'R2', moduleID: 'M1', fromKey: '3', toKey: '3' },
    })
  })

  it('sends an empty group for the ungrouped column rather than omitting it', async () => {
    const w = await mountWith({ ...board, group: '' }, cards)

    dragList(w).vm.$emit('add', { item: moved, index: 0, fromKey: '7' })
    await flushPromises()

    expect(organizeArgs(recordExec.mock.calls[0][0]).group).toBe('')
  })

  it('re-reads the column after the move, so the card shows what the server stored', async () => {
    const w = await mountWith(board, cards)

    // The card arrives holding the key it had in the column it left; the
    // column that took it only knows better once it has read the record back.
    const stale = { recordID: 'R9', values: { title: 'Nine', status: 'todo' } }
    list.mockClear()
    list.mockResolvedValueOnce({
      set: [{ recordID: 'R9', values: { title: 'Nine', status: 'done', pos: '0' } }],
    })

    dragList(w).vm.$emit('add', { item: stale, index: 0, fromKey: '7' })
    await flushPromises()

    expect(list).toHaveBeenCalledTimes(1)
    expect(cardRecords(w).map(r => r.values.status)).toEqual(['done'])
  })

  it('tells the board which two columns the card moved between', async () => {
    const w = await mountWith(board, cards)

    dragList(w).vm.$emit('add', { item: moved, index: 0, fromKey: '7' })
    await flushPromises()

    expect(emit).toHaveBeenCalledWith('refetch-records', {
      organized: { recordID: 'R9', moduleID: 'M1', fromKey: '7', toKey: '3' },
    })
  })
})

// A drop used to broadcast a bare refetch-records, and every block on the page
// answered it with a spinner — eight of them on an eight-column board, for a
// move two of those columns already showed.
describe('RecordOrganizerBlock refresh on someone else’s move', () => {
  const board = { groupField: 'status', group: 'done', positionField: 'pos' }

  const organized = extra => ({ organized: { recordID: 'R9', moduleID: 'M1', ...extra } })

  it('reloads on a bare refetch, as a saved record still requires', async () => {
    await mountWith(board)
    list.mockClear()

    onRefetch()
    await flushPromises()

    expect(list).toHaveBeenCalledTimes(1)
  })

  it('leaves the two columns the card moved between alone', async () => {
    await mountWith(board)
    list.mockClear()

    onRefetch(organized({ fromKey: '3', toKey: '9' }))
    await flushPromises()
    onRefetch(organized({ fromKey: '9', toKey: '3' }))
    await flushPromises()

    expect(list).not.toHaveBeenCalled()
  })

  it('refreshes a column the move did not touch without blanking it', async () => {
    const w = await mountWith(board, [
      { recordID: 'R1', values: { title: 'One', pos: '1' } },
      { recordID: 'R2', values: { title: 'Two', pos: '2' } },
    ])
    list.mockClear()

    let release
    list.mockReturnValueOnce(new Promise(res => (release = res)))

    onRefetch(organized({ fromKey: '8', toKey: '9' }))
    await flushPromises()

    // mid-flight: the cards are still on screen, not replaced by a spinner
    expect(list).toHaveBeenCalledTimes(1)
    expect(w.findAll('[data-organizer-card]')).toHaveLength(2)

    release({ set: [] })
    await flushPromises()
  })

  it('moves the card back when the server refuses it', async () => {
    const w = await mountWith(board, [{ recordID: 'R1', values: { title: 'One', pos: '1' } }])
    recordExec.mockRejectedValueOnce(new Error('nope'))

    dragList(w).vm.$emit('add', { item: { recordID: 'R9', values: {} }, index: 0, fromKey: '7' })
    await flushPromises()

    expect(emit).toHaveBeenCalledWith('refetch-records', {
      organized: { recordID: 'R9', moduleID: 'M1', fromKey: '7', toKey: '3', failed: true },
    })

    // and the two ends are exactly who re-reads, so the card lands back
    list.mockClear()
    onRefetch(organized({ fromKey: '7', toKey: '3', failed: true }))
    await flushPromises()
    expect(list).toHaveBeenCalledTimes(1)

    list.mockClear()
    onRefetch(organized({ fromKey: '8', toKey: '9', failed: true }))
    await flushPromises()
    expect(list).not.toHaveBeenCalled()
  })
})
