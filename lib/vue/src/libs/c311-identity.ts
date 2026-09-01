export const passwordPolicy = Object.freeze({
  minLength: 12,
  maxLength: 128,
  minimumCharacterClasses: 3,
})

export function validatePassword (value: string): string[] {
  const errors: string[] = []
  if (value.length < passwordPolicy.minLength) errors.push('too-short')
  if (value.length > passwordPolicy.maxLength) errors.push('too-long')
  const classes = [/[A-Z]/, /[a-z]/, /[0-9]/, /[^A-Za-z0-9]/].filter(pattern => pattern.test(value)).length
  if (classes < passwordPolicy.minimumCharacterClasses) errors.push('character-classes')
  return errors
}

export function resetTokenFromLocation (location: Location): string {
  const params = new URLSearchParams(location.search || '')
  const token = params.get('token') || ''
  const browserHistory = typeof window !== 'undefined' && window.history
    ? window.history
    : typeof history !== 'undefined' ? history : undefined
  if (token && browserHistory && typeof browserHistory.replaceState === 'function') {
    params.delete('token')
    const query = params.toString()
    browserHistory.replaceState(null, '', `${location.pathname}${query ? `?${query}` : ''}${location.hash || ''}`)
  }
  return token
}
