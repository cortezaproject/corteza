const re = /^[A-Za-z][0-9A-Za-z_\-.]*[A-Za-z0-9]$/

export const isValid = (h: string) => re.test(h)

// Used for state
export function handleState (h: string) {
  if (h) {
    return isValid(h) ? null : false
  }
}
