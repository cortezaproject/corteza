/**
 * URL utility functions for the Url field viewer.
 */

/**
 * Processes a raw URL string and returns a URL object.
 * Prepends https:// if no protocol is present.
 */
function makeURL(url: string): URL {
  if (!/^https?:\/\//i.test(url)) {
    url = 'https://' + url
  }
  return new URL(url)
}

/**
 * Removes URL's fragment (#hash)
 */
export function trimUrlFragment(url: string): string {
  if (!url) return url
  const u = makeURL(url)
  u.hash = ''
  return u.toString()
}

/**
 * Removes URL's query string (?params)
 */
export function trimUrlQuery(url: string): string {
  if (!url) return url
  const u = makeURL(url)
  u.search = ''
  return u.toString()
}

/**
 * Removes URL's path
 */
export function trimUrlPath(url: string): string {
  if (!url) return url
  const u = makeURL(url)
  u.pathname = ''
  return u.toString()
}

/**
 * Forces the URL to use https protocol
 */
export function onlySecureUrl(url: string): string {
  if (!url) return url
  const u = makeURL(url)
  u.protocol = 'https'
  return u.toString()
}
