<template>
  <FriaScenarioEditorShell :project="project" :disabled="disabled" v-slot="{ scenario, update }">
    <div class="flex flex-col gap-6 max-w-4xl">
      <p class="text-sm text-muted-color">{{ $t('fria.vectors.intro') }}</p>

      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-3">
        <FriaToggleCard
          v-for="v in AI_HARM_VECTORS"
          :key="v.key"
          :label="$t(v.labelKey)"
          :description="$t(v.descriptionKey)"
          :icon="AI_HARM_VECTOR_ICONS[v.key]"
          :selected="scenario.harmVectors.includes(v.key)"
          :disabled="disabled"
          @toggle="toggle(scenario, update, v.key)"
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
  </FriaScenarioEditorShell>
</template>

<script setup>
import FriaScenarioEditorShell from '@/sections/project/components/wizard/fria/FriaScenarioEditorShell.vue'
import FriaToggleCard from '@/sections/project/components/wizard/fria/FriaToggleCard.vue'
import { AI_HARM_VECTOR_ICONS } from '@/sections/project/config/friaScenario'
import { AI_HARM_VECTORS } from '@/sections/project/config/friaTaxonomies'

defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

function toggle(scenario, update, key) {
  const list = scenario.harmVectors
  update({
    harmVectors: list.includes(key) ? list.filter(k => k !== key) : [...list, key],
  })
}
</script>
