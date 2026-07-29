<template>
  <!-- Icon-only is the same tinted badge with the label taken out, for places
       that carry a status beside something already named — a step row in the
       nav, a tab. The wording moves into the tooltip rather than disappearing. -->
  <span
    v-if="iconOnly"
    v-tooltip.bottom="chipLabel"
    class="inline-flex items-center justify-center w-6 h-6 rounded-md ring-1 shrink-0"
    :class="[preset.bg, preset.ring, preset.text]"
  >
    <i :class="[preset.icon, 'text-xs']" />
  </span>
  <CTag
    v-else
    :icon="preset.icon"
    :bg="preset.bg"
    :ring="preset.ring"
    :text="preset.text"
    :label="chipLabel"
    :small="small"
  />
</template>

<script setup>
// The one status indicator this section renders — wizard step headers, the
// Publish tab and the revision switcher each had their own Tag/pill/severity
// map before this existed, so the same state could read green in one place and
// grey in another. CTag (lib) is the dumb pill; the palette lives here,
// because these statuses are the project section's vocabulary.
//
// Two axes end up in the same tag, deliberately:
//   - REVIEW status, the wizard's per-step and per-publish governance flag
//     (draft / submitted / changes-requested / approved). Session-local — see
//     stores/projects.js, it never persists.
//   - LIFECYCLE status, the project/revision's own persisted `status`
//     (draft / active / published / archived / suspended / deprecated).
// They share a vocabulary at the edges ('draft' means "not reviewed yet" and
// "never published" both), so one table covers both rather than two components
// the caller has to choose between.
//
// `label` overrides the wording without touching the colour — several call
// sites need their own phrasing for a shared state ("Not published" rather than
// the lifecycle "Draft" on the Publish header, "Reviewed" for a finished
// stage), and the point of this component is that they can't drift on COLOUR
// while doing it.
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  status: { type: String, default: '' },
  // Wording override; the status still picks the colour and icon.
  label: { type: String, default: '' },
  // Passed through to CTag — for dropdown items and table cells.
  small: { type: Boolean, default: false },
  // Badge with no wording, the label moving to a tooltip (step nav, tabs).
  iconOnly: { type: Boolean, default: false },
})

const { t, te } = useI18n()

// Alpha tints, matching the idiom the step-nav badges and publish stage
// numbers already use: a /10 fill inside a /30 ring, so both themes work.
const TONES = {
  neutral: { bg: 'bg-emphasis', ring: 'ring-surface', text: 'text-muted-color' },
  info: { bg: 'bg-primary/10', ring: 'ring-primary/30', text: 'text-primary' },
  warn: { bg: 'bg-amber-500/10', ring: 'ring-amber-500/30', text: 'text-amber-500' },
  success: { bg: 'bg-green-500/10', ring: 'ring-green-500/30', text: 'text-green-500' },
}

// The four review states first, then the lifecycle statuses that ride the same
// tones. Add a status here, not another map at a call site.
const PRESETS = {
  draft: { tone: 'neutral', icon: 'pi pi-pencil', labelKey: 'project.status.draft' },
  submitted: {
    tone: 'info',
    icon: 'pi pi-clock',
    labelKey: 'project.governance.status.waitingApproval',
  },
  'changes-requested': {
    tone: 'warn',
    icon: 'pi pi-exclamation-circle',
    labelKey: 'project.governance.status.changesRequested',
  },
  // Approved is green like active but says so: a revision can sit approved for
  // days before anyone publishes it, and calling that "Active" would claim it
  // is live.
  approved: {
    tone: 'success',
    icon: 'pi pi-check-circle',
    labelKey: 'project.governance.status.approved',
  },
  active: { tone: 'success', icon: 'pi pi-bolt', labelKey: 'project.status.active' },
  published: { tone: 'success', icon: 'pi pi-check-circle', labelKey: 'project.status.published' },
  archived: { tone: 'neutral', icon: 'pi pi-inbox', labelKey: 'project.status.archived' },
  suspended: { tone: 'warn', icon: 'pi pi-pause-circle', labelKey: 'project.status.suspended' },
  deprecated: { tone: 'neutral', icon: 'pi pi-ban', labelKey: 'project.status.deprecated' },
}

const preset = computed(() => {
  const cfg = PRESETS[props.status] || {
    tone: 'neutral',
    icon: 'pi pi-circle',
    labelKey: `project.status.${props.status}`,
  }
  return { ...TONES[cfg.tone], icon: cfg.icon, labelKey: cfg.labelKey }
})

// An unknown status falls back to its raw value rather than rendering a
// missing-key path — better a stray enum on screen than 'project.status.xyz'.
const chipLabel = computed(() => {
  if (props.label) return props.label
  const key = preset.value.labelKey
  return te(key) ? t(key) : props.status
})
</script>
