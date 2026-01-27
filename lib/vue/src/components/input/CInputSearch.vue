<template>
  <InputGroup>
    <IconField>
      <InputText
        :model-value="modelValue"
        @update:model-value="$emit('update:modelValue', $event)"
        :placeholder="placeholder"
        class="w-full"
      />
      <InputIcon
        v-show="submittable ? modelValue : true"
        :class="getSearchIconClass()"
        @click="clearSearch"
      />
    </IconField>

    <InputGroupAddon v-if="submittable">
      <Button
        icon="pi pi-search"
        severity="secondary"
        variant="text"
        class="text-primary w-14"
        @click="$emit('search')"
      />
    </InputGroupAddon>
  </InputGroup>
</template>

<script setup>
import IconField from 'primevue/iconfield'
import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: 'Search...',
  },
  submittable: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'search'])

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
