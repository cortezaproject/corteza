import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { ref, computed, toValue } from 'vue'

// A system field is a property of the record, not one of its values. This file
// pins both directions of that: what the editor is handed, and where what it
// emits is written.

const route = { name: 'page.record', params: { recordID: 'R1' } }

vi.mock('vue-router', () => ({ useRoute: () => route }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ t: k => k }) }))
vi.mock('./PageBlock.vue', () => ({
  default: { props: ['block', 'namespace'], template: '<div><slot /></div>' },
}))

const moduleStore = { getByID: () => module_ }
const recordStore = { findByID: vi.fn(), update: vi.fn(rec => Promise.resolve(rec)) }

vi.mock('@planetcrust/human-vue', () => ({
  useDeferredBusy: src => computed(() => !!toValue(src)),
  useModuleStore: () => moduleStore,
  useRecordStore: () => recordStore,
  components: {
    CFieldViewer: {
      props: ['field', 'record'],
      template: '<div class="viewer" :data-field="field.name" :data-owner="record.ownedBy" />',
    },
    CFieldEditor: {
      props: ['field', 'modelValue'],
      emits: ['update:modelValue'],
      template:
        '<div class="editor" :data-field="field.name" :data-value="String(modelValue)"' +
        " @click=\"$emit('update:modelValue', 'U-next')\" />",
    },
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

const OWNED_BY = { name: 'ownedBy', label: 'Owned by', kind: 'User', isSystem: true }
const CREATED_AT = { name: 'createdAt', label: 'Created at', kind: 'DateTime', isSystem: true }
const TITLE = { fieldID: 'F-title', name: 'title', label: 'Title', kind: 'String' }

let module_
let record

beforeEach(() => {
  module_ = {
    moduleID: 'M1',
    namespaceID: 'N1',
    fields: [TITLE],
    filterFields: names => [TITLE, OWNED_BY, CREATED_AT].filter(f => names.includes(f.name)),
  }
  record = {
    recordID: 'R1',
    values: { title: 'Hello' },
    ownedBy: 'U-owner',
    createdAt: '2026-08-19T16:35:00Z',
    canManageOwnerOnRecord: true,
    setValue: vi.fn(function (name, value) {
      record.values[name] = value
    }),
    serialize: () => ({ recordID: 'R1', values: record.values }),
  }
})

function mountBlock({ mode = 'edit' } = {}) {
  return mount(RecordBlock, {
    props: {
      block: {
        blockID: 'B1',
        kind: 'Record',
        options: { fields: ['title', 'ownedBy', 'createdAt'] },
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
        FormField: { template: '<div><slot :invalid="false" :error="null" /></div>' },
      },
      directives: { tooltip: {} },
      mocks: { $t: k => k },
      provide: {
        $ComposeAPI: {},
        $SystemAPI: { expressionEvaluate: vi.fn() },
        $Auth: { user: { userID: 'U1' } },
        $toast: { toastSuccess: vi.fn(), toastErrorHandler: () => vi.fn() },
        recordViewContext: {
          mode: computed(() => mode),
          record: ref(record),
          isNew: computed(() => false),
          isSaving: ref(false),
        },
      },
    },
  })
}

describe('RecordBlock system fields', () => {
  it('hands the owner editor the record property, not a value', async () => {
    const w = mountBlock()
    await flushPromises()

    const editor = w.find('.editor[data-field="ownedBy"]')
    expect(editor.exists()).toBe(true)
    expect(editor.attributes('data-value')).toBe('U-owner')
  })

  it('writes an edited owner onto the record', async () => {
    const w = mountBlock()
    await flushPromises()

    await w.find('.editor[data-field="ownedBy"]').trigger('click')

    expect(record.ownedBy).toBe('U-next')
    // values is what the save serializes as module fields; an owner stored
    // there is stored under a name nothing reads
    expect(record.values.ownedBy).toBeUndefined()
    expect(record.setValue).not.toHaveBeenCalledWith('ownedBy', expect.anything())
  })

  it('leaves a module field going through values', async () => {
    const w = mountBlock()
    await flushPromises()

    await w.find('.editor[data-field="title"]').trigger('click')

    expect(record.setValue).toHaveBeenCalledWith('title', 'U-next')
    expect(record.title).toBeUndefined()
  })

  it('draws a read-only system field with its viewer', async () => {
    const w = mountBlock()
    await flushPromises()

    const viewer = w.find('.viewer[data-field="createdAt"]')
    expect(viewer.exists()).toBe(true)
    expect(w.find('.editor[data-field="createdAt"]').exists()).toBe(false)
  })
})
