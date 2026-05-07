<template>
  <div
    v-if="!roles.length"
    class="text-muted-color p-3 border rounded-lg bg-highlight text-center text-sm"
  >
    {{ $t('system.authclients.editor.info.security.noRoles') }}
  </div>
  <div v-else class="flex flex-col border rounded-lg divide-y bg-surface">
    <div
      v-for="role in roles"
      :key="role.roleID"
      class="flex items-center justify-between p-2 px-3"
    >
      <span class="font-medium">{{ role.name || role.handle || role.roleID }}</span>
      <Button
        icon="pi pi-trash"
        severity="danger"
        text
        rounded
        size="small"
        @click="$emit('remove', role)"
      />
    </div>
  </div>
</template>

<script setup>
defineProps({
  roles: {
    type: Array,
    required: true,
  },
})

defineEmits(['remove'])
</script>
