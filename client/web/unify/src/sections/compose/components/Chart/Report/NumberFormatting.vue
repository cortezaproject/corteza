<template>
  <Fieldset :legend="$t('chart.edit.formatting.presetFormats.label')" class="mt-4">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <CFormGroup :label="$t('chart.edit.formatting.prefix.label')" :input-id="`${idPrefix}Prefix`">
        <InputText
          :id="`${idPrefix}Prefix`"
          v-model="model.prefix"
          :placeholder="$t('chart.edit.formatting.prefix.placeholder')"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup :label="$t('chart.edit.formatting.suffix.label')" :input-id="`${idPrefix}Suffix`">
        <InputText
          :id="`${idPrefix}Suffix`"
          v-model="model.suffix"
          :placeholder="$t('chart.edit.formatting.suffix.placeholder')"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('chart.edit.formatting.presetFormats.label')"
        :input-id="`${idPrefix}PresetFormat`"
      >
        <Select
          :id="`${idPrefix}PresetFormat`"
          v-model="model.presetFormat"
          :options="formatOptions"
          option-label="text"
          option-value="value"
          class="w-full"
        />
        <template v-if="model.presetFormat" #description>
          <span class="whitespace-pre-line">
            {{ $t(`chart.edit.formatting.presetFormats.description.${model.presetFormat}`) }}
          </span>
        </template>
      </CFormGroup>

      <CFormGroup :label="$t('chart.edit.formatting.format.label')" :input-id="`${idPrefix}Format`">
        <InputText
          :id="`${idPrefix}Format`"
          v-model="model.format"
          :disabled="model.presetFormat !== 'custom'"
          :placeholder="$t('chart.edit.formatting.format.placeholder')"
          class="w-full"
        />
      </CFormGroup>
    </div>
  </Fieldset>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// Number presentation for a metric or an axis. A model rather than a plain
// prop because the object is edited in place, the way the rest of the report
// draft is — the parent keeps the same object and sees the writes.
const model = defineModel({ type: Object, required: true })

defineProps({
  // Distinguishes this instance's input ids from every other formatting block
  // on the page — a chart can show one per metric plus one for the axis.
  idPrefix: {
    type: String,
    required: true,
  },
})

const formatOptions = [
  { value: 'custom', text: t('chart.edit.formatting.presetFormats.options.custom') },
  { value: 'accounting', text: t('chart.edit.formatting.presetFormats.options.accounting') },
]
</script>
