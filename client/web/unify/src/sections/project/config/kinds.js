// Per-resource-kind visual config (icon + colors), reused from the project_old
// concept and extended for the new kinds (connection, chatbot). Tailwind classes
// are kept literal so JIT picks them up.

import { STEPS } from '@/sections/project/config/pipeline'

export const KIND_CONFIG = {
  module: {
    label: 'Modules',
    icon: 'pi pi-database',
    text: 'text-indigo-600 dark:text-indigo-400',
    bg: 'bg-indigo-50 dark:bg-indigo-950/40',
    ring: 'ring-indigo-200 dark:ring-indigo-800/60',
    stroke: '#4f46e5',
  },
  connection: {
    label: 'Connections',
    icon: 'pi pi-link',
    text: 'text-teal-600 dark:text-teal-400',
    bg: 'bg-teal-50 dark:bg-teal-950/40',
    ring: 'ring-teal-200 dark:ring-teal-800/60',
    stroke: '#0d9488',
  },
  automation: {
    label: 'Automations',
    icon: 'pi pi-bolt',
    text: 'text-orange-600 dark:text-orange-400',
    bg: 'bg-orange-50 dark:bg-orange-950/40',
    ring: 'ring-orange-200 dark:ring-orange-800/60',
    stroke: '#ea580c',
  },
  agent: {
    label: 'Agents',
    icon: 'pi pi-sparkles',
    text: 'text-fuchsia-600 dark:text-fuchsia-400',
    bg: 'bg-fuchsia-50 dark:bg-fuchsia-950/40',
    ring: 'ring-fuchsia-200 dark:ring-fuchsia-800/60',
    stroke: '#c026d3',
  },
  chatbot: {
    label: 'Chatbots',
    icon: 'pi pi-comments',
    text: 'text-rose-600 dark:text-rose-400',
    bg: 'bg-rose-50 dark:bg-rose-950/40',
    ring: 'ring-rose-200 dark:ring-rose-800/60',
    stroke: '#e11d48',
  },
  page: {
    label: 'Pages',
    icon: 'pi pi-window-maximize',
    text: 'text-sky-600 dark:text-sky-400',
    bg: 'bg-sky-50 dark:bg-sky-950/40',
    ring: 'ring-sky-200 dark:ring-sky-800/60',
    stroke: '#0284c7',
  },
  chart: {
    label: 'Charts',
    icon: 'pi pi-chart-bar',
    text: 'text-amber-600 dark:text-amber-400',
    bg: 'bg-amber-50 dark:bg-amber-950/40',
    ring: 'ring-amber-200 dark:ring-amber-800/60',
    stroke: '#d97706',
  },
  role: {
    label: 'Roles',
    icon: 'pi pi-id-card',
    text: 'text-violet-600 dark:text-violet-400',
    bg: 'bg-violet-50 dark:bg-violet-950/40',
    ring: 'ring-violet-200 dark:ring-violet-800/60',
    stroke: '#7c3aed',
  },
  user: {
    label: 'Users',
    icon: 'pi pi-user',
    text: 'text-emerald-600 dark:text-emerald-400',
    bg: 'bg-emerald-50 dark:bg-emerald-950/40',
    ring: 'ring-emerald-200 dark:ring-emerald-800/60',
    stroke: '#059669',
  },
  // Not a resource step kind; used for Group nodes in the relationship graph.
  group: {
    label: 'Groups',
    icon: 'pi pi-folder',
    text: 'text-slate-700 dark:text-slate-200',
    bg: 'bg-slate-100 dark:bg-slate-800',
    ring: 'ring-slate-300 dark:ring-slate-600',
    stroke: '#475569',
  },
}

const FALLBACK = {
  label: 'Resources',
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
