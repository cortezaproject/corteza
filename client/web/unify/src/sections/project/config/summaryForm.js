// Field schema for the Project Summary governance step. Drives GovernanceForm.
// Section titles mirror the groupings in the source spec.

export const SUMMARY_SCHEMA = [
  {
    title: 'General',
    fields: [
      { key: 'systemName', label: 'System Name', type: 'text', fullWidth: true },
      { key: 'version', label: 'Version', type: 'number' },
      { key: 'releaseDate', label: 'Latest Production Release Date', type: 'datetime' },
      {
        key: 'releaseStage',
        label: 'Release Stage',
        type: 'select',
        options: ['Dev', 'Test (Preview)', 'Production'],
        default: 'Dev',
      },
      {
        key: 'orgRole',
        label: 'Organisation Role',
        type: 'multiselect',
        options: ['Provider', 'Deployer', 'Importer', 'Distributor', 'Authorised Representative'],
        default: ['Provider'],
      },
      { key: 'intendedPurpose', label: 'Intended Purpose', type: 'textarea', fullWidth: true },
      {
        key: 'primaryUsers',
        label: 'Primary Users',
        type: 'multiselect',
        options: ['Employees', 'Public', 'Customers', 'Government', 'Vulnerable Groups'],
        fullWidth: true,
      },
      {
        key: 'vulnerableGroups',
        label: 'Vulnerable Groups',
        type: 'multiselect',
        options: [
          'Children and Young People',
          'Persons with Disabilities',
          'Elderly Persons',
          'Socio-economically disadvantaged persons',
          'Migrants, asylum seekers and minorities',
          'Patients and people dependent on healthcare or social care',
          'Workers and jobseekers in asymmetrical power relationships',
          'Persons with low digital or AI literacy',
          'Any group whose specific traits can be exploited or who is systematically at risk in context',
        ],
        fullWidth: true,
      },
      {
        key: 'vulnerableGroupPrecisions',
        label: 'Vulnerable Group Additional Precisions',
        type: 'text',
        placeholder: 'e.g. Tenants in Social Housing',
        fullWidth: true,
      },
    ],
  },
  {
    title: 'Geography & GPAI',
    fields: [
      {
        key: 'geoDeployment',
        label: 'Geographic Deployment',
        type: 'select',
        options: ['EU', 'Non-EU', 'Both'],
        default: 'EU',
      },
      {
        key: 'geoUse',
        label: 'Geographic Use',
        type: 'select',
        options: ['EU', 'Non-EU', 'Both'],
        default: 'EU',
      },
      {
        key: 'gpaiUsed',
        label: 'Does the system use or embed a general purpose AI (GPAI) model?',
        type: 'select',
        options: ['Yes', 'No', 'Unknown'],
        default: 'Yes',
      },
      {
        key: 'orgIsGpaiProvider',
        label: 'Is your organisation the provider of that general purpose AI (GPAI) model?',
        type: 'select',
        options: ['Yes', 'No', 'Unknown'],
        default: 'No',
      },
      {
        key: 'sourcing',
        label: 'Open-Source or Proprietary',
        type: 'select',
        options: ['Open-Source', 'Proprietary', 'Mixed'],
        default: 'Mixed',
      },
    ],
  },
  {
    title: 'Deployment context',
    fields: [
      {
        key: 'publiclyAccessible',
        label: 'Is the system used in publicly accessible spaces?',
        type: 'select',
        options: ['Yes', 'No'],
        default: 'No',
      },
      {
        key: 'publicAuthority',
        label: 'Is the system used by or on behalf of a public authority?',
        type: 'select',
        options: ['Yes', 'No', 'Unknown'],
        default: 'No',
      },
      {
        key: 'safetyComponent',
        label: 'Is the system a safety component of a regulated product?',
        type: 'select',
        options: ['Yes', 'No', 'Unknown'],
        default: 'No',
      },
      { key: 'regulatedProductCategory', label: 'Regulated product category', type: 'text', fullWidth: true },
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
