import defaultAppIcon from '../assets/default-app.png'
import adminAreaIcon from '../assets/admin-area.png'
import discoveryIcon from '../assets/discovery.png'
import lowCodeCrmIcon from '../assets/low-code-crm-app.png'
import lowCodePlatformIcon from '../assets/low-code-platform.png'
import lowCodeServiceIcon from '../assets/low-code-service-solution-app.png'
import privacyIcon from '../assets/privacy.png'
import reporterIcon from '../assets/reporter.png'
import videoConferenceIcon from '../assets/video-conference.png'
import workflowsIcon from '../assets/workflows.png'

export const appIconMap: Record<string, string> = {
  'applications/default-app.png': defaultAppIcon,
  'applications/admin-area.png': adminAreaIcon,
  'applications/discovery.png': discoveryIcon,
  'applications/low-code-crm-app.png': lowCodeCrmIcon,
  'applications/low-code-platform.png': lowCodePlatformIcon,
  'applications/low-code-service-solution-app.png': lowCodeServiceIcon,
  'applications/privacy.png': privacyIcon,
  'applications/reporter.png': reporterIcon,
  'applications/video-conference.png': videoConferenceIcon,
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
