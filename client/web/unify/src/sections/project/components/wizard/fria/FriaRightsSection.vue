<template>
  <!-- § 4 of the scenario editor — content moved verbatim from the old
       fria-rights wizard step (components/wizard/steps/FriaRightsStep.vue,
       removed); see FriaScenarioEditor.vue. -->
  <FriaScenarioSection
    :number="number"
    :title="$t('fria.sections.rights.title')"
    :description="$t('fria.sections.rights.description')"
  >
    <div class="flex flex-col gap-5">
      <!-- Chapter filter — a non-colour affordance (icon + roman numeral +
           name) rather than the design mockup's per-chapter colour dot: see
           config/friaScenario.js's RIGHTS_CHAPTER_ICONS comment. A validated
           7-hue chapter palette is a worthwhile follow-up, not done here. -->
      <div class="flex flex-wrap gap-2">
        <button
          type="button"
          class="rounded-full border-2 px-3.5 py-1.5 text-xs font-semibold uppercase tracking-wide transition-colors"
          :class="
            chapterFilter === 'all'
              ? 'bg-primary border-primary text-primary-contrast'
              : 'border-surface hover:border-primary/60'
          "
          @click="chapterFilter = 'all'"
        >
          {{ $t('general.label.all') }}
        </button>
        <button
          v-for="c in RIGHTS_CHAPTERS"
          :key="c.key"
          type="button"
          class="rounded-full border-2 px-3.5 py-1.5 text-xs font-semibold uppercase tracking-wide transition-colors inline-flex items-center gap-1.5"
          :class="
            chapterFilter === c.key
              ? 'bg-primary border-primary text-primary-contrast'
              : 'border-surface hover:border-primary/60'
          "
          @click="chapterFilter = c.key"
        >
          <i class="pi text-xs" :class="RIGHTS_CHAPTER_ICONS[c.key]" />
          <span v-if="c.charterTitle">{{ c.charterTitle }} ·</span>
          {{ $t(c.labelKey) }}
        </button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
        <button
          v-for="r in visibleRights"
          :key="r.key"
          type="button"
          class="text-left rounded-lg border-2 p-4 transition-colors relative"
          :class="[
            scenario.rights.includes(r.key)
              ? 'border-primary bg-primary/5'
              : 'border-surface hover:border-primary/60',
            disabled ? 'cursor-not-allowed opacity-70' : 'cursor-pointer',
          ]"
          :disabled="disabled"
          :aria-pressed="scenario.rights.includes(r.key)"
          @click="toggle(r.key)"
        >
          <div class="flex items-start gap-2 mb-1 pr-6">
            <span
              class="mt-0.5 shrink-0 inline-flex items-center justify-center w-5 h-5 rounded-full bg-emphasis"
            >
              <i class="pi text-xs text-muted-color" :class="RIGHTS_CHAPTER_ICONS[r.chapter]" />
            </span>
            <span
              class="font-medium text-sm leading-snug"
              :class="scenario.rights.includes(r.key) ? 'text-primary' : ''"
            >
              {{ $t(r.labelKey) }}
            </span>
          </div>
          <div class="text-xs text-muted-color ml-7 mb-1.5">{{ r.article }}</div>
          <p class="text-xs text-muted-color leading-relaxed ml-7">{{ $t(r.descriptionKey) }}</p>
          <i
            v-if="scenario.rights.includes(r.key)"
            class="pi pi-check-circle text-primary absolute top-3 right-3"
          />
        </button>
      </div>
    </div>
  </FriaScenarioSection>
</template>

<script setup>
import FriaScenarioSection from './FriaScenarioSection.vue'
import { RIGHTS_CHAPTER_ICONS } from '@/sections/project/config/friaScenario'
import { FUNDAMENTAL_RIGHTS, RIGHTS_CHAPTERS } from '@/sections/project/config/friaTaxonomies'
import { computed, ref } from 'vue'

const props = defineProps({
  number: { type: [Number, String], required: true },
  scenario: { type: Object, required: true },
  update: { type: Function, required: true },
  disabled: { type: Boolean, default: false },
})

const chapterFilter = ref('all')
const visibleRights = computed(() =>
  chapterFilter.value === 'all'
    ? FUNDAMENTAL_RIGHTS
    : FUNDAMENTAL_RIGHTS.filter(r => r.chapter === chapterFilter.value),
)

function toggle(key) {
  const list = props.scenario.rights
  props.update({ rights: list.includes(key) ? list.filter(k => k !== key) : [...list, key] })
}
</script>
