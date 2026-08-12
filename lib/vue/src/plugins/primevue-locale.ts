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
// Week start is deliberately left at PrimeVue's default. Intl's `getWeekInfo`
// is the only way to derive it and it is missing from Firefox and from Node, so
// wiring it up would give the same user a different calendar per browser and
// leave the code untestable here. It wants a locale table and a decision of its
// own.

// A date whose parts are all distinguishable — day 22, month 11 — so a locale
// that happens to put them in the same position cannot be misread.
const SAMPLE = new Date(Date.UTC(2026, 10, 22))

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
 * The partial PrimeVue `locale` config for a language. PrimeVue deep-merges it
 * over its defaults, so every key left out here keeps its built-in value.
 */
export function primeVueLocale(locale: string): Record<string, unknown> {
  return { dateFormat: localeDateFormat(locale) }
}
