// New Event form — config-driven field sets per category, consumed by the
// GovernanceForm renderer (schema = array of sections; each field: key,
// labelKey, type, options?, placeholder?). Field labels are i18n keys. Values
// persist to the per-project event resources via the events store; owner fields
// are user references resolved at render time.
const p = 'project.dashboard.event.f.'

// Shared option sets. Statuses are fixed; users are resolved at render time
// from the project's user directory (see NewEventDialog userOptions) — owner
// and approver fields carry `source: 'users'` instead of static options.
const STATUS = ['Open', 'In Progress', 'Ready to Test', 'Completed']
const STATUS_REVIEW = ['Open', 'In Progress', 'Completed']
const SEVERITY = ['Critical', 'Serious', 'Major', 'Minor', 'Informational']
const RISK = ['Critical', 'High', 'Medium', 'Low']

// A user-reference select — options/labels/values are injected from the project
// user directory at render time (backend stores these as user IDs).
const user = (key, required = false) => ({
  key,
  labelKey: `${p}${key}`,
  type: 'select',
  source: 'users',
  filter: true,
  required,
})

// Severity/risk selects (shared labels with the list columns).
const severity = { key: 'severity', labelKey: 'project.dashboard.columns.severity', type: 'select', options: SEVERITY }
const risk = { key: 'risk', labelKey: 'project.dashboard.columns.risk', type: 'select', options: RISK }

// Title text input — shown first in every category form (spans both columns).
const title = { key: 'title', labelKey: 'project.dashboard.columns.title', type: 'text', full: true, placeholder: 'Short, descriptive title…', required: true }

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
    title,
    { key: 'incidentType', labelKey: `${p}incidentType`, type: 'select', options: ['Serious Incident', 'Critical', 'Major', 'Minor', 'Informational'], required: true },
    { key: 'groupSystem', labelKey: `${p}groupSystem`, type: 'select', options: ['AI Model Layer', 'Data Pipeline', 'API Gateway', 'User Interface', 'Infrastructure', 'Access Control'] },
    severity,
    risk,
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Describe the incident in detail…', required: true },
    { key: 'riskIssue', labelKey: `${p}riskIssue`, type: 'textarea', placeholder: 'Identify the risk, impact, likelihood, and affected parties…' },
    { key: 'changeRequired', labelKey: `${p}changeRequired`, type: 'textarea', placeholder: 'What corrective action or change is required?…' },
    { key: 'riskChange', labelKey: `${p}riskChange`, type: 'textarea', placeholder: 'Assess the risk of implementing the proposed change…' },
    user('issueOwner', true),
    user('changeOwner'),
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    user('changeApprovedBy'),
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
    { key: 'completedDate', labelKey: `${p}completedDate`, type: 'date' },
  ] }],
  feature: [{ fields: [
    title,
    { key: 'featureType', labelKey: `${p}featureType`, type: 'select', options: ['Platform', 'Data Model', 'Automation', 'Agent', 'Chatbot', 'Other'], required: true },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    severity,
    risk,
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Describe the feature or update…', required: true },
    { key: 'riskFeature', labelKey: `${p}riskFeature`, type: 'textarea', placeholder: 'Assess the risk introduced by this feature…' },
    { key: 'changeRequired', labelKey: `${p}changeRequired`, type: 'textarea', placeholder: 'Any mitigations or changes required?…' },
    { key: 'riskChange', labelKey: `${p}riskChange`, type: 'textarea', placeholder: 'Risk of the corrective change…' },
    user('featureOwner', true),
    user('changeOwner'),
    user('changeApprovedBy'),
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
  ] }],
  privacy: [{ fields: [
    title,
    { key: 'requestType', labelKey: `${p}requestType`, type: 'select', options: ['Access (Export)', 'Rectification', 'Deletion'], required: true },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    severity,
    risk,
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Describe the privacy request…', required: true },
    { key: 'riskAssessment', labelKey: `${p}riskAssessment`, type: 'textarea', placeholder: 'GDPR / AI Act risk considerations…' },
    { key: 'changeRequired', labelKey: `${p}changeRequired`, type: 'textarea', placeholder: 'Data deletion steps, systems to update…' },
    { key: 'riskChange', labelKey: `${p}riskChange`, type: 'textarea', placeholder: 'Any downstream impact of fulfilling the request…' },
    user('requestOwner', true),
    user('changeOwner'),
    user('changeApprovedBy'),
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date', required: true },
  ] }],
  task: [{ fields: [
    title,
    { key: 'taskName', labelKey: `${p}taskName`, type: 'select', options: ['RBAC Review', 'AI Model Change', 'Risk Register Update', 'Compliance Audit', 'Custom…'], required: true },
    { key: 'taskType', labelKey: `${p}taskType`, type: 'select', options: ['Review/Action', 'Audit', 'Change'] },
    severity,
    risk,
    { key: 'description', labelKey: `${p}description`, type: 'textarea', placeholder: 'Detail what needs to be done…' },
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS },
    user('owner', true),
    user('changeOwner'),
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
    { key: 'completedDate', labelKey: `${p}completedDate`, type: 'date' },
  ] }],
  review: [{ fields: [
    title,
    { key: 'reviewType', labelKey: `${p}reviewType`, type: 'select', options: ['Periodic Compliance Review', 'AI Risk Review', 'Access Control Review', 'Model Performance Review'], required: true },
    { key: 'reviewFrequency', labelKey: `${p}reviewFrequency`, type: 'select', options: ['Monthly', 'Quarterly', 'Bi-Annual', 'Annual'] },
    { key: 'scope', labelKey: `${p}scope`, type: 'textarea', placeholder: 'What will be reviewed and what outcomes are expected…' },
    user('reviewer', true),
    { key: 'status', labelKey: `${p}status`, type: 'select', options: STATUS_REVIEW },
    { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
    user('approvedBy'),
  ] }],
}
