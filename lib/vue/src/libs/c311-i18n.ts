export const C311_LOCALES = ['en', 'es', 'vi'] as const
export type C311Locale = typeof C311_LOCALES[number]

export const C311_MESSAGES: Record<C311Locale, Record<string, string>> = {
  en: {
    'accessibility.skipToMain': 'Skip to main content',
    'dirty.leave': 'You have unsaved changes. Leave this page?',
    'status.loading': 'Loading',
    'status.loading.message': 'Loading data.',
    'status.empty': 'No matching records',
    'status.empty.message': 'Zero matching records were found.',
    'status.permission-denied': 'Access denied',
    'status.forbidden': 'Access denied',
    'status.forbidden.message': 'You do not have access to this information.',
    'status.not-found': 'Not found',
    'status.not-found.message': 'The requested information could not be found.',
    'status.validation-error': 'Check your information',
    'status.validation-error.message': 'Correct the highlighted fields and try again.',
    'status.retryable-error': 'Temporarily unavailable',
    'status.retryable-error.message': 'The service is temporarily unavailable.',
    'status.terminal-error': 'Unable to complete this operation',
    'status.terminal-error.message': 'The operation could not be completed.',
    'status.versionConflict': 'Current server version: {{version}}. Reload before trying again.',
    'help.public.request.submit': 'Submit a service request.',
    'help.public.request.lookup': 'Look up the status of a service request.',
    'help.staff.request.triage': 'Review and classify a request.',
    'help.staff.request.reassign': 'Reassign a request with a reason.',
    'help.staff.request.bulk-update': 'Apply an update to selected requests.',
    'help.admin.workflow.author': 'Create or update an approved workflow.',
    'help.staff.report.create': 'Create a report from permitted records.',
    'help.admin.branding.publish': 'Preview and publish approved branding.',
  },
  es: {
    'accessibility.skipToMain': 'Saltar al contenido principal',
    'dirty.leave': 'Tiene cambios sin guardar. ¿Salir de esta página?',
  },
  vi: {
    'accessibility.skipToMain': 'Chuyển đến nội dung chính',
    'dirty.leave': 'Bạn có thay đổi chưa lưu. Rời khỏi trang này?',
  },
}

export function installC311Translations (i18next: any): void {
  if (!i18next?.addResourceBundle) return
  C311_LOCALES.forEach(locale => i18next.addResourceBundle(locale, 'c311', C311_MESSAGES[locale], true, true))
  if (i18next.options) {
    i18next.options.fallbackLng = 'en'
    i18next.options.supportedLngs = C311_LOCALES
  }
}

export function persistC311Locale (locale: C311Locale, actorID?: string): void {
  if (typeof localStorage === 'undefined') return
  localStorage.setItem(`c311.locale.${actorID || 'anonymous'}`, locale)
}

export function readC311Locale (actorID?: string): C311Locale | null {
  if (typeof localStorage === 'undefined') return null
  const value = localStorage.getItem(`c311.locale.${actorID || 'anonymous'}`)
  return C311_LOCALES.includes(value as C311Locale) ? value as C311Locale : null
}
