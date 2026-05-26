// Data sensitivity / classification levels. In real Corteza these are
// configurable Sensitivity Levels; for the POC we use a standard scheme.

export const SENSITIVITY_LEVELS = [
  { id: 'public', label: 'Public', severity: 'secondary' },
  { id: 'internal', label: 'Internal', severity: 'info' },
  { id: 'confidential', label: 'Confidential', severity: 'warn' },
  { id: 'restricted', label: 'Restricted', severity: 'danger' },
]

// For selectors: an explicit "Unclassified" choice instead of a clear button.
export const SENSITIVITY_OPTIONS = [{ id: null, label: 'Unclassified' }, ...SENSITIVITY_LEVELS]

export const sensitivity = id => SENSITIVITY_LEVELS.find(s => s.id === id)
export const sensitivityLabel = id => sensitivity(id)?.label || 'Unspecified'
