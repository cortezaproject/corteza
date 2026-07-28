<template>
  <div class="flex flex-col gap-6">
    <template v-if="scenario">
      <!-- Scenario switcher — lets a member move between scenarios without
           dropping back to the list step, since all five section-editor
           steps (see config/pipeline.js) share this same header. -->
      <header class="flex items-center gap-3 flex-wrap">
        <Select
          :model-value="scenario.id"
          :options="scenarioOptions"
          option-label="label"
          option-value="value"
          class="w-full sm:w-80"
          @update:model-value="setActiveScenarioId"
        />
        <FriaSeverityBadge :severity="scenario.severity" />
        <span class="ml-auto flex items-center gap-2">
          <Button
            v-if="!disabled"
            :label="$t('fria.editor.newScenario')"
            icon="pi pi-plus"
            text
            size="small"
            @click="createNew"
          />
          <Button
            :label="$t('fria.editor.allScenarios')"
            icon="pi pi-arrow-left"
            text
            severity="secondary"
            size="small"
            @click="goToScenarios"
          />
        </span>
      </header>

      <slot :scenario="scenario" :update="update" />
    </template>

    <CEmptyState v-else class="flex flex-col items-center gap-3 py-10">
      <p>{{ $t('fria.editor.noActiveScenario') }}</p>
      <div class="flex items-center gap-2">
        <Button
          v-if="!disabled"
          :label="$t('fria.editor.newScenario')"
          icon="pi pi-plus"
          size="small"
          @click="createNew"
        />
        <Button
          :label="$t('fria.editor.allScenarios')"
          icon="pi pi-table"
          severity="secondary"
          text
          size="small"
          @click="goToScenarios"
        />
      </div>
    </CEmptyState>
  </div>
</template>

<script setup>
import { useFriaActiveScenario } from '@/sections/project/composables/useFriaActiveScenario'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import FriaSeverityBadge from './FriaSeverityBadge.vue'

const { CEmptyState } = components
const { t } = useI18n()

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { activeScenarioId, setActiveScenarioId, openScenario, goToScenarios } =
  useFriaActiveScenario()

const scenarios = computed(() => store.friaScenariosFor(props.project.projectID))
const scenario = computed(() => scenarios.value.find(s => s.id === activeScenarioId.value) || null)

const scenarioOptions = computed(() =>
  scenarios.value.map((s, i) => ({
    value: s.id,
    label: s.title?.trim() || t('fria.scenarios.untitled', { index: i + 1 }),
  })),
)

function update(patch) {
  store.updateFriaScenario(props.project.projectID, scenario.value.id, patch)
}

function createNew() {
  const id = store.createFriaScenario(props.project.projectID)
  openScenario(id)
}
</script>
