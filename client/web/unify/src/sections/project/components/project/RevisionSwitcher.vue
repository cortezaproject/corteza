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
      :role="hasMenu ? 'button' : undefined"
      :tabindex="hasMenu ? 0 : undefined"
      class="inline-flex shrink-0 items-center gap-1 rounded px-1"
      :class="hasMenu ? 'cursor-pointer hover:bg-emphasis' : ''"
      :aria-label="hasMenu ? $t('project.wizard.revisionSwitcher.trigger') : undefined"
      v-tooltip.bottom="hasMenu ? $t('project.wizard.revisionSwitcher.trigger') : undefined"
      @click="hasMenu && toggleMenu($event)"
      @keydown.enter.stop.prevent="hasMenu && toggleMenu($event)"
      @keydown.space.stop.prevent="hasMenu && toggleMenu($event)"
    >
      <Tag :value="triggerLabel" severity="secondary" class="!text-xs" />
      <!-- Where the open revision stands, stated on the control that names it —
           and by the same derivation its dropdown entry uses, so the trigger and
           the row it corresponds to always read identically. Suppressed on the
           dashboard, where the trigger names the chain-wide Dashboard rather
           than a revision. -->
      <StatusChip
        v-if="!onDashboard && project"
        :status="currentTag.status"
        :label="currentTag.label"
        small
      />
      <!-- No chevron when there is nowhere to go: a lone revision that can't be
           branched and has no dashboard leaves the menu with nothing in it but
           the entry you are already on, so the whole control degrades to a
           label (see hasMenu). -->
      <i v-if="hasMenu" class="pi pi-chevron-down text-xs text-muted-color" />
    </span>

    <!-- Popup menu: the chain-wide Dashboard first, then the revision chain
         (see stores/projects.js listRevisions/revisionsFor), each entry
         showing its 1-based version + lifecycle status, then a separator and
         the "New revision from this one" action. Which entry is marked
         current is derived from the active route (see onDashboard/menuItems
         below), not from a prop, so the exact same popup reads correctly
         whichever surface mounted it. -->
    <Menu v-if="hasMenu" ref="menuRef" :model="menuItems" popup>
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
          <StatusChip :status="item.tag.status" :label="item.tag.label" small />
        </a>
        <a
          v-else-if="item.kind === 'new-revision'"
          v-ripple
          v-bind="itemProps.action"
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
import StatusChip from '@/sections/project/components/project/StatusChip.vue'
import { chainHasPublished, projectStatusTag } from '@/sections/project/config/publishState'
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

// User-facing versions are 1-based (the original live project is v1), so we
// display the backend revision + 1. Just the version, exactly as the dropdown
// entries name it — the trigger used to append "draft" and turn amber for an
// unpublished revision, which now says twice, worse, what the project status
// tag beside it states properly (see Wizard.vue's topbar).
const versionLabel = computed(() =>
  t('project.dashboard.version', { number: (props.project?.revision ?? 0) + 1 }),
)

// The trigger always shows whichever entry is current: the chain-wide
// Dashboard label while on a dashboard route, the open revision's version
// label otherwise.
const triggerLabel = computed(() =>
  onDashboard.value ? t('project.dashboard.title') : versionLabel.value,
)

// One derivation for every revision this component shows, trigger and dropdown
// alike (see config/publishState.js). The approval standing is read off each
// revision's own row, so a sibling revision reports where ITS review actually
// stands — including one submitted by somebody else, in another session, which
// the previous session-local reading could never see.
const statusTagFor = rev =>
  projectStatusTag({
    status: rev?.status,
    publishStatus: rev ? store.publishApprovalStatus(rev.projectID) : '',
  })
const currentTag = computed(() => statusTagFor(props.project))

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

// Tone/icon/label per status all live in StatusChip now — including
// `deprecated`, which publish stamps on the outgoing revision and is therefore
// the single most common status to meet in a chain of any age.

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

// "New revision from this one" is offered only when it would actually work —
// the entry is absent, not disabled-with-a-reason, when the backend would
// reject createRevision: the revision to branch from is not active, or the
// chain already has a draft (the backend allows only one draft per chain).
// onCreateRevision's catch below is still the backstop for the race where
// someone else created a draft first.
//
// "This one" means different things on the two surfaces: in the Wizard it's
// literally the open revision (you branch off what you're looking at); the
// dashboard has no single open revision (it's chain-wide), so it branches off
// whichever revision in the chain is currently active.
const branchSource = computed(() =>
  onDashboard.value ? revisionChain.value.find(r => r.status === 'active') || null : props.project,
)
// canReviseProject is the backend's own answer (project.revise, carried on the
// payload); the status/draft conditions mirror the rest of what CreateRevision
// enforces. Asking the server rather than assuming means a user without the
// permission is never offered a branch that would come back as a 403.
const canCreateRevision = computed(
  () =>
    branchSource.value?.status === 'active' &&
    branchSource.value?.canReviseProject &&
    !revisionChain.value.some(r => r.status === 'draft'),
)

// Whether the control is a menu at all. Everything it can offer beyond the
// entry you are already on: a dashboard to switch to, a sibling revision, or a
// new revision to branch. With none of those the popup would list only the
// current revision, so the trigger drops its chevron and stops being clickable
// (see the template) rather than opening a menu of one.
const hasMenu = computed(
  () =>
    chainHasPublished(props.project) || revisionChain.value.length > 1 || canCreateRevision.value,
)

const menuItems = computed(() => {
  const items = [
    // The dashboard only exists once the chain has published something to
    // report on (see config/publishState.js) — until then this switcher offers
    // the revisions alone.
    ...(chainHasPublished(props.project)
      ? [
          {
            kind: 'dashboard',
            label: t('project.dashboard.title'),
            current: onDashboard.value,
            command: goToDashboard,
          },
        ]
      : []),
    ...revisionChain.value.map(rev => ({
      kind: 'revision',
      label: t('project.dashboard.version', { number: rev.revision + 1 }),
      tag: statusTagFor(rev),
      current: !onDashboard.value && rev.projectID === props.project?.projectID,
      command: () => goToRevision(rev.projectID),
    })),
  ]
  if (canCreateRevision.value) {
    items.push({ separator: true })
    items.push({
      kind: 'new-revision',
      label: t('project.wizard.revisionSwitcher.newRevision'),
      command: onCreateRevision,
    })
  }
  return items
})

const creatingRevision = ref(false)
async function onCreateRevision() {
  if (creatingRevision.value || !canCreateRevision.value || !branchSource.value) return
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
