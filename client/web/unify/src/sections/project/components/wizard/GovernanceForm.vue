<template>
  <div class="flex flex-col gap-6">
    <div
      v-for="(section, i) in schema"
      :key="section.titleKey || i"
    >
      <div class="grid grid-cols-1 gap-x-6 gap-y-4" :class="{ 'sm:grid-cols-2': columns === 2 }">
        <CFormGroup
          v-for="field in section.fields"
          :key="field.key"
          :label="$t(field.labelKey)"
          :required="!!field.required"
          :description="field.descriptionKey ? $t(field.descriptionKey) : undefined"
          :class="{ 'sm:col-span-2': columns === 2 && (field.full || field.type === 'textarea') }"
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
          <DatePicker
            v-else-if="field.type === 'date'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
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
  // 1 = single column (governance forms); 2 = two-column grid where textareas
  // (and fields flagged `full`) span both columns.
  columns: { type: Number, default: 1 },
})
const emit = defineEmits(['update:modelValue'])

function update(key, value) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>
