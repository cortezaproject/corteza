// Maps action-log resource types onto the project's own resource kinds, so an
// event about an agent wears the same icon and colour as that agent everywhere
// else in the project (KIND_CONFIG in ./kinds is the single source of truth for
// kind icons; category badges come from config/categories — do NOT invent icons
// here beyond the two project-level marks below, which nothing else owns).
//
// This mapping lives in the project section rather than in the admin ActionLog
// vocab: the vocab is shared vocabulary (labels), while icon/colour is project
// presentation, and admin must not depend on project config.
//
// Children map to their parent's kind — a record, a module field and a revision
// are all things you think of as "the module". Anything unmapped (settings,
// queues, templates…) falls through to kindConfig's neutral FALLBACK, which is
// correct: those are not project resources.

import { CATEGORY_CONFIG } from './categories'

const RESOURCE_KIND = {
  // Compose — records/fields/revisions are all "module" to a reader
  'corteza::compose:module': 'module',
  'corteza::compose:module-field': 'module',
  'corteza::compose:record': 'module',
  'corteza::compose:record-revision': 'module',
  'corteza::compose:page': 'page',
  'corteza::compose:page-layout': 'page',
  'corteza::compose:chart': 'chart',

  // System
  'corteza::system:agent': 'agent',
  'corteza::system:ai-conversation': 'agent',
  'corteza::system:chatbot': 'chatbot',
  'corteza::system:chatbot-session': 'chatbot',
  'corteza::system:role': 'role',
  'corteza::system:role-member': 'role',
  'corteza::system:user': 'user',
  'corteza::system:user-group': 'group',
  'corteza::system:connection': 'connection',
  'corteza::system:configured-connection': 'connection',
  'corteza::system:dal-connection': 'connection',

  // Automation — every flavour reads as "automation"
  'corteza::automation': 'automation',
  'corteza::automation:workflow': 'automation',
  'corteza::automation:trigger': 'automation',
  'corteza::automation:session': 'automation',
  'corteza::automation:ng-automation': 'automation',
}

// eventKind resolves an action-log resource TYPE (no trailing /id — use
// resourceType() first) to a project kind, or '' when it is not one.
export function eventKind(resourceType) {
  return RESOURCE_KIND[resourceType] || ''
}

// Colour overrides for events that went wrong, shaped to spread over a
// KIND_CONFIG (so the kind's own icon is kept and only the tone changes).
// Severity is the syslog uint8: 0–3 error-ish, 4 warning, 5–7 routine.
//
// Routine returns null — the event keeps its kind's colours. Only a real problem
// repaints the square, which is what makes it worth noticing.
const ALARM = {
  text: 'text-red-600 dark:text-red-400',
  bg: 'bg-red-50 dark:bg-red-950/40',
  ring: 'ring-red-200 dark:ring-red-800/60',
}

const WARN = {
  text: 'text-amber-600 dark:text-amber-400',
  bg: 'bg-amber-50 dark:bg-amber-950/40',
  ring: 'ring-amber-200 dark:ring-amber-800/60',
}

export function severityTone(severity) {
  if (severity <= 3) return ALARM
  if (severity === 4) return WARN
  return null
}

// The project's own mark for project/membership events — pi-sitemap (the
// wizard/overview's project glyph) in the indigo family, matching the events
// accent (EVENTS_COLOR) rather than any single category.
const PROJECT_BADGE = {
  icon: 'pi pi-sitemap',
  text: 'text-indigo-600 dark:text-indigo-400',
  bg: 'bg-indigo-50 dark:bg-indigo-950/40',
  ring: 'ring-indigo-200 dark:ring-indigo-800/60',
}

// Backlog items are category-neutral (they span all five) — same neutral
// th-large mark BacklogItemDialog/BacklogView use.
const BACKLOG_BADGE = { icon: 'pi pi-th-large', bg: 'bg-emphasis', ring: 'ring-surface', text: 'text-color' }

// Project-dashboard resources wear the dashboard's own badges (the category
// badge for category events, the backlog/project marks above) — these take
// precedence over RESOURCE_KIND/kindConfig in the timeline.
const RESOURCE_BADGE = {
  'corteza::system:project': PROJECT_BADGE,
  'corteza::system:project-member': PROJECT_BADGE,
  'corteza::system:project-group': PROJECT_BADGE,
  'corteza::system:project-group-entry': PROJECT_BADGE,
  'corteza::system:project-backlog-item': BACKLOG_BADGE,
  'corteza::system:project-incident': CATEGORY_CONFIG.incident.badge,
  'corteza::system:project-feature': CATEGORY_CONFIG.feature.badge,
  'corteza::system:project-privacy': CATEGORY_CONFIG.privacy.badge,
  'corteza::system:project-task': CATEGORY_CONFIG.task.badge,
  'corteza::system:project-review': CATEGORY_CONFIG.review.badge,
}

// eventBadge resolves an action-log resource TYPE to a ready badge config
// (KindIcon-shaped), or null when the resource has no project-specific mark
// (callers then fall back to kindConfig(eventKind(...))).
export function eventBadge(resourceType) {
  return RESOURCE_BADGE[resourceType] || null
}
