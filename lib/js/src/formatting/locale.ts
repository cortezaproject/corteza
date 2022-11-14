/**
 * Returns current language
 *
 * This temporary solution returns an empty array;
 * this will cause Intl functions to format strings and numbers in the current (by-browser) language
 */
export function currentLanguage (): string|string[] {
  return []
}
