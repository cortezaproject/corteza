/**
 * Version utility functions
 */

/**
 * Parses a version string into major and minor components
 * @param {string} version - Version string (e.g., "2024.1.0")
 * @returns {Object} Object with year and month properties
 */
export function parseVersion (version) {
  if (!version) return { year: '', month: '' }

  const [year, month] = version.split('.')
  return { year, month }
}

/**
 * Generates a documentation URL for the given path using the current version
 * @param {string} path - Documentation path (e.g., "integrator-guide/automation/workflows/index.html")
 * @returns {string} Full documentation URL
 */
export function getDocumentationURL (path) {
  // eslint-disable-next-line no-undef
  const { year, month } = parseVersion(VERSION)
  return `https://docs.cortezaproject.org/corteza-docs/${year}.${month}/${path}`
}
