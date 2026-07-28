<template>
  <FriaScenarioEditorShell :project="project" :disabled="disabled" v-slot="{ scenario, update }">
    <div class="flex flex-col gap-6 max-w-4xl">
      <CFormGroup :label="$t('fria.parties.impactedLabel')" required>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="p in IMPACTED_PARTIES"
            :key="p.key"
            type="button"
            class="rounded-full border-2 px-4 py-2 text-xs font-medium transition-colors"
            :class="[
              scenario.impactedParties.includes(p.key)
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-surface hover:border-primary/60',
              disabled ? 'cursor-not-allowed opacity-70' : 'cursor-pointer',
            ]"
            :disabled="disabled"
            :aria-pressed="scenario.impactedParties.includes(p.key)"
            @click="toggle(scenario, update, 'impactedParties', p.key)"
          >
            {{ $t(p.labelKey) }}
          </button>
        </div>
      </CFormGroup>

      <hr class="border-surface" />

      <CFormGroup :label="$t('fria.parties.vulnerableGroupsLabel')">
        <div class="flex items-start gap-3 rounded-lg border border-surface bg-emphasis p-4 mb-1">
          <i class="pi pi-exclamation-triangle text-amber-600 dark:text-amber-400 mt-0.5" />
          <p class="text-xs leading-relaxed">
            <strong class="block font-semibold mb-1">
              {{ $t('fria.parties.proxyCalloutTitle') }}
            </strong>
            {{ $t('fria.parties.proxyCalloutBody') }}
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-2">
          <button
            v-for="g in VULNERABLE_GROUPS"
            :key="g.key"
            type="button"
            class="flex items-start gap-3 rounded-md border-2 px-4 py-3 text-left transition-colors"
            :class="[
              scenario.vulnerableGroups.includes(g.key)
                ? 'border-primary bg-primary/10'
                : 'border-surface hover:border-primary/60',
              disabled ? 'cursor-not-allowed opacity-70' : 'cursor-pointer',
            ]"
            :disabled="disabled"
            :aria-pressed="scenario.vulnerableGroups.includes(g.key)"
            @click="toggle(scenario, update, 'vulnerableGroups', g.key)"
          >
            <span
              class="mt-0.5 shrink-0 w-4 h-4 rounded-sm border-2 flex items-center justify-center"
              :class="
                scenario.vulnerableGroups.includes(g.key)
                  ? 'bg-primary border-primary'
                  : 'border-surface-400'
              "
            >
              <i
                v-if="scenario.vulnerableGroups.includes(g.key)"
                class="pi pi-check text-xs text-primary-contrast"
              />
            </span>
            <span>
              <span
                class="block text-xs font-medium"
                :class="scenario.vulnerableGroups.includes(g.key) ? 'text-primary' : ''"
              >
                {{ $t(g.labelKey) }}
              </span>
              <span class="block text-xs text-muted-color mt-0.5 leading-snug">
                {{ $t(g.sublabelKey) }}
              </span>
            </span>
          </button>
        </div>
      </CFormGroup>

      <CFormGroup :label="$t('fria.parties.notesLabel')">
        <Textarea
          :model-value="scenario.vulnerableGroupsNotes"
          rows="3"
          auto-resize
          fluid
          :disabled="disabled"
          :placeholder="$t('fria.parties.notesPlaceholder')"
          @update:model-value="v => update({ vulnerableGroupsNotes: v })"
        />
      </CFormGroup>
    </div>
  </FriaScenarioEditorShell>
</template>

<script setup>
import FriaScenarioEditorShell from '@/sections/project/components/wizard/fria/FriaScenarioEditorShell.vue'
import { IMPACTED_PARTIES, VULNERABLE_GROUPS } from '@/sections/project/config/friaTaxonomies'

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
