<template>
  <!-- The FRIA scenario editor: ONE scrolling page stacking all five design-
       mockup sections rather than five separate wizard steps (see
       config/pipeline.js).
       Mounted by components/wizard/steps/FriaScenariosStep.vue in place of
       the scenario list whenever a scenario is open (existing or new — see
       composables/useFriaActiveScenario.js), keyed by scenario id so a
       switch between scenarios always remounts fresh. -->
  <div class="flex flex-col">
    <template v-if="draft">
      <!-- Everything above the action bar scrolls; the bar does not. The step
           hands this component the whole pane, unpadded (views/Wizard.vue
           stepOwnsBody), so the padding the sections need lives here. -->
      <div class="flex-1 min-h-0 overflow-y-auto p-4 flex flex-col gap-5">
        <header class="flex items-center gap-3 flex-wrap">
          <Button
            icon="pi pi-arrow-left"
            :label="$t('fria.editor.allScenarios')"
            text
            severity="secondary"
            size="small"
            @click="cancel"
          />
          <h3 class="text-base font-medium flex-1 min-w-0 truncate">
            {{ draft.title?.trim() || $t('fria.editor.untitledScenario') }}
          </h3>
          <FriaSeverityBadge :severity="draft.severity" />
        </header>

        <!-- Which AI system this scenario assesses. Sits ABOVE the five
             numbered sections rather than inside one because it is not part
             of describing the harm — it is what the harm is being described
             ABOUT. Art. 27 attaches a FRIA to a specific high-risk AI system,
             and a project can hold several, so a scenario with no system
             assesses nothing in particular. It is the one field Save blocks
             on (see friaScenarioSaveable). -->
        <section
          class="rounded-xl border p-4 flex flex-col gap-2"
          :class="draft.aiSystemID ? 'border-surface bg-surface' : 'border-primary/40 bg-primary/5'"
        >
          <label class="text-sm font-medium" for="fria-ai-system">
            {{ $t('fria.editor.aiSystem.label') }}
          </label>
          <Select
            id="fria-ai-system"
            :model-value="draft.aiSystemID"
            :options="aiSystemOptions"
            option-label="label"
            option-value="value"
            :disabled="disabled"
            :placeholder="$t('fria.editor.aiSystem.placeholder')"
            :empty-message="$t('fria.editor.aiSystem.none')"
            fluid
            @update:model-value="v => update({ aiSystemID: v })"
          />
          <small class="text-muted-color">{{ $t('fria.editor.aiSystem.hint') }}</small>
        </section>

        <FriaHarmSection :number="1" :scenario="draft" :update="update" :disabled="disabled" />
        <FriaTriggerSection :number="2" :scenario="draft" :update="update" :disabled="disabled" />
        <FriaPartiesSection :number="3" :scenario="draft" :update="update" :disabled="disabled" />
        <FriaRightsSection :number="4" :scenario="draft" :update="update" :disabled="disabled" />
        <FriaVectorsSection :number="5" :scenario="draft" :update="update" :disabled="disabled" />
      </div>

      <!-- Explicit save/cancel — edits above only ever touch the local draft
           (see the `draft` ref below); nothing reaches the store until Save,
           and Cancel just drops the draft. The panel's own bottom edge, not a
           sticky bar over the form: it is a sibling of the scrolling area
           rather than a box inside it, so it spans the full width, stays put,
           and nothing passes behind it. -->
      <div
        class="shrink-0 border-t border-surface bg-surface px-4 py-3 flex items-center justify-between gap-3"
      >
        <div class="flex items-center gap-3 min-w-0">
          <div class="w-28 h-1.5 rounded-full bg-emphasis overflow-hidden shrink-0">
            <div
              class="h-full bg-primary rounded-full transition-all"
              :style="{ width: `${(progress.filled / progress.total) * 100}%` }"
            />
          </div>
          <span class="text-xs text-muted-color whitespace-nowrap truncate">
            {{ $t('fria.editor.sectionsProgress', progress) }}
          </span>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="cancel"
          />
          <!-- Blocked until an AI system is chosen. Everything else in a
               scenario can be filled in over time; this one cannot be
               deferred, because without it there is nothing the assessment
               is about. The hint says why rather than leaving a dead
               button. -->
          <span v-if="!disabled && !saveable" class="text-xs text-muted-color whitespace-nowrap">
            {{ $t('fria.editor.aiSystem.required') }}
          </span>
          <Button
            v-if="!disabled"
            :label="$t('fria.editor.saveScenario')"
            icon="pi pi-check"
            size="small"
            :disabled="!saveable"
            :loading="saving"
            @click="save"
          />
        </div>
      </div>
    </template>

    <!-- Scenario id in the URL doesn't resolve to anything (removed, or a
         stale/bad deep link) — offer a way back rather than a blank page. -->
    <CEmptyState v-else class="flex flex-col items-center gap-3 p-4 py-10">
      <p>{{ $t('fria.editor.noActiveScenario') }}</p>
      <Button
        :label="$t('fria.editor.allScenarios')"
        icon="pi pi-arrow-left"
        severity="secondary"
        text
        size="small"
        @click="cancel"
      />
    </CEmptyState>
  </div>
</template>

<script setup>
import FriaHarmSection from './FriaHarmSection.vue'
import FriaPartiesSection from './FriaPartiesSection.vue'
import FriaRightsSection from './FriaRightsSection.vue'
import FriaSeverityBadge from './FriaSeverityBadge.vue'
import FriaTriggerSection from './FriaTriggerSection.vue'
import FriaVectorsSection from './FriaVectorsSection.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import {
  cloneFriaScenario,
  friaScenarioCompletion,
  friaScenarioSaveable,
  newFriaScenario,
} from '@/sections/project/config/friaScenario'
import { components } from '@planetcrust/human-vue'
import Select from 'primevue/select'
import { computed, onMounted, ref, watch } from 'vue'

const { CEmptyState } = components

const props = defineProps({
  project: { type: Object, required: true },
  // The scenario being edited, or null to create a new one. The parent
  // (FriaScenariosStep.vue) resolves the 'new' sentinel from
  // useFriaActiveScenario's isCreatingScenario down to this prop, and keys
  // this component by scenario id so a switch always remounts with a fresh
  // draft rather than reusing stale local state.
  scenarioId: { type: String, default: null },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['saved', 'cancelled'])

const store = useProjectsStore()

// Seeded ONCE at mount — the local draft every section component below
// reads and patches via `update`. An existing scenario is cloned (so edits
// never touch the store in place before Save); a new one starts from the
// blank shape. `null` means the scenarioId prop named a scenario that no
// longer exists — see the CEmptyState fallback above.
const existing = props.scenarioId
  ? store.friaScenario(props.project.projectID, props.scenarioId)
  : null
const draft = ref(
  props.scenarioId ? (existing ? cloneFriaScenario(existing) : null) : newFriaScenario(),
)

// Re-seed once the scenario arrives, if it was not cached at setup. Since
// scenarios became persisted, a refresh or a shared link straight into
// ?scenario=<id> reaches this component before the parent's load resolves —
// the draft would seed to null and the editor would show "not found"
// permanently, for a scenario that exists. Only fills a null draft, so it can
// never clobber edits in progress.
watch(
  () => store.friaScenario(props.project.projectID, props.scenarioId),
  found => {
    if (found && !draft.value) draft.value = cloneFriaScenario(found)
  },
)

function update(patch) {
  draft.value = { ...draft.value, ...patch }
}

// Live completion for the action bar's progress indicator, recomputed as the
// user edits sections above. Reuses the exact same required-field definition
// as the scenario list's Complete/In Progress badge (FriaScenariosStep.vue)
// so the two views can never disagree — see friaScenarioCompletion's own
// comment for which fields count. Only rendered while `draft` is truthy (see
// the template's v-if above), so draft.value is never null here.
const progress = computed(() => friaScenarioCompletion(draft.value))

// The AI systems this project defines, for the picker above. Loaded here
// because the Govern tab's ai-systems step may never have been opened in this
// session — the list would otherwise be empty and look like "none defined"
// when several exist.
onMounted(() => store.loadAiSystems(props.project.projectID))

const aiSystemOptions = computed(() =>
  store.aiSystemsFor(props.project.projectID).map(s => ({
    value: s.id,
    label: s.name || s.handle,
  })),
)

// Deliberately NOT part of `progress`: completion measures how much of the
// assessment has been written, this measures whether it is addressable at all.
// See config/friaScenario.js for why the two are kept apart.
const saveable = computed(() => friaScenarioSaveable(draft.value))

// Commits the draft to the store in one shot — the only place any of this
// screen's edits reach stores/projects.js. See that file's FRIA risk
// scenarios section for why nothing runs until this point.
// Awaited since scenarios became persisted: 'saved' must not fire until the
// write lands, or the list re-renders from a cache the server has not
// confirmed. `saving` keeps the button from being pressed twice.
const saving = ref(false)

async function save() {
  // Guarded here too, not just on the button: a scenario with no AI system
  // must never reach the store, whatever path calls this.
  if (!draft.value || !saveable.value || saving.value) return

  saving.value = true
  try {
    if (props.scenarioId) {
      await store.updateFriaScenario(props.project.projectID, props.scenarioId, draft.value)
    } else {
      await store.createFriaScenario(props.project.projectID, draft.value)
    }
    emit('saved')
  } finally {
    saving.value = false
  }
}

// Discards the draft — nothing to undo in the store, since nothing was ever
// written there (see the header comment above and useFriaActiveScenario.js's
// NEW_SCENARIO_SENTINEL note for the cancelled-create case specifically).
function cancel() {
  emit('cancelled')
}
</script>
