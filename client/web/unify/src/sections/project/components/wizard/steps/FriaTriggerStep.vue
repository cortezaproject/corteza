<template>
  <FriaScenarioEditorShell :project="project" :disabled="disabled" v-slot="{ scenario, update }">
    <div class="flex flex-col gap-6 max-w-4xl">
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
            @toggle="toggle(scenario, update, 'triggerTypes', c.key)"
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
  </FriaScenarioEditorShell>
</template>

<script setup>
import FriaScenarioEditorShell from '@/sections/project/components/wizard/fria/FriaScenarioEditorShell.vue'
import FriaToggleCard from '@/sections/project/components/wizard/fria/FriaToggleCard.vue'
import { TRIGGER_ICONS } from '@/sections/project/config/friaScenario'
import { TRIGGER_CONDITIONS } from '@/sections/project/config/friaTaxonomies'

defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

function toggle(scenario, update, field, key) {
  const list = scenario[field]
  update({
    [field]: list.includes(key) ? list.filter(k => k !== key) : [...list, key],
  })
}
</script>
