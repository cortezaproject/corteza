import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mountWithContext } from '@planetcrust/human-test-utils'

// The component reads from the pinia record/module stores directly (inject
// pattern) and delegates label rendering to CFieldViewer. Mock the stores so we
// can drive resolvedRecords / labelFieldDef, and stub CFieldViewer to surface
// the value it would render.
const getByID = vi.fn()
const resolveRecordLabels = vi.fn().mockResolvedValue(undefined)
const moduleGetByID = vi.fn()

vi.mock('../../../stores/useRecordStore', () => ({
  useRecordStore: () => ({ getByID, resolveRecordLabels }),
}))
vi.mock('../../../stores/useModuleStore', () => ({
  useModuleStore: () => ({ getByID: moduleGetByID }),
}))
vi.mock('../../../stores/usePageStore', () => ({
  usePageStore: () => ({ set: [] }),
}))

import CFieldRecordViewer from './CFieldRecordViewer.vue'

const CFieldViewerStub = {
  name: 'CFieldViewer',
  props: ['field', 'record', 'namespace', 'disableClick', 'valueOnly'],
  template: '<span class="cfv">{{ record.values[field.name] }}</span>',
}

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'ref', isMulti: false, isSystem: false, options: { moduleID: 'm1' }, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function mountViewer(props: Record<string, unknown>) {
  return mountWithContext(CFieldRecordViewer, { namespace: { namespaceID: 'ns1' } }, {
    props: { namespace: { namespaceID: 'ns1' }, ...props },
    global: { stubs: { CFieldViewer: CFieldViewerStub } },
  })
}

describe('CFieldRecordViewer', () => {
  beforeEach(() => {
    getByID.mockReset()
    resolveRecordLabels.mockClear()
    moduleGetByID.mockReset()
    // Default module with two fields; first is the implicit label field.
    moduleGetByID.mockReturnValue({ fields: [{ name: 'title' }, { name: 'email' }] })
  })

  describe('no resolved record', () => {
    it('renders empty when no value', () => {
      getByID.mockReturnValue(null)
      const wrapper = mountViewer({ field: field(), record: record() })
      expect(wrapper.findAll('span')).toHaveLength(0)
    })

    it('shows recordID as fallback when record not in store', () => {
      getByID.mockReturnValue(null)
      const wrapper = mountViewer({ field: field(), record: record({ ref: 'rec-abc' }) })
      expect(wrapper.find('span').text()).toBe('rec-abc')
    })

    it('renders multiple IDs as fallback', () => {
      getByID.mockReturnValue(null)
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ ref: ['r1', 'r2'] }),
      })
      const spans = wrapper.findAll('span')
      expect(spans).toHaveLength(2)
      expect(spans[0].text()).toContain('r1')
      expect(spans[1].text()).toContain('r2')
    })
  })

  describe('with resolved record', () => {
    it('renders the implicit first field via CFieldViewer', () => {
      getByID.mockReturnValue({ recordID: 'r1', values: { title: 'My Record', email: 'a@test.com' } })
      const wrapper = mountViewer({ field: field(), record: record({ ref: 'r1' }) })
      expect(wrapper.find('.cfv').text()).toBe('My Record')
    })

    it('uses labelField from field options', () => {
      getByID.mockReturnValue({ recordID: 'r1', values: { title: 'My Record', email: 'a@test.com' } })
      const wrapper = mountViewer({
        field: field({ options: { moduleID: 'm1', labelField: 'email' } }),
        record: record({ ref: 'r1' }),
      })
      expect(wrapper.find('.cfv').text()).toBe('a@test.com')
    })

    it('falls back to recordID when module has no field definition', () => {
      moduleGetByID.mockReturnValue(null)
      getByID.mockReturnValue({ recordID: 'r1', values: { title: 'My Record' } })
      const wrapper = mountViewer({ field: field(), record: record({ ref: 'r1' }) })
      expect(wrapper.find('span').text()).toBe('r1')
    })
  })

  describe('delimiter', () => {
    it('uses default ", " delimiter between multi values', () => {
      getByID.mockImplementation((id: string) => ({ recordID: id, values: { title: id === 'a' ? 'Rec A' : 'Rec B' } }))
      const wrapper = mountViewer({
        field: field({ isMulti: true }),
        record: record({ ref: ['a', 'b'] }),
      })
      expect(wrapper.text()).toContain('Rec A')
      expect(wrapper.text()).toContain('Rec B')
      expect(wrapper.text()).toContain(',')
    })

    it('uses custom delimiter from field options', () => {
      getByID.mockImplementation((id: string) => ({ recordID: id, values: { title: id } }))
      const wrapper = mountViewer({
        field: field({ isMulti: true, options: { moduleID: 'm1', multiDelimiter: ' | ' } }),
        record: record({ ref: ['a', 'b'] }),
      })
      expect(wrapper.text()).toContain('|')
    })
  })

  describe('disableClick', () => {
    it('does not add link styling when disableClick is true', () => {
      getByID.mockReturnValue({ recordID: 'r1', values: { title: 'T' } })
      const wrapper = mountViewer({
        field: field(),
        record: record({ ref: 'r1' }),
        disableClick: true,
      })
      expect(wrapper.find('span').classes()).not.toContain('cursor-pointer')
    })

    it('adds link styling when disableClick is false and record has ID', () => {
      getByID.mockReturnValue({ recordID: 'r1', values: { title: 'T' } })
      const wrapper = mountViewer({
        field: field(),
        record: record({ ref: 'r1' }),
        disableClick: false,
      })
      expect(wrapper.find('span').classes()).toContain('cursor-pointer')
    })
  })

  describe('label resolution', () => {
    it('asks the store to resolve missing labels on mount', () => {
      getByID.mockReturnValue(null)
      mountViewer({ field: field(), record: record({ ref: 'r9' }) })
      expect(resolveRecordLabels).toHaveBeenCalledWith({
        namespaceID: 'ns1',
        moduleID: 'm1',
        recordIDs: ['r9'],
      })
    })
  })
})
