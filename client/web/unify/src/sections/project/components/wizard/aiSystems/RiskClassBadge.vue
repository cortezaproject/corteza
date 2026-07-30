<template>
  <span
    class="text-xs px-2 py-0.5 rounded-full font-medium whitespace-nowrap"
    :class="RISK_CLASS_BADGE_CLASSES[riskClass] || UNCLASSIFIED_CLASSES"
  >
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  riskClass: { type: String, default: null },
})

const { t } = useI18n()

// Same worst→least warm family already used across this section (EventBadge's
// SEVERITY tints, RiskPips, and config/friaScenario.js's FRIA_SEVERITY_BADGE_
// CLASSES): red → orange → amber → yellow, from Tailwind's own predefined
// palette. No new hexes, so config/chartColors.js's standing dataviz-validator
// rule does not come into play — this is the same "worst is red" concept those
// already validated.
const RISK_CLASS_BADGE_CLASSES = {
  prohibited: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300',
  high: 'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-300',
  limited: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
  minimal: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/15 dark:text-yellow-300',
}

// Deliberately NEUTRAL, not the `minimal` tint. An unclassified AI system has
// not been assessed; showing it in the lowest-risk colour would state a
// conclusion nobody reached.
const UNCLASSIFIED_CLASSES = 'bg-surface text-muted-color border border-surface'

const label = computed(() =>
  props.riskClass
    ? t(`fria.riskClasses.${props.riskClass}.label`)
    : t('fria.aiSystems.unclassified'),
)
</script>
