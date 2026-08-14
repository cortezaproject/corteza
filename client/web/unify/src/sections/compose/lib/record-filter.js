import { compose } from '@planetcrust/human-js'
import moment from 'moment'

export const nonQueryableFieldNames = ['recordID']
export const nonQueryableFieldKinds = [
  'Number',
  'Record',
  'User',
  'Bool',
  'DateTime',
  'File',
  'Geometry',
]

/**
 * Escapes special characters for Human QL string literals
 *
 * Human QL uses backslash-based escaping within string literals (delimited by single quotes).
 * The QL lexer recognizes only two escape sequences:
 * - \' for single quote
 * - \\ for backslash
 *
 * For LIKE patterns, additional characters need escaping:
 * - Backslashes need double escaping (once for QL, once for SQL LIKE)
 * - % and _ are SQL LIKE wildcards and must be escaped to match literally
 */
export function escapeQlString(str, isLikePattern = false) {
  let result = String(str)
    // Escape backslashes first (must be done before quotes)
    .replace(/\\/g, isLikePattern ? '\\\\\\\\' : '\\\\')
    // Escape single quotes with backslash (QL syntax, not SQL's '')
    .replace(/'/g, "\\'")

  if (isLikePattern) {
    result = result
      // Escape % → \\% (the \\ becomes \ in QL, then \% in SQL means literal %)
      .replace(/%/g, '\\\\%')
      // Escape _ → \\_ (the \\ becomes \ in QL, then \_ in SQL means literal _)
      .replace(/_/g, '\\\\_')
  }

  return result
}

// How a date value arrives from the filter UI.
//
// CInputDateTime (lib/vue) emits one of three shapes, depending on the field's
// options: Date.toISOString() for a datetime field ('2024-06-29T12:30:00.000Z'),
// YYYY-MM-DD when it is date-only, and HH:mm:ss when it is time-only. A shape
// no branch matches makes getFieldFilter return undefined and
// getRecordListFilterSql drop the condition silently.
//
// Date-only is tested BEFORE datetime: moment's ISO_8601 accepts a bare
// YYYY-MM-DD too, and letting it through the timestamp branch would resolve the
// value to midnight and lose the day-level comparison that keeps an end-of-day
// bound inclusive.
const asDateOnly = value => moment(value, 'YYYY-MM-DD', true)
const asDateTime = value => moment(value, moment.ISO_8601, true)
const asTimeOnly = value => moment(value, ['HH:mm:ss', 'HH:mm'], true)

// Generate record list sql query string based on filter object input

export function getRecordListFilterSql(filter) {
  let query = ''

  let existsPreviousElement = false

  filter.forEach(f => {
    if (f.name && f.operator) {
      const fieldFilter = getFieldFilter(f.name, f.kind, f.value, f.operator)

      if (fieldFilter) {
        if (existsPreviousElement) {
          query += ` ${f.condition} `
        }

        query += getFieldFilter(f.name, f.kind, f.value, f.operator)
        existsPreviousElement = true
      }
    }
  })

  // Clean up query by grouping OR conditions and trimming whitespace
  // Split query into AND parts, wrap any OR conditions in parentheses,
  // trim whitespace from each part, then rejoin with AND
  query = query
    .split(' AND ')
    .map(q => {
      if (q.includes('OR')) {
        return `(${q.trim()})`
      }

      return q.trim()
    })
    .join(' AND ')

  return query ? `(${query})` : query
}

// Operators that test for the presence of a value and so carry none of their
// own. The UI hides the value editor for these, and getFieldFilter builds the
// condition from the field kind alone.
export const IS_EMPTY = 'IS EMPTY'
export const IS_NOT_EMPTY = 'IS NOT EMPTY'

export function isValuelessOperator(op) {
  return [IS_EMPTY, IS_NOT_EMPTY].includes(op)
}

// Kinds whose value lands in a text column, where "empty" has two spellings: a
// record that never held a value has no row at all (IS NULL), while one whose
// value was cleared keeps its row with an empty string. Only testing IS NULL
// finds the first and misses the second, which is the more common of the two —
// clearing a field is something users do, never setting one happens once.
//
// Every other kind is stored in a typed column (numeric, boolean, timestamptz,
// bigint) where comparing against '' is a hard postgres error, e.g.
// `invalid input syntax for type numeric: ""`. Those get the NULL test alone.
//
// An allowlist rather than a denylist on purpose: a kind not named here falls
// to the form that cannot error, so a new field kind fails safe.
const emptyStringKinds = ['String', 'Email', 'Url', 'Select']

// Helper function that creates a query for a specific field kind
export function getFieldFilter(name, kind, query = '', operator = '=') {
  const numQuery = Number.parseFloat(query)

  // Before every other branch: these ignore the value entirely, and the Bool
  // branch below would otherwise read their absent value as `false`.
  if (isValuelessOperator(operator)) {
    const wantEmpty = operator === IS_EMPTY

    if (emptyStringKinds.includes(kind)) {
      // `!=` alone is the whole not-empty test: SQL makes `NULL != ''` unknown,
      // so a missing row is excluded without a second condition. Spelling it
      // out as `(x IS NOT NULL) AND (x != '')` returns the same rows but puts
      // an ' AND ' inside one field's condition, which getRecordListFilterSql
      // splits on when it groups OR conditions.
      return wantEmpty ? `((${name} IS NULL) OR (${name} = ''))` : `(${name} != '')`
    }

    return wantEmpty ? `(${name} IS NULL)` : `(${name} IS NOT NULL)`
  }

  const build = (op, left, right) => {
    switch (op.toUpperCase()) {
      case '!=':
      case 'NOT LIKE':
        return `((${left} ${op} ${right}) OR (${left} IS NULL))`
      case 'IN':
      case 'NOT IN':
        // flip left/right for IN/NOT IN
        return `(${right} ${op} ${left})`
      default:
        return `(${left} ${op} ${right})`
    }
  }

  // Boolean should search for literal values. Example `${name} = true` or just `${name}
  // At the moment it doesn't seem to be working as intended

  if (kind === 'Bool') {
    const operation = operator === '=' ? 'is' : 'is not'
    const boolQuery = toBoolean(query)

    if (boolQuery) {
      return `(${name} ${operation} true)`
    } else {
      return `((${name} ${operation} false) OR (${name} IS NULL))`
    }
  }

  // Take care of special case where query is undefined and its not a Bool field
  if (!query && query !== 0) {
    if (['=', 'IN'].includes(operator)) {
      return `(${name} IS NULL)`
    } else if (['!=', 'NOT IN'].includes(operator)) {
      return `(${name} IS NOT NULL)`
    }

    return undefined
  }

  if (['Number'].includes(kind)) {
    if (['BETWEEN', 'NOT BETWEEN'].includes(operator)) {
      return build(operator, name, `${query.start} ${query.end}`)
    } else if (!isNaN(numQuery)) {
      return build(operator, name, `'${numQuery}'`)
    }
  }

  if (['DateTime'].includes(kind)) {
    // Callers hold dates as Date objects (compose.Record.createdAt is one). A
    // Date stringifies as 'Fri Aug 14 2026 13:19:57 GMT+0200 (Central European
    // Summer Time)', which postgres rejects outright, and moment reports a Date
    // valid against any format — so the strict date-only parse below claims it
    // and drops the time. ISO first, and both branches read a string.
    const asIso = v => (v instanceof Date || moment.isMoment(v) ? moment(v).toISOString() : v)

    if (['BETWEEN', 'NOT BETWEEN'].includes(operator)) {
      query = { ...query, start: asIso(query.start), end: asIso(query.end) }
    } else {
      query = asIso(query)
    }

    const dataFmtEntry = date =>
      `TIMESTAMP(DATE_FORMAT('${date.format()}', '%Y-%m-%dT%H:%i:%s.%f+00:00'))`

    if (['BETWEEN', 'NOT BETWEEN'].includes(operator)) {
      const startDate = asDateOnly(query.start)
      const endDate = asDateOnly(query.end)

      if (startDate.isValid() && endDate.isValid()) {
        // Compare on the field's date part so a date-only range is inclusive of the
        // whole end day; otherwise DATE('end') resolves to midnight and records
        // logged during the end day itself are excluded.
        return build(operator, `DATE(${name})`, `DATE('${query.start}') DATE('${query.end}')`)
      }

      const startDateTime = asDateTime(query.start)
      const endDateTime = asDateTime(query.end)

      if (startDateTime.isValid() && endDateTime.isValid()) {
        return build(
          operator,
          `TIMESTAMP(DATE_FORMAT(${name}, '%Y-%m-%dT%H:%i:%s.%f+00:00'))`,
          `${dataFmtEntry(startDateTime)} ${dataFmtEntry(endDateTime)}`,
        )
      }

      const startTime = asTimeOnly(query.start)
      const endTime = asTimeOnly(query.end)

      if (startTime.isValid() && endTime.isValid()) {
        // Compare on the field's time part, as the date branch does with DATE():
        // the value is stored as a timestamp, and postgres refuses
        // "timestamp >= time" outright.
        return build(operator, `TIME(${name})`, `TIME('${query.start}') TIME('${query.end}')`)
      }

      // Special case where between object is invalid
      return undefined
    } else {
      // Build different querries if date, time or datetime
      const date = asDateOnly(query)
      const dateTime = asDateTime(query)
      const time = asTimeOnly(query)

      // @note tweaking the template a bit:
      // * adding %f to include fractions; mysql sometimes forces them when formatting date
      // * changing Z to +00:00
      // * doing the same for time-only fields
      if (date.isValid()) {
        // Compare on the field's date part so date-only operators (e.g. <=) include
        // the whole day instead of resolving the value to midnight.
        return build(operator, `DATE(${name})`, `DATE('${query}')`)
      } else if (dateTime.isValid()) {
        return build(
          operator,
          `TIMESTAMP(DATE_FORMAT(${name}, '%Y-%m-%dT%H:%i:%s.%f+00:00'))`,
          dataFmtEntry(dateTime),
        )
      } else if (time.isValid()) {
        return build(operator, `TIME(${name})`, `TIME('${query}')`)
      }
    }
  }

  // Since userID and recordID must be numbers, we check if query is number to avoid wrong queries
  if (['User', 'Record'].includes(kind) && !isNaN(numQuery)) {
    return build(operator, name, `'${escapeQlString(query)}'`)
  }

  if (['String', 'Url', 'Select', 'Email'].includes(kind)) {
    if (operator === 'LIKE' || operator === 'NOT LIKE') {
      const escapedQuery = escapeQlString(query, true)
        // Convert user's * to % wildcard AFTER escaping (so literal % is already \%)
        .replace(/\*/g, '%')

      return build(operator, name, `'%${escapedQuery}%'`)
    }

    return build(operator, name, `'${escapeQlString(query, false)}'`)
  }
}

// Helper to determine if and value for given bool query
// == is intentional
const toBoolean = v => {
  if (v == 'false' || v == 0) {
    return false
  }
  if (v == 'true' || v == 1) {
    return true
  }

  return undefined
}

export function convertRecordListFilter(groupFilter = []) {
  // Group filters by name
  const filtersByName = groupFilter
    .sort((a, b) => a.name.localeCompare(b.name))
    .reduce((acc, filter) => {
      if (!acc[filter.name]) {
        acc[filter.name] = []
      }

      acc[filter.name].push(filter)

      return acc
    }, {})

  // Rebuild filter array with proper conditions
  return Object.entries(filtersByName).flatMap(([_, filters], groupIndex) =>
    filters.map((filter, idx) => ({
      ...filter,
      condition: idx === 0 ? (groupIndex === 0 ? '' : 'AND') : 'OR',
    })),
  )
}

// Takes fields and prefilter and record list filter and merges them into query string
// ie: Return records that have strings in columns (fields) we're showing that start with <query> in case
//     of text or are exactly the same in case of numbers
export function queryToFilter(
  searchQuery = '',
  prefilter = '',
  fields = [],
  recordListFilter = [],
) {
  searchQuery = (searchQuery || '').trim()

  // Create query for search string
  if (searchQuery || searchQuery === 0) {
    searchQuery = fields
      .filter(
        f =>
          f.isQueryable !== false &&
          !nonQueryableFieldNames.includes(f.name) &&
          !nonQueryableFieldKinds.includes(f.kind),
      )
      .map(f => getFieldFilter(f.name, f.kind, searchQuery, 'LIKE'))
      .filter(q => !!q)
      .join(' OR ')

    searchQuery = searchQuery ? `(${searchQuery})` : ''
  }

  // Build query with proper parenthesization for mixed AND/OR
  // Groups connected by AND have higher precedence, so we wrap consecutive ANDs
  const groups = recordListFilter
    .map(({ filter = [], groupCondition }) => ({
      sql: getRecordListFilterSql(filter),
      condition: groupCondition || 'OR',
    }))
    .filter(({ sql }) => sql)

  // Group consecutive ANDs together, then join with ORs for proper precedence
  const orSegments = []
  let andGroup = []

  for (let i = 0; i < groups.length; i++) {
    const { sql, condition } = groups[i]

    if (i === 0 || condition === 'AND') {
      andGroup.push(sql)
    } else {
      // OR encountered - flush AND group as its own segment
      orSegments.push(andGroup.length > 1 ? `(${andGroup.join(' AND ')})` : andGroup[0])
      andGroup = [sql]
    }
  }
  // Flush remaining
  if (andGroup.length) {
    orSegments.push(andGroup.length > 1 ? `(${andGroup.join(' AND ')})` : andGroup[0])
  }

  let recordListFilterSql = orSegments.join(' OR ')

  // When there is a top-level OR (more than one segment), wrap the combined
  // filter so it is grouped before being AND-joined with the prefilter/search
  // query. Without this, SQL operator precedence (AND binds tighter than OR)
  // lets the last OR branch escape the prefilter:
  // `prefilter AND (a) OR (b)` => `(prefilter AND (a)) OR (b)`.
  if (orSegments.length > 1) {
    recordListFilterSql = `(${recordListFilterSql})`
  }

  return [prefilter, recordListFilterSql, searchQuery].filter(f => f).join(' AND ')
}

// Evaluates the given prefilter. Allows JS template literal expressions
// such as id = ${recordID}
export function evaluatePrefilter(prefilter, { record, user, recordID, ownerID, userID }) {
  return compose.interpolateTemplate(prefilter, { record, user, recordID, ownerID, userID })
}

// Reports whether an author-typed template reads the page record. Such a
// template can only be evaluated where a record exists — off a record page
// (page builder, list page) `${record.values.x}` throws on the missing record
// and `${ownerID}`/`${recordID}` have nothing to resolve to, so callers skip
// the feed/metric/list instead of evaluating it.
export function usesRecordVariables(template = '') {
  return !!template && (template.includes('${record') || template.includes('${ownerID}'))
}

// Removes char from end of string
export function trimChar(text = '', char = '') {
  if (text.substring(text.length - char.length, text.length) === char) {
    text = text.substring(0, text.length - char.length)
  }
  return text
}

// Helper function that checks if field name is included in a filter
export function isFieldInFilter(fieldName, filter = '') {
  if (!fieldName || !filter) return

  return filter.includes(fieldName)
}

export function formatActiveFilterOperator(op) {
  const operators = {
    '=': 'equal',
    '!=': 'notEqual',
    IN: 'in',
    'NOT IN': 'notIn',
    '>': 'greaterThan',
    '<': 'lessThan',
    LIKE: 'like',
    'NOT LIKE': 'notLike',
    BETWEEN: 'between',
    'NOT BETWEEN': 'notBetween',
    [IS_EMPTY]: 'isEmpty',
    [IS_NOT_EMPTY]: 'isNotEmpty',
  }

  return operators[op] || 'is'
}

export function isBetweenOperator(op) {
  return ['BETWEEN', 'NOT BETWEEN'].includes(op)
}

// Browser-storage keys for one record list's saved filter and saved presets.
//
// The pageID is part of the key because a blockID is only unique WITHIN its
// page: the server numbers a page's blocks from 1 (server/compose/service/
// page.go), so the first record list of every page is blockID 1. Keying on the
// blockID alone put all of those lists in one bucket — a filter set on one
// page's list came back on an unrelated page's list, against a different
// module, silently hiding its records. A pageID is globally unique, so
// page+block names exactly one record list.
export function recordListFilterStorageKey(pageID, blockID) {
  return `recordListFilter-${pageID || '0'}-${blockID}`
}

export function recordListPresetsStorageKey(pageID, blockID) {
  return `recordListFilterPresets-${pageID || '0'}-${blockID}`
}
