import { describe, it, expect } from 'vitest'
import { trimUrlFragment, trimUrlQuery, trimUrlPath, onlySecureUrl } from './url'

describe('trimUrlFragment', () => {
  it('removes the fragment', () => {
    expect(trimUrlFragment('https://www.example.com/page?foo=bar#section')).toBe(
      'https://www.example.com/page?foo=bar',
    )
  })

  it('passes empty values through', () => {
    expect(trimUrlFragment('')).toBe('')
  })
})

describe('trimUrlQuery', () => {
  it('removes the query string', () => {
    expect(trimUrlQuery('https://www.example.com/page?foo=bar#section')).toBe(
      'https://www.example.com/page#section',
    )
  })
})

describe('trimUrlPath', () => {
  it('trims down to the domain, dropping path, query and fragment', () => {
    expect(trimUrlPath('https://www.example.com/page/employee?foo=bar#section')).toBe(
      'https://www.example.com',
    )
  })

  it('keeps the port', () => {
    expect(trimUrlPath('https://www.example.com:8443/page?foo=bar')).toBe(
      'https://www.example.com:8443',
    )
  })

  it('keeps the scheme of insecure URLs', () => {
    expect(trimUrlPath('http://www.example.com/page#section')).toBe('http://www.example.com')
  })

  it('assumes https when no protocol is given', () => {
    expect(trimUrlPath('www.example.com/page?foo=bar')).toBe('https://www.example.com')
  })

  it('drops credentials', () => {
    expect(trimUrlPath('https://user:pass@www.example.com/page')).toBe('https://www.example.com')
  })

  it('passes empty values through', () => {
    expect(trimUrlPath('')).toBe('')
  })
})

describe('onlySecureUrl', () => {
  it('upgrades http to https', () => {
    expect(onlySecureUrl('http://www.example.com/page')).toBe('https://www.example.com/page')
  })
})
