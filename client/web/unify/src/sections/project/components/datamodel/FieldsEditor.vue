<template>
  <div class="cform-table rounded-lg border border-surface overflow-hidden">
  <CFormList
    v-model="fields"
    :columns="columns"
    empty-message="No fields yet."
    :hide-remove="disabled"
  >
    <template #row="{ item: f }">
      <InputText
        v-model="f.name"
        placeholder="Field name"
        size="small"
        fluid
        :disabled="disabled"
      />
      <Select
        v-model="f.type"
        :options="FIELD_TYPES"
        option-label="label"
        option-value="id"
        size="small"
        fluid
        :disabled="disabled"
        @update:model-value="onType(f)"
      />
      <Select
        v-model="f.sensitivity"
        :options="SENSITIVITY_OPTIONS"
        option-label="label"
        option-value="id"
        size="small"
        fluid
        :disabled="disabled"
      />
      <div class="flex justify-center">
        <ToggleSwitch v-model="f.required" :disabled="disabled" />
      </div>
    </template>

    <!-- Type-specific settings (e.g. a record reference's target module) -->
    <template #extra="{ item: f }">
      <Fieldset
        v-if="isRecordRef(f.type)"
        legend="Record"
        class="mt-1"
        :pt="{ content: { class: '!p-3' } }"
      >
        <CFormGroup label="Target module">
          <Select
            v-model="f.targetModuleId"
            :options="moduleOptions"
            option-label="name"
            option-value="id"
            size="small"
            placeholder="Select module"
            class="w-72"
            :disabled="disabled"
          />
        </CFormGroup>
      </Fieldset>
    </template>
  </CFormList>
  </div>
</template>

<script setup>
import { FIELD_TYPES, isRecordRef } from '@/sections/project/config/fieldTypes'
import { SENSITIVITY_OPTIONS } from '@/sections/project/config/sensitivity'

// Staged field list (two-way; the dialog owns add, CFormList owns remove).
const fields = defineModel({ type: Array, required: true })

defineProps({
  // Modules a Record field can point at.
  moduleOptions: { type: Array, default: () => [] },
  disabled: { type: Boolean, default: false },
})

const columns = [
  { label: 'Name', width: 'minmax(12rem, 1fr)' },
  { label: 'Type', width: '13rem' },
  { label: 'Sensitivity', width: '13rem' },
  { label: 'Required', width: '6rem', headerClass: 'text-center' },
]

// Clear the target when a field is no longer a record reference.
function onType(f) {
  if (!isRecordRef(f.type)) f.targetModuleId = null
}
</script>

<style scoped>
/* Flatten CFormList's card rows into a single bounded table. */
.cform-table :deep(.flex.flex-col.gap-3) {
  gap: 0;
}
/* Header bar */
.cform-table :deep(.grid.pt-3) {
  padding: 0.5rem 0.75rem;
  background: var(--p-content-hover-background, rgba(0, 0, 0, 0.03));
  border-bottom: 1px solid var(--p-content-border-color);
}
/* Rows: drop the per-card border/shadow/radius, separate with a divider */
.cform-table :deep(.shadow-sm) {
  border: 0;
  border-radius: 0;
  box-shadow: none;
  border-top: 1px solid var(--p-content-border-color);
}
.cform-table :deep(.shadow-sm:first-of-type) {
  border-top: 0;
}
/* Empty-state box already sits inside our border */
.cform-table :deep(.bg-highlight) {
  border: 0;
}
</style>
