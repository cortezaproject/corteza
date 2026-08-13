import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CFieldFileEditor from './CFieldFileEditor.vue'

// processFiles only reads name/size/type, so a plain object is enough — jsdom's
// File derives its size from the content, which makes multi-megabyte fixtures awkward.
function file(name: string, size: number) {
  return { name, size, type: 'application/pdf' }
}

function fieldWithMaxSize(maxSize: number) {
  return { name: 'attachment', kind: 'File', isMulti: true, options: { maxSize } }
}

async function drop(maxSize: number, ...files: ReturnType<typeof file>[]) {
  const wrapper = mount(CFieldFileEditor, {
    props: { field: fieldWithMaxSize(maxSize), modelValue: [] },
  })
  await flushPromises()

  await wrapper.find('.border-dashed').trigger('drop', { dataTransfer: { files } })
  await flushPromises()

  return wrapper
}

describe('CFieldFileEditor max file size', () => {
  it('reads the maxSize option as megabytes, not bytes', async () => {
    // 3 MB against a 5 MB limit — read as bytes this would reject everything
    // larger than five bytes.
    const wrapper = await drop(5, file('report.pdf', 3_000_000))

    expect(wrapper.text()).toContain('report.pdf')
    expect(wrapper.text()).not.toContain('exceeds the maximum size')
  })

  it('rejects a file over the limit and states the limit in MB', async () => {
    const wrapper = await drop(5, file('huge.pdf', 6_000_000))

    expect(wrapper.text()).toContain('exceeds the maximum size of 5.0 MB')
    expect(wrapper.text()).not.toContain('pending')
  })

  it('applies no limit when maxSize is unset', async () => {
    const wrapper = await drop(0, file('huge.pdf', 6_000_000))

    expect(wrapper.text()).toContain('huge.pdf')
    expect(wrapper.text()).not.toContain('exceeds the maximum size')
  })
})
