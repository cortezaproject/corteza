<template>
  <div class="h-full flex">
    <!-- Group selector -->
    <aside class="w-52 shrink-0 h-full border-r border-surface bg-surface flex flex-col">
      <div class="p-3 border-b border-surface">
        <span class="block text-xs uppercase tracking-wider text-muted-color">Groups</span>
      </div>
      <div class="flex-1 overflow-auto px-2 py-2">
        <button
          type="button"
          class="w-full flex items-center gap-2 px-2 py-1.5 rounded-md text-left transition-colors"
          :class="selectedId === '__all__' ? 'bg-primary/10 text-primary' : 'hover:bg-emphasis text-color'"
          @click="selectedId = '__all__'"
        >
          <i class="pi pi-th-large" />
          <span class="text-sm font-medium">All groups</span>
          <span class="ml-auto text-xs text-muted-color">{{ groups.length }}</span>
        </button>

        <div v-if="groups.length" class="border-t border-surface my-2" />

        <button
          v-for="group in groups"
          :key="group.id"
          type="button"
          class="w-full flex flex-col items-start gap-0.5 px-2 py-2 rounded-md text-left transition-colors"
          :class="selectedId === group.id ? 'bg-primary/10 text-primary' : 'hover:bg-emphasis text-color'"
          @click="selectedId = group.id"
        >
          <div class="flex items-center gap-2 w-full">
            <i class="pi pi-folder text-sm" />
            <span class="text-sm font-medium truncate flex-1">{{ group.name || 'Untitled' }}</span>
            <span class="text-xs text-muted-color">{{ (group.resourceIds || []).length }}</span>
          </div>
          <span v-if="roleLabel(group)" class="text-xs text-muted-color truncate w-full">
            {{ roleLabel(group) }}
          </span>
        </button>

        <div v-if="!groups.length" class="px-2 py-1.5 text-sm text-muted-color italic">
          No groups yet.
        </div>
      </div>
    </aside>

    <!-- Group editor -->
    <div class="flex-1 min-w-0 overflow-auto p-6">
      <div>
        <template v-if="selectedId === '__all__'">
          <div class="flex flex-col gap-3">
            <GroupCard
              v-for="group in groups"
              :key="group.id"
              :group="group"
              :project="project"
              :disabled="disabled"
              @update="patch => store.updateGroup(project.id, group.id, patch)"
              @remove="store.removeGroup(project.id, group.id)"
            />
            <NewGroupForm v-if="!disabled" @create="onCreate" />
          </div>
        </template>

        <template v-else-if="focusedGroup">
          <GroupCard
            :group="focusedGroup"
            :project="project"
            :disabled="disabled"
            @update="patch => store.updateGroup(project.id, focusedGroup.id, patch)"
            @remove="onRemoveFocused"
          />
        </template>

        <div v-else class="text-sm text-muted-color italic">Group not found.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import GroupCard from '@/sections/project/components/group/GroupCard.vue'
import NewGroupForm from '@/sections/project/components/group/NewGroupForm.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, watch } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  // Selected group id ('__all__' for the overview). Two-way bound so the
  // resource graph panel can follow the selection.
  modelValue: { type: String, default: '__all__' },
  // Locks editing (e.g. once the step is submitted/approved at its gate).
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const store = useProjectsStore()

const groups = computed(() => props.project?.groups || [])
const selectedId = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const focusedGroup = computed(() =>
  selectedId.value === '__all__' ? null : groups.value.find(g => g.id === selectedId.value),
)

// Fall back to "All" if the selected group is deleted.
watch(groups, list => {
  if (selectedId.value === '__all__') return
  if (!list.some(g => g.id === selectedId.value)) selectedId.value = '__all__'
})

const roleName = id => (props.project?.resources || []).find(r => r.id === id)?.name || id

const roleLabel = group => {
  const ids = group.roleIds || []
  if (!ids.length) return ''
  if (ids.length === 1) return roleName(ids[0])
  return `${roleName(ids[0])} +${ids.length - 1}`
}

const onCreate = group => {
  const id = store.addGroup(props.project.id, group)
  if (id) selectedId.value = id
}

const onRemoveFocused = () => {
  if (!focusedGroup.value) return
  store.removeGroup(props.project.id, focusedGroup.value.id)
  selectedId.value = '__all__'
}
</script>
