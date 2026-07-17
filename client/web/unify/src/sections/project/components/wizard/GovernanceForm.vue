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
            :placeholder="fieldPlaceholder(field)"
            :invalid="submitted && !!errors[field.key]"
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <Textarea
            v-else-if="field.type === 'textarea'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            :placeholder="fieldPlaceholder(field)"
            :invalid="submitted && !!errors[field.key]"
            rows="3"
            auto-resize
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <InputNumber
            v-else-if="field.type === 'number'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            :invalid="submitted && !!errors[field.key]"
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <DatePicker
            v-else-if="field.type === 'datetime'"
            :model-value="modelValue[field.key]"
            :disabled="disabled"
            :invalid="submitted && !!errors[field.key]"
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
            :invalid="submitted && !!errors[field.key]"
            show-button-bar
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <Select
            v-else-if="field.type === 'select'"
            :model-value="modelValue[field.key]"
            :options="field.options"
            :option-label="field.optionLabel"
            :option-value="field.optionValue"
            :filter="!!field.filter"
            :disabled="disabled"
            :invalid="submitted && !!errors[field.key]"
            fluid
            @update:model-value="update(field.key, $event)"
          >
            <!-- Risk keeps its list identity — the 5-pip meter — with the
                 label alongside (a bare meter is unreadable in a select);
                 every other badge-hinted field renders its pill. -->
            <template v-if="field.badge" #option="{ option }">
              <span v-if="field.badge === 'risk'" class="flex items-center gap-2">
                <RiskPips :level="option" />
                <span class="text-sm">{{ option }}</span>
              </span>
              <span v-else class="flex items-center">
                <EventBadge :value="option" :variant="field.badge" size="md" />
              </span>
            </template>
            <template v-if="field.badge" #value="{ value, placeholder }">
              <!-- Blank selection falls back to the placeholder text instead
                   of rendering an empty badge pill (possible on older records
                   whose optional severity/risk was never set). -->
              <span v-if="value && field.badge === 'risk'" class="flex items-center gap-2 h-full">
                <RiskPips :level="value" />
                <span class="text-sm">{{ value }}</span>
              </span>
              <span v-else-if="value" class="flex items-center h-full">
                <EventBadge :value="value" :variant="field.badge" size="md" />
              </span>
              <span v-else class="text-muted-color">{{ placeholder }}</span>
            </template>
          </Select>
          <MultiSelect
            v-else-if="field.type === 'multiselect'"
            :model-value="modelValue[field.key]"
            :options="field.options"
            :disabled="disabled"
            :invalid="submitted && !!errors[field.key]"
            display="chip"
            filter
            fluid
            @update:model-value="update(field.key, $event)"
          />
          <ValidationMessage :message="submitted ? errors[field.key] : ''" />
        </CFormGroup>
      </div>
    </div>
  </div>
</template>

<script setup>
import EventBadge from '@/sections/project/components/dashboard/EventBadge.vue'
import RiskPips from '@/sections/project/components/dashboard/RiskPips.vue'
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  schema: { type: Array, required: true },
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
  // 1 = single column (governance forms); 2 = two-column grid where textareas
  // (and fields flagged `full`) span both columns.
  columns: { type: Number, default: 1 },
  // Per-field validation messages, keyed by field.key. Only shown once
  // `submitted` is true, matching the section's other forms (error styling
  // stays quiet until the user actually tries to submit).
  errors: { type: Object, default: () => ({}) },
  submitted: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

function update(key, value) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

// Prefer the i18n placeholderKey idiom; fall back to a plain `placeholder`
// string for any schema that still carries literal text.
function fieldPlaceholder(field) {
  return field.placeholderKey ? t(field.placeholderKey) : field.placeholder
}
</script>
