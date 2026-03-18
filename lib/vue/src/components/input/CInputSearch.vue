<template>
  <IconField>
    <InputText
      :model-value="modelValue"
      @update:model-value="$emit('update:modelValue', $event)"
      :placeholder="placeholder"
      :size="size"
      class="w-full"
    />
    <InputIcon :class="getSearchIconClass()" @click="clearSearch" />
  </IconField>
</template>

<script setup>


const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: 'Search...',
  },
  size: {
    type: String,
    default: undefined,
  },
})

const emit = defineEmits(['update:modelValue'])

const getSearchIconClass = () => {
  return !props.modelValue
    ? 'pi pi-search text-primary'
    : 'pi pi-times cursor-pointer hover:text-primary'
}

const clearSearch = () => {
  if (props.modelValue) {
    emit('update:modelValue', '')
  }
}
</script>
