<template>
  <Dialog
    :visible="visible"
    modal
    :header="$t('system.user-groups.delete-blocked.title', { name })"
    class="w-full max-w-lg"
    data-testid="user-group-delete-blocked"
    @update:visible="$emit('update:visible', $event)"
  >
    <div class="flex flex-col gap-3">
      <p>
        {{
          $t('system.user-groups.delete-blocked.summary', {
            members: $t('system.user-groups.delete-blocked.members', members),
            childGroups: $t('system.user-groups.delete-blocked.child-groups', childGroups),
          })
        }}
      </p>
      <p class="text-sm text-muted-color">
        {{ $t('system.user-groups.delete-blocked.how-to') }}
      </p>
    </div>

    <template #footer>
      <Button :label="$t('general.label.close')" @click="$emit('update:visible', false)" />
    </template>
  </Dialog>
</template>

<script setup>
defineProps({
  visible: { type: Boolean, default: false },
  name: { type: String, default: '' },
  members: { type: Number, default: 0 },
  childGroups: { type: Number, default: 0 },
})

defineEmits(['update:visible'])
</script>
