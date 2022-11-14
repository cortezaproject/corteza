const re = /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$/

export const isValid = (h) => re.test(h)

// Used for state
export function handleState (h) {
  if (!h || h.length === 0) {
    return null
  }

  return isValid(h) ? null : false
}
