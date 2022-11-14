/* eslint-disable @typescript-eslint/ban-ts-ignore */
import qs from 'qs'

const baseQsConfig = {
  arrayFormat: 'brackets',
  encode: false,
}

/**
 * Little URL helper
 *
 * We need it to handle relative URLs, especially ones w/o schema
 */
export function Make ({ url = '', query = {}, hash = '', ref = window.location.toString(), config = {} }): string {
  let u

  if (/^http(s)?:\/\//.test(url)) {
    u = new URL(url)
  } else if (/^\/\//.test(url)) {
    u = new URL(ref)
    u.href = `${u.protocol}${url}`
  } else {
    // Construct full relative URL from path
    u = new URL(ref)
    u.pathname = url
  }

  if (hash) {
    u.hash = hash
  }

  // TypeScript somehow thinks that 'brackets' is not a string.
  // @ts-ignore
  u.search = qs.stringify(query, {
    ...baseQsConfig,
    ...config,
  })

  return u.toString()
}
