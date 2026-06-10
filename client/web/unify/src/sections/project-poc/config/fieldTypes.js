// Module field types — 1:1 with Corteza Compose module field kinds
// (see lib/vue/src/components/field/registry.ts). `id` is the compose kind.
// `Record` references another module (relationship); `User` references a user.

export const FIELD_TYPES = [
  { id: 'String', label: 'String' },
  { id: 'Number', label: 'Number' },
  { id: 'Bool', label: 'Yes / No (Bool)' },
  { id: 'DateTime', label: 'Date / Time' },
  { id: 'Select', label: 'Select' },
  { id: 'Email', label: 'Email' },
  { id: 'Url', label: 'URL' },
  { id: 'User', label: 'User' },
  { id: 'Record', label: 'Record (reference)' },
  { id: 'File', label: 'File' },
  { id: 'Geometry', label: 'Geometry' },
]

export const RECORD_KIND = 'Record'
export const isRecordRef = type => type === RECORD_KIND

export const fieldTypeLabel = id => FIELD_TYPES.find(t => t.id === id)?.label || id
