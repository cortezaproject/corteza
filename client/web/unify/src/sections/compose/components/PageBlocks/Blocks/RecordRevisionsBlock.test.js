import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { compose } from '@planetcrust/human-js'

// A revision read has three outcomes the block used to render identically: a
// history that is empty, a module that keeps no history at all, and a request
// that failed. It says which, and it renders the changes through the same field
// viewers the record page uses rather than raw keys and raw IDs.

let module_ = null

vi.mock('@planetcrust/human-vue', () => ({
  components: {
    CFieldViewer: {
      props: ['field', 'record'],
      template: '<span class="viewer">{{ field.name }}</span>',
    },
  },
  useModuleStore: () => ({ getByID: () => module_ }),
}))

const RecordRevisionsBlock = (await import('./RecordRevisionsBlock.vue')).default

function makeModule(enabled = true) {
  return new compose.Module({
    moduleID: '77',
    namespaceID: '1',
    config: { recordRevisions: { enabled } },
    fields: [
      { name: 'title', label: 'Title', kind: 'String' },
      { name: 'tags', label: 'Tags', kind: 'String', isMulti: true },
    ],
  })
}

function mountBlock(fetch, { record = { recordID: '1', revision: 2 } } = {}) {
  return mount(RecordRevisionsBlock, {
    props: {
      block: { kind: 'RecordRevisions', options: { preload: true }, fetch },
      page: { moduleID: '77' },
      namespace: { namespaceID: '1' },
      record,
    },
    global: {
      provide: { $ComposeAPI: {}, $eventBus: { on: () => () => {} } },
      mocks: { $t: k => k },
      stubs: {
        PageBlock: { template: '<div><slot /></div>' },
        ProgressSpinner: true,
        DataTable: { props: ['value'], template: '<table><slot /></table>' },
        Column: true,
        Dialog: true,
        Button: true,
      },
    },
  })
}

beforeEach(() => {
  module_ = makeModule()
})

describe('RecordRevisionsBlock', () => {
  it('says the read failed instead of reporting an empty history', async () => {
    const fetch = vi.fn(() => Promise.reject(new Error('record.errors.revisionsDisabledOnModule')))
    const w = mountBlock(fetch)
    await flushPromises()

    expect(w.text()).toContain('block.recordRevisions.viewer.errors.load-failed')
    expect(w.text()).not.toContain('block.recordRevisions.viewer.errors.no-revisions')
  })

  it('still reports an empty history as empty', async () => {
    const w = mountBlock(vi.fn(() => Promise.resolve([])))
    await flushPromises()

    expect(w.text()).toContain('block.recordRevisions.viewer.errors.no-revisions')
    expect(w.text()).not.toContain('block.recordRevisions.viewer.errors.load-failed')
  })

  it('does not offer a read the module keeps no history for', async () => {
    module_ = makeModule(false)
    const fetch = vi.fn(() => Promise.resolve([]))
    const w = mountBlock(fetch)
    await flushPromises()

    expect(w.text()).toContain('block.recordRevisions.viewer.errors.disabled-on-module')
    expect(fetch).not.toHaveBeenCalled()
  })

  it('labels a change with the module field label and carries its definition', async () => {
    const fetch = vi.fn(() =>
      Promise.resolve([
        {
          revision: 2,
          operation: 'updated',
          timestamp: '2026-01-01T00:00:00Z',
          userID: '5',
          changes: [
            { key: 'title', old: ['one'], new: ['two'] },
            { key: 'gone', new: ['x'] },
          ],
        },
      ]),
    )
    const w = mountBlock(fetch)
    await flushPromises()

    const [rev] = w.vm.revisions
    expect(rev.changes.map(c => c.label)).toEqual(['Title', 'gone'])
    expect(rev.changes[0].field.name).toBe('title')
    // A key the module no longer defines has no viewer to render it with.
    expect(rev.changes[1].field).toBe(null)
  })

  it('builds both sides of a revision as records the field viewers can read', async () => {
    const fetch = vi.fn(() =>
      Promise.resolve([
        {
          revision: 2,
          operation: 'updated',
          timestamp: '2026-01-01T00:00:00Z',
          userID: '5',
          changes: [
            { key: 'title', old: ['one'], new: ['two'] },
            { key: 'tags', old: ['a', 'b'], new: ['c'] },
          ],
        },
      ]),
    )
    const w = mountBlock(fetch)
    await flushPromises()

    const [rev] = w.vm.revisions
    expect(rev.oldRecord.values.title).toBe('one')
    expect(rev.newRecord.values.title).toBe('two')
    // A multi-value field keeps its array; a single-value one is unwrapped.
    expect(rev.oldRecord.values.tags).toEqual(['a', 'b'])
    expect(rev.newRecord.values.tags).toEqual(['c'])
  })

  it('leaves the created side without a record when only one side has values', async () => {
    const fetch = vi.fn(() =>
      Promise.resolve([
        {
          revision: 1,
          operation: 'created',
          timestamp: '2026-01-01T00:00:00Z',
          userID: '5',
          changes: [{ key: 'title', new: ['one'] }],
        },
      ]),
    )
    const w = mountBlock(fetch)
    await flushPromises()

    const [rev] = w.vm.revisions
    expect(rev.oldRecord).toBe(null)
    expect(rev.newRecord.values.title).toBe('one')
  })
})
