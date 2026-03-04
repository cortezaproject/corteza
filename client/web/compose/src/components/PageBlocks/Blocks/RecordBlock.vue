<template>
  <PageBlock :block="block">
    <div v-if="loading" class="flex items-center justify-center h-full p-5">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else-if="!fieldModule" class="p-5 text-muted-color italic">
      {{ $t('block.record.noModule') }}
    </div>

    <div v-else-if="!activeRecord" class="p-5 text-muted-color italic">
      {{ $t('block.record.noRecord') }}
    </div>

    <div v-else class="p-4 overflow-y-auto" :class="layoutClass">
      <div
        v-for="field in visibleFields"
        :key="field.fieldID || field.name"
        class="field-item"
        :class="fieldContainerClass"
      >
        <label
          v-if="field.kind !== 'Bool' || field.options?.switch || !isEditing"
          class="text-sm font-semibold text-primary mb-1.5 block"
        >
          {{ field.label || field.name }}
          <span v-if="field.isRequired" class="text-red-500">*</span>
        </label>

        <div class="field-value text-color">
          <!-- Editor (edit/create mode) -->
          <template v-if="isEditing && field.canReadRecordValue !== false">
            <FormField :name="field.name" v-slot="{ invalid, error }">
              <CFieldEditor
                :field="field"
                :namespace="namespace"
                :model-value="getFieldValue(field)"
                @update:model-value="setFieldValue(field, $event)"
              />
              <Message v-if="invalid" severity="error" size="small" variant="simple">
                {{ error?.message }}
              </Message>
            </FormField>
          </template>

          <!-- Viewer (view mode) -->
          <CFieldViewer
            v-else-if="field.canReadRecordValue !== false"
            :field="field"
            :record="activeRecord"
            :namespace="namespace"
          />

          <span v-else class="text-muted-color italic text-sm">
            {{ $t('block.field.noPermission') }}
          </span>
        </div>
      </div>
    </div>
  </PageBlock>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { components } from '@cortezaproject/corteza-vue-next'
const { CFieldViewer, CFieldEditor } = components
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

// Inject edit context from RecordView (may be null on non-record pages)
const ctx = inject('recordViewContext', null)

const loading = ref(false)
const localRecord = ref(null)

const options = computed(() => props.block.options || {})

// Whether we're in edit or create mode
const isEditing = computed(() => {
  return !!ctx && ctx.mode.value !== 'view'
})

// Resolve the module: use block's referenceModuleID if set, otherwise fall back to the page's moduleID
const fieldModule = computed(() => {
  const moduleID = options.value.referenceModuleID || props.page?.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// The active record: use context record in edit/create mode, otherwise local
const activeRecord = computed(() => {
  if (isEditing.value) {
    return ctx.record.value
  }
  return localRecord.value
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
  if (fieldModule.value.filterFields) {
    return fieldModule.value.filterFields(configuredFields)
  }
  // Fallback: filter module fields to only those listed in block options
  const names = configuredFields.map(f => f.name ?? f)
  return (fieldModule.value.fields || []).filter(f => names.includes(f.name))
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

// Field value helpers for edit mode
function getFieldValue(field) {
  const r = ctx?.record?.value
  if (!r) return field.isMulti ? [] : ''
  const val = r.values[field.name]
  if (val === undefined || val === null) return field.isMulti ? [] : ''
  return val
}

function setFieldValue(field, value) {
  const r = ctx?.record?.value
  if (!r) return
  r.setValue(field.name, value)
}

// Load the record when in view mode and module or route changes
async function loadRecord() {
  // In edit/create mode, record comes from context — no need to load
  if (isEditing.value) return
  if (!fieldModule.value) return

  const recordID = route.params?.recordID
  if (!recordID || recordID === '0') return

  loading.value = true

  try {
    localRecord.value = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: fieldModule.value.moduleID,
      recordID,
    })
  } catch (e) {
    console.error('Failed to load record for Record block:', e)
    localRecord.value = null
  } finally {
    loading.value = false
  }
}

watch(
  () => [fieldModule.value?.moduleID, route.params?.recordID, isEditing.value],
  () => loadRecord(),
  { immediate: true },
)
</script>
