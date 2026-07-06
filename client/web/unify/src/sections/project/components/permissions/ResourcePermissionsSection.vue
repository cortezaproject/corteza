<template>
  <!-- Per-resource access, embedded in a resource detail dialog. Roles are the
       columns; the scoped matrix shows just this one resource's toggles. Mirrors
       the Role/User dialogs' inline permissions section. -->
  <CFormGroup :label="$t('project.permissions.resourceLabel')">
    <p class="text-sm text-muted-color mb-2">{{ $t('project.permissions.resourceHint') }}</p>
    <ProjectPermissionMatrix
      v-if="resourceId && roles.length"
      :project="project"
      :roles="roles"
      :scope="{ kind, id: resourceId }"
    />
    <p v-else class="text-sm text-muted-color italic">
      {{ $t('project.permissions.resourceNoRoles') }}
    </p>
  </CFormGroup>
</template>

<script setup>
import ProjectPermissionMatrix from '@/sections/project/components/permissions/ProjectPermissionMatrix.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, watch } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  // Matrix KINDS kind: 'module' | 'page' | 'agent' | 'chatbot' | 'automation'.
  kind: { type: String, required: true },
  resourceId: { type: String, default: null },
})

const store = useProjectsStore()

// The project's access roles = the column set. Resource detail dialogs don't
// load roles on their own, so pull them in when the section mounts / the project
// changes — without this both the columns and the traced cell states are empty.
const roles = computed(() => store.rolesFor(props.project?.id) || [])

watch(
  () => props.project?.id,
  id => {
    if (id) store.loadRoles(id)
  },
  { immediate: true },
)
</script>
