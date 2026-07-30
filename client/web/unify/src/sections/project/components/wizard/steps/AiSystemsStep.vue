<template>
  <!-- Opening a system replaces this whole body with its editor; going back
       returns here. Same list <-> editor swap the fria-scenarios step uses,
       driven by the `aiSystem` query param (composables/useActiveAiSystem.js)
       so an open system stays deep-linkable and back/forward still work.
       Keyed by id so switching straight from one system to another always
       remounts against fresh state. -->
  <AiSystemEditor
    v-if="activeAiSystemId"
    :key="activeAiSystemId"
    :project="project"
    :ai-system-id="activeAiSystemId"
    :disabled="disabled"
    @closed="closeAiSystem"
  />

  <div v-else class="flex flex-col gap-6">
    <!-- Stat row — mirrors the scenario step's, so the two Govern surfaces
         read as one family. -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium">{{ systems.length }}</div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.aiSystems.stats.total') }}
        </div>
      </div>
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium text-red-600 dark:text-red-400">
          {{ countOfClass('prohibited') }}
        </div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.riskClasses.prohibited.label') }}
        </div>
      </div>
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium text-orange-600 dark:text-orange-400">
          {{ countOfClass('high') }}
        </div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.riskClasses.high.label') }}
        </div>
      </div>
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium text-primary">{{ unclassifiedCount }}</div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.aiSystems.stats.unclassified') }}
        </div>
      </div>
    </div>

    <!-- A prohibited system is not a risk to manage, it is a practice the Act
         forbids outright (Art. 5). Surfaced as a banner rather than a row
         badge because it should stop the build, not decorate it. -->
    <Message v-if="countOfClass('prohibited') > 0" severity="error" :closable="false">
      {{ $t('fria.aiSystems.prohibitedWarning') }}
    </Message>

    <div class="flex flex-col gap-3">
      <article
        v-for="system in systems"
        :key="system.id"
        class="rounded-xl border border-surface bg-surface overflow-hidden"
      >
        <button
          type="button"
          class="w-full text-left px-4 py-3 flex items-start gap-3 hover:bg-emphasis transition-colors"
          @click="openAiSystem(system.id)"
        >
          <i class="pi pi-sitemap text-primary mt-1" />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="font-medium truncate">
                {{ system.name || $t('fria.aiSystems.untitled') }}
              </span>
              <RiskClassBadge :risk-class="system.riskClass" />
              <span
                v-if="friaRequiredForRiskClass(system.riskClass)"
                class="text-xs px-2 py-0.5 rounded-full bg-primary/10 text-primary"
              >
                {{ $t('fria.aiSystems.friaRequired') }}
              </span>
            </div>
            <p v-if="system.intendedPurpose" class="text-sm text-muted-color mt-1 line-clamp-2">
              {{ system.intendedPurpose }}
            </p>
            <p v-else class="text-sm text-muted-color italic mt-1">
              {{ $t('fria.aiSystems.noPurpose') }}
            </p>
          </div>
          <span class="text-xs text-muted-color shrink-0 mt-1">
            {{ $t('fria.aiSystems.scenarioCount', scenarioCountFor(system.id)) }}
          </span>
        </button>
      </article>

      <div v-if="!systems.length" class="text-sm text-muted-color italic px-1">
        {{ $t('fria.aiSystems.empty') }}
      </div>

      <!-- Inline create. No 'new' sentinel route: a system's membership rows
           are keyed by its ID, so the row has to exist before an editor can
           mean anything (see useActiveAiSystem.js). -->
      <div v-if="!disabled" class="rounded-lg border border-dashed border-surface p-3 bg-emphasis">
        <InputText
          v-model="draftName"
          :placeholder="$t('fria.aiSystems.namePlaceholder')"
          class="w-full mb-2"
          size="small"
          fluid
          @keydown.enter="create"
        />
        <div class="flex items-center justify-between gap-3">
          <span class="text-xs text-muted-color">{{ $t('fria.aiSystems.createHint') }}</span>
          <Button
            icon="pi pi-plus"
            :label="$t('fria.aiSystems.add')"
            severity="secondary"
            size="small"
            :disabled="!draftName.trim() || creating"
            :loading="creating"
            @click="create"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useActiveAiSystem } from '@/sections/project/composables/useActiveAiSystem'
import { friaRequiredForRiskClass } from '@/sections/project/config/aiSystemClasses'
import AiSystemEditor from '@/sections/project/components/wizard/aiSystems/AiSystemEditor.vue'
import RiskClassBadge from '@/sections/project/components/wizard/aiSystems/RiskClassBadge.vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { activeAiSystemId, openAiSystem, closeAiSystem } = useActiveAiSystem()

const projectId = computed(() => props.project?.projectID)
const systems = computed(() => store.aiSystemsFor(projectId.value))

const draftName = ref('')
const creating = ref(false)

onMounted(() => {
  store.loadAiSystems(projectId.value)
  // The per-system scenario count reads these; without the load it would
  // report zero everywhere, which reads as "unassessed" rather than "not
  // loaded yet".
  store.loadFriaScenarios(projectId.value)
})

function countOfClass(key) {
  return systems.value.filter(s => s.riskClass === key).length
}

// Unclassified is called out on its own because an AI system with no risk
// class is not "minimal risk" — it is unassessed, and the two must never look
// the same on a compliance surface.
const unclassifiedCount = computed(() => systems.value.filter(s => !s.riskClass).length)

// How many risk scenarios name this system, from the store's scenario cache
// (loaded in onMounted above).
function scenarioCountFor(aiSystemId) {
  return store.friaScenariosFor(projectId.value).filter(s => s.aiSystemID === aiSystemId).length
}

async function create() {
  const name = draftName.value.trim()
  if (!name || creating.value) return

  creating.value = true
  try {
    const id = await store.createAiSystem(projectId.value, { name })
    draftName.value = ''
    openAiSystem(id)
  } finally {
    creating.value = false
  }
}
</script>
