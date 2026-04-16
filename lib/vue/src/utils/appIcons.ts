import defaultAppIcon from '../assets/default-app.png'
import adminAreaIcon from '../assets/admin-area.png'
import agenticIcon from '../assets/agentic.png'
import namespacesIcon from '../assets/namespaces.png'
import projectsIcon from '../assets/projects.png'
import taqIcon from '../assets/taq.png'
import workflowsIcon from '../assets/workflows.png'

export const appIconMap: Record<string, string> = {
  'applications/default-app.png': defaultAppIcon,
  'applications/admin-area.png': adminAreaIcon,
  'applications/agentic.png': agenticIcon,
  'applications/namespaces.png': namespacesIcon,
  'applications/projects.png': projectsIcon,
  'applications/taq.png': taqIcon,
  'applications/workflows.png': workflowsIcon,
}

export { defaultAppIcon }

/**
 * Resolves the logo/icon URL for an application.
 * Handles bundled icons, API paths, and absolute URLs.
 */
export function resolveAppLogoUrl(app: any, apiBaseUrl: string): string {
  const logo = app?.unify?.logo || app?.unify?.icon || ''
  if (!logo) return defaultAppIcon

  // Check if it's a known bundled app icon
  if (appIconMap[logo]) return appIconMap[logo]

  // Handle API paths like /api/system/...
  const apiSystem = '/api/system'
  const base = new URL(apiBaseUrl, window.location.origin).toString()
  if (logo.startsWith(apiSystem)) {
    return base.substring(0, base.length - apiSystem.length) + logo
  }

  if (logo.startsWith('/') || logo.startsWith('http')) return logo
  return base.substring(0, base.length - apiSystem.length) + '/' + logo
}
