/**
 * String utility functions
 */

/**
 * Converts a camelCase string to Title Case with spaces
 * @param {string} str - The camelCase string to convert
 * @returns {string} The converted title case string
 */
export function camelToTitle (str) {
  if (!str) return ''

  return str
    .split(/(?=[A-Z])/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}
