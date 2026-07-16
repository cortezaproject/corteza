// Per-resource-kind visual config (icon + colors), reused from the project_old
// concept and extended for the new kinds (connection, chatbot). Tailwind classes
// are kept literal so JIT picks them up. `labelKey`/`singularKey` are i18n keys
// (plural / singular); components resolve them with $t for display.

import { STEPS } from '@/sections/project/config/pipeline'

export const KIND_CONFIG = {
  module: {
    labelKey: 'project.kinds.module.plural',
    singularKey: 'project.kinds.module.single',
    icon: 'pi pi-database',
    text: 'text-indigo-600 dark:text-indigo-400',
    bg: 'bg-indigo-50 dark:bg-indigo-950/40',
    ring: 'ring-indigo-200 dark:ring-indigo-800/60',
    stroke: '#4f46e5',
  },
  connection: {
    labelKey: 'project.kinds.connection.plural',
    singularKey: 'project.kinds.connection.single',
    icon: 'pi pi-link',
    text: 'text-teal-600 dark:text-teal-400',
    bg: 'bg-teal-50 dark:bg-teal-950/40',
    ring: 'ring-teal-200 dark:ring-teal-800/60',
    stroke: '#0d9488',
  },
  automation: {
    labelKey: 'project.kinds.automation.plural',
    singularKey: 'project.kinds.automation.single',
    icon: 'pi pi-bolt',
    text: 'text-orange-600 dark:text-orange-400',
    bg: 'bg-orange-50 dark:bg-orange-950/40',
    ring: 'ring-orange-200 dark:ring-orange-800/60',
    stroke: '#ea580c',
  },
  agent: {
    labelKey: 'project.kinds.agent.plural',
    singularKey: 'project.kinds.agent.single',
    icon: 'pi pi-sparkles',
    text: 'text-fuchsia-600 dark:text-fuchsia-400',
    bg: 'bg-fuchsia-50 dark:bg-fuchsia-950/40',
    ring: 'ring-fuchsia-200 dark:ring-fuchsia-800/60',
    stroke: '#c026d3',
  },
  chatbot: {
    labelKey: 'project.kinds.chatbot.plural',
    singularKey: 'project.kinds.chatbot.single',
    icon: 'pi pi-comments',
    text: 'text-rose-600 dark:text-rose-400',
    bg: 'bg-rose-50 dark:bg-rose-950/40',
    ring: 'ring-rose-200 dark:ring-rose-800/60',
    stroke: '#e11d48',
  },
  page: {
    labelKey: 'project.kinds.page.plural',
    singularKey: 'project.kinds.page.single',
    icon: 'pi pi-window-maximize',
    text: 'text-sky-600 dark:text-sky-400',
    bg: 'bg-sky-50 dark:bg-sky-950/40',
    ring: 'ring-sky-200 dark:ring-sky-800/60',
    stroke: '#0284c7',
  },
  chart: {
    labelKey: 'project.kinds.chart.plural',
    singularKey: 'project.kinds.chart.single',
    icon: 'pi pi-chart-bar',
    text: 'text-amber-600 dark:text-amber-400',
    bg: 'bg-amber-50 dark:bg-amber-950/40',
    ring: 'ring-amber-200 dark:ring-amber-800/60',
    stroke: '#d97706',
  },
  role: {
    labelKey: 'project.kinds.role.plural',
    singularKey: 'project.kinds.role.single',
    icon: 'pi pi-id-card',
    text: 'text-violet-600 dark:text-violet-400',
    bg: 'bg-violet-50 dark:bg-violet-950/40',
    ring: 'ring-violet-200 dark:ring-violet-800/60',
    stroke: '#7c3aed',
  },
  user: {
    labelKey: 'project.kinds.user.plural',
    singularKey: 'project.kinds.user.single',
    icon: 'pi pi-user',
    text: 'text-emerald-600 dark:text-emerald-400',
    bg: 'bg-emerald-50 dark:bg-emerald-950/40',
    ring: 'ring-emerald-200 dark:ring-emerald-800/60',
    stroke: '#059669',
  },
  // Not a resource step kind; used for Group nodes in the relationship graph.
  group: {
    labelKey: 'project.kinds.group.plural',
    singularKey: 'project.kinds.group.single',
    icon: 'pi pi-folder',
    text: 'text-surface-700 dark:text-surface-200',
    bg: 'bg-surface-100 dark:bg-surface-800',
    ring: 'ring-surface-300 dark:ring-surface-600',
    stroke: '#475569',
  },
}

const FALLBACK = {
  labelKey: 'project.kinds.fallback.plural',
  singularKey: 'project.kinds.fallback.single',
  icon: 'pi pi-circle',
  text: 'text-color',
  bg: 'bg-emphasis',
  ring: 'ring-surface',
  stroke: '#6b7280',
}

export const kindConfig = kind => KIND_CONFIG[kind] || FALLBACK

// Resource kinds in pipeline order (one per resource-collection step). Only
// modules are backend-backed today; other kinds rejoin with their steps.
export const RESOURCE_KINDS = STEPS.filter(s => s.type === 'resource').map(s => s.kind)

// Full system overview, in the conceptual pipeline's order. The resource
// metrics strip always shows every kind — it's an overview of the whole
// system, not just the steps built so far.
export const OVERVIEW_KINDS = [
  'module',
  'connection',
  'automation',
  'agent',
  'chatbot',
  'page',
  'chart',
  'role',
  'user',
]

// Semantic layers for the graph's layer selector. Each preset toggles its set
// of kinds together; the per-kind chips toggle them individually. Roles/users
// are deliberately NOT a preset here — they ride the separate "show access"
// overlay (RBAC edges), handled apart from these node layers.
export const LAYERS = [
  { key: 'data', labelKey: 'project.graph.layer.data', kinds: ['connection', 'module'] },
  { key: 'logic', labelKey: 'project.graph.layer.logic', kinds: ['automation', 'agent'] },
  {
    key: 'experience',
    labelKey: 'project.graph.layer.experience',
    kinds: ['chatbot', 'page', 'chart'],
  },
]

// Every kind toggled by the layer presets/chips (i.e. everything except the
// access overlay). Used to seed the graph view's visible set so all node
// layers start on.
export const NODE_LAYER_KINDS = LAYERS.flatMap(l => l.kinds)

// Roles/users surface through the "show access" overlay only, off by default.
export const ACCESS_KINDS = ['role', 'user']
