<template>
  <!-- The FRIA scenario editor: ONE scrolling page stacking all five design-
       mockup sections, replacing what used to be five separate wizard steps
       (fria-harm/trigger/parties/rights/vectors — see config/pipeline.js).
       Mounted by components/wizard/steps/FriaScenariosStep.vue in place of
       the scenario list whenever a scenario is open (existing or new — see
       composables/useFriaActiveScenario.js), keyed by scenario id so a
       switch between scenarios always remounts fresh. -->
  <div class="flex flex-col gap-5">
    <template v-if="draft">
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

      <FriaHarmSection :number="1" :scenario="draft" :update="update" :disabled="disabled" />
      <FriaTriggerSection :number="2" :scenario="draft" :update="update" :disabled="disabled" />
      <FriaPartiesSection :number="3" :scenario="draft" :update="update" :disabled="disabled" />
      <FriaRightsSection :number="4" :scenario="draft" :update="update" :disabled="disabled" />
      <FriaVectorsSection :number="5" :scenario="draft" :update="update" :disabled="disabled" />

      <!-- Explicit save/cancel — edits above only ever touch the local
           draft (see the `draft` ref below); nothing reaches the store
           until Save, and Cancel just drops the draft. Sticky so it stays
           reachable at the bottom of a long scroll, mirroring the design
           mockup's sticky action bar. Negative-margined past the p-4 padding
           of its scrolling ancestor (views/Wizard.vue's step content pane)
           so it reads as a real full-bleed bottom bar — square, top border
           only, opaque bg-surface — rather than an inset floating card that
           scrolled content shows around and underneath. -->
      <div
        class="sticky bottom-0 -mx-4 -mb-4 border-t border-surface bg-surface px-4 py-3 flex items-center justify-between gap-3"
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
          <Button
            v-if="!disabled"
            :label="$t('fria.editor.saveScenario')"
            icon="pi pi-check"
            size="small"
            @click="save"
          />
        </div>
      </div>
    </template>

    <!-- Scenario id in the URL doesn't resolve to anything (removed, or a
         stale/bad deep link) — offer a way back rather than a blank page. -->
    <CEmptyState v-else class="flex flex-col items-center gap-3 py-10">
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
  newFriaScenario,
} from '@/sections/project/config/friaScenario'
import { components } from '@planetcrust/human-vue'
import { computed, ref } from 'vue'

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

function update(patch) {
  draft.value = { ...draft.value, ...patch }
}

// Live completion for the sticky bar's progress indicator, recomputed as the
// user edits sections above. Reuses the exact same required-field definition
// as the scenario list's Complete/In Progress badge (FriaScenariosStep.vue)
// so the two views can never disagree — see friaScenarioCompletion's own
// comment for which fields count. Only rendered while `draft` is truthy (see
// the template's v-if above), so draft.value is never null here.
const progress = computed(() => friaScenarioCompletion(draft.value))

// Commits the draft to the store in one shot — the only place any of this
// screen's edits reach stores/projects.js. See that file's FRIA risk
// scenarios section for why nothing runs until this point.
function save() {
  if (!draft.value) return
  if (props.scenarioId) {
    store.updateFriaScenario(props.project.projectID, props.scenarioId, draft.value)
  } else {
    store.createFriaScenario(props.project.projectID, draft.value)
  }
  emit('saved')
}

// Discards the draft — nothing to undo in the store, since nothing was ever
// written there (see the header comment above and useFriaActiveScenario.js's
// NEW_SCENARIO_SENTINEL note for the cancelled-create case specifically).
function cancel() {
  emit('cancelled')
}
</script>
