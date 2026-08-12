// Date presentation for PrimeVue, derived from the signed-in user's language.
//
// PrimeVue ships one hard-coded default — 'mm/dd/yy' — and the app never
// configured it, so every date input in Human read as US-formatted whatever the
// user's language. Nothing here picks a house format: the ordering and
// separators come from Intl for the locale, so a user on `sl` sees
// `20. 08. 2026` and one on `en-US` still sees `08/20/2026`.
//
// PrimeVue's date tokens are the jQuery-UI set, not moment's: `dd` zero-padded
// day, `mm` zero-padded month, `yy` FOUR-digit year (`y` is the two-digit one).
//
// Week start comes from the table below rather than from Intl's `getWeekInfo`,
// which is missing from Firefox and from Node: reading it at runtime would give
// the same user a different calendar in different browsers, and leave this
// untestable in CI.

// A date whose parts are all distinguishable — day 22, month 11 — so a locale
// that happens to put them in the same position cannot be misread.
const SAMPLE = new Date(Date.UTC(2026, 10, 22))

// CLDR week data, keyed by region. Extracted from Chromium's ICU on 2026-08-12
// (`new Intl.Locale('und-XX').getWeekInfo().firstDay` over every ISO region)
// rather than written from memory, which had UAE on Saturday — it moved to
// Monday after the 2022 workweek change — and China on Sunday, which CLDR has
// never said. Regenerate the same way if it needs refreshing; obsolete ISO codes
// (BU, RH, YD) are dropped.
//
// Only the regions that are NOT Monday-first are listed: Monday is CLDR's own
// `001` default and covers 200-odd of them.
const SUNDAY_FIRST = new Set(
  `AG AS BD BR BS BT BW BZ CA CO DM DO ET GT GU HK HN ID IL IN IS JM JP KE KH KR
   LA MH MM MO MT MX MZ NI NP PA PE PH PK PR PT PY SA SG SV TH TT TW UM US VE VI
   WS YE ZA ZW`.split(/\s+/),
)
const SATURDAY_FIRST = new Set('AF BH DJ DZ EG IQ IR JO KW LY OM QA SD SY'.split(' '))
const FRIDAY_FIRST = new Set(['MV'])

// PrimeVue counts Sunday..Saturday as 0..6.
const SUNDAY = 0
const MONDAY = 1
const FRIDAY = 5
const SATURDAY = 6

/**
 * PrimeVue `dateFormat` for a locale, e.g. 'mm/dd/yy' for en-US,
 * 'dd. mm. yy' for sl. Falls back to PrimeVue's own default if Intl cannot
 * resolve the locale.
 */
export function localeDateFormat(locale: string): string {
  let parts: Intl.DateTimeFormatPart[]
  try {
    parts = new Intl.DateTimeFormat(locale, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      timeZone: 'UTC',
    }).formatToParts(SAMPLE)
  } catch {
    return 'mm/dd/yy'
  }

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

/**
 * PrimeVue `firstDayOfWeek` (0 = Sunday) for a locale. The language alone
 * decides it: CLDR keys week data by region, so a bare `sl` is widened to its
 * likely region (`sl-Latn-SI`) first — the same widening that makes `en` mean
 * en-US, which is what PrimeVue already assumed.
 */
export function localeFirstDayOfWeek(locale: string): number {
  let region: string | undefined
  try {
    region = new Intl.Locale(locale).maximize().region
  } catch {
    // Unresolvable, so match what localeDateFormat falls back to: en-US.
    return SUNDAY
  }
  if (!region) return SUNDAY
  if (SUNDAY_FIRST.has(region)) return SUNDAY
  if (SATURDAY_FIRST.has(region)) return SATURDAY
  if (FRIDAY_FIRST.has(region)) return FRIDAY
  return MONDAY
}

/**
 * The partial PrimeVue `locale` config for a language. PrimeVue deep-merges it
 * over its defaults, so every key left out here keeps its built-in value.
 */
export function primeVueLocale(locale: string): Record<string, unknown> {
  return {
    dateFormat: localeDateFormat(locale),
    firstDayOfWeek: localeFirstDayOfWeek(locale),
  }
}
