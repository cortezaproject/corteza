<template>
  <div class="flex flex-col gap-6">
    <!-- Module selector -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-surface-500 text-sm">{{ $t('field.kind.record.moduleLabel') }}</label>
      <Select
        v-model="field.options.moduleID"
        :options="moduleOptions"
        option-label="name"
        option-value="moduleID"
        :placeholder="$t('field.kind.record.modulePlaceholder')"
        filter
        class="w-full"
        @update:model-value="onModuleChange"
      />
    </div>

    <template v-if="selectedModule">
      <!-- Label field -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-surface-500 text-sm">{{ $t('field.kind.record.moduleField') }}</label>
        <Select
          v-model="field.options.labelField"
          :options="fieldOptions"
          option-label="text"
          option-value="value"
          :placeholder="$t('field.kind.record.pickField')"
          show-clear
          class="w-full"
          @update:model-value="onLabelFieldChange"
        />
      </div>

      <!-- Query fields -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-surface-500 text-sm">{{ $t('field.kind.record.queryFieldsLabel') }}</label>
        <MultiSelect
          v-model="field.options.queryFields"
          :options="queryFieldOptions"
          option-label="text"
          option-value="value"
          :placeholder="$t('field.kind.record.queryFieldsPlaceholder')"
          filter
          class="w-full"
        />
      </div>

      <!-- Prefilter -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-surface-500 text-sm">{{ $t('field.kind.record.prefilterLabel') }}</label>
        <Textarea
          v-model="field.options.prefilter"
          :placeholder="$t('field.kind.record.prefilterPlaceholder')"
          rows="3"
          class="w-full"
        />
        <small class="text-muted-color">
          {{ $t('field.kind.record.prefilterFootnote', ['${record.values.fieldName}', '${recordID}', '${ownerID}', '${userID}']) }}
        </small>
      </div>
    </template>

    <!-- Multi-value select type -->
    <template v-if="field.isMulti">
      <div>
        <label class="font-medium text-surface-500 text-sm block mb-3">
          {{ $t('field.kind.select.optionType.label') }}
        </label>
        <div class="flex flex-col gap-2">
          <div
            v-for="opt in selectTypeOptions"
            :key="opt.value"
            class="flex items-center gap-2"
          >
            <RadioButton
              :input-id="`recordSelectType-${opt.value}`"
              v-model="field.options.selectType"
              :value="opt.value"
              @update:model-value="onSelectTypeChange"
            />
            <label :for="`recordSelectType-${opt.value}`" class="cursor-pointer">{{ opt.label }}</label>
          </div>
        </div>
      </div>

      <div v-if="showAllowDuplicates" class="flex items-center gap-2">
        <Checkbox
          :model-value="!field.options.isUniqueMultiValue"
          inputId="recordAllowDuplicates"
          :binary="true"
          @update:model-value="field.options.isUniqueMultiValue = !$event"
        />
        <label for="recordAllowDuplicates" class="cursor-pointer">{{ $t('field.kind.select.allow-duplicates') }}</label>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const moduleStore = useModuleStore()

// Non-queryable field kinds (same as in lib/js/src/compose/types/module-field/base.ts)
const nonQueryableFieldKinds = ['Number', 'Record', 'User', 'Bool', 'DateTime', 'File', 'Geometry']

const moduleOptions = computed(() => moduleStore.set || [])

const selectedModule = computed(() => {
  if (!props.field.options.moduleID || props.field.options.moduleID === '0') return null
  return moduleStore.getByID(props.field.options.moduleID) || null
})

const fieldOptions = computed(() => {
  if (!selectedModule.value) return []
  return selectedModule.value.fields
    .map(f => ({ value: f.name, text: f.label || f.name, kind: f.kind }))
    .sort((a, b) => a.text.localeCompare(b.text))
})

const queryFieldOptions = computed(() =>
  fieldOptions.value.filter(f => !nonQueryableFieldKinds.includes(f.kind)),
)

const duplicatesAllowedTypes = ['default', 'each']
const showAllowDuplicates = computed(() =>
  duplicatesAllowedTypes.includes(props.field.options.selectType),
)

const selectTypeOptions = computed(() => [
  { value: 'default', label: t('field.kind.select.optionType.default') },
  { value: 'multiple', label: t('field.kind.select.optionType.multiple') },
  { value: 'each', label: t('field.kind.select.optionType.each') },
])

function onModuleChange() {
  props.field.options.labelField = ''
  props.field.options.queryFields = []
  props.field.options.prefilter = ''
}

function onLabelFieldChange() {
  props.field.options.queryFields = []
  props.field.options.prefilter = ''
}

function onSelectTypeChange(val) {
  if (!duplicatesAllowedTypes.includes(val)) {
    props.field.options.isUniqueMultiValue = true
  }
}
</script>
