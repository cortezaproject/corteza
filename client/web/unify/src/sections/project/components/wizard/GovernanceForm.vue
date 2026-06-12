<template>
  <div class="flex flex-col gap-6">
    <div
      v-for="(section, i) in schema"
      :key="section.title || i"
    >
      <div class="grid grid-cols-1 gap-x-6 gap-y-4">
        <CFormGroup
          v-for="field in section.fields"
          :key="field.key"
          :label="field.label"
          :description="field.description"
        >
          <InputText
            v-if="field.type === 'text'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            :placeholder="field.placeholder"
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <Textarea
            v-else-if="field.type === 'textarea'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            :placeholder="field.placeholder"
            rows="3"
            auto-resize
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <InputNumber
            v-else-if="field.type === 'number'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <DatePicker
            v-else-if="field.type === 'datetime'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            show-time
            hour-format="24"
            show-button-bar
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <Select
            v-else-if="field.type === 'select'"
            :model-value="modelValue[field.key]"
            :options="field.options"
            :disabled="disabled"
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <MultiSelect
            v-else-if="field.type === 'multiselect'"
            :model-value="modelValue[field.key]"
            :options="field.options"
            :disabled="disabled"
            display="chip"
            filter
            fluid
            @update:model-value="update(field.key, $event)"
          />
        </CFormGroup>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  schema: { type: Array, required: true },
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

function update(key, value) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>
