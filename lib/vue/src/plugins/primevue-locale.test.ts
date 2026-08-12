import { describe, it, expect } from 'vitest'
import { localeDateFormat, localeFirstDayOfWeek, primeVueLocale } from './primevue-locale'

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

// PrimeVue counts Sunday..Saturday as 0..6, where Intl's own week data counts
// Monday..Sunday as 1..7 — an off-by-one here silently shifts every calendar.
describe('localeFirstDayOfWeek', () => {
  it('starts the week on Sunday for the US', () => {
    expect(localeFirstDayOfWeek('en-US')).toBe(0)
    expect(localeFirstDayOfWeek('en')).toBe(0)
  })

  it('starts the week on Monday across Europe', () => {
    for (const l of ['en-GB', 'de-DE', 'sl', 'fr', 'lt-LT', 'hu-HU', 'pl']) {
      expect(localeFirstDayOfWeek(l)).toBe(1)
    }
  })

  it('starts the week on Saturday where CLDR says so', () => {
    expect(localeFirstDayOfWeek('ar-EG')).toBe(6)
    expect(localeFirstDayOfWeek('fa-IR')).toBe(6)
  })

  it('starts the week on Friday in the Maldives', () => {
    expect(localeFirstDayOfWeek('dv-MV')).toBe(5)
  })

  it('follows the region, not the language', () => {
    // pt widens to pt-BR (Sunday); pt-PT is Sunday too per CLDR, while a
    // Portuguese speaker in Angola gets Monday.
    expect(localeFirstDayOfWeek('pt-PT')).toBe(0)
    expect(localeFirstDayOfWeek('pt-AO')).toBe(1)
    // en is US-flavoured by default but not in Ireland.
    expect(localeFirstDayOfWeek('en-IE')).toBe(1)
  })

  it('reflects the CLDR changes that memory gets wrong', () => {
    // The UAE moved to a Mon-Fri week in 2022 and CLDR followed; China has
    // always been Monday-first. Both are easy to mis-remember, and both were
    // wrong in the first draft of the table.
    expect(localeFirstDayOfWeek('ar-AE')).toBe(1)
    expect(localeFirstDayOfWeek('zh-CN')).toBe(1)
  })

  it('matches the date-format fallback for an unusable locale', () => {
    expect(localeFirstDayOfWeek('not a locale')).toBe(0)
  })
})

describe('primeVueLocale', () => {
  it('carries the date format and the week start, and nothing else', () => {
    // Every other key must be absent rather than undefined: PrimeVue deep-merges
    // this over its defaults, and an explicit undefined would erase one.
    expect(primeVueLocale('en-GB')).toEqual({ dateFormat: 'dd/mm/yy', firstDayOfWeek: 1 })
    expect(primeVueLocale('en-US')).toEqual({ dateFormat: 'mm/dd/yy', firstDayOfWeek: 0 })
  })
})
