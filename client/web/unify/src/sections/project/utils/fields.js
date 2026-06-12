// Compose field names must be handles; the human label carries the display
// name. Shared by the store (marshalling) and the field editors (e.g. the
// Record label-field picker must predict the machine name of unsaved fields).
export const fieldName = label =>
  (label || 'field')
    .trim()
    .replace(/[^a-zA-Z0-9_]+/g, '_')
    .replace(/^([0-9])/, 'f$1')
    .replace(/(^_+|_+$)/g, '') || 'field'
