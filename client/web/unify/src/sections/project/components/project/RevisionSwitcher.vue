<template>
  <template v-if="project">
    <!-- Revision switcher trigger — shared by the Wizard header and the
         dashboard topbar (DashboardLayout.vue mounts this exact component),
         replacing what used to be two divergent bits of chrome: the wizard's
         own version Tag + New-revision menu, and the dashboard's static
         version Tag + "view overview" button. role="button"/tabindex mirror
         the icon-trigger idiom used for other popup Menus in this app (e.g.
         TabsBlock.vue's tab menu). -->
    <span
      role="button"
      tabindex="0"
      class="inline-flex items-center gap-1 cursor-pointer rounded px-1 hover:bg-emphasis"
      :aria-label="$t('project.wizard.revisionSwitcher.trigger')"
      v-tooltip.bottom="$t('project.wizard.revisionSwitcher.trigger')"
      @click="toggleMenu"
      @keydown.enter.stop.prevent="toggleMenu"
      @keydown.space.stop.prevent="toggleMenu"
    >
      <Tag :value="triggerLabel" :severity="triggerSeverity" class="!text-xs" />
      <i class="pi pi-chevron-down text-xs text-muted-color" />
    </span>

    <!-- Popup menu: the chain-wide Dashboard first, then the revision chain
         (see stores/projects.js listRevisions/revisionsFor), each entry
         showing its 1-based version + lifecycle status, then a separator and
         the "New revision from this one" action. Which entry is marked
         current is derived from the active route (see onDashboard/menuItems
         below), not from a prop, so the exact same popup reads correctly
         whichever surface mounted it. -->
    <Menu ref="menuRef" :model="menuItems" popup>
      <template #item="{ item, props: itemProps }">
        <a
          v-if="item.kind === 'dashboard'"
          v-ripple
          v-bind="itemProps.action"
          class="flex items-center gap-3"
          :class="{ 'font-semibold': item.current }"
        >
          <i class="pi pi-check text-primary text-xs" :class="{ invisible: !item.current }" />
          <i class="pi pi-gauge text-xs text-muted-color" />
          <span class="flex-1">{{ item.label }}</span>
        </a>
        <a
          v-else-if="item.kind === 'revision'"
          v-ripple
          v-bind="itemProps.action"
          class="flex items-center gap-3"
          :class="{ 'font-semibold': item.current }"
        >
          <i class="pi pi-check text-primary text-xs" :class="{ invisible: !item.current }" />
          <span class="flex-1">{{ item.label }}</span>
          <Tag
            :value="$t(`project.status.${item.status}`)"
            :severity="revisionStatusSeverity(item.status)"
            class="!text-xs"
          />
        </a>
        <a
          v-else-if="item.kind === 'new-revision'"
          v-ripple
          v-bind="itemProps.action"
          v-tooltip.bottom="item.tooltip"
          class="flex items-center gap-2"
        >
          <i class="pi pi-plus text-xs" />
          <span>{{ item.label }}</span>
        </a>
      </template>
    </Menu>
  </template>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

// Shared revision-chain switcher. `project` is whichever project row the
// mounting view already has loaded — the open revision in the Wizard, or the
// chain member the dashboard route resolved (DashboardLayout.vue's own
// `project` computed) — used to seed the trigger label/severity and to
// resolve the chain (stores/projects.js revisionsFor, keyed by chain root).
const props = defineProps({
  project: { type: Object, default: null },
})

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const store = useProjectsStore()
const $toast = inject('$toast')

// The dashboard is chain-wide, so nothing about "current" is tied to
// `project`'s own id there — it's current purely by being on the
// project.overview route (or one of its child routes: events/category/etc).
const onDashboard = computed(
  () => typeof route.name === 'string' && route.name.startsWith('project.overview'),
)

// A live (published) project has a dashboard to switch to; drafts are
// wizard-only. Also drives the trigger's severity (mirrors the dashboard
// topbar's old version Tag) and the version label's wording.
const isLive = computed(() => ['active', 'published'].includes(props.project?.status))

// User-facing versions are 1-based (the original live project is v1), so we
// display the backend revision + 1. An unpublished project is flagged as a
// draft (e.g. "v1 draft").
const versionLabel = computed(() =>
  t(isLive.value ? 'project.dashboard.version' : 'project.dashboard.versionDraft', {
    number: (props.project?.revision ?? 0) + 1,
  }),
)

// The trigger always shows whichever entry is current: the chain-wide
// Dashboard label while on a dashboard route, the open revision's version
// label otherwise.
const triggerLabel = computed(() =>
  onDashboard.value ? t('project.dashboard.title') : versionLabel.value,
)
const triggerSeverity = computed(() => (!onDashboard.value && !isLive.value ? 'warn' : 'secondary'))

const menuRef = ref()
function toggleMenu(event) {
  menuRef.value?.toggle(event)
}

const revisionChain = computed(() => {
  if (!props.project) return []
  const chain = store.revisionsFor(props.project.projectID)
  return (chain.length ? chain : [props.project]).slice().sort((a, b) => a.revision - b.revision)
})

// Load the chain whenever the mounting view's project changes — both surfaces
// need it, so the switcher owns this fetch itself rather than each parent
// duplicating it. A failure just leaves the switcher showing only the
// passed-in project (see revisionChain's fallback above); it isn't critical
// enough to toast.
watch(
  () => props.project?.projectID,
  id => {
    if (!id) return
    store.listRevisions(id).catch(err => console.error('Failed to load project revisions', err))
  },
  { immediate: true },
)

// Severities mirror ProjectList.vue / ProjectSidebar.vue's project-status Tag
// mapping (active/published live, draft in review, suspended flagged).
const REVISION_STATUS_SEVERITY = {
  active: 'success',
  published: 'success',
  draft: 'info',
  suspended: 'warn',
  archived: 'secondary',
}
const revisionStatusSeverity = status => REVISION_STATUS_SEVERITY[status] || 'secondary'

// Selecting another revision navigates to its own wizard; selecting the
// already-open one (from the wizard) is a harmless no-op.
function goToRevision(projectId) {
  if (!onDashboard.value && String(projectId) === String(props.project?.projectID)) return
  router.push({ name: 'project.wizard', params: { projectId } })
}

// Selecting the Dashboard entry always targets the chain ROOT, so the
// dashboard's URL is the same canonical one regardless of which chain member
// happened to be open when the switcher was used.
function goToDashboard() {
  if (onDashboard.value) return
  const rootId = props.project?.rootProjectID || props.project?.projectID
  router.push({ name: 'project.overview', params: { projectId: rootId } })
}

// AGREED BEHAVIOUR: prevent rather than fail. The "New revision from this
// one" entry is disabled — with a tooltip explaining why — whenever the
// backend would reject createRevision: no revision to branch from is
// currently active, or the chain already has a draft (the backend allows
// only one draft per chain). onCreateRevision's catch below is only a
// backstop for the race where someone else created a draft first.
//
// "This one" means different things on the two surfaces: in the Wizard it's
// literally the open revision (you branch off what you're looking at); the
// dashboard has no single open revision (it's chain-wide), so it branches off
// whichever revision in the chain is currently active.
const branchSource = computed(() =>
  onDashboard.value ? revisionChain.value.find(r => r.status === 'active') || null : props.project,
)
const currentNotActive = computed(() => branchSource.value?.status !== 'active')
const chainHasDraft = computed(() => revisionChain.value.some(r => r.status === 'draft'))
const newRevisionDisabled = computed(() => currentNotActive.value || chainHasDraft.value)
const newRevisionDisabledReason = computed(() => {
  if (currentNotActive.value)
    return t('project.wizard.revisionSwitcher.newRevisionDisabledNotActive')
  if (chainHasDraft.value)
    return t('project.wizard.revisionSwitcher.newRevisionDisabledDraftExists')
  return ''
})

const menuItems = computed(() => {
  const items = [
    {
      kind: 'dashboard',
      label: t('project.dashboard.title'),
      current: onDashboard.value,
      command: goToDashboard,
    },
    ...revisionChain.value.map(rev => ({
      kind: 'revision',
      label: t('project.dashboard.version', { number: rev.revision + 1 }),
      status: rev.status,
      current: !onDashboard.value && rev.projectID === props.project?.projectID,
      command: () => goToRevision(rev.projectID),
    })),
  ]
  items.push({ separator: true })
  items.push({
    kind: 'new-revision',
    label: t('project.wizard.revisionSwitcher.newRevision'),
    disabled: newRevisionDisabled.value,
    tooltip: newRevisionDisabled.value ? newRevisionDisabledReason.value : '',
    command: onCreateRevision,
  })
  return items
})

const creatingRevision = ref(false)
async function onCreateRevision() {
  if (creatingRevision.value || newRevisionDisabled.value || !branchSource.value) return
  creatingRevision.value = true
  try {
    const draft = await store.createRevision(branchSource.value.projectID)
    router.push({ name: 'project.wizard', params: { projectId: draft.projectID } })
  } catch (err) {
    $toast.toastErrorHandler(t('project.wizard.revisionSwitcher.toastCreateFailed'))(err)
  } finally {
    creatingRevision.value = false
  }
}
</script>
