// Readable names for action-log entries, from the action log's own
// vocabulary so the dashboard and the log page agree on every label.
import { actionLabel, resourceID, resourceTypeLabel } from '../../views/system/ActionLog/vocab'

function capitalize(s) {
  return s ? s[0].toUpperCase() + s.slice(1) : s
}

// "reloadDALModels" → "reload DAL models": an action the vocabulary does not
// know is split on its camel humps, keeping its acronyms.
export function humanizeAction(action) {
  return String(action || '')
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .replace(/([A-Z]+)([A-Z][a-z])/g, '$1 $2')
    .replace(/[-_]+/g, ' ')
    .trim()
}

// "corteza::compose:namespace/*" + "lookup" → "Namespace lookup"
export function describeEvent(resource, action) {
  let noun = resourceTypeLabel(resource) || ''
  // a type the vocabulary does not know ("system:auth") reads as its last part
  if (noun.includes(':'))
    noun = capitalize(noun.slice(noun.lastIndexOf(':') + 1).replace(/[-_]+/g, ' '))
  const label = actionLabel(action) || ''
  const verb =
    label === action
      ? humanizeAction(action)
          .toLowerCase()
          .replace(/\b(dal|api|ai|id|llm|taq)\b/g, m => m.toUpperCase())
      : label.toLowerCase()
  return capitalize(`${noun} ${verb}`.trim())
}

export { resourceID, resourceTypeLabel }
