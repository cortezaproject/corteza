<template>
  <article class="bg-surface rounded-xl border border-surface overflow-hidden">
    <header class="px-4 py-2 flex items-center gap-2 border-b border-surface">
      <i class="pi pi-folder text-primary" />
      <InputText
        :model-value="group.name"
        :disabled="disabled"
        class="!border-0 !shadow-none focus:!ring-0 !bg-transparent !p-0 font-medium text-sm flex-1"
        @update:model-value="v => emit('update', { name: v })"
      />
      <Chip
        :label="`${(group.roleIds || []).length} roles`"
        icon="pi pi-id-card"
        class="!text-xs !py-0"
      />
      <Chip
        :label="`${(group.resourceIds || []).length} linked`"
        icon="pi pi-link"
        class="!text-xs !py-0"
      />
      <Button
        v-if="!disabled"
        icon="pi pi-times"
        severity="secondary"
        text
        rounded
        size="small"
        @click="emit('remove')"
      />
    </header>

    <Textarea
      :model-value="group.description"
      :disabled="disabled"
      rows="2"
      placeholder="Description (optional)"
      class="!border-0 !shadow-none focus:!ring-0 !bg-transparent w-full px-4 pt-3 !text-sm text-color resize-none"
      auto-resize
      @update:model-value="v => emit('update', { description: v })"
    />

    <GroupKindRow
      :project="project"
      kind="role"
      :linked-ids="group.roleIds || []"
      :disabled="disabled"
      @link="linkRole"
      @unlink="unlinkRole"
    />

    <template v-for="kind in RESOURCE_KINDS" :key="kind">
      <GroupKindRow
        :project="project"
        :kind="kind"
        :linked-ids="resourceIdsOfKind(kind)"
        :disabled="disabled"
        @link="linkResource"
        @unlink="unlinkResource"
      />
    </template>

    <div class="pb-3" />
  </article>
</template>

<script setup>
import GroupKindRow from '@/sections/project-poc/components/group/GroupKindRow.vue'
import { computed } from 'vue'

// Resource kinds a group can link (roles are handled in their own row above).
const RESOURCE_KINDS = ['module', 'page', 'chart', 'automation', 'agent', 'chatbot']

const props = defineProps({
  group: { type: Object, required: true },
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update', 'remove'])

const resourceKindById = computed(() => {
  const map = new Map()
  for (const r of props.project?.resources || []) map.set(r.id, r.kind)
  return map
})

const resourceIdsOfKind = kind =>
  (props.group.resourceIds || []).filter(id => resourceKindById.value.get(id) === kind)

const linkRole = id => {
  if ((props.group.roleIds || []).includes(id)) return
  emit('update', { roleIds: [...(props.group.roleIds || []), id] })
}

const unlinkRole = id => {
  emit('update', { roleIds: (props.group.roleIds || []).filter(x => x !== id) })
}

const linkResource = id => {
  if ((props.group.resourceIds || []).includes(id)) return
  emit('update', { resourceIds: [...(props.group.resourceIds || []), id] })
}

const unlinkResource = id => {
  emit('update', { resourceIds: (props.group.resourceIds || []).filter(x => x !== id) })
}
</script>
