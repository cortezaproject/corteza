/**
 * Resolves the logo/icon URL for an application.
 * Handles bundled icons, API paths, and absolute URLs.
 * iconMap is provided by the host app (e.g. unify) so the lib stays asset-free.
 */
export function resolveAppLogoUrl(
  app: any,
  apiBaseUrl: string,
  iconMap: Record<string, string> = {},
): string {
  const logo = app?.unify?.logo || app?.unify?.icon || ''

  if (!logo) return iconMap['applications/default-app.png'] || ''

  if (iconMap[logo]) return iconMap[logo]

  const apiSystem = '/api/system'
  const base = new URL(apiBaseUrl, window.location.origin).toString()
  if (logo.startsWith(apiSystem)) {
    return base.substring(0, base.length - apiSystem.length) + logo
  }

  if (logo.startsWith('/') || logo.startsWith('http')) return logo
  return base.substring(0, base.length - apiSystem.length) + '/' + logo
}
