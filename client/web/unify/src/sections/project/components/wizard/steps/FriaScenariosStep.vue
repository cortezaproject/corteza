<template>
  <!-- Opening a scenario (existing or new) replaces this whole body with the
       single-page editor; going back returns here. See
       composables/useFriaActiveScenario.js for the `scenario` query param
       driving the swap, and FriaScenarioEditor.vue for the editor itself.
       Keyed by activeScenarioId so switching straight from one scenario to
       another (e.g. via browser back/forward) always remounts with a fresh
       local draft rather than reusing stale editor state. -->
  <FriaScenarioEditor
    v-if="activeScenarioId"
    :key="activeScenarioId"
    :project="project"
    :scenario-id="isCreatingScenario ? null : activeScenarioId"
    :disabled="disabled"
    @saved="closeScenario"
    @cancelled="closeScenario"
  />

  <div v-else class="flex flex-col gap-6">
    <!-- Stat row -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium">{{ stats.total }}</div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.scenarios.stats.total') }}
        </div>
      </div>
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium text-red-600 dark:text-red-400">
          {{ stats.critical }}
        </div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.severity.critical.label') }}
        </div>
      </div>
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium text-orange-600 dark:text-orange-400">
          {{ stats.high }}
        </div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.severity.high.label') }}
        </div>
      </div>
      <div class="rounded-xl border border-surface bg-surface p-4">
        <div class="text-2xl font-medium text-primary">{{ stats.complete }}</div>
        <div class="text-xs text-muted-color uppercase tracking-wide">
          {{ $t('fria.scenarios.stats.complete') }}
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex items-center gap-3 flex-wrap">
      <IconField class="w-64">
        <InputIcon class="pi pi-search" />
        <InputText v-model="search" fluid :placeholder="$t('fria.scenarios.searchPlaceholder')" />
      </IconField>
      <SelectButton
        v-model="severityFilter"
        :options="SEVERITY_FILTERS"
        option-label="label"
        option-value="value"
        :allow-empty="false"
      />
      <Button
        v-if="!disabled"
        class="ml-auto"
        :label="$t('fria.scenarios.newScenario')"
        icon="pi pi-plus"
        size="small"
        @click="addScenario"
      />
    </div>

    <!-- Grid -->
    <div v-if="scenarios.length" class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <div
        v-for="s in filtered"
        :key="s.id"
        class="rounded-xl border border-surface bg-surface p-5 flex flex-col gap-3 hover:shadow-md transition-shadow"
      >
        <div class="flex items-start gap-3">
          <FriaSeverityBadge :severity="s.severity" class="mt-0.5" />
          <h3 class="font-medium leading-snug flex-1 min-w-0">
            {{ s.title?.trim() || $t('fria.scenarios.untitled', { index: indexOf(s) + 1 }) }}
          </h3>
        </div>

        <div v-if="s.impactedParties.length || s.triggerTypes.length" class="flex flex-wrap gap-2">
          <span
            v-if="s.impactedParties.length"
            class="text-xs text-muted-color bg-emphasis rounded-full px-2.5 py-1"
          >
            <i class="pi pi-users text-xs mr-1" />
            {{ partyLabel(s.impactedParties[0]) }}
            <template v-if="s.impactedParties.length > 1">
              +{{ s.impactedParties.length - 1 }}
            </template>
          </span>
          <span
            v-if="s.triggerTypes.length"
            class="text-xs text-muted-color bg-emphasis rounded-full px-2.5 py-1"
          >
            {{ triggerLabel(s.triggerTypes[0]) }}
            <template v-if="s.triggerTypes.length > 1">+{{ s.triggerTypes.length - 1 }}</template>
          </span>
        </div>

        <p class="text-xs text-muted-color leading-relaxed line-clamp-3">
          {{ s.description || $t('fria.scenarios.noDescription') }}
        </p>

        <div
          v-if="s.rights.length || s.harmVectors.length || s.vulnerableGroups.length"
          class="flex flex-wrap gap-1.5"
        >
          <span
            v-if="s.rights.length"
            class="text-xs font-medium rounded-full px-2 py-0.5 bg-primary/10 text-primary"
          >
            {{ $t('fria.scenarios.tagRights', s.rights.length) }}
          </span>
          <span
            v-if="s.harmVectors.length"
            class="text-xs font-medium rounded-full px-2 py-0.5 bg-orange-500/10 text-orange-600 dark:text-orange-400"
          >
            {{ $t('fria.scenarios.tagVectors', s.harmVectors.length) }}
          </span>
          <span
            v-if="s.vulnerableGroups.length"
            class="text-xs font-medium rounded-full px-2 py-0.5 bg-purple-500/10 text-purple-600 dark:text-purple-400"
          >
            {{ $t('fria.scenarios.tagVulnerableGroups', s.vulnerableGroups.length) }}
          </span>
        </div>

        <div class="mt-1 pt-3 border-t border-surface flex items-center justify-between">
          <span
            class="text-xs flex items-center gap-1.5"
            :class="
              completion(s).complete
                ? 'text-emerald-600 dark:text-emerald-400'
                : 'text-amber-600 dark:text-amber-400'
            "
          >
            <i :class="completion(s).complete ? 'pi pi-check-circle' : 'pi pi-circle'" />
            {{
              completion(s).complete
                ? $t('fria.scenarios.statusComplete')
                : $t('fria.scenarios.statusInProgress', completion(s))
            }}
          </span>
          <div class="flex items-center gap-1">
            <Button
              :label="$t('general.label.edit')"
              text
              size="small"
              @click="openScenario(s.id)"
            />
            <Button
              v-if="!disabled"
              icon="pi pi-trash"
              text
              rounded
              severity="secondary"
              size="small"
              @click="removeScenario(s)"
            />
          </div>
        </div>
      </div>

      <!-- New scenario card -->
      <button
        v-if="!disabled"
        type="button"
        class="rounded-xl border-2 border-dashed border-surface hover:border-primary hover:bg-primary/5 transition-colors flex flex-col items-center justify-center gap-2 p-8 min-h-48 text-muted-color hover:text-primary"
        @click="addScenario"
      >
        <i class="pi pi-plus text-xl" />
        <span class="text-sm font-medium">{{ $t('fria.scenarios.documentNew') }}</span>
      </button>
    </div>

    <CEmptyState v-else class="flex flex-col items-center gap-3 py-10">
      <p>{{ $t('fria.scenarios.empty') }}</p>
      <Button
        v-if="!disabled"
        :label="$t('fria.scenarios.newScenario')"
        icon="pi pi-plus"
        size="small"
        @click="addScenario"
      />
    </CEmptyState>
  </div>
</template>

<script setup>
import FriaScenarioEditor from '@/sections/project/components/wizard/fria/FriaScenarioEditor.vue'
import FriaSeverityBadge from '@/sections/project/components/wizard/fria/FriaSeverityBadge.vue'
import { useFriaActiveScenario } from '@/sections/project/composables/useFriaActiveScenario'
import { friaScenarioCompletion } from '@/sections/project/config/friaScenario'
import { IMPACTED_PARTIES, TRIGGER_CONDITIONS } from '@/sections/project/config/friaTaxonomies'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const { CEmptyState } = components
const { confirmDelete } = useConfirmDelete()

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { activeScenarioId, isCreatingScenario, openScenario, openNewScenario, closeScenario } =
  useFriaActiveScenario()

const scenarios = computed(() => store.friaScenariosFor(props.project.projectID))
const indexOf = s => scenarios.value.indexOf(s)

const partyLabel = key => t(IMPACTED_PARTIES.find(p => p.key === key)?.labelKey || '')
const triggerLabel = key => t(TRIGGER_CONDITIONS.find(tr => tr.key === key)?.labelKey || '')
const completion = s => friaScenarioCompletion(s)

const stats = computed(() => {
  const list = scenarios.value
  return {
    total: list.length,
    critical: list.filter(s => s.severity === 'critical').length,
    high: list.filter(s => s.severity === 'high').length,
    complete: list.filter(s => friaScenarioCompletion(s).complete).length,
  }
})

const search = ref('')
const severityFilter = ref('all')
const SEVERITY_FILTERS = computed(() => [
  { label: t('general.label.all'), value: 'all' },
  { label: t('fria.severity.critical.label'), value: 'critical' },
  { label: t('fria.severity.high.label'), value: 'high' },
  { label: t('fria.severity.medium.label'), value: 'medium' },
  { label: t('fria.severity.low.label'), value: 'low' },
])

const filtered = computed(() =>
  scenarios.value.filter(s => {
    if (severityFilter.value !== 'all' && s.severity !== severityFilter.value) return false
    if (!search.value.trim()) return true
    return (s.title || '').toLowerCase().includes(search.value.trim().toLowerCase())
  }),
)

// Opens the editor in create mode WITHOUT touching the store — the new
// scenario only becomes real if the editor's Save is used (see
// FriaScenarioEditor.vue and useFriaActiveScenario.js's NEW_SCENARIO_SENTINEL).
function addScenario() {
  openNewScenario()
}

function removeScenario(s) {
  confirmDelete({
    header: t('fria.scenarios.removeConfirm.header'),
    message: t('fria.scenarios.removeConfirm.message', {
      title: s.title?.trim() || t('fria.scenarios.untitled', { index: indexOf(s) + 1 }),
    }),
    onConfirm: () => store.removeFriaScenario(props.project.projectID, s.id),
  })
}
</script>
