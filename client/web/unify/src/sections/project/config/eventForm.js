// New Event form — config-driven field sets per category, consumed by the
// GovernanceForm renderer (schema = array of sections; each field: key,
// labelKey, type, options?, placeholder?). Field labels are i18n keys; option
// values and placeholders are plain strings (mock/demo data for the scaffold).
// No backend: submitting just toasts and navigates to the Events view.
const p = 'project.dashboard.event.f.'

// Shared option sets (mock users/statuses until wired to real stores).
const STATUS = ['Open', 'In Progress', 'Ready to Test', 'Completed']
const STATUS_REVIEW = ['Open', 'In Progress', 'Completed']
const USERS = ['N. McCarthy', 'S. Roche', 'A. Dupont', 'M. Laurent', 'P. Martin']
const APPROVERS = ['N. McCarthy', 'P. Martin']

// The five event categories shown in step 1 (label/desc are i18n keys, icon is
// a PrimeIcons class). Order mirrors the reference demo.
export const EVENT_CATEGORIES = [
  { key: 'incident', icon: 'pi-exclamation-triangle', color: 'text-red-500' },
  { key: 'feature', icon: 'pi-sparkles', color: 'text-blue-500' },
  { key: 'privacy', icon: 'pi-shield', color: 'text-purple-500' },
  { key: 'task', icon: 'pi-check-square', color: 'text-emerald-500' },
  { key: 'review', icon: 'pi-sync', color: 'text-amber-500' },
]

// Per-category schemas (single section each). GovernanceForm renders these.
export const EVENT_FORMS = {
  incident: [{ fields: [
    { key: 'incidentType', labelKey: `${p}incidentType`, type: 'select', options: ['Serious Incident', 'Critical', 'Major', 'Minor', 'Informational'], required: true },
    { key: 'groupSystem', labelKey: `${p}groupSystem`, type: 'select', options: ['AI Model Layer', 'Data Pipeline', 'API Gateway', 'User Interface', 'Infrastructure', 'Access Control'] },
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Describe the incident in detail…', required: true },
    { key: 'riskIssue', labelKey: `${p}riskIssue`, type: 'textarea', placeholder: 'Identify the risk, impact, likelihood, and affected parties…' },
    { key: 'changeRequired', labelKey: `${p}changeRequired`, type: 'textarea', placeholder: 'What corrective action or change is required?…' },
    { key: 'riskChange', labelKey: `${p}riskChange`, type: 'textarea', placeholder: 'Assess the risk of implementing the proposed change…' },
    { key: 'issueOwner', labelKey: `${p}issueOwner`, type: 'select', options: USERS, required: true },
    { key: 'changeOwner', labelKey: `${p}changeOwner`, type: 'select', options: USERS },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    { key: 'changeApprovedBy', labelKey: `${p}changeApprovedBy`, type: 'select', options: APPROVERS },
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
    { key: 'completedDate', labelKey: `${p}completedDate`, type: 'date' },
  ] }],
  feature: [{ fields: [
    { key: 'featureType', labelKey: `${p}featureType`, type: 'select', options: ['Platform', 'Data Model', 'Automation', 'Agent', 'Chatbot', 'Other'], required: true },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Describe the feature or update…', required: true },
    { key: 'riskFeature', labelKey: `${p}riskFeature`, type: 'textarea', placeholder: 'Assess the risk introduced by this feature…' },
    { key: 'changeRequired', labelKey: `${p}changeRequired`, type: 'textarea', placeholder: 'Any mitigations or changes required?…' },
    { key: 'riskChange', labelKey: `${p}riskChange`, type: 'textarea', placeholder: 'Risk of the corrective change…' },
    { key: 'featureOwner', labelKey: `${p}featureOwner`, type: 'select', options: USERS, required: true },
    { key: 'changeOwner', labelKey: `${p}changeOwner`, type: 'select', options: USERS },
    { key: 'changeApprovedBy', labelKey: `${p}changeApprovedBy`, type: 'select', options: APPROVERS },
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
  ] }],
  privacy: [{ fields: [
    { key: 'requestType', labelKey: `${p}requestType`, type: 'select', options: ['Access (Export)', 'Rectification', 'Deletion'], required: true },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Describe the privacy request…', required: true },
    { key: 'riskAssessment', labelKey: `${p}riskAssessment`, type: 'textarea', placeholder: 'GDPR / AI Act risk considerations…' },
    { key: 'changeRequired', labelKey: `${p}changeRequired`, type: 'textarea', placeholder: 'Data deletion steps, systems to update…' },
    { key: 'riskChange', labelKey: `${p}riskChange`, type: 'textarea', placeholder: 'Any downstream impact of fulfilling the request…' },
    { key: 'requestOwner', labelKey: `${p}requestOwner`, type: 'select', options: USERS, required: true },
    { key: 'changeOwner', labelKey: `${p}changeOwner`, type: 'select', options: USERS },
    { key: 'changeApprovedBy', labelKey: `${p}changeApprovedBy`, type: 'select', options: APPROVERS },
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date', required: true },
  ] }],
  task: [{ fields: [
    { key: 'taskName', labelKey: `${p}taskName`, type: 'select', options: ['RBAC Review', 'AI Model Change', 'Risk Register Update', 'Compliance Audit', 'Custom…'], required: true },
    { key: 'taskType', labelKey: `${p}taskType`, type: 'select', options: ['Review/Action', 'Audit', 'Change'] },
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Detail what needs to be done…' },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    { key: 'owner', labelKey: `${p}owner`, type: 'select', options: USERS, required: true },
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
    { key: 'completedDate', labelKey: `${p}completedDate`, type: 'date' },
  ] }],
  review: [{ fields: [
    { key: 'reviewType', labelKey: `${p}reviewType`, type: 'select', options: ['Periodic Compliance Review', 'AI Risk Review', 'Access Control Review', 'Model Performance Review'], required: true },
    { key: 'reviewFrequency', labelKey: `${p}reviewFrequency`, type: 'select', options: ['Monthly', 'Quarterly', 'Bi-Annual', 'Annual'] },
    { key: 'scope', labelKey: `${p}scope`, type: 'textarea', placeholder: 'What will be reviewed and what outcomes are expected…' },
    { key: 'reviewer', labelKey: `${p}reviewer`, type: 'select', options: USERS, required: true },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS_REVIEW },
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
    { key: 'approvedBy', labelKey: `${p}approvedBy`, type: 'select', options: APPROVERS },
  ] }],
}
