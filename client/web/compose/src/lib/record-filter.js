/**
 * Evaluates prefilter string by replacing template variables.
 *
 * Supported variables:
 * - ${recordID} — current record ID
 * - ${ownerID} — record owner ID
 * - ${userID} — current user ID
 * - ${record.values.fieldName} — record field value
 * - ${user.name} — user name
 * - ${user.email} — user email
 *
 * @param {string} prefilter - filter string with ${...} variables
 * @param {Object} context
 * @param {Object} [context.record] - current record
 * @param {Object} [context.user] - current user
 * @param {string} [context.recordID] - record ID
 * @param {string} [context.ownerID] - record owner ID
 * @param {string} [context.userID] - user ID
 * @returns {string} - evaluated filter string
 */
export function evaluatePrefilter(prefilter, { record, user, recordID, ownerID, userID } = {}) {
  if (!prefilter) return prefilter

  return prefilter.replace(/\$\{([^}]+)\}/g, (match, expr) => {
    const trimmed = expr.trim()

    if (trimmed === 'recordID') return recordID || '0'
    if (trimmed === 'ownerID') return ownerID || '0'
    if (trimmed === 'userID') return userID || '0'

    // ${record.values.fieldName}
    if (trimmed.startsWith('record.values.')) {
      const fieldName = trimmed.slice('record.values.'.length)
      if (record && record.values) {
        const val = record.values[fieldName]
        return val !== undefined && val !== null ? String(val) : ''
      }
      return ''
    }

    // ${record.recordID}, ${record.ownedBy}, etc.
    if (trimmed.startsWith('record.')) {
      const prop = trimmed.slice('record.'.length)
      if (record && record[prop] !== undefined) {
        return String(record[prop])
      }
      return ''
    }

    // ${user.name}, ${user.email}, etc.
    if (trimmed.startsWith('user.')) {
      const prop = trimmed.slice('user.'.length)
      if (user && user[prop] !== undefined) {
        return String(user[prop])
      }
      return ''
    }

    // Unknown variable — leave as-is
    return match
  })
}

/**
 * Checks if a field name is used in a filter string.
 *
 * @param {string} fieldName
 * @param {string} filter
 * @returns {boolean}
 */
export function isFieldInFilter(fieldName, filter = '') {
  if (!fieldName || !filter) return false
  return filter.includes(fieldName)
}
