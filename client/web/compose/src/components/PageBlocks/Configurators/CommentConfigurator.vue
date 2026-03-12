<template>
  <div class="flex flex-col gap-4">
    <!-- Module -->
    <div class="flex flex-col gap-1">
      <label class="font-semibold text-primary text-sm">
        {{ $t('block.general.module') }}
      </label>
      <Select
        v-model="options.moduleID"
        :options="modules"
        option-label="name"
        option-value="moduleID"
        :placeholder="$t('block.comment.module.placeholder')"
        filter
        class="w-full"
      />
    </div>

    <template v-if="selectedModule">
      <!-- Prefilter -->
      <div class="flex flex-col gap-1">
        <label class="font-semibold text-primary text-sm">
          {{ $t('block.recordList.record.prefilterLabel') }}
        </label>
        <InputText
          v-model.trim="options.filter"
          :placeholder="$t('block.recordList.record.prefilterPlaceholder')"
          class="w-full"
        />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- Title field -->
        <div class="flex flex-col gap-1">
          <label class="font-semibold text-primary text-sm">
            {{ $t('block.comment.titleField.label') }}
          </label>
          <Select
            v-model="options.titleField"
            :options="stringFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
          <small class="text-muted-color">{{ $t('block.comment.titleField.footnote') }}</small>
        </div>

        <!-- Content field -->
        <div class="flex flex-col gap-1">
          <label class="font-semibold text-primary text-sm">
            {{ $t('block.comment.contentField.label') }}
          </label>
          <Select
            v-model="options.contentField"
            :options="stringFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
          <small class="text-muted-color">{{ $t('block.comment.contentField.footnote') }}</small>
        </div>

        <!-- Reply field -->
        <div class="flex flex-col gap-1">
          <label class="font-semibold text-primary text-sm">
            {{ $t('block.comment.replyField.label') }}
          </label>
          <Select
            v-model="options.replyField"
            :options="recordFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
          <small class="text-muted-color">{{ $t('block.comment.replyField.footnote') }}</small>
        </div>

        <!-- Reference field -->
        <div class="flex flex-col gap-1">
          <label class="font-semibold text-primary text-sm">
            {{ $t('block.comment.referenceField.label') }}
          </label>
          <Select
            v-model="options.referenceField"
            :options="recordFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
          <small class="text-muted-color">{{ $t('block.comment.referenceField.footnote') }}</small>
        </div>

        <!-- Attachment field -->
        <div class="flex flex-col gap-1">
          <label class="font-semibold text-primary text-sm">
            {{ $t('block.comment.attachmentField.label') }}
          </label>
          <Select
            v-model="options.attachmentField"
            :options="fileFields"
            :option-label="fieldLabel"
            option-value="name"
            :placeholder="$t('general.label.none')"
            show-clear
            class="w-full"
          />
          <small class="text-muted-color">{{ $t('block.comment.attachmentField.footnote') }}</small>
        </div>

        <!-- Sort direction -->
        <div class="flex flex-col gap-1">
          <label class="font-semibold text-primary text-sm">
            {{ $t('block.comment.sortDirection.label') }}
          </label>
          <Select
            v-model="options.sortDirection"
            :options="sortDirections"
            option-label="label"
            option-value="value"
            class="w-full"
          />
          <small class="text-muted-color">{{ $t('block.comment.sortDirection.footnote') }}</small>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useModuleStore } from '@/stores/module'

const { t } = useI18n()

const props = defineProps({
  block: { type: Object, required: true },
  namespace: { type: Object, default: () => ({}) },
  page: { type: Object, default: () => ({}) },
  record: { type: Object, default: undefined },
})

const moduleStore = useModuleStore()

const options = computed(() => props.block.options)

const modules = computed(() => moduleStore.set || [])

const selectedModule = computed(() =>
  modules.value.find(m => m.moduleID === options.value.moduleID),
)

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
