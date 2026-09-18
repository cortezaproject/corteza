import { describe, it, expect, vi, beforeEach } from 'vitest'
import { flushPromises } from '@vue/test-utils'
import { mountWithContext, createTestPinia } from '@planetcrust/human-test-utils'
import CFieldFileViewer from './CFieldFileViewer.vue'

function field(overrides: Record<string, unknown> = {}) {
  return { name: 'docs', kind: 'File', isSystem: false, isMulti: true, options: {}, ...overrides }
}

function attachment(attachmentID: string, name: string) {
  return {
    attachmentID,
    name,
    url: `/namespace/1/attachment/record/${attachmentID}/original/${name}?sign=x&userID=1`,
    meta: { original: { size: 10, mimetype: 'text/plain' } },
  }
}

// A read that answers only when told to, so a test can see which reads are
// under way at the same time.
function deferredReads() {
  const pending = new Map<string, { resolve: (v: unknown) => void; reject: (e: unknown) => void }>()
  const attachmentRead = vi.fn(
    ({ attachmentID }: { attachmentID: string }) =>
      new Promise((resolve, reject) => pending.set(attachmentID, { resolve, reject })),
  )
  return { attachmentRead, pending }
}

function mount(values: Record<string, unknown>, composeAPI: Record<string, unknown>) {
  return mountWithContext(
    CFieldFileViewer,
    { composeAPI: { baseURL: '', ...composeAPI } },
    { props: { field: field(), record: { values }, namespace: { namespaceID: '1' } } },
  )
}

describe('CFieldFileViewer', () => {
  beforeEach(() => {
    createTestPinia()
  })

  it('reads every attachment at once and lists them in value order', async () => {
    const { attachmentRead, pending } = deferredReads()
    const wrapper = mount({ docs: ['11', '12', '13'] }, { attachmentRead })
    await flushPromises()

    expect(attachmentRead).toHaveBeenCalledTimes(3)

    pending.get('13')!.resolve(attachment('13', 'c.txt'))
    pending.get('11')!.resolve(attachment('11', 'a.txt'))
    pending.get('12')!.resolve(attachment('12', 'b.txt'))
    await flushPromises()

    const text = wrapper.text()
    expect(text.indexOf('a.txt')).toBeLessThan(text.indexOf('b.txt'))
    expect(text.indexOf('b.txt')).toBeLessThan(text.indexOf('c.txt'))
  })

  it('shows the ID of an attachment that cannot be read', async () => {
    const attachmentRead = vi.fn(({ attachmentID }: { attachmentID: string }) =>
      attachmentID === '21'
        ? Promise.reject(new Error('gone'))
        : Promise.resolve(attachment(attachmentID, 'ok.txt')),
    )
    const wrapper = mount({ docs: ['21', '22'] }, { attachmentRead })
    await flushPromises()

    expect(wrapper.text()).toContain('21')
    expect(wrapper.text()).toContain('ok.txt')
  })

  it('keeps the newer attachments when an older read answers last', async () => {
    const { attachmentRead, pending } = deferredReads()
    const wrapper = mount({ docs: ['31'] }, { attachmentRead })
    await flushPromises()

    await wrapper.setProps({ record: { values: { docs: ['32'] } } })
    await flushPromises()

    pending.get('32')!.resolve(attachment('32', 'new.txt'))
    await flushPromises()
    pending.get('31')!.resolve(attachment('31', 'old.txt'))
    await flushPromises()

    expect(wrapper.text()).toContain('new.txt')
    expect(wrapper.text()).not.toContain('old.txt')
  })
})
