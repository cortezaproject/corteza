<template>
  <div class="flex items-center gap-1">
    <ColorPicker
      v-if="modelValue"
      :model-value="modelValue"
      format="hex"
      @update:model-value="$emit('update:modelValue', $event)"
    />
    <button
      v-else
      type="button"
      class="w-6 h-6 rounded border-2 border-dashed border-surface-300 bg-transparent cursor-pointer"
      :title="$t('general.label.none')"
      @click="$emit('update:modelValue', 'cccccc')"
    />
    <Button
      v-if="modelValue"
      icon="pi pi-times"
      text
      rounded
      size="small"
      severity="secondary"
      class="!w-5 !h-5 !p-0"
      @click="$emit('update:modelValue', '')"
    />
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

defineEmits(['update:modelValue'])
</script>
