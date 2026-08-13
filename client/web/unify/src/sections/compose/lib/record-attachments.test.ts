import { describe, it, expect } from 'vitest'
import { mergeAttachmentIDs } from './record-attachments'

describe('mergeAttachmentIDs', () => {
  it('appends to a multi-value field', () => {
    expect(mergeAttachmentIDs(['a'], ['b', 'c'], true)).toEqual(['a', 'b', 'c'])
  })

  it('normalises a single stored ID to an array', () => {
    expect(mergeAttachmentIDs('a', ['b'], true)).toEqual(['a', 'b'])
  })

  it('drops empty entries already in the value', () => {
    expect(mergeAttachmentIDs(['a', '', null], ['b'], true)).toEqual(['a', 'b'])
  })

  it('treats an empty value as no attachments', () => {
    expect(mergeAttachmentIDs('', ['b'], true)).toEqual(['b'])
    expect(mergeAttachmentIDs(undefined, ['b'], true)).toEqual(['b'])
  })

  // Appending on a single-value field hands the server two references, which it
  // rejects as an invalid reference format; Record.setValue keeps only the first
  // entry, so the file just uploaded would be silently discarded instead.
  it('replaces the attachment on a single-value field', () => {
    expect(mergeAttachmentIDs('a', ['b'], false)).toEqual(['b'])
  })

  it('keeps the existing attachment when nothing was uploaded', () => {
    expect(mergeAttachmentIDs('a', [], false)).toEqual(['a'])
  })
})
