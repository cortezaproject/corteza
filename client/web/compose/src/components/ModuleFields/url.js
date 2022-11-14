/**
 * Processes a given URL and returns a URL object
 * @param {String} url URL to process
 * @returns {URL}
 */
function makeURL (url) {
  if (!(/^http(s)?:\/\//.test(url))) {
    url = 'https://' + url
  }
  return new URL(url)
}

/**
 * Removes URL's fragment
 * @param {String} url URL to process
 * @returns {String}
 */
export function trimUrlFragment (url) {
  if (!url) {
    return url
  }

  url = makeURL(url)
  url.hash = ''
  return url.toString()
}

/**
 * Removes URL's query
 * @param {String} url URL to process
 * @returns {String}
 */
export function trimUrlQuery (url) {
  if (!url) {
    return url
  }

  url = makeURL(url)
  url.search = ''
  return url.toString()
}

/**
 * Removes URL's path
 * @param {String} url URL to process
 * @returns {String}
 */
export function trimUrlPath (url) {
  if (!url) {
    return url
  }

  url = makeURL(url)
  url.pathname = ''
  return url.toString()
}

/**
 * Assures that the given URL will provide a https schema
 * @param {String} url URL to process
 * @returns {String}
 */
export function onlySecureUrl (url) {
  if (!url) {
    return url
  }

  url = makeURL(url)
  url.protocol = 'https'
  return url.toString()
}
