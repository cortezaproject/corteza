<template>
  <!-- § 5 of the scenario editor — content moved verbatim from the old
       fria-vectors wizard step (components/wizard/steps/FriaVectorsStep.vue,
       removed); see FriaScenarioEditor.vue. -->
  <FriaScenarioSection
    :number="number"
    :title="$t('fria.sections.vectors.title')"
    :description="$t('fria.sections.vectors.description')"
  >
    <div class="flex flex-col gap-6 max-w-4xl">
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-3">
        <FriaToggleCard
          v-for="v in AI_HARM_VECTORS"
          :key="v.key"
          :label="$t(v.labelKey)"
          :description="$t(v.descriptionKey)"
          :icon="AI_HARM_VECTOR_ICONS[v.key]"
          :selected="scenario.harmVectors.includes(v.key)"
          :disabled="disabled"
          @toggle="toggle(v.key)"
        />
      </div>

      <CFormGroup :label="$t('fria.vectors.descriptionLabel')">
        <Textarea
          :model-value="scenario.harmVectorsDescription"
          rows="4"
          auto-resize
          fluid
          :disabled="disabled"
          :placeholder="$t('fria.vectors.descriptionPlaceholder')"
          @update:model-value="v => update({ harmVectorsDescription: v })"
        />
      </CFormGroup>
    </div>
  </FriaScenarioSection>
</template>

<script setup>
import FriaScenarioSection from './FriaScenarioSection.vue'
import FriaToggleCard from './FriaToggleCard.vue'
import { AI_HARM_VECTOR_ICONS } from '@/sections/project/config/friaScenario'
import { AI_HARM_VECTORS } from '@/sections/project/config/friaTaxonomies'

const props = defineProps({
  number: { type: [Number, String], required: true },
  scenario: { type: Object, required: true },
  update: { type: Function, required: true },
  disabled: { type: Boolean, default: false },
})

function toggle(key) {
  const list = props.scenario.harmVectors
  props.update({
    harmVectors: list.includes(key) ? list.filter(k => k !== key) : [...list, key],
  })
}
</script>
