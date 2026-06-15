// Field schema for the Project Summary governance step. Drives GovernanceForm.
// Section titles mirror the groupings in the source spec. `titleKey`/`labelKey`/
// `descriptionKey` are i18n keys; GovernanceForm resolves them with $t.
// The multiselect `options` are AI Act enum values that persist verbatim — they
// stay literal (untranslated), like a status or mode value.

export const SUMMARY_SCHEMA = [
  {
    titleKey: 'project.summaryForm.sections.general',
    fields: [
      { key: 'systemName', labelKey: 'project.summaryForm.fields.systemName', type: 'text' },
      {
        key: 'orgRole',
        labelKey: 'project.summaryForm.fields.orgRole.label',
        descriptionKey: 'project.summaryForm.fields.orgRole.description',
        type: 'multiselect',
        options: ['Provider', 'Deployer', 'Importer', 'Distributor', 'Authorised Representative'],
        default: ['Provider'],
      },
      {
        key: 'intendedPurpose',
        labelKey: 'project.summaryForm.fields.intendedPurpose.label',
        descriptionKey: 'project.summaryForm.fields.intendedPurpose.description',
        type: 'textarea',
      },
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
