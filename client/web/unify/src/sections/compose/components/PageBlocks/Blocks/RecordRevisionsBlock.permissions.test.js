import { describe, it, expect, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'

// Reading a record's revisions is its own RBAC operation, and the server answers
// it on the record payload. An empty revisions table and a refused one look the
// same, so the block says which it is rather than reporting "no revisions".

vi.mock('@planetcrust/human-vue', () => ({
  components: { CFieldViewer: { template: '<span />' } },
  useModuleStore: () => ({ getByID: () => null }),
}))

const fetch = vi.fn(() => Promise.resolve([{ revision: 1, operation: 'create', changes: [] }]))

function mountBlock(record) {
  return mount(RecordRevisionsBlock, {
    props: {
      block: { kind: 'RecordRevisions', options: { preload: true }, fetch },
      page: { moduleID: '0' },
      record,
    },
    global: {
      provide: { $ComposeAPI: {}, $eventBus: { on: () => () => {} } },
      mocks: { $t: k => k },
      stubs: {
        PageBlock: { template: '<div><slot /></div>' },
        ProgressSpinner: true,
        DataTable: { template: '<table><slot /></table>' },
        Column: true,
        Dialog: true,
        Button: true,
      },
    },
  })
}

const RecordRevisionsBlock = (await import('./RecordRevisionsBlock.vue')).default

describe('RecordRevisionsBlock permission', () => {
  it('says so when the server refused the revision read', async () => {
    const w = mountBlock({ recordID: '1', revision: 2, canSearchRevisions: false })
    await flushPromises()
    expect(w.text()).toContain('block.noPermission')
    expect(fetch).not.toHaveBeenCalled()
  })

  it('loads revisions when the permission is granted', async () => {
    fetch.mockClear()
    const w = mountBlock({ recordID: '1', revision: 2, canSearchRevisions: true })
    await flushPromises()
    expect(w.text()).not.toContain('block.noPermission')
    expect(fetch).toHaveBeenCalled()
  })

  // A record that never carried the flag — a raw label-cache row, or a page
  // still loading — is not a refusal, and hiding the block there would read as
  // a permission problem the user does not have.
  it('does not read a missing flag as a refusal', async () => {
    fetch.mockClear()
    const w = mountBlock({ recordID: '1', revision: 2 })
    await flushPromises()
    expect(w.text()).not.toContain('block.noPermission')
    expect(fetch).toHaveBeenCalled()
  })
})
