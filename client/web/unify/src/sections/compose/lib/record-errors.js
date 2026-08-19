// What a record save's issues are, in terms a form can put on screen.
//
// The record endpoints answer a failed save with one detail per issue, each
// naming its field in `meta.field`; the top-level message is only a count
// ("2 issue(s) found") and names nothing. Showing that count is the same as
// showing nothing, so the details are what a caller places.

// Messages arrive with their placeholders unfilled — the server sends
// `The value "{{value}}" already exists in another record` and puts the value
// in the detail's own meta. Filling them here is what makes the message read.
const PLACEHOLDER = /{{\s*(\w+)\s*}}/g

/**
 * The Error behind a failure the API reported in the body of a 200, or null if
 * the body reported none.
 *
 * A generated client method runs this for its caller. An upload does not: it
 * posts its form data through the raw axios instance, so it never reaches
 * `stdResolve` and has to read the body itself — and `error` is an object, so
 * handing it straight to `new Error` yields the literal '[object Object]'
 * where the server's own sentence should be.
 */
export function apiError(data, fallback = 'unknown error') {
  const err = data?.error
  if (!err) return null
  if (typeof err === 'string') return new Error(err)

  const out = new Error(err.message || fallback)
  if (Array.isArray(err.details) && err.details.length > 0) out.details = err.details
  return out
}

/** One issue's message, with its placeholders filled from its meta. */
export function detailMessage(detail) {
  const message = detail?.message
  if (!message) return ''

  return message.replace(PLACEHOLDER, (whole, key) => {
    const fill = detail.meta?.[key]
    return fill === undefined || fill === null || fill === '' ? whole : String(fill)
  })
}

/**
 * Splits a failed save into what belongs beside a field and what does not.
 *
 * `canShow` answers whether a field is on screen right now: a layout need not
 * place every field of its module, and an error written against one it omits
 * is an error nobody ever sees. Those are named by label in `general` instead,
 * alongside the issues that belong to no field at all.
 */
export function partitionSaveErrors(err, { canShow = () => true, labelOf = name => name } = {}) {
  const fieldErrors = {}
  // One issue per offending record, so a value taken twice is reported twice in
  // the same words. Repeating a sentence tells the reader nothing.
  const general = new Set()

  for (const detail of err?.details ?? []) {
    const message = detailMessage(detail)
    if (!message) continue

    const field = detail.meta?.field
    if (!field) {
      general.add(message)
    } else if (canShow(field)) {
      const seen = fieldErrors[field]
      if (!seen) fieldErrors[field] = message
      else if (!seen.split('\n').includes(message)) fieldErrors[field] = `${seen}\n${message}`
    } else {
      general.add(`${labelOf(field)}: ${message}`)
    }
  }

  return { fieldErrors, general: [...general] }
}

/**
 * Issues carried by a save that SUCCEEDED — a duplicate-detection rule that is
 * not strict warns rather than refuses, and the payload reports it in
 * `valueErrors` instead of as an error.
 *
 * The record is saved and the form has moved on by the time these are read, so
 * there is no field left to sit beside: each is named by its field's label.
 */
export function saveWarnings(saved, { labelOf = name => name } = {}) {
  const out = new Set()

  for (const detail of saved?.valueErrors?.set ?? []) {
    const message = detailMessage(detail)
    if (!message) continue

    const field = detail.meta?.field
    out.add(field ? `${labelOf(field)}: ${message}` : message)
  }

  return [...out]
}

/** A field's label, for a message that has to name the field in its text. */
export function fieldLabeller(mod) {
  return name => mod?.fields?.find(f => f.name === name)?.label || name
}

/**
 * Which fields a record form has on screen, collected from the blocks drawing
 * them — a layout places blocks, a block places a selection of fields, and
 * field conditions narrow that again, so no single place knows the answer.
 *
 * Blocks register a getter rather than a value: the set is read at the moment a
 * save fails, by which time a condition may have moved a field out of reach.
 */
export function displayedFieldRegistry() {
  const reporters = new Set()

  return {
    register(fn) {
      reporters.add(fn)
      return () => reporters.delete(fn)
    },

    names() {
      const names = new Set()
      for (const report of reporters) {
        try {
          report().forEach(name => names.add(name))
        } catch (e) {
          console.error('A block failed to report its fields:', e)
        }
      }
      return names
    },
  }
}
