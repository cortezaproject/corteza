import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref, computed, toValue } from 'vue'

// Field conditions are evaluated server-side against the record every block on
// the page shares. Two consequences this file pins: a conditioned field must not
// appear before the evaluation answers for it, and clearing a hidden field's
// value must not happen where there is no save to carry it.

const route = { name: 'page.record', params: { recordID: 'R1' } }

vi.mock('vue-router', () => ({ useRoute: () => route }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('./PageBlock.vue', () => ({
  default: { props: ['block', 'namespace'], template: '<div><slot /></div>' },
}))

const moduleStore = { getByID: () => module_ }
const recordStore = { findByID: vi.fn(), update: vi.fn(rec => Promise.resolve(rec)) }

vi.mock('@planetcrust/human-vue', () => ({
  // No deferral under test: the spinner's timing is pinned in lib/vue
  useDeferredBusy: src => computed(() => !!toValue(src)),
  useModuleStore: () => moduleStore,
  useRecordStore: () => recordStore,
  components: {
    CFieldViewer: { props: ['field', 'record'], template: '<div class="viewer" />' },
    CFieldEditor: { props: ['field', 'modelValue'], template: '<div class="editor" />' },
  },
}))

vi.mock('@planetcrust/human-js', () => ({
  compose: {
    Record: class {
      constructor() {
        this.values = {}
      }
    },
  },
}))

import RecordBlock from './RecordBlock.vue'

const SECRET = 'F-secret'

let module_
let record
let expressionEvaluate
let deferred

function makeRecord(values) {
  return { recordID: 'R1', values, serialize: () => ({ recordID: 'R1', values }) }
}

beforeEach(() => {
  module_ = {
    moduleID: 'M1',
    namespaceID: 'N1',
    fields: [
      { fieldID: 'F-title', name: 'title', label: 'Title', kind: 'String' },
      { fieldID: SECRET, name: 'secret', label: 'Secret', kind: 'String' },
    ],
  }
  record = makeRecord({ title: 'Hello', secret: 'CLASSIFIED' })
  // Held open so the pre-evaluation render can be inspected
  deferred = {}
  deferred.promise = new Promise(resolve => {
    deferred.resolve = resolve
  })
  expressionEvaluate = vi.fn(() => deferred.promise)
})

// The record the block is handed, reassignable so a test can swap it the way a
// record page does when it navigates to the next record.
let ctxRecord

function mountBlock({ mode = 'view', clearOnHide = false, condition = 'false', adoptSaved } = {}) {
  ctxRecord = ref(record)
  return mount(RecordBlock, {
    props: {
      block: {
        blockID: 'B1',
        kind: 'Record',
        options: {
          fields: [{ name: 'title' }, { name: 'secret' }],
          fieldConditions: [{ field: SECRET, condition, clearOnHide }],
        },
      },
      namespace: { namespaceID: 'N1' },
      page: { pageID: 'P1', moduleID: 'M1' },
    },
    global: {
      stubs: {
        Button: true,
        Select: true,
        ProgressSpinner: true,
        Skeleton: true,
        // scoped slot: the editor branch reads { invalid, error } off it
        FormField: { template: '<div><slot :invalid="false" :error="null" /></div>' },
      },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: {},
        $SystemAPI: { expressionEvaluate },
        $Auth: { user: { userID: 'U1' } },
        $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        recordViewContext: {
          mode: computed(() => mode),
          record: ctxRecord,
          isNew: computed(() => false),
          isSaving: ref(false),
          adoptSaved,
        },
      },
    },
  })
}

const renders = (w, name) => w.html().toLowerCase().includes(name)

describe('RecordBlock field conditions', () => {
  it('withholds a conditioned field until the evaluation answers for it', async () => {
    const w = mountBlock()
    await flushPromises()

    // The request is in flight — 'secret' must not have been on screen yet
    expect(expressionEvaluate).toHaveBeenCalled()
    expect(renders(w, 'secret')).toBe(false)
    expect(renders(w, 'title')).toBe(true)

    deferred.resolve({ [SECRET]: true })
    await flushPromises()

    expect(renders(w, 'secret')).toBe(true)
  })

  it('covers itself while a replacement record waits for its own answers', async () => {
    const w = mountBlock({ condition: 'true' })
    deferred.resolve({ [SECRET]: true })
    await flushPromises()
    expect(renders(w, 'title')).toBe(true)

    // A second record arrives and its evaluation is held open
    deferred = {}
    deferred.promise = new Promise(resolve => {
      deferred.resolve = resolve
    })
    expressionEvaluate.mockImplementation(() => deferred.promise)

    ctxRecord.value = {
      recordID: 'R2',
      values: { title: 'Hello', secret: 'ALSO CLASSIFIED' },
      serialize: () => ({ recordID: 'R2', values: {} }),
    }
    await flushPromises()

    // Not one field is on screen: the fields already carry R1's answers, so the
    // block covers itself rather than showing R2 under them
    expect(renders(w, 'title')).toBe(false)
    expect(renders(w, 'secret')).toBe(false)

    deferred.resolve({ [SECRET]: false })
    await flushPromises()

    expect(renders(w, 'title')).toBe(true)
    expect(renders(w, 'secret')).toBe(false)
  })

  it('shows the first record it is given without covering itself', async () => {
    // Nothing on screen to protect yet, so only the conditioned field waits
    const w = mountBlock()
    await flushPromises()

    expect(renders(w, 'title')).toBe(true)
    expect(renders(w, 'secret')).toBe(false)
  })

  it('evaluates immediately rather than behind the typing debounce', async () => {
    mountBlock()
    await flushPromises()

    // No timer advance: a first paint that waits 300ms is the flash we removed
    expect(expressionEvaluate).toHaveBeenCalledTimes(1)
  })

  it('releases the field when there is nothing to evaluate it with', async () => {
    const w = mount(RecordBlock, {
      props: {
        block: {
          blockID: 'B1',
          kind: 'Record',
          options: { fields: [{ name: 'title' }, { name: 'secret' }], fieldConditions: [] },
        },
        namespace: { namespaceID: 'N1' },
        page: { pageID: 'P1', moduleID: 'M1' },
      },
      global: {
        stubs: {
          Button: true,
          Select: true,
          ProgressSpinner: true,
          Skeleton: true,
          // scoped slot: the editor branch reads { invalid, error } off it
          FormField: { template: '<div><slot :invalid="false" :error="null" /></div>' },
        },
        directives: { tooltip: {} },
        mocks: { $t: k => k },
        provide: {
          $ComposeAPI: {},
          $SystemAPI: { expressionEvaluate },
          $Auth: { user: { userID: 'U1' } },
          $toast: null,
          recordViewContext: {
            mode: computed(() => 'view'),
            record: ref(record),
            isNew: computed(() => false),
            isSaving: ref(false),
          },
        },
      },
    })
    await flushPromises()

    expect(expressionEvaluate).not.toHaveBeenCalled()
    expect(renders(w, 'secret')).toBe(true)
  })

  it('does not clear a hidden field in view mode', async () => {
    mountBlock({ mode: 'view', clearOnHide: true })
    await flushPromises()
    deferred.resolve({ [SECRET]: false })
    await flushPromises()

    // Nothing is being saved, so the write would only take the value off the
    // screen — including out of any other block showing the same field.
    expect(record.values.secret).toBe('CLASSIFIED')
  })

  it('clears a hidden field in edit mode', async () => {
    mountBlock({ mode: 'edit', clearOnHide: true })
    await flushPromises()
    deferred.resolve({ [SECRET]: false })
    await flushPromises()

    expect(record.values.secret).toBeUndefined()
  })

  it("hands an inline save's response back to the record view", async () => {
    // Inline edit is the one save path the record view does not run itself, so
    // without this the view keeps the record it had before the write.
    const saved = makeRecord({ title: 'Hello', secret: 'CLASSIFIED', computed: 'by-server' })
    recordStore.update.mockResolvedValueOnce(saved)
    const adoptSaved = vi.fn()

    const w = mountBlock({ adoptSaved })
    await flushPromises()
    deferred.resolve({ [SECRET]: true })
    await flushPromises()

    await w.vm.saveInlineEdits()

    expect(adoptSaved).toHaveBeenCalledWith(saved)
  })

  it('leaves the value alone in edit mode when clearOnHide is off', async () => {
    mountBlock({ mode: 'edit', clearOnHide: false })
    await flushPromises()
    deferred.resolve({ [SECRET]: false })
    await flushPromises()

    expect(record.values.secret).toBe('CLASSIFIED')
  })
})
