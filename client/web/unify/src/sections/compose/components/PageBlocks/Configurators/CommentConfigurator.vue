<template>
  <div class="flex flex-col gap-4">
    <CFormGroup :label="$t('block.general.module')">
      <Select
        v-model="options.moduleID"
        :options="modules"
        option-label="name"
        option-value="moduleID"
        :placeholder="$t('block.comment.module.placeholder')"
        filter
        class="w-full"
      />
    </CFormGroup>

    <template v-if="selectedModule">
      <CFormGroup :label="$t('block.recordList.record.prefilterLabel')">
        <CInputExpression
          ref="filterInput"
          v-model="options.filter"
          dialect="ql"
          :scope="scope"
          :query-fields="queryFields"
          :placeholder="$t('block.recordList.record.prefilterPlaceholder')"
        />
        <CExpressionHint :scope="scope" @insert="filterInput?.insert($event)" />
      </CFormGroup>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <CFormGroup
          :label="$t('block.comment.titleField.label')"
          :description="$t('block.comment.titleField.footnote')"
        >
          <Select
            v-model="options.titleField"
            :options="stringFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.comment.contentField.label')"
          :description="$t('block.comment.contentField.footnote')"
        >
          <Select
            v-model="options.contentField"
            :options="stringFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.comment.replyField.label')"
          :description="$t('block.comment.replyField.footnote')"
        >
          <Select
            v-model="options.replyField"
            :options="recordFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.comment.referenceField.label')"
          :description="$t('block.comment.referenceField.footnote')"
        >
          <Select
            v-model="options.referenceField"
            :options="recordFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.comment.attachmentField.label')"
          :description="$t('block.comment.attachmentField.footnote')"
        >
          <Select
            v-model="options.attachmentField"
            :options="fileFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.comment.sortDirection.label')"
          :description="$t('block.comment.sortDirection.footnote')"
        >
          <Select
            v-model="options.sortDirection"
            :options="sortDirections"
            option-label="label"
            option-value="value"
            class="w-full"
          />
        </CFormGroup>

        <CFormGroup
          :label="$t('block.comment.reactionsField.label')"
          :description="$t('block.comment.reactionsField.footnote')"
        >
          <Select
            v-model="options.reactionsField"
            :options="stringFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
        </CFormGroup>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@planetcrust/human-vue'
import { useExpressionScope } from '@/sections/compose/composables/useExpressionScope'

const { t } = useI18n()

const props = defineProps({
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const block = inject('blockDraft')

const moduleStore = useModuleStore()

const options = computed(() => block.value.options)

const modules = computed(() => moduleStore.set || [])

const selectedModule = computed(() =>
  modules.value.find(m => m.moduleID === options.value.moduleID),
)

const filterInput = ref(null)
const { scope, queryFields } = useExpressionScope({
  page: computed(() => props.page),
  queryModule: computed(() => selectedModule.value),
})

const moduleFields = computed(() => {
  if (!selectedModule.value) return []
  return [...selectedModule.value.fields].sort((a, b) =>
    (a.label || a.name).localeCompare(b.label || b.name),
  )
})

const stringFields = computed(() =>
  moduleFields.value.filter(f => f.kind === 'String' && !f.isMulti),
)

const recordFields = computed(() => moduleFields.value.filter(f => f.kind === 'Record'))

const fileFields = computed(() => moduleFields.value.filter(f => f.kind === 'File'))

const sortDirections = computed(() => [
  { label: t('block.comment.sortDirection.asc'), value: 'asc' },
  { label: t('block.comment.sortDirection.desc'), value: 'desc' },
])

function fieldLabel(f) {
  return `${f.label || f.name} (${f.kind})`
}

// Auto-detect common field names when module changes
watch(
  () => options.value.moduleID,
  () => {
    if (!selectedModule.value) return
    options.value.titleField = ''
    options.value.contentField = ''
    options.value.referenceField = ''
    options.value.attachmentField = ''

    moduleFields.value.forEach(f => {
      if (f.name === 'Content') options.value.contentField = 'Content'
      if (f.name === 'Reference') options.value.referenceField = 'Reference'
      if (f.name === 'Attachments') options.value.attachmentField = 'Attachments'
    })
  },
)

// Default sort direction
if (!options.value.sortDirection) {
  options.value.sortDirection = 'desc'
}
</script>
