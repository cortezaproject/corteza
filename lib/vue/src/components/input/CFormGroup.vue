<template>
  <FormField v-if="name" :name="name" class="flex flex-col gap-2" v-slot="$field">
    <div
      v-if="label || $slots.label || $slots.actions || description || $slots.description"
      class="flex flex-col gap-1"
    >
      <div v-if="label || $slots.label || $slots.actions" class="flex items-center gap-2">
        <label
          v-if="label || $slots.label"
          :for="inputId || name"
          class="font-medium text-primary text-sm"
        >
          <slot name="label">{{ label }}</slot>
          <span v-if="required" class="text-red-500">*</span>
        </label>
        <slot name="actions" />
      </div>
      <small v-if="description || $slots.description" class="text-muted-color">
        <slot name="description">{{ description }}</slot>
      </small>
    </div>
    <slot />
    <Message v-if="$field?.invalid" severity="error" size="small" variant="simple">
      {{ $field.error?.message }}
    </Message>
  </FormField>

  <div v-else class="flex flex-col gap-2">
    <div
      v-if="label || $slots.label || $slots.actions || description || $slots.description"
      class="flex flex-col gap-1"
    >
      <div v-if="label || $slots.label || $slots.actions" class="flex items-center gap-2">
        <label
          v-if="label || $slots.label"
          :for="inputId || undefined"
          class="font-medium text-primary text-sm"
        >
          <slot name="label">{{ label }}</slot>
          <span v-if="required" class="text-red-500">*</span>
        </label>
        <slot name="actions" />
      </div>
      <small v-if="description || $slots.description" class="text-muted-color">
        <slot name="description">{{ description }}</slot>
      </small>
    </div>
    <slot />
  </div>
</template>

<script setup>
defineProps({
  name: { type: String, default: '' },
  label: { type: String, default: '' },
  description: { type: String, default: '' },
  required: { type: Boolean, default: false },
  inputId: { type: String, default: '' },
})
</script>
