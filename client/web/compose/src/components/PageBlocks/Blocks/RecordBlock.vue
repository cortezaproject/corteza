<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full p-5">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else-if="!fieldModule" class="p-5 text-muted-color italic">
      {{ $t('block.record.noModule') }}
    </div>

    <div v-else-if="!record" class="p-5 text-muted-color italic">
      {{ $t('block.record.noRecord') }}
    </div>

    <div v-else class="p-5" :class="layoutClass">
      <div
        v-for="field in visibleFields"
        :key="field.fieldID || field.name"
        class="field-item"
        :class="fieldContainerClass"
      >
        <label class="text-sm font-semibold text-primary mb-1.5 block">
          {{ field.label || field.name }}
        </label>

        <div class="field-value text-color">
          <CFieldViewer
            v-if="field.canReadRecordValue !== false"
            :field="field"
            :record="record"
            :namespace="namespace"
          />
          <span v-else class="text-muted-color italic text-sm">
            {{ $t('field.noPermission') }}
          </span>
        </div>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { components } from '@cortezaproject/corteza-vue-next'
const { CFieldViewer } = components
import { useModuleStore } from '@/stores/module'
import { useRecordStore } from '@/stores/record'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  page: {
    type: Object,
    default: () => ({}),
  },
})

const route = useRoute()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const loading = ref(false)
const record = ref(null)

const options = computed(() => props.block.options || {})

// Resolve the module: use block's referenceModuleID if set, otherwise fall back to the page's moduleID
const fieldModule = computed(() => {
  const moduleID = options.value.referenceModuleID || props.page?.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// Determine which fields to display
const visibleFields = computed(() => {
  if (!fieldModule.value) return []

  const configuredFields = options.value.fields || []

  if (configuredFields.length === 0) {
    // No fields configured — show all module fields (non-system)
    return fieldModule.value.fields || []
  }

  // Show only configured fields in order
  return fieldModule.value.filterFields
    ? fieldModule.value.filterFields(configuredFields)
    : configuredFields
})

// Layout class based on block options
const layoutClass = computed(() => {
  const layout = options.value.recordFieldLayoutOption || 'default'
  const classes = {
    default: 'flex flex-col gap-5',
    noWrap: 'flex gap-6',
    wrap: 'grid grid-cols-2 gap-5',
  }
  return classes[layout] || classes.default
})

const fieldContainerClass = computed(() => {
  const layout = options.value.recordFieldLayoutOption || 'default'
  if (layout === 'noWrap') return 'min-w-[20rem]'
  return ''
})

// Load the record when module or route changes
async function loadRecord() {
  if (!fieldModule.value) return

  const recordID = route.params?.recordID
  if (!recordID) return

  loading.value = true

  try {
    record.value = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: fieldModule.value.moduleID,
      recordID,
    })
  } catch (e) {
    console.error('Failed to load record for Record block:', e)
    record.value = null
  } finally {
    loading.value = false
  }
}

watch(
  () => [fieldModule.value?.moduleID, route.params?.recordID],
  () => loadRecord(),
  { immediate: true },
)
</script>
