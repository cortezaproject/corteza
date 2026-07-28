import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

// Resolves a work item's `revisionID` (events and backlog items both carry
// one — see stores/events.js#load / stores/backlogItems.js#load, and
// server/system/types/project_incident.gen.go et al) to display info for the
// dashboard's chain-wide list views (BacklogView, CategoryView), which show
// items filed against ANY revision in the chain side by side and so need to
// say which one each belongs to (project.intent.md "Dashboards"). Items with
// no revision are legitimately unassigned — `unassigned: true` lets callers
// render that clearly rather than leaving the cell blank.
//
// Reads the same revisionsFor(chain root) cache the shared RevisionSwitcher
// populates (components/project/RevisionSwitcher.vue, mounted once by
// DashboardLayout.vue) — reactive, so this fills in once that fetch resolves
// without either view needing to trigger its own.
export function useRevisionLabel() {
  const { t } = useI18n()
  const route = useRoute()
  const store = useProjectsStore()

  const chainRootId = computed(() => {
    const p = store.findById(route.params.projectId)
    return p?.rootProjectID || route.params.projectId
  })
  const chain = computed(() => store.revisionsFor(chainRootId.value))

  // Same severities the revision switcher's menu items use, so a status tint
  // reads consistently wherever a revision shows up.
  const REVISION_STATUS_SEVERITY = {
    active: 'success',
    published: 'success',
    draft: 'info',
    suspended: 'warn',
    archived: 'secondary',
  }

  function revisionInfo(revisionId) {
    const id = String(revisionId || '')
    if (!id || id === '0') {
      return { unassigned: true, label: t('project.dashboard.revision.unassigned') }
    }
    const rev = chain.value.find(r => r.projectID === id)
    if (!rev) {
      // Not in the loaded chain (still loading, or a stale/removed
      // revision) — fall back to the raw reference rather than hiding it.
      return { unassigned: false, label: `#${id}`, severity: 'secondary' }
    }
    return {
      unassigned: false,
      label: t('project.dashboard.version', { number: rev.revision + 1 }),
      severity: REVISION_STATUS_SEVERITY[rev.status] || 'secondary',
    }
  }

  return { revisionInfo }
}
