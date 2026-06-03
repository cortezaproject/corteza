/**
 * Constraint utility functions
 */

/**
 * Converts a constraint name to a human-readable label
 * @param {string} name - The constraint name (e.g., "user.email")
 * @returns {string} The formatted label (e.g., "User Email")
 */
export function getConstraintNameLabel (name) {
  if (!name) return ''

  return name
    .split('.')
    .map(s => {
      const first = s[0] || ''
      const rest = s.slice(1) || ''
      return first.toUpperCase() + rest.toLowerCase()
    })
    .join(' ')
}
