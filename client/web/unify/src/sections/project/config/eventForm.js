// New Event form — config-driven field sets per category, consumed by the
// GovernanceForm renderer (schema = array of sections; each field: key,
// labelKey, type, options?, placeholderKey?, badge?, default?). Field labels
// and placeholders are i18n keys — never plain English (see resolution in
// GovernanceForm.vue). Values persist to the per-project event resources via
// the events store; owner fields are user references resolved at render time.
const p = 'project.dashboard.event.f.'
const pp = 'project.dashboard.event.placeholder.'

// Shared option sets. Statuses are fixed; users are resolved at render time
// from the project's user directory (see NewEventDialog userOptions) — owner
// and approver fields carry `source: 'users'` instead of static options.
// All six work-item types (the five categories below, plus backlog items)
// share this one four-value set — AGREED INTENT (2026-07-28): review used to
// carry a shorter STATUS_REVIEW (no "Ready to Test"), but the Manage & Monitor
// board renders one shared column set across every type, so review's schema
// now uses the same STATUS list as everything else. Status is a free-text
// string in the backend, so this is a frontend-only change.
const STATUS = ['Open', 'In Progress', 'Ready to Test', 'Completed']
// Exported so other schemas that share the same lifecycle (backlog items —
// see components/dashboard/BacklogItemDialog.vue) don't redeclare the list.
export const EVENT_STATUS = STATUS
const SEVERITY = ['Critical', 'Serious', 'Major', 'Minor', 'Informational']
const RISK = ['Critical', 'High', 'Medium', 'Low', 'Very Low', 'None']

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

// Severity/risk/status selects (shared labels with the list columns). `badge`
// tells GovernanceForm to render the option/value slots with EventBadge
// (variant = the badge key) instead of plain text, matching how these fields
// already render in the list/drawer. `default` seeds NewEventDialog's
// create-mode model (see buildDefaults there) so these required fields never
// open on an empty selection; EventDetailDialog (edit) ignores it and seeds
// from the record instead.
const severity = {
  key: 'severity',
  labelKey: 'project.dashboard.columns.severity',
  type: 'select',
  options: SEVERITY,
  badge: 'severity',
  default: 'Minor',
}
const risk = {
  key: 'risk',
  labelKey: 'project.dashboard.columns.risk',
  type: 'select',
  options: RISK,
  badge: 'risk',
  default: 'Low',
}
// Status is repeated per category below (each category's schema pulls in its
// own `status(...)` field) so it's a factory rather than a single shared
// const like severity/risk — every category currently passes the same
// STATUS list (see the note above).
const status = options => ({
  key: 'status',
  labelKey: `${p}status`,
  type: 'select',
  options,
  badge: 'status',
  default: 'Open',
})

// Title text input — shown first in every category form (spans both columns).
const title = {
  key: 'title',
  labelKey: 'project.dashboard.columns.title',
  type: 'text',
  full: true,
  placeholderKey: `${pp}title`,
  required: true,
}

// Revision-reference select — the item's `revisionID` (0/null = unassigned).
// Work items carry a revision the way an issue carries a milestone
// (project.intent.md "Dashboards"): filed against the project, assigned to a
// revision, reassignable. `source: 'revisions'` mirrors the `user()` factory's
// `source: 'users'` mechanism — GovernanceForm itself stays option-agnostic;
// each consuming dialog (NewEventDialog/EventDetailDialog/BacklogItemDialog)
// injects the project's revision chain as options at render time (see their
// resolvedSchema + composables/useRevisionOptions.js), including an explicit
// "Unassigned" entry (value null) so the field can clear an assignment, not
// just set one. `default: null` seeds create-mode's model to that entry
// whenever the field is shown (see the dialogs' `allowRevisionSelect` prop —
// board-scoped create flows leave it unshown so their existing
// context-assigns-the-revision behaviour is untouched; edit mode always shows
// it, since reassigning an existing item is the point of this field).
export const REVISION_FIELD = {
  key: 'revisionID',
  labelKey: `${p}revision`,
  type: 'select',
  source: 'revisions',
  default: null,
}

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
  incident: [
    {
      fields: [
        title,
        {
          key: 'incidentType',
          labelKey: `${p}incidentType`,
          type: 'select',
          options: ['Serious Incident', 'Critical', 'Major', 'Minor', 'Informational'],
          required: true,
        },
        {
          key: 'groupSystem',
          labelKey: `${p}groupSystem`,
          type: 'select',
          options: [
            'AI Model Layer',
            'Data Pipeline',
            'API Gateway',
            'User Interface',
            'Infrastructure',
            'Access Control',
          ],
        },
        severity,
        risk,
        {
          key: 'description',
          labelKey: `${p}description`,
          type: 'textarea',
          placeholderKey: `${pp}incident.description`,
          required: true,
        },
        {
          key: 'riskIssue',
          labelKey: `${p}riskIssue`,
          type: 'textarea',
          placeholderKey: `${pp}incident.riskIssue`,
        },
        {
          key: 'changeRequired',
          labelKey: `${p}changeRequired`,
          type: 'textarea',
          placeholderKey: `${pp}incident.changeRequired`,
        },
        {
          key: 'riskChange',
          labelKey: `${p}riskChange`,
          type: 'textarea',
          placeholderKey: `${pp}incident.riskChange`,
        },
        user('issueOwner', true),
        user('changeOwner'),
        status(STATUS),
        user('changeApprovedBy'),
        { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
        { key: 'completedDate', labelKey: `${p}completedDate`, type: 'date' },
        REVISION_FIELD,
      ],
    },
  ],
  feature: [
    {
      fields: [
        title,
        {
          key: 'featureType',
          labelKey: `${p}featureType`,
          type: 'select',
          options: ['Platform', 'Data Model', 'Automation', 'Agent', 'Chatbot', 'Other'],
          required: true,
        },
        status(STATUS),
        severity,
        risk,
        {
          key: 'description',
          labelKey: `${p}description`,
          type: 'textarea',
          placeholderKey: `${pp}feature.description`,
          required: true,
        },
        {
          key: 'riskFeature',
          labelKey: `${p}riskFeature`,
          type: 'textarea',
          placeholderKey: `${pp}feature.riskFeature`,
        },
        {
          key: 'changeRequired',
          labelKey: `${p}changeRequired`,
          type: 'textarea',
          placeholderKey: `${pp}feature.changeRequired`,
        },
        {
          key: 'riskChange',
          labelKey: `${p}riskChange`,
          type: 'textarea',
          placeholderKey: `${pp}feature.riskChange`,
        },
        user('featureOwner', true),
        user('changeOwner'),
        user('changeApprovedBy'),
        { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
        REVISION_FIELD,
      ],
    },
  ],
  privacy: [
    {
      fields: [
        title,
        {
          key: 'requestType',
          labelKey: `${p}requestType`,
          type: 'select',
          options: ['Access (Export)', 'Rectification', 'Deletion'],
          required: true,
        },
        status(STATUS),
        severity,
        risk,
        {
          key: 'description',
          labelKey: `${p}description`,
          type: 'textarea',
          placeholderKey: `${pp}privacy.description`,
          required: true,
        },
        {
          key: 'riskAssessment',
          labelKey: `${p}riskAssessment`,
          type: 'textarea',
          placeholderKey: `${pp}privacy.riskAssessment`,
        },
        {
          key: 'changeRequired',
          labelKey: `${p}changeRequired`,
          type: 'textarea',
          placeholderKey: `${pp}privacy.changeRequired`,
        },
        {
          key: 'riskChange',
          labelKey: `${p}riskChange`,
          type: 'textarea',
          placeholderKey: `${pp}privacy.riskChange`,
        },
        user('requestOwner', true),
        user('changeOwner'),
        user('changeApprovedBy'),
        { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date', required: true },
        REVISION_FIELD,
      ],
    },
  ],
  task: [
    {
      fields: [
        title,
        // `taskName` (a fixed preset list) was dropped from the form — it
        // duplicated `title` without adding meaning. The backend column still
        // exists and old rows keep their value; nothing writes it any more.
        {
          key: 'taskType',
          labelKey: `${p}taskType`,
          type: 'select',
          options: ['Review/Action', 'Audit', 'Change'],
        },
        severity,
        risk,
        {
          key: 'description',
          labelKey: `${p}description`,
          type: 'textarea',
          placeholderKey: `${pp}task.description`,
        },
        status(STATUS),
        user('owner'),
        user('changeOwner'),
        { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
        { key: 'completedDate', labelKey: `${p}completedDate`, type: 'date' },
        REVISION_FIELD,
      ],
    },
  ],
  review: [
    {
      fields: [
        title,
        {
          key: 'reviewType',
          labelKey: `${p}reviewType`,
          type: 'select',
          options: [
            'Periodic Compliance Review',
            'AI Risk Review',
            'Access Control Review',
            'Model Performance Review',
          ],
          required: true,
        },
        {
          key: 'reviewFrequency',
          labelKey: `${p}reviewFrequency`,
          type: 'select',
          options: ['Monthly', 'Quarterly', 'Bi-Annual', 'Annual'],
        },
        {
          key: 'scope',
          labelKey: `${p}scope`,
          type: 'textarea',
          placeholderKey: `${pp}review.scope`,
        },
        user('reviewer', true),
        status(STATUS),
        { key: 'dateDue', labelKey: `${p}dateDue`, type: 'date' },
        user('approvedBy'),
        REVISION_FIELD,
      ],
    },
  ],
}
