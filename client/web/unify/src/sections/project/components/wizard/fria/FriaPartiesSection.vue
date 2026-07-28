<template>
  <!-- § 3 of the scenario editor — content moved verbatim from the old
       fria-parties wizard step (components/wizard/steps/FriaPartiesStep.vue,
       removed); see FriaScenarioEditor.vue. -->
  <FriaScenarioSection
    :number="number"
    :title="$t('fria.sections.parties.title')"
    :description="$t('fria.sections.parties.description')"
  >
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
            @click="toggle('impactedParties', p.key)"
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
            @click="toggle('vulnerableGroups', g.key)"
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
  </FriaScenarioSection>
</template>

<script setup>
import FriaScenarioSection from './FriaScenarioSection.vue'
import { IMPACTED_PARTIES, VULNERABLE_GROUPS } from '@/sections/project/config/friaTaxonomies'

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
