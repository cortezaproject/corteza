/**
 * The href for a registry application's `unify.url`.
 *
 * The registry stores that url exactly as it was typed, so it arrives in three
 * shapes: an address (`https://www.google.com`), a path this shell serves
 * (`admin/system/labels`), or a bare host (`www.google.com`). A scheme settles
 * the first and a leading slash the second; what is left is decided by a dot in
 * the first segment — a shell path never carries one, a host always does.
 *
 * A bare host has to be given a scheme: left alone the browser reads it as a
 * relative path and the link lands on the shell's catch-all instead of the
 * site it names.
 */
export function resolveAppUrl(url: string | null | undefined): string {
  const raw = String(url ?? '').trim()
  if (!raw) return ''

  // A path this shell serves, or a protocol-relative `//host`.
  if (raw.startsWith('/')) return raw

  // Already carries a scheme.
  if (/^[a-z][a-z0-9+.-]*:/i.test(raw)) return raw

  return raw.split('/')[0].includes('.') ? 'https://' + raw : '/' + raw
}

/**
 * Whether the given `unify.url` names a path this shell serves, as opposed to
 * an address somewhere else. `//host` is protocol-relative, so a leading slash
 * alone does not settle it.
 */
export function isAppUrlLocal(url: string | null | undefined): boolean {
  const href = resolveAppUrl(url)
  return href.startsWith('/') && !href.startsWith('//')
}
