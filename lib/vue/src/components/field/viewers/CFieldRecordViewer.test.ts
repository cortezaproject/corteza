import { describe, it, expect, vi } from 'vitest'
import { mountWithContext } from '@planetcrust/human-test-utils'
import CFieldRecordViewer from './CFieldRecordViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'ref', isMulti: false, isSystem: false, options: {}, ...overrides }
}

function record(values: Record<string, unknown> = {}) {
  return { values }
}

function makeRecordStore(entries: Record<string, unknown>[] = []) {
  const map = new Map(entries.map(r => [(r as any).recordID, r]))
  return {
    getByID: (id: string) => map.get(id) || null,
    resolveRecordLabels: vi.fn().mockResolvedValue(undefined),
  }
}

describe('CFieldRecordViewer', () => {
  describe('no recordStore', () => {
    it('renders empty when no value', () => {
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: null }, {
        props: { field: field(), record: record() },
      })
      expect(wrapper.findAll('span')).toHaveLength(0)
    })

    it('shows recordID as fallback when no store', () => {
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: null }, {
        props: { field: field(), record: record({ ref: 'rec-abc' }) },
      })
      expect(wrapper.find('span').text()).toBe('rec-abc')
    })

    it('renders multiple IDs as fallback', () => {
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: null }, {
        props: { field: field({ isMulti: true }), record: record({ ref: ['r1', 'r2'] }) },
      })
      const spans = wrapper.findAll('span')
      expect(spans).toHaveLength(2)
      expect(spans[0].text()).toContain('r1')
      expect(spans[1].text()).toContain('r2')
    })
  })

  describe('with recordStore — object values (compose.Record shape)', () => {
    it('shows first field value from object values', () => {
      const store = makeRecordStore([
        { recordID: 'r1', values: { title: 'My Record', other: 'extra' } },
      ])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: { field: field(), record: record({ ref: 'r1' }) },
      })
      expect(wrapper.find('span').text()).toBe('My Record')
    })

    it('uses labelField from field options (object values)', () => {
      const store = makeRecordStore([
        { recordID: 'r1', values: { name: 'Alice', email: 'alice@test.com' } },
      ])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: {
          field: field({ options: { labelField: 'email' } }),
          record: record({ ref: 'r1' }),
        },
      })
      expect(wrapper.find('span').text()).toBe('alice@test.com')
    })
  })

  describe('with recordStore — array values (raw API shape)', () => {
    it('shows first value from array values', () => {
      const store = makeRecordStore([
        { recordID: 'r2', values: [{ name: 'title', value: 'Raw Record' }] },
      ])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: { field: field(), record: record({ ref: 'r2' }) },
      })
      expect(wrapper.find('span').text()).toBe('Raw Record')
    })

    it('uses labelField from field options (array values)', () => {
      const store = makeRecordStore([
        {
          recordID: 'r3',
          values: [
            { name: 'name', value: 'Bob' },
            { name: 'email', value: 'bob@test.com' },
          ],
        },
      ])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: {
          field: field({ options: { labelField: 'email' } }),
          record: record({ ref: 'r3' }),
        },
      })
      expect(wrapper.find('span').text()).toBe('bob@test.com')
    })
  })

  describe('delimiter', () => {
    it('uses default ", " delimiter between multi values', () => {
      const store = makeRecordStore([
        { recordID: 'a', values: { title: 'Rec A' } },
        { recordID: 'b', values: { title: 'Rec B' } },
      ])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: {
          field: field({ isMulti: true }),
          record: record({ ref: ['a', 'b'] }),
        },
      })
      expect(wrapper.text()).toContain('Rec A')
      expect(wrapper.text()).toContain('Rec B')
      expect(wrapper.text()).toContain(', ')
    })

    it('uses custom delimiter from field options', () => {
      const store = makeRecordStore([
        { recordID: 'a', values: { title: 'X' } },
        { recordID: 'b', values: { title: 'Y' } },
      ])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: {
          field: field({ isMulti: true, options: { multiDelimiter: ' | ' } }),
          record: record({ ref: ['a', 'b'] }),
        },
      })
      expect(wrapper.text()).toContain(' | ')
    })
  })

  describe('disableClick', () => {
    it('does not add record-link class when disableClick is true', () => {
      const store = makeRecordStore([{ recordID: 'r1', values: { title: 'T' } }])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: {
          field: field(),
          record: record({ ref: 'r1' }),
          disableClick: true,
        },
      })
      expect(wrapper.find('span').classes()).not.toContain('record-link')
    })

    it('adds record-link class when disableClick is false and record has ID', () => {
      const store = makeRecordStore([{ recordID: 'r1', values: { title: 'T' } }])
      const wrapper = mountWithContext(CFieldRecordViewer, { recordStore: store }, {
        props: {
          field: field(),
          record: record({ ref: 'r1' }),
          disableClick: false,
        },
      })
      expect(wrapper.find('span').classes()).toContain('record-link')
    })
  })
})
