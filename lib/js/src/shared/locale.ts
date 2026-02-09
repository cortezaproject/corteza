/**
 * Locale-related helper functions
 */

/**
 * Get the first day of the week for a given locale.
 * Returns 0 for Sunday, 1 for Monday, etc.
 *
 * @param locale - BCP 47 language tag (e.g., 'en-US', 'hu-HU')
 * @returns The first day of the week (0-6)
 */
export function getWeekStartDay (locale: string): number {
  const l = new Intl.Locale(locale)

  // @ts-ignore - getWeekInfo is a newer API
  if (l.getWeekInfo) {
    // @ts-ignore
    return (l.getWeekInfo().firstDay % 7)
  }

  // @ts-ignore - weekInfo is deprecated but still supported in some browsers
  if (l.weekInfo && l.weekInfo.firstDay !== undefined) {
    // @ts-ignore
    return (l.weekInfo.firstDay % 7)
  }

  return 0
}
