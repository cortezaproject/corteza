<template>
  <template v-if="project">
    <!-- Members dialog opener — the project topbar's "Members" button (see
         MembersDialog.vue, mounted below). Shared by the Wizard header and
         the dashboard topbar so both surfaces get working Members without
         each mounting its own dialog. -->
    <Button
      :label="$t('project.wizard.toolbar.members')"
      icon="pi pi-users"
      size="small"
      severity="secondary"
      outlined
      @click="membersOpen = true"
    />

    <!-- Open the project's compose namespace — only once the project has one
         (project.hasNamespace, resolved from namespaceID; see
         stores/projects.js#hasNamespace). CRouterLinkButton (a real
         RouterLink under a PrimeVue Button) mirrors the idiom already used
         for every other namespace link in this section (DataModelStep,
         Page/ModuleDetailDialog): native anchor semantics (new-tab,
         copy-link, correct href) over a click handler + router.push. -->
    <CRouterLinkButton
      v-if="project.hasNamespace"
      :to="{ name: 'namespace.view', params: { slug: project.namespaceID } }"
      :label="$t('project.viewProject')"
      icon="pi pi-external-link"
      size="small"
    />

    <MembersDialog v-model:visible="membersOpen" :project="project" />
  </template>
</template>

<script setup>
import MembersDialog from '@/sections/project/components/project/MembersDialog.vue'
import { components } from '@planetcrust/human-vue'
import { ref, watch } from 'vue'

// Shared project topbar tools: Members (opens MembersDialog, owned here so
// every mounting surface gets working Members without duplicating the
// dialog) + "View project" (the compose namespace link, when the project has
// one). Mounted by the Wizard header and DashboardLayout, each inside their
// own Teleport to="#topbar-title" block beside RevisionSwitcher
// (components/project/RevisionSwitcher.vue) — same extraction shape.
const { CRouterLinkButton } = components

const props = defineProps({
  project: { type: Object, default: null },
  // Opt-in one-shot trigger to open the Members dialog: Wizard.vue derives
  // this from the `?new=1` just-created-project deep link (ProjectList.vue
  // sets the query flag; Wizard.vue strips it via router.replace once
  // consumed). Any transition to a truthy value (re)opens the dialog; the
  // dashboard never passes this.
  autoOpenMembers: { type: Boolean, default: false },
})

const membersOpen = ref(false)
watch(
  () => props.autoOpenMembers,
  val => {
    if (val) membersOpen.value = true
  },
  { immediate: true },
)
</script>
