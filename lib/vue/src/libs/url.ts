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

interface AppURLParams {
  // Web application handle (admin, compose, workflow, ...)
  app: string;
  // Path inside that application, without a leading slash
  path?: string;
  // Base URL of all web applications; defaults to window.CortezaWebapp,
  // which the server sets when it serves the applications
  webapp?: string;
  // Base of the current application (its <base href>), used when the webapp
  // base is not configured; defaults to document.baseURI
  base?: string;
}

/**
 * Builds a link to another web application
 *
 * Applications can be served under a path prefix (HTTP_BASE_URL), so a
 * root-absolute path like /workflow/... would miss it. The link is derived
 * from the configured webapp base or, failing that, from the current
 * application's <base href> which the server rewrites to <prefix>/<app>/.
 */
export function MakeAppURL ({ app, path = '', webapp, base }: AppURLParams): string {
  if (webapp === undefined) {
    // @ts-ignore
    webapp = typeof window !== 'undefined' ? window.CortezaWebapp : undefined
  }

  if (base === undefined) {
    base = typeof document !== 'undefined' ? document.baseURI : undefined
  }

  const rel = `${app}/${path}`.replace(/\/+$/, '')

  if (webapp) {
    const root = webapp.replace(/\/+$/, '')
    return /^(https?:)?\/\//.test(root)
      ? new URL(`${root}/${rel}`).toString()
      : Make({ url: `${root}/${rel}`, ref: base })
  }

  // the current app lives one level below the webapp base
  return new URL(`../${rel}`, base).toString()
}
