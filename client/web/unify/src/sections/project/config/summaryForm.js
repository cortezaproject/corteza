// Field schema for the Project Summary governance step. Drives GovernanceForm.
// Section titles mirror the groupings in the source spec.

export const SUMMARY_SCHEMA = [
  {
    title: 'General',
    fields: [
      { key: 'systemName', label: 'Project Name', type: 'text' },
      {
        key: 'orgRole',
        label: 'Organisation Role',
        type: 'multiselect',
        options: ['Provider', 'Deployer', 'Importer', 'Distributor', 'Authorised Representative'],
        default: ['Provider'],
      },
      { key: 'intendedPurpose', label: 'Project Objective', type: 'textarea' },
    ],
  },
]

// Build the default values object from the schema.
export function summaryDefaults() {
  const out = {}
  for (const section of SUMMARY_SCHEMA) {
    for (const f of section.fields) {
      if ('default' in f) out[f.key] = f.default
      else out[f.key] = f.type === 'multiselect' ? [] : null
    }
  }
  return out
}
