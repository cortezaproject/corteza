<template>
  <!-- § 2 of the scenario editor — content moved verbatim from the old
       fria-trigger wizard step (components/wizard/steps/FriaTriggerStep.vue,
       removed); see FriaScenarioEditor.vue. -->
  <FriaScenarioSection
    :number="number"
    :title="$t('fria.sections.trigger.title')"
    :description="$t('fria.sections.trigger.description')"
  >
    <div class="flex flex-col gap-6">
      <CFormGroup :label="$t('fria.trigger.typesLabel')" required>
        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-3">
          <FriaToggleCard
            v-for="c in TRIGGER_CONDITIONS"
            :key="c.key"
            :label="$t(c.labelKey)"
            :description="$t(c.descriptionKey)"
            :icon="TRIGGER_ICONS[c.key]"
            :selected="scenario.triggerTypes.includes(c.key)"
            :disabled="disabled"
            @toggle="toggle('triggerTypes', c.key)"
          />
        </div>
      </CFormGroup>

      <CFormGroup :label="$t('fria.trigger.contextLabel')">
        <Textarea
          :model-value="scenario.triggerDescription"
          rows="4"
          auto-resize
          fluid
          :disabled="disabled"
          :placeholder="$t('fria.trigger.contextPlaceholder')"
          @update:model-value="v => update({ triggerDescription: v })"
        />
      </CFormGroup>
    </div>
  </FriaScenarioSection>
</template>

<script setup>
import FriaScenarioSection from './FriaScenarioSection.vue'
import FriaToggleCard from './FriaToggleCard.vue'
import { TRIGGER_ICONS } from '@/sections/project/config/friaScenario'
import { TRIGGER_CONDITIONS } from '@/sections/project/config/friaTaxonomies'

const props = defineProps({
  number: { type: [Number, String], required: true },
  scenario: { type: Object, required: true },
  update: { type: Function, required: true },
  disabled: { type: Boolean, default: false },
})

function toggle(field, key) {
  const list = props.scenario[field]
  props.update({
    [field]: list.includes(key) ? list.filter(k => k !== key) : [...list, key],
  })
}
</script>
