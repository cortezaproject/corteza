import { describe, it, expect, afterEach } from 'vitest'
import {
  localeDateFormat,
  localeFirstDayOfWeek,
  primeVueLocale,
  resolveLocale,
} from './primevue-locale'

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

  it('emits only PrimeVue tokens, never a moment-style one', () => {
    // 'YYYY'/'DD' would be silently wrong rather than an error, so pin the
    // vocabulary: nothing but dd, mm, yy and separators comes out.
    for (const l of ['en-US', 'en-GB', 'de-DE', 'sl', 'lt-LT', 'hu-HU', 'ja-JP']) {
      expect(localeDateFormat(l).replace(/dd|mm|yy/g, '')).not.toMatch(/[dmyoDM@]/)
    }
  })
})

describe('resolveLocale', () => {
  it('keeps a locale Intl can use', () => {
    expect(resolveLocale('en-GB')).toBe('en-GB')
    expect(resolveLocale('sl')).toBe('sl')
  })

  it('falls back to the runtime locale rather than to English', () => {
    // The whole point of the fallback: a user who set no UI language still gets
    // their own browser's date conventions, not en-US.
    const runtime = new Intl.DateTimeFormat().resolvedOptions().locale
    expect(resolveLocale(undefined)).toBe(runtime)
    expect(resolveLocale('not a locale')).toBe(runtime)
  })
})

// This runtime (Node 22) exposes week data as the legacy `weekInfo` property and
// has no `getWeekInfo()` method at all — probing only for the method is what
// made an engine that HAS the data look like one that does not, and cost a
// hand-written CLDR table before it was noticed.
describe('localeFirstDayOfWeek', () => {
  const proto = Intl.Locale.prototype as unknown as Record<string, unknown>
  const descriptors = {
    getWeekInfo: Object.getOwnPropertyDescriptor(proto, 'getWeekInfo'),
    weekInfo: Object.getOwnPropertyDescriptor(proto, 'weekInfo'),
  }

  const stub = (firstDay: number | undefined, via: 'getWeekInfo' | 'weekInfo') => {
    delete proto.getWeekInfo
    delete proto.weekInfo
    if (firstDay === undefined) return
    const value = via === 'getWeekInfo' ? () => ({ firstDay }) : { firstDay }
    Object.defineProperty(proto, via, { value, configurable: true, writable: true })
  }

  afterEach(() => {
    for (const [name, d] of Object.entries(descriptors)) {
      delete proto[name]
      if (d) Object.defineProperty(proto, name, d)
    }
  })

  it('reads the real week data this runtime carries', () => {
    // No stubbing: proves the live path works, not just the mapping.
    expect(localeFirstDayOfWeek('en-US')).toBe(0)
    expect(localeFirstDayOfWeek('en-GB')).toBe(1)
    expect(localeFirstDayOfWeek('sl')).toBe(1)
    expect(localeFirstDayOfWeek('ar-EG')).toBe(6)
    // dv is a language this runtime's ICU cannot format with. Week start must
    // still follow MV, which it only does because the raw tag is used here
    // rather than the format-negotiated one.
    expect(localeFirstDayOfWeek('dv-MV')).toBe(5)
  })

  it('renumbers Intl Sunday (7) to PrimeVue Sunday (0)', () => {
    // Intl counts Monday..Sunday as 1..7, PrimeVue Sunday..Saturday as 0..6. The
    // only value where those disagree is Sunday, and getting it wrong shifts
    // every calendar by a day.
    stub(7, 'getWeekInfo')
    expect(localeFirstDayOfWeek('en-US')).toBe(0)
  })

  it('passes the other days through unchanged', () => {
    for (const firstDay of [1, 5, 6]) {
      stub(firstDay, 'getWeekInfo')
      expect(localeFirstDayOfWeek('de-DE')).toBe(firstDay)
    }
  })

  it('accepts either spelling', () => {
    stub(1, 'weekInfo')
    expect(localeFirstDayOfWeek('de-DE')).toBe(1)
    stub(1, 'getWeekInfo')
    expect(localeFirstDayOfWeek('de-DE')).toBe(1)
  })

  it('returns undefined where the engine has neither', () => {
    stub(undefined, 'getWeekInfo')
    // PrimeVue's own default then stands, which is the intended degradation.
    expect(localeFirstDayOfWeek('de-DE')).toBeUndefined()
    expect(primeVueLocale('de-DE')).not.toHaveProperty('firstDayOfWeek')
  })
})

describe('primeVueLocale', () => {
  it('names the days and months in the locale, Sunday first', () => {
    // Sunday-first regardless of where the week starts: PrimeVue indexes these
    // arrays by getDay(), and rotates them itself.
    const sl = primeVueLocale('sl')
    expect(sl.dayNames).toEqual([
      'nedelja',
      'ponedeljek',
      'torek',
      'sreda',
      'četrtek',
      'petek',
      'sobota',
    ])
    expect((sl.monthNames as string[])[0]).toBe('januar')
    expect((sl.dayNamesMin as string[]).length).toBe(7)
    expect((sl.monthNamesShort as string[]).length).toBe(12)
  })

  it('still speaks English for en', () => {
    const en = primeVueLocale('en-US')
    expect((en.dayNames as string[])[0]).toBe('Sunday')
    expect((en.monthNames as string[])[0]).toBe('January')
    expect(en.dateFormat).toBe('mm/dd/yy')
  })
})
