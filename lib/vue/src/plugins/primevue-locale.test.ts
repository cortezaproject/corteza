import { describe, it, expect } from 'vitest'
import { localeDateFormat, primeVueLocale } from './primevue-locale'

// PrimeVue's tokens are jQuery-UI's, not moment's: dd = padded day, mm = padded
// month, yy = FOUR-digit year. Getting `yy` wrong is the easy mistake and it is
// silent — it renders a plausible-looking two-digit year.
describe('localeDateFormat', () => {
  it('keeps the US order for en-US', () => {
    expect(localeDateFormat('en-US')).toBe('mm/dd/yy')
  })

  it('puts the day first for en-GB', () => {
    expect(localeDateFormat('en-GB')).toBe('dd/mm/yy')
  })

  it('follows a dot-separated European locale, spacing included', () => {
    expect(localeDateFormat('de-DE')).toBe('dd.mm.yy')
    expect(localeDateFormat('sl')).toBe('dd. mm. yy')
  })

  it('orders year first where the locale does', () => {
    expect(localeDateFormat('lt-LT')).toBe('yy-mm-dd')
    expect(localeDateFormat('hu-HU')).toBe('yy. mm. dd.')
  })

  it('falls back to the PrimeVue default for an unusable locale', () => {
    expect(localeDateFormat('not a locale')).toBe('mm/dd/yy')
  })

  it('emits only PrimeVue tokens, never a moment-style one', () => {
    // 'YYYY'/'DD' would be silently wrong rather than an error, so pin the
    // vocabulary: nothing but dd, mm, yy and separators comes out.
    for (const l of ['en-US', 'en-GB', 'de-DE', 'sl', 'lt-LT', 'hu-HU', 'ja-JP']) {
      const f = localeDateFormat(l)
      expect(f.replace(/dd|mm|yy/g, '')).not.toMatch(/[dmyoDM@]/)
    }
  })
})

describe('primeVueLocale', () => {
  it('carries the date format and nothing else', () => {
    // Every other key must be absent rather than undefined: PrimeVue deep-merges
    // this over its defaults, and an explicit undefined would erase one.
    expect(primeVueLocale('en-GB')).toEqual({ dateFormat: 'dd/mm/yy' })
  })
})
