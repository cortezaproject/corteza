<template>
  <div class="flex items-center gap-1">
    <IconField class="flex-1">
      <InputIcon class="pi pi-link" />
      <InputText
        :model-value="label"
        class="w-full"
        :size="size === 'small' ? 'small' : undefined"
        @update:model-value="emit('update:label', $event)"
        @focus="emit('click')"
        @click="emit('click')"
      />
    </IconField>
    <Button
      icon="pi pi-times"
      text
      rounded
      size="small"
      severity="secondary"
      @click="emit('clear')"
    />
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: 'Reference',
  },
  size: {
    type: String,
    default: 'normal',
  },
})

const emit = defineEmits(['click', 'clear', 'update:label'])
</script>

<style scoped>
:deep(.p-inputicon) {
  color: var(--p-primary-color);
}
</style>
