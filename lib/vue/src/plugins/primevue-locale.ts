// Date presentation for PrimeVue, taken from Intl rather than from any data of
// our own. PrimeVue ships one hard-coded default — US ordering, English names,
// Sunday-first — and the app never configured it.
//
// The UI language and the date locale are separate concerns: someone reading an
// English interface still expects their own date order. So when the user has set
// no explicit language we resolve to the BROWSER's locale, region and all,
// rather than falling back to 'en'.
//
// PrimeVue's date tokens are the jQuery-UI set, not moment's: `dd` zero-padded
// day, `mm` zero-padded month, `yy` FOUR-digit year (`y` is the two-digit one).

// A date whose parts are all distinguishable — day 22, month 11 — so a locale
// that happens to put them in the same position cannot be misread.
const SAMPLE = new Date(Date.UTC(2026, 10, 22))

// 2026-01-04 is a Sunday, so +0..+6 walks Sunday..Saturday — the order PrimeVue
// indexes day names in, whatever the week actually starts on.
const SUNDAY = Date.UTC(2026, 0, 4)
const DAY_MS = 86400000

/**
 * The locale actually in force: the requested one when Intl can resolve it,
 * otherwise the runtime's own. A bare language is left bare — Intl applies
 * CLDR's likely region itself, so 'en' already formats and starts its week the
 * US way without us spelling that out.
 */
export function resolveLocale(locale?: string): string {
  try {
    return new Intl.DateTimeFormat(locale).resolvedOptions().locale
  } catch {
    return new Intl.DateTimeFormat().resolvedOptions().locale
  }
}

/** PrimeVue `dateFormat`, e.g. 'mm/dd/yy' for en-US, 'dd. mm. yy' for sl-SI. */
export function localeDateFormat(locale?: string): string {
  const parts = new Intl.DateTimeFormat(resolveLocale(locale), {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    timeZone: 'UTC',
  }).formatToParts(SAMPLE)

  return parts
    .map(p => {
      if (p.type === 'day') return 'dd'
      if (p.type === 'month') return 'mm'
      if (p.type === 'year') return 'yy'
      // Everything else is a separator. PrimeVue's formatDate reads d/m/y/o/D/M/@/!
      // as tokens wherever they appear, so a separator carrying one has to be
      // quoted — 'de' in a Spanish long form would otherwise become a day and a
      // stray 'e'.
      return /[dmyoDM@!']/.test(p.value) ? `'${p.value.replace(/'/g, "''")}'` : p.value
    })
    .join('')
}

function names(
  locale: string,
  count: number,
  at: (i: number) => Date,
  opts: Intl.DateTimeFormatOptions,
): string[] {
  const fmt = new Intl.DateTimeFormat(locale, { ...opts, timeZone: 'UTC' })
  return Array.from({ length: count }, (_, i) => fmt.format(at(i)))
}

const dayAt = (i: number) => new Date(SUNDAY + i * DAY_MS)
const monthAt = (i: number) => new Date(Date.UTC(2026, i, 1))

/**
 * PrimeVue `firstDayOfWeek` (0 = Sunday), or undefined where the engine carries
 * no week data — `getWeekInfo` is ES2023 and not everywhere yet, so PrimeVue's
 * own default stands in until the browser ships it. Deliberately not backed by a
 * hand-kept CLDR table: that data goes stale silently and nobody notices.
 */
export function localeFirstDayOfWeek(locale?: string): number | undefined {
  try {
    // The RAW tag, not the resolved one. resolveLocale negotiates against the
    // formatting data the runtime actually ships, so a language its ICU lacks
    // collapses to en-US and takes the region with it — 'dv-MV' would answer
    // Sunday instead of Friday. Week data needs no formatting data, only a
    // parseable tag.
    //
    // Both spellings matter: `getWeekInfo()` is the standardised one, `weekInfo`
    // the property engines shipped first — Node 22 has only the latter, and
    // probing for the method alone makes a runtime that HAS the data look like
    // one that does not.
    const l = new Intl.Locale(locale ?? resolveLocale()) as Intl.Locale & {
      getWeekInfo?: () => { firstDay: number }
      weekInfo?: { firstDay: number }
    }
    const firstDay = l.getWeekInfo?.().firstDay ?? l.weekInfo?.firstDay
    if (!firstDay) return undefined
    // Intl counts Monday..Sunday as 1..7; PrimeVue counts Sunday..Saturday as 0..6.
    return firstDay === 7 ? 0 : firstDay
  } catch {
    return undefined
  }
}

/**
 * The partial PrimeVue `locale` config for a language. PrimeVue deep-merges it
 * over its defaults, so every key left out here keeps its built-in value — which
 * is why `firstDayOfWeek` is omitted rather than sent as undefined when unknown.
 *
 * The button and aria strings ('Today', 'Clear', 'Choose Date') are NOT set
 * here: they are app chrome, not locale data, and belong in the `human-webapp`
 * i18n bundle with everything else the user reads.
 */
export function primeVueLocale(locale?: string): Record<string, unknown> {
  const resolved = resolveLocale(locale)
  // Names and format need the resolved tag (they need data the runtime ships);
  // week start takes the raw one, see localeFirstDayOfWeek.
  const firstDayOfWeek = localeFirstDayOfWeek(locale)

  return {
    dateFormat: localeDateFormat(resolved),
    dayNames: names(resolved, 7, dayAt, { weekday: 'long' }),
    dayNamesShort: names(resolved, 7, dayAt, { weekday: 'short' }),
    dayNamesMin: names(resolved, 7, dayAt, { weekday: 'narrow' }),
    monthNames: names(resolved, 12, monthAt, { month: 'long' }),
    monthNamesShort: names(resolved, 12, monthAt, { month: 'short' }),
    ...(firstDayOfWeek === undefined ? {} : { firstDayOfWeek }),
  }
}
