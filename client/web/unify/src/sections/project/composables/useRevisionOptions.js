import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

// Select options for a work item's `revisionID` field (config/eventForm.js
// REVISION_FIELD, consumed by NewEventDialog/EventDetailDialog/
// BacklogItemDialog) — the project's revision chain, 1-based and oldest
// first, plus an explicit "Unassigned" entry (value null) so the field can
// clear an assignment as well as set one. Work items carry a revision the way
// an issue carries a milestone (project.intent.md "Dashboards").
//
// Mirrors useRevisionLabel.js's chain-root resolution (route.params.projectId,
// falling back through rootProjectID) so both read the exact same chain
// regardless of whether the mounting route is the dashboard (root) or the
// wizard (an open revision) — READS stores/projects.js only, never mutates it.
export function useRevisionOptions() {
  const { t } = useI18n()
  const route = useRoute()
  const store = useProjectsStore()

  const chainRootId = computed(() => {
    const p = store.findById(route.params.projectId)
    return p?.rootProjectID || route.params.projectId
  })

  // Defensive refetch: RevisionSwitcher (mounted once by DashboardLayout.vue /
  // Wizard.vue) already populates revisionsFor(chainRoot) in the common case,
  // but an editor opened before that resolves would otherwise show a
  // stale/empty chain. Cheap and idempotent, so call freely (e.g. whenever a
  // dialog opens) — mirrors RevisionSwitcher's own watch-driven fetch.
  function ensureLoaded() {
    const id = chainRootId.value
    if (id) {
      store.listRevisions(id).catch(err => console.error('Failed to load project revisions', err))
    }
  }

  const revisionOptions = computed(() => {
    const chain = store.revisionsFor(chainRootId.value)
    const opts = chain
      .slice()
      .sort((a, b) => a.revision - b.revision)
      .map(rev => ({
        label: t('project.dashboard.version', { number: rev.revision + 1 }),
        value: rev.projectID,
      }))
    return [{ label: t('project.dashboard.revision.unassigned'), value: null }, ...opts]
  })

  return { revisionOptions, ensureLoaded }
}
