import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

// In-memory mock store for the published-project dashboard's category pages
// (incident / feature / privacy / task / review). There is NO backend: this is
// a POC data layer. Creating an item pushes it here so the CResourceList, the
// KPI row, the charts AND the nav badge all update live off the same ref.
//
// The list is global (not per-project) — deliberately, for the scaffold. Event
// objects are loose bags: { id, category, title, status, dateDue, ...form } —
// they carry whatever the per-category form (config/eventForm.js EVENT_FORMS)
// produced, plus the four fields every consumer relies on.

// ID prefix per category. The numeric suffix is zero-padded per category.
const PREFIXES = {
  incident: 'INC',
  feature: 'FEA',
  privacy: 'PRV',
  task: 'TSK',
  review: 'REV',
}

// Seed data — ~4 realistic events per category, using the option VALUES from
// EVENT_FORMS so they render identically to freshly-created ones. Statuses are
// spread across the lifecycle; dateDue mixes past (some overdue) and future.
// Each seed carries its category's type field and owner field (see the field
// keys in config/eventForm.js).
// Enrichment fields mirror the reference demo's events table so the per-category
// lists can show the full column set: `description` (title subtitle), `severity`
// (level badge), `risk` (level badge), `changeOwner`, `backlog` (linked BL ids),
// and the approval field (changeApprovedBy / approvedBy).
const SEED = [
  // --- incident (type=incidentType, owner=issueOwner) ---------------------------
  { id: 'INC-0001', category: 'incident', title: 'Model output leaked PII in export', description: 'Systematic PII leak in the record export path', status: 'In Progress', dateDue: '2026-06-28', incidentType: 'Serious Incident', severity: 'Critical', risk: 'Critical', issueOwner: 'N. McCarthy', changeOwner: 'A. Dupont', backlog: ['BL-014', 'BL-015'], changeApprovedBy: 'P. Martin', groupSystem: 'AI Model Layer' },
  { id: 'INC-0002', category: 'incident', title: 'API gateway 5xx spike under load', description: 'Gateway returning 5xx above threshold at peak', status: 'Open', dateDue: '2026-07-18', incidentType: 'Major', severity: 'Major', risk: 'High', issueOwner: 'S. Roche', changeOwner: 'N. McCarthy', backlog: ['BL-017'], changeApprovedBy: 'P. Martin', groupSystem: 'API Gateway' },
  { id: 'INC-0003', category: 'incident', title: 'Stale data in nightly pipeline', description: 'Nightly ETL served stale partitions', status: 'Completed', dateDue: '2026-05-30', incidentType: 'Minor', severity: 'Minor', risk: 'Low', issueOwner: 'A. Dupont', changeOwner: 'A. Dupont', backlog: ['BL-006', 'BL-007'], changeApprovedBy: 'P. Martin', groupSystem: 'Data Pipeline' },
  { id: 'INC-0004', category: 'incident', title: 'Access control misconfiguration', description: 'Role grants wider than intended on a module', status: 'Ready to Test', dateDue: '2026-08-05', incidentType: 'Critical', severity: 'Serious', risk: 'High', issueOwner: 'M. Laurent', changeOwner: 'N. McCarthy', backlog: ['BL-020'], changeApprovedBy: 'P. Martin', groupSystem: 'Access Control' },

  // --- feature (type=featureType, owner=featureOwner) ---------------------------
  { id: 'FEA-0001', category: 'feature', title: 'Bulk record import for data model', description: 'CSV importer with field mapping', status: 'In Progress', dateDue: '2026-07-22', featureType: 'Data Model', severity: 'Minor', risk: 'Low', featureOwner: 'S. Roche', changeOwner: 'A. Dupont', backlog: ['BL-008', 'BL-009'], changeApprovedBy: 'P. Martin' },
  { id: 'FEA-0002', category: 'feature', title: 'Slack notification automation', description: 'Notify channels on record events', status: 'Open', dateDue: '2026-08-14', featureType: 'Automation', severity: 'Minor', risk: 'Low', featureOwner: 'A. Dupont', changeOwner: 'N. McCarthy', backlog: ['BL-021'], changeApprovedBy: 'P. Martin' },
  { id: 'FEA-0003', category: 'feature', title: 'Support agent handoff flow', description: 'Escalate chatbot sessions to human agents', status: 'Ready to Test', dateDue: '2026-06-20', featureType: 'Agent', severity: 'Major', risk: 'Medium', featureOwner: 'P. Martin', changeOwner: 'N. McCarthy', backlog: ['BL-010', 'BL-011', 'BL-012'], changeApprovedBy: 'P. Martin' },
  { id: 'FEA-0004', category: 'feature', title: 'Onboarding chatbot scenarios', description: 'Guided onboarding conversation flows', status: 'Completed', dateDue: '2026-05-12', featureType: 'Chatbot', severity: 'Minor', risk: 'Low', featureOwner: 'N. McCarthy', changeOwner: 'A. Dupont', backlog: ['BL-005'], changeApprovedBy: 'P. Martin' },

  // --- privacy (type=requestType, owner=requestOwner) ---------------------------
  { id: 'PRV-0001', category: 'privacy', title: 'Data subject access request — export', description: 'Full data export for a data subject', status: 'Open', dateDue: '2026-07-01', requestType: 'Access (Export)', severity: 'Informational', risk: 'Low', requestOwner: 'M. Laurent', changeOwner: 'S. Roche', backlog: ['BL-019'], changeApprovedBy: '' },
  { id: 'PRV-0002', category: 'privacy', title: 'Right to erasure — closed account', description: 'Erase data for a deactivated account', status: 'In Progress', dateDue: '2026-06-15', requestType: 'Deletion', severity: 'Serious', risk: 'High', requestOwner: 'N. McCarthy', changeOwner: 'S. Roche', backlog: ['BL-016'], changeApprovedBy: '' },
  { id: 'PRV-0003', category: 'privacy', title: 'Rectify inaccurate contact record', description: 'Correct contact details on request', status: 'Completed', dateDue: '2026-05-20', requestType: 'Rectification', severity: 'Minor', risk: 'Low', requestOwner: 'S. Roche', changeOwner: 'S. Roche', backlog: ['BL-013'], changeApprovedBy: 'P. Martin' },
  { id: 'PRV-0004', category: 'privacy', title: 'Bulk deletion request — vendor data', description: 'Delete vendor-provided datasets', status: 'Open', dateDue: '2026-08-30', requestType: 'Deletion', severity: 'Major', risk: 'Medium', requestOwner: 'A. Dupont', changeOwner: 'S. Roche', backlog: ['BL-022'], changeApprovedBy: '' },

  // --- task (type=taskType, owner=owner) ----------------------------------------
  { id: 'TSK-0001', category: 'task', title: 'Quarterly RBAC review', description: 'Review role grants across modules', status: 'Open', dateDue: '2026-07-25', taskType: 'Review/Action', taskName: 'RBAC Review', severity: 'Minor', risk: 'Low', owner: 'P. Martin', changeOwner: 'M. Laurent', backlog: ['BL-018'] },
  { id: 'TSK-0002', category: 'task', title: 'AI model version rollout', description: 'Promote the new model to production', status: 'In Progress', dateDue: '2026-06-30', taskType: 'Change', taskName: 'AI Model Change', severity: 'Major', risk: 'Medium', owner: 'A. Dupont', changeOwner: 'A. Dupont', backlog: ['BL-023'] },
  { id: 'TSK-0003', category: 'task', title: 'Annual compliance audit', description: 'External audit of controls', status: 'Ready to Test', dateDue: '2026-09-10', taskType: 'Audit', taskName: 'Compliance Audit', severity: 'Minor', risk: 'Low', owner: 'N. McCarthy', changeOwner: 'N. McCarthy', backlog: ['BL-024'] },
  { id: 'TSK-0004', category: 'task', title: 'Risk register refresh', description: 'Update the risk register entries', status: 'Completed', dateDue: '2026-05-05', taskType: 'Review/Action', taskName: 'Risk Register Update', severity: 'Minor', risk: 'Low', owner: 'M. Laurent', changeOwner: 'M. Laurent', backlog: ['BL-004'] },

  // --- review (type=reviewType, owner=reviewer) ---------------------------------
  { id: 'REV-0001', category: 'review', title: 'Q3 periodic compliance review', description: 'Quarterly compliance checkpoint', status: 'Open', dateDue: '2026-07-31', reviewType: 'Periodic Compliance Review', reviewFrequency: 'Quarterly', reviewer: 'N. McCarthy', approvedBy: 'P. Martin' },
  { id: 'REV-0002', category: 'review', title: 'AI risk review — new model', description: 'Assess risks of a new model release', status: 'In Progress', dateDue: '2026-06-18', reviewType: 'AI Risk Review', reviewFrequency: 'Monthly', reviewer: 'P. Martin', approvedBy: 'P. Martin' },
  { id: 'REV-0003', category: 'review', title: 'Access control review', description: 'Review access grants and roles', status: 'Completed', dateDue: '2026-05-28', reviewType: 'Access Control Review', reviewFrequency: 'Bi-Annual', reviewer: 'S. Roche', approvedBy: 'N. McCarthy' },
  { id: 'REV-0004', category: 'review', title: 'Model performance review', description: 'Evaluate model performance metrics', status: 'Open', dateDue: '2026-08-22', reviewType: 'Model Performance Review', reviewFrequency: 'Annual', reviewer: 'A. Dupont', approvedBy: 'P. Martin' },
]

export const useEventsStore = defineStore('events', () => {
  // Single source of truth. New items are UNSHIFTed, so the array is naturally
  // newest-first; every getter derives from it and stays reactive.
  const events = ref([...SEED])

  // Items for one category, newest first (array order already is, thanks to
  // unshift-on-add + seeds treated as the initial backlog).
  const byCategory = computed(() => cat => events.value.filter(e => e.category === cat))

  const countByCategory = computed(() => cat => byCategory.value(cat).length)

  // KPI trio for a category. `open` = not Completed; `overdue` = due in the past
  // and not Completed (bad/blank dates are guarded and never count as overdue).
  const kpis = computed(() => cat => {
    const items = byCategory.value(cat)
    const now = new Date()
    let open = 0
    let overdue = 0
    for (const e of items) {
      const isOpen = e.status !== 'Completed'
      if (isOpen) open++
      const due = e.dateDue ? new Date(e.dateDue) : null
      if (isOpen && due && !Number.isNaN(due.getTime()) && due < now) overdue++
    }
    return { total: items.length, open, overdue }
  })

  // Group a category's items by one field and count each bucket. Missing/blank
  // values collapse into an em-dash bucket. Shaped for the bar charts:
  // [{ label, value }].
  const breakdown = computed(() => (cat, field) => {
    const counts = new Map()
    for (const e of byCategory.value(cat)) {
      const label = String(e[field] || '—')
      counts.set(label, (counts.get(label) || 0) + 1)
    }
    return [...counts.entries()].map(([label, value]) => ({ label, value }))
  })

  // Next zero-padded id for a category, from the current max numeric suffix.
  function nextId(cat) {
    const prefix = PREFIXES[cat] || cat.slice(0, 3).toUpperCase()
    let max = 0
    for (const e of events.value) {
      if (e.category !== cat) continue
      const n = parseInt(String(e.id).split('-')[1], 10)
      if (Number.isFinite(n) && n > max) max = n
    }
    return `${prefix}-${String(max + 1).padStart(4, '0')}`
  }

  // Create an event from a form payload and push it to the front. Title falls
  // back to a description snippet, then a generic label; status defaults to
  // 'Open'. The whole payload is spread so every form field is retained.
  function add(cat, payload = {}) {
    const event = {
      ...payload,
      id: nextId(cat),
      category: cat,
      title:
        payload.title ||
        (payload.description ? String(payload.description).slice(0, 80) : '(untitled)'),
      status: payload.status || 'Open',
      // Enrichment columns the form doesn't collect — kept present so the list
      // cells render cleanly (empty until a real backend supplies them).
      backlog: payload.backlog || [],
    }
    events.value.unshift(event)
    return event
  }

  return {
    events,
    byCategory,
    countByCategory,
    kpis,
    breakdown,
    add,
  }
})
