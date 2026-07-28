<template>
  <!-- One event in the log. Reads as a sentence — who did what to which resource
       — and then answers the question the view exists for, inline: what changed.
       Everything else is deliberately quiet so the diff is the only thing with a
       voice.

       No connector rail: a rail asserts that consecutive rows are one thread,
       and they are not — an agent created at 12:02 and a user updated at 12:30
       share nothing but a day. Order is already carried by the day headings.

       Severity is not tinted for routine events either: in practice nearly every
       row is `info`, so colouring them all would spend the row's one accent on
       noise and leave a genuine error with nothing left to shout with. -->
  <!-- The whole row toggles the details disclosure (when `collapsible`) — the
       chevron next to the title is the visible affordance, the row is the hit
       target. Inner interactive elements (resource filter, more-changes,
       links) stop propagation so they don't double as a toggle. The row owns
       its own padding — no negative-margin bleed, so it can never be clipped
       by (or poke out of) whatever container hosts it. -->
  <div
    class="group flex gap-3 rounded-md px-2 py-2 transition-colors hover:bg-emphasis motion-reduce:transition-none"
    :class="{ 'cursor-pointer': collapsible }"
    @click="collapsible && (open = !open)"
  >
    <!-- KindIcon owns the kind icon-square recipe (bg + ring + coloured glyph),
         so an event about an agent wears the same mark that agent wears
         everywhere else in the project. Project-dashboard resources (project,
         backlog items, the five event categories) wear their own category/
         backlog/project badges via eventBadge(); everything else outside the
         project's kinds (settings, queues, templates…) falls through to
         kindConfig's neutral fallback, which is honest: they are not project
         resources. -->
    <KindIcon :config="iconConfig" size="lg" class="mt-0.5" />

    <div class="min-w-0 flex-1 flex flex-col gap-1">
      <!-- The sentence. Clicking the resource type filters by it. break-words
           keeps long unbroken tokens (raw resource strings, snowflake IDs)
           from forcing the row — and narrow containers like the overview's
           activity card — to overflow horizontally. -->
      <p class="text-sm text-color break-words">
        <UserCell v-if="actorName" :name="actorName" avatar class="align-middle" />
        <span v-else class="text-muted-color">{{ t('project.dashboard.activity.system') }}</span>
        <span class="mx-1">{{ verb }}</span>
        <span
          v-tooltip.top="data.resource"
          role="button"
          tabindex="0"
          class="font-medium cursor-pointer rounded-sm hover:underline focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          :class="resourceActive ? 'text-primary' : 'text-primary-500'"
          @click.stop="emit('filter', 'resource', typeValue)"
          @keydown.enter.prevent="emit('filter', 'resource', typeValue)"
          @keydown.space.prevent="emit('filter', 'resource', typeValue)"
        >
          {{ typeLabel }}
        </span>
        <span v-if="subject" class="font-medium">&nbsp;{{ subject }}</span>
        <!-- Disclosure chevron — sits right after the sentence so it reads as
             part of the row's headline. Stops propagation: the row itself also
             toggles, and both firing would cancel out. -->
        <button
          v-if="collapsible"
          type="button"
          class="ml-1.5 shrink-0 align-middle cursor-pointer rounded-sm text-muted-color hover:text-color focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          :aria-expanded="open"
          :aria-controls="detailsID"
          :aria-label="t('project.dashboard.activity.details')"
          @click.stop="open = !open"
        >
          <i class="pi text-xs" :class="open ? 'pi-chevron-up' : 'pi-chevron-down'" />
        </button>
      </p>

      <!-- The signature: what actually changed, without a click. -->
      <EventDiff v-if="changeCount" :delta="data.delta" :limit="open ? 0 : INLINE_DIFF_LIMIT" />

      <button
        v-if="collapsible && hiddenChanges"
        type="button"
        class="self-start text-xs text-primary-500 hover:underline cursor-pointer rounded-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
        @click.stop="open = true"
      >
        {{ t('project.dashboard.activity.moreChanges', hiddenChanges) }}
      </button>

      <!-- Error is the exception colour is saved for. -->
      <p v-if="data.error" class="text-xs text-red-600 break-all dark:text-red-400">
        {{ data.error }}
      </p>

      <!-- Provenance last: you read what happened, then when and from where. -->
      <p class="flex flex-wrap items-center gap-x-2 gap-y-0.5 text-xs text-muted-color">
        <span v-tooltip.top="absoluteTime">{{ relativeTime }}</span>
        <template v-if="originText">
          <span aria-hidden="true">·</span>
          <span>{{ originText }}</span>
        </template>
      </p>

      <!-- Everything the row deliberately did not say. One disclosure, not three:
           meta and the request plumbing are the same kind of thing — evidence you
           only want when you are already investigating. -->
      <div
        v-if="open"
        :id="detailsID"
        role="region"
        :aria-label="t('project.dashboard.activity.details')"
        class="grid grid-cols-1 gap-x-4 gap-y-1 pt-1 sm:grid-cols-2"
        @click.stop
      >
        <div v-for="row in detailRows" :key="row.k" class="flex gap-2 min-w-0 text-xs">
          <span class="text-muted-color shrink-0">{{ row.k }}</span>
          <span class="font-mono break-all">{{ row.v }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { filters } from '@planetcrust/human-vue'
import {
  actionVerb,
  originLabel,
  resourceID,
  resourceType,
  resourceTypeLabel,
} from '@/sections/admin/views/system/ActionLog/vocab'
import { kindConfig } from '@/sections/project/config/kinds'
import { eventBadge, eventKind, severityTone } from '@/sections/project/config/eventKinds'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import EventDiff from './EventDiff.vue'
import UserCell from './UserCell.vue'

const { locFullDateTime } = filters
const { t, locale } = useI18n()

// Three fields answers the common update without turning the row into a table.
const INLINE_DIFF_LIMIT = 3

// Why this disclosure is hand-rolled rather than PrimeVue's Accordion
// (AccordionPanel/AccordionHeader/AccordionContent, as used in taq's
// ReferencePanel.vue):
//
// Accordion's model is "the header IS the toggle" — AccordionHeader renders a
// <button> around its default slot. This row's leading line contains its own
// interactive element (the resource type, which filters by it), and an
// interactive descendant of a <button> is invalid HTML and breaks focus order.
// ReferencePanel gets away with it because its header is only an icon + labels.
//
// The escapes both cost more than they return: `asChild` hands back
// {class, active, a11yAttrs, onClick} and we wire it ourselves anyway, while a
// chevron-only AccordionHeader needs pt overrides to strip Aura's 1.125rem
// padding / 600 weight / full-width button, plus the open state lifted to the
// list as a value array.
//
// So we borrow the part that matters — the aria-expanded + aria-controls/id
// contract — and skip the wrapper. Motion is skipped deliberately: expansion is
// instant, which is what you want when scanning an audit log.

const props = defineProps({
  data: { type: Object, required: true },
  // Resolved display name of the actor; empty for system-originated events.
  actorName: { type: String, default: '' },
  // Currently active resource filter, so the sentence can show it as engaged.
  activeResource: { type: String, default: '' },
  // Compact hosts (the overview's recent-activity card) turn the disclosure
  // off entirely: no chevron, no row toggle, no more-changes reveal — the row
  // is a static summary and the full view is one "View all" away.
  collapsible: { type: Boolean, default: true },
})

const emit = defineEmits(['filter'])

const open = ref(false)

// Pairs the toggle with its region (aria-controls -> id). actionID is unique per
// event, so it is stable across re-renders and re-sorts.
const detailsID = computed(() => `event-details-${props.data.actionID}`)

const typeValue = computed(() => resourceType(props.data.resource))

// The kind's own icon + colours, repainted only when the event went wrong. The
// tone spreads OVER the config so the icon survives — a failed agent update
// should still look like an agent, just alarmed.
const iconConfig = computed(() => {
  const base = eventBadge(typeValue.value) || kindConfig(eventKind(typeValue.value))
  const tone = severityTone(props.data.severity)
  return tone ? { ...base, ...tone } : base
})
const typeLabel = computed(() => resourceTypeLabel(props.data.resource))
const verb = computed(() => actionVerb(props.data.action))
const originText = computed(() =>
  props.data.requestOrigin ? originLabel(props.data.requestOrigin) : '',
)
const resourceActive = computed(() => props.activeResource === typeValue.value)
const changeCount = computed(() => (props.data.delta || []).length)
const hiddenChanges = computed(() =>
  open.value ? 0 : Math.max(0, changeCount.value - INLINE_DIFF_LIMIT),
)

// The event's own name for the affected resource. Meta is keyed `<thing>.name` /
// `<thing>.handle` / `<thing>.email` (see the *_actions.yaml props), so only ever
// take the key whose prefix matches THIS resource: a record event carries both
// `record.ID` and `module.name`, and borrowing the latter would label the row
// "Record Module A". Falls back to the trailing ID, which is unhelpful but true.
const NAME_FIELDS = ['name', 'handle', 'email', 'slug']

const subject = computed(() => {
  const meta = props.data.meta || {}
  const short = typeValue.value.split(':').pop() || ''
  const camel = short.replace(/-([a-z])/g, (_, c) => c.toUpperCase())

  for (const field of NAME_FIELDS) {
    const v = meta[`${camel}.${field}`] || meta[`${short}.${field}`]
    if (v) return v
  }

  return resourceID(props.data.resource)
})

const absoluteTime = computed(() => locFullDateTime(props.data.timestamp))

// Relative time via Intl — locale-aware and dependency-free. Coarse by design:
// an audit log is read recency-first, and "3 days ago" is more useful at a
// glance than a timestamp you have to subtract in your head.
const relativeTime = computed(() => {
  const ts = props.data.timestamp ? new Date(props.data.timestamp) : null
  if (!ts || Number.isNaN(ts.getTime())) return ''

  const diff = Math.round((ts.getTime() - Date.now()) / 1000)
  const abs = Math.abs(diff)
  const rtf = new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' })

  if (abs < 60) return rtf.format(Math.round(diff), 'second')
  if (abs < 3600) return rtf.format(Math.round(diff / 60), 'minute')
  if (abs < 86400) return rtf.format(Math.round(diff / 3600), 'hour')
  if (abs < 2592000) return rtf.format(Math.round(diff / 86400), 'day')
  // Beyond a month, a date is easier to place than "2 months ago".
  return absoluteTime.value
})

// Meta the sentence already spoke is dropped — repeating `user.email` under a row
// that names the user is noise.
const metaRows = computed(() => {
  const meta = props.data.meta || {}
  return Object.keys(meta)
    .filter(k => !NAME_FIELDS.some(f => k.endsWith(`.${f}`)))
    .map(k => ({ k, v: String(meta[k]) }))
    .filter(r => r.v !== '')
})

const detailRows = computed(() => {
  const d = props.data
  const base = [
    { k: t('system.actionlog.list.details.id'), v: d.actionID },
    { k: t('system.actionlog.list.details.timestamp'), v: absoluteTime.value },
    { k: t('system.actionlog.list.details.resource'), v: d.resource },
    { k: t('system.actionlog.list.details.action'), v: d.action },
    { k: t('system.actionlog.list.details.actorID'), v: d.actorID },
    { k: t('system.actionlog.list.details.actorIPAddr'), v: d.actorIPAddr },
    { k: t('system.actionlog.list.details.requestID'), v: d.requestID },
  ]
  return [...base, ...metaRows.value].filter(r => r.v !== undefined && r.v !== null && r.v !== '')
})
</script>
