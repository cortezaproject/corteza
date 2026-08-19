import { describe, it, expect, vi, afterEach } from 'vitest'
import { useFileUpload } from './useFileUpload'

const file = () => new File(['x'], 'logo.png', { type: 'image/png' })

function respond(body: string, { ok = true, status = 200 } = {}) {
  // Both readers, as a Response has: with only text() the old implementation
  // would fail on the mock rather than on what it got wrong.
  const fetchMock = vi.fn().mockResolvedValue({
    ok,
    status,
    text: async () => body,
    json: async () => JSON.parse(body),
  })
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

afterEach(() => vi.unstubAllGlobals())

describe('uploadFileRaw', () => {
  it('returns the response payload of a successful upload', async () => {
    respond(JSON.stringify({ response: { attachmentID: '42' } }))
    const { uploadFileRaw, uploadError } = useFileUpload()

    await expect(uploadFileRaw(file(), { url: '/upload' })).resolves.toEqual({ attachmentID: '42' })
    expect(uploadError.value).toBe('')
  })

  it('asks for JSON, so the server does not answer in plain text', async () => {
    const fetchMock = respond(JSON.stringify({ response: {} }))
    const { uploadFileRaw } = useFileUpload()

    await uploadFileRaw(file(), { url: '/upload', token: 't' })

    const { headers } = fetchMock.mock.calls[0][1]
    expect(headers.Accept).toBe('application/json')
    expect(headers.Authorization).toBe('Bearer t')
  })

  // The API answers a failed upload with HTTP 200 and an error payload. Read as
  // a success, that payload becomes the attachment the caller stores.
  it('fails on an error payload served under HTTP 200', async () => {
    respond(JSON.stringify({ error: { message: 'could not process image' } }))
    const { uploadFileRaw, uploadError } = useFileUpload()

    await expect(uploadFileRaw(file(), { url: '/upload' })).rejects.toThrow(
      'could not process image',
    )
    expect(uploadError.value).toBe('could not process image')
  })

  it('reports the message of a plain-text error, not the parse failure', async () => {
    respond(
      'Error: could not process image\n' +
        '----------------------------------------\n' +
        'resource: system:attachment\n' +
        'type:     failedToProcessImage\n',
    )
    const { uploadFileRaw, uploadError } = useFileUpload()

    await expect(uploadFileRaw(file(), { url: '/upload' })).rejects.toThrow(
      'could not process image',
    )
    expect(uploadError.value).toBe('could not process image')
  })

  it('names the file when the failure carries no message', async () => {
    respond('', { ok: false, status: 500 })
    const { uploadFileRaw, uploadError } = useFileUpload()

    await expect(uploadFileRaw(file(), { url: '/upload' })).rejects.toThrow(
      'Upload failed for logo.png',
    )
    expect(uploadError.value).toBe('Upload failed for logo.png')
  })

  it('clears the uploading flag whichever way the upload ends', async () => {
    respond(JSON.stringify({ error: { message: 'nope' } }))
    const { uploadFileRaw, uploading } = useFileUpload()

    await uploadFileRaw(file(), { url: '/upload' }).catch(() => {})

    expect(uploading.value).toBe(false)
  })
})
