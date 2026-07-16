<template>
  <div class="flex flex-col gap-8">
    <GovernanceForm
      :schema="SUMMARY_SCHEMA"
      :model-value="modelValue"
      :disabled="disabled"
      @update:model-value="$emit('update:modelValue', $event)"
    />

    <!-- Read-only reference: how @Human Governance satisfies each Article 17 QMS aspect. -->
    <section class="rounded-xl border border-surface bg-surface overflow-hidden">
      <header class="flex items-start gap-3 px-5 py-4 border-b border-surface">
        <span class="shrink-0 mt-0.5 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <i class="pi pi-verified text-lg" />
        </span>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold uppercase tracking-wide text-muted-color">{{ $t('project.qms.article17') }}</p>
          <h3 class="text-lg font-medium leading-tight">{{ $t('project.qms.title') }}</h3>
        </div>
      </header>

      <ol class="divide-y divide-surface">
        <li v-for="item in QMS_ARTICLE_17" :key="item.label" class="px-5 py-4">
          <p class="text-sm leading-relaxed">{{ item.body }}</p>
          <ul v-if="item.points?.length" class="mt-2.5 flex flex-col gap-2">
            <li v-for="(p, i) in item.points" :key="i" class="flex gap-2.5 text-sm text-muted-color leading-relaxed">
              <span class="w-1.5 h-1.5 rounded-full bg-primary/60 mt-[0.45rem] shrink-0" />
              <span>
                {{ pointText(p) }}
                <a
                  v-if="pointUrl(p)"
                  :href="pointUrl(p)"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-primary hover:underline break-all"
                >{{ pointUrl(p) }}</a>
              </span>
            </li>
          </ul>
        </li>
      </ol>
    </section>
  </div>
</template>

<script setup>
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { SUMMARY_SCHEMA } from '@/sections/project/config/summaryForm'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue'])

const { t } = useI18n()

// A bullet is either a plain string or { text, url } when it carries a link.
const pointText = p => (typeof p === 'string' ? p : p.text)
const pointUrl = p => (typeof p === 'string' ? null : p.url)

const qms = (label, point) => t(`project.qms.items.${label}.${point}`)

// EU AI Act Article 17 — Quality Management System. Each aspect is paired with
// how @Human Governance delivers it. Reference content only (not an input). The
// text lives in the locale bundle (project.qms.items.*); the structure and the
// (non-translatable) source URLs stay here.
const QMS_ARTICLE_17 = computed(() => [
  {
    label: 'a',
    body: qms('a', 'body'),
    points: [qms('a', 'point1'), qms('a', 'point2'), qms('a', 'point3')],
  },
  { label: 'b', body: qms('b', 'body'), points: [] },
  {
    label: 'c',
    body: qms('c', 'body'),
    points: [qms('c', 'point1'), qms('c', 'point2'), qms('c', 'point3'), qms('c', 'point4')],
  },
  { label: 'd', body: qms('d', 'body'), points: [] },
  {
    label: 'e',
    body: qms('e', 'body'),
    points: [
      { text: qms('e', 'point1Text'), url: 'https://github.com/crusttech' },
      { text: qms('e', 'point2Text'), url: 'https://docs.planetcrust.com' },
      qms('e', 'point3'),
      qms('e', 'point4'),
    ],
  },
  {
    label: 'f',
    body: qms('f', 'body'),
    points: [
      qms('f', 'point1'),
      qms('f', 'point2'),
      qms('f', 'point3'),
      qms('f', 'point4'),
      qms('f', 'point5'),
    ],
  },
  { label: 'g', body: qms('g', 'body'), points: [] },
  { label: 'h', body: qms('h', 'body'), points: [] },
  { label: 'i', body: qms('i', 'body'), points: [] },
  { label: 'j', body: qms('j', 'body'), points: [] },
  { label: 'k', body: qms('k', 'body'), points: [] },
  { label: 'l', body: qms('l', 'body'), points: [qms('l', 'point1')] },
  { label: 'm', body: qms('m', 'body'), points: [qms('m', 'point1')] },
])
</script>
