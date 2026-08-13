import { describe, it, expect } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CFieldFileEditor from './CFieldFileEditor.vue'

// processFiles only reads name/size/type, so a plain object is enough — jsdom's
// File derives its size from the content, which makes multi-megabyte fixtures awkward.
function file(name: string, size: number, type = 'application/pdf') {
  return { name, size, type }
}

// Stands in for the Settings plugin, which serves the server's struct rather
// than the kv keys the settings editor writes.
function settings(attachments: Record<string, unknown>) {
  return {
    get: (k: string) =>
      attachments[k.replace('compose.Record.Attachments.', '') as keyof typeof attachments],
  }
}

type DropOptions = {
  options?: Record<string, unknown>
  global?: Record<string, unknown>
}

async function drop({ options = {}, global }: DropOptions, ...files: ReturnType<typeof file>[]) {
  const wrapper = mount(CFieldFileEditor, {
    props: {
      field: { name: 'attachment', kind: 'File', isMulti: true, options },
      modelValue: [],
    },
    global: { provide: global ? { $Settings: settings(global) } : {} },
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
    const wrapper = await drop({ options: { maxSize: 5 } }, file('report.pdf', 3_000_000))

    expect(wrapper.text()).toContain('report.pdf')
    expect(wrapper.text()).not.toContain('exceeds the maximum size')
  })

  it('rejects a file over the limit and states the limit in MB', async () => {
    const wrapper = await drop({ options: { maxSize: 5 } }, file('huge.pdf', 6_000_000))

    expect(wrapper.text()).toContain('exceeds the maximum size of 5.0 MB')
    expect(wrapper.text()).not.toContain('pending')
  })

  it('falls back to the system-wide limit when maxSize is unset', async () => {
    const wrapper = await drop({ global: { MaxSize: 10 } }, file('huge.pdf', 12_000_000))

    expect(wrapper.text()).toContain('exceeds the maximum size of 10.0 MB')
  })

  it('lets the field option override the system-wide limit', async () => {
    const wrapper = await drop(
      { options: { maxSize: 5 }, global: { MaxSize: 10 } },
      file('huge.pdf', 6_000_000),
    )

    expect(wrapper.text()).toContain('exceeds the maximum size of 5.0 MB')
  })

  it('applies no limit when neither the field nor the settings set one', async () => {
    const wrapper = await drop({ global: { MaxSize: 0 } }, file('huge.pdf', 6_000_000))

    expect(wrapper.text()).toContain('huge.pdf')
    expect(wrapper.text()).not.toContain('exceeds the maximum size')
  })

  it('applies no limit when the settings are unavailable', async () => {
    const wrapper = await drop({}, file('huge.pdf', 6_000_000))

    expect(wrapper.text()).toContain('huge.pdf')
    expect(wrapper.text()).not.toContain('exceeds the maximum size')
  })
})

describe('CFieldFileEditor mime types', () => {
  it('falls back to the system-wide allow list when mimetypes is unset', async () => {
    const wrapper = await drop(
      { global: { Mimetypes: ['image/png', 'image/jpeg'] } },
      file('notes.pdf', 1000, 'application/pdf'),
    )

    expect(wrapper.text()).toContain('unsupported file type')
  })

  it('lets the field option override the system-wide allow list', async () => {
    const wrapper = await drop(
      { options: { mimetypes: 'application/pdf' }, global: { Mimetypes: ['image/png'] } },
      file('notes.pdf', 1000, 'application/pdf'),
    )

    expect(wrapper.text()).toContain('notes.pdf')
    expect(wrapper.text()).not.toContain('unsupported file type')
  })

  it('accepts any type when the system-wide allow list is empty', async () => {
    const wrapper = await drop({ global: { Mimetypes: [] } }, file('notes.pdf', 1000))

    expect(wrapper.text()).toContain('notes.pdf')
    expect(wrapper.text()).not.toContain('unsupported file type')
  })
})
