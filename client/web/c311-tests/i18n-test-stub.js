export const C311_LOCALES = ['en', 'es', 'vi']
export const C311_MESSAGES = {
  en: { 'status.loading': 'Loading' },
  es: { 'status.loading': 'Cargando' },
  vi: { 'status.loading': 'Đang tải' },
}

export function persistC311Locale (locale, actorID) {
  localStorage.setItem(`c311.locale.${actorID || 'anonymous'}`, locale)
}

export function readC311Locale (actorID) {
  const value = localStorage.getItem(`c311.locale.${actorID || 'anonymous'}`)
  return C311_LOCALES.includes(value) ? value : null
}
