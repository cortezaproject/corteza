<template>
  <div class="flex flex-col gap-8">
    <GovernanceForm
      :schema="SUMMARY_SCHEMA"
      :model-value="modelValue"
      :disabled="disabled"
      @update:model-value="$emit('update:modelValue', $event)"
    />

    <!-- Read-only reference: how @Human Governance satisfies each Article 17 QMS aspect. -->
    <section class="rounded-xl border border-surface bg-surface-50 dark:bg-surface-950 overflow-hidden">
      <header class="flex items-start gap-3 px-5 py-4 border-b border-surface">
        <span class="shrink-0 mt-0.5 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
          <i class="pi pi-verified text-lg" />
        </span>
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold uppercase tracking-wide text-muted-color">Article 17</p>
          <h3 class="text-lg font-medium leading-tight">Quality Management System</h3>
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

defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue'])

// A bullet is either a plain string or { text, url } when it carries a link.
const pointText = p => (typeof p === 'string' ? p : p.text)
const pointUrl = p => (typeof p === 'string' ? null : p.url)

// EU AI Act Article 17 — Quality Management System. Each aspect is paired with
// how @Human Governance delivers it. Reference content only (not an input).
const QMS_ARTICLE_17 = [
  {
    label: 'a',
    body: '@Human Governance is an AI Compliance-as-Infrastructure tooling layer to deliver an out-of-the-box strategy for AI regulatory compliance, including compliance with conformity assessment procedures and procedures for the management of modifications to the high-risk AI systems.',
    points: [
      'Conformity Assessment Procedures are documented and implemented in the Conformity Assessment section.',
      'Management of modifications to a high-risk AI system in production must go through a gated change management procedure documented in the Incident, Issue and Change Management section.',
      'By following the full gated development process for AI Systems, EU AI Act Annex IV documentation will be programmatically delivered and updated (after change management).',
    ],
  },
  {
    label: 'b',
    body: 'Techniques, procedures and systematic actions to be used for the design, design control and design verification of high-risk AI systems.',
    points: [],
  },
  {
    label: 'c',
    body: 'Techniques, procedures and systematic actions to be used for the development, quality control and quality assurance of high-risk AI systems.',
    points: [
      '@Human Governance delivers a wizard-like user journey for building high-risk AI systems.',
      'The design and development process for any high-risk AI system project is gated for governance verification at multiple strategic points: 1. Executive Authority is required to sign off the project after the Project Summary, Quality Management Systems and Project Members have been defined; 2. etc.',
      'After an AI System Project is in production, changes can only be made via a request, approval and review system.',
      'Multi-factor authentication is enforced for all Project Members.',
    ],
  },
  {
    label: 'd',
    body: 'Examination, test and validation procedures to be carried out before, during and after the development of the high-risk AI system, and the frequency with which they have to be carried out.',
    points: [],
  },
  {
    label: 'e',
    body: 'Technical specifications, including standards, to be applied and, where the relevant harmonised standards are not applied in full or do not cover all of the relevant requirements set out in Section 2, the means to be used to ensure that the high-risk AI system complies with those requirements.',
    points: [
      { text: '@Human Governance is 100% open-source and can be inspected at', url: 'https://github.com/crusttech' },
      { text: '@Human Governance documentation is 100% open and can be inspected at', url: 'https://docs.planetcrust.com' },
      'High-level technical architecture for any AI System Project is available to download in the project documentation section.',
      'Configuration files for any individual AI System Project are available to download in the project documentation section.',
    ],
  },
  {
    label: 'f',
    body: 'Systems and procedures for data management, including data acquisition, data collection, data analysis, data labelling, data storage, data filtration, data mining, data aggregation, data retention and any other operation regarding the data that is performed before and for the purpose of the placing on the market or the putting into service of high-risk AI systems.',
    points: [
      'A Data Access Layer (DAL) is used to connect to third-party Postgres, MySQL (and variants) and Microsoft SQL Server databases. The DAL may also connect to third-party databases via API. The third-party database schema must be programmatically mapped to the @Human data schema and field types.',
      'An Application Connection layer allows connection to third-party web applications via API using standard REST methods. This can be bidirectional and all connections must be programmatically mapped to the @Human data schema and field types.',
      'All fields in any data model must be labelled for their Sensitivity level.',
      'All local data is stored in a Postgres database. Data at rest should be encrypted.',
      'The Application layer only provides soft-deletes of data. Hard (i.e. permanent) data deletion must be provided at the database layer.',
    ],
  },
  {
    label: 'g',
    body: 'Each AI System Project has its own documented Risk Management System section.',
    points: [],
  },
  {
    label: 'h',
    body: 'The setup, implementation and maintenance of a post-market monitoring system are documented in the Monitoring section.',
    points: [],
  },
  {
    label: 'i',
    body: 'Procedures related to the reporting of a serious incident are documented in the Incident, Issue and Change Management section.',
    points: [],
  },
  {
    label: 'j',
    body: 'The handling of communication with national competent authorities, other relevant authorities, including those providing or supporting the access to data, notified bodies, other operators, customers or other interested parties.',
    points: [],
  },
  {
    label: 'k',
    body: 'Systems and procedures for record-keeping of all relevant documentation and information.',
    points: [],
  },
  {
    label: 'l',
    body: 'Resource management, including security-of-supply related measures.',
    points: ['Documented in the Resource Management section.'],
  },
  {
    label: 'm',
    body: 'An accountability framework setting out the responsibilities of the management and other staff with regard to all the aspects listed in this Quality Management System.',
    points: ['Documented in the Members section.'],
  },
]
</script>
