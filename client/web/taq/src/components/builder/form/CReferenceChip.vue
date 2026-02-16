<template>
  <div class="flex items-center gap-2">
    <Tag
      :value="label"
      severity="primary"
      class="text-sm flex-1 cursor-pointer py-2"
      icon="pi pi-link"
      @click="emit('click')"
    />
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
})

const emit = defineEmits(['click', 'clear'])
</script>
