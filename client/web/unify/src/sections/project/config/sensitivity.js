// Data sensitivity / classification levels. These are the standard scheme the
// store seeds as real DAL sensitivity-level resources when missing (see
// ensureStandardLevels in stores/projects.js); `id` is the resource handle.
// `labelKey` is an i18n key; components resolve it with $t for display.

export const SENSITIVITY_LEVELS = [
  { id: 'public', labelKey: 'project.sensitivity.public', severity: 'success' },
  { id: 'internal', labelKey: 'project.sensitivity.internal', severity: 'info' },
  { id: 'confidential', labelKey: 'project.sensitivity.confidential', severity: 'warn' },
  { id: 'restricted', labelKey: 'project.sensitivity.restricted', severity: 'danger' },
]

// For selectors: an explicit "Unclassified" choice instead of a clear button.
export const SENSITIVITY_OPTIONS = [
  { id: null, labelKey: 'project.sensitivity.unclassified' },
  ...SENSITIVITY_LEVELS,
]

export const sensitivity = id => SENSITIVITY_LEVELS.find(s => s.id === id)
export const sensitivityLabelKey = id => sensitivity(id)?.labelKey || 'project.sensitivity.unspecified'
