<template>
  <div class="h-full flex flex-col min-h-0">
    <!-- The step description lives in the wizard header (STEP_BLURB); no in-body repeat. -->
    <div class="flex-1 min-h-0 overflow-auto p-4">
      <ProjectPermissionMatrix :project="project" :roles="roles" :disabled="disabled" />
    </div>
  </div>
</template>

<script setup>
import ProjectPermissionMatrix from '@/sections/project/components/permissions/ProjectPermissionMatrix.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')

const roles = computed(() => store.rolesFor(props.project.id))

// Load the roles and every kind of resource the matrix charts. The matrix builds
// its tree from these store caches and evaluates effective access itself.
async function refresh(id) {
  if (!id) return
  try {
    await Promise.all([
      store.loadRoles(id),
      store.loadResources(id),
      store.loadPages(id),
      store.loadAutomations(id),
      store.loadAgents(id),
      store.loadChatbots(id),
      store.loadConnections(id),
    ])
  } catch (err) {
    $toast.toastErrorHandler(t('project.permissions.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(() => props.project.id, refresh)
</script>
