import { describe, expect, it } from 'vitest'
import { isAppUrlLocal, resolveAppUrl } from './appUrl'

describe('resolveAppUrl', () => {
  it('leaves a path this shell serves alone', () => {
    expect(resolveAppUrl('/project')).toBe('/project')
    expect(resolveAppUrl('/admin/system/labels')).toBe('/admin/system/labels')
  })

  it('gives a slash-less shell path its leading slash', () => {
    expect(resolveAppUrl('admin/')).toBe('/admin/')
    expect(resolveAppUrl('admin/system/labels')).toBe('/admin/system/labels')
  })

  it('leaves an address that carries a scheme alone', () => {
    expect(resolveAppUrl('https://www.google.com')).toBe('https://www.google.com')
    expect(resolveAppUrl('http://example.com/x')).toBe('http://example.com/x')
    expect(resolveAppUrl('//example.com')).toBe('//example.com')
  })

  it('gives a bare host a scheme rather than reading it as a path', () => {
    // Without this the anchor is relative and the link lands on the shell's
    // catch-all instead of google.
    expect(resolveAppUrl('www.google.com')).toBe('https://www.google.com')
    expect(resolveAppUrl('example.com/some/page')).toBe('https://example.com/some/page')
  })

  it('answers empty for no url', () => {
    expect(resolveAppUrl('')).toBe('')
    expect(resolveAppUrl(undefined)).toBe('')
    expect(resolveAppUrl(null)).toBe('')
  })
})

describe('isAppUrlLocal', () => {
  it('is true only for a path this shell serves', () => {
    expect(isAppUrlLocal('admin/')).toBe(true)
    expect(isAppUrlLocal('/project')).toBe(true)

    expect(isAppUrlLocal('www.google.com')).toBe(false)
    expect(isAppUrlLocal('https://www.google.com')).toBe(false)
    // Protocol-relative: a leading slash that is not a local path.
    expect(isAppUrlLocal('//example.com')).toBe(false)
    expect(isAppUrlLocal('')).toBe(false)
  })
})
