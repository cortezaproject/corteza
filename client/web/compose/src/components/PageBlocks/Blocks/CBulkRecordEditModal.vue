<template>
  <Dialog
    :visible="visible"
    modal
    :header="modalTitle || $t('block.recordList.bulkRecord.title')"
    :style="{ width: '50vw' }"
    :breakpoints="{ '960px': '75vw', '641px': '90vw' }"
    @update:visible="$emit('update:visible', $event)"
    @hide="onModalHide"
  >
    <div class="flex flex-col gap-4 mt-2">
      <div
        v-for="(fieldName, index) in fields"
        :key="fieldName"
        class="flex flex-col gap-1 relative"
      >
        <div class="flex items-center gap-2 mb-1">
          <label class="text-sm font-medium text-primary">
            {{ getFieldLabel(getField(fieldName)) }}
          </label>
          <Button
            icon="pi pi-trash"
            text
            severity="danger"
            size="small"
            class="h-6 w-6 p-0"
            @click="fields.splice(index, 1)"
          />
        </div>
        <CFieldEditor
          :field="getField(fieldName)"
          :namespace="namespace"
          :module="module"
          :model-value="getFieldValue(fieldName)"
          @update:model-value="setFieldValue(fieldName, $event)"
        />
      </div>

      <Divider v-if="fields.length" class="!m-0" />

      <Select
        v-model="selectedField"
        :options="availableModuleFields"
        option-label="label"
        option-value="name"
        :placeholder="getFieldSelectorPlaceholder"
        class="w-full"
        @update:model-value="addField"
        filter
      >
        <template #option="{ option }">
          {{ option.label || option.name }}
        </template>
      </Select>
    </div>

    <template #footer>
      <div class="flex flex-col w-full gap-4 mt-2">
        <div class="flex justify-between items-center w-full">
          <Button
            :label="$t('general.label.reset')"
            severity="secondary"
            text
            :disabled="processing"
            @click="onReset"
          />
          <div class="flex gap-2">
            <Button
              :label="$t('general.label.cancel')"
              severity="secondary"
              text
              @click="$emit('update:visible', false)"
            />
            <Button
              :label="$t('general.label.save')"
              severity="primary"
              :disabled="!fields.length || processing"
              :loading="processing"
              @click="handleBulkUpdate"
            />
          </div>
        </div>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { ref, computed, watch, inject } from 'vue'
import { compose } from '@cortezaproject/corteza-js-next'
import { components } from '@cortezaproject/corteza-vue-next'
import { useI18n } from 'vue-i18n'

const { CFieldEditor } = components

const props = defineProps({
  visible: { type: Boolean, required: true },
  namespace: { type: Object, required: true },
  module: { type: Object, required: true },
  selectedFields: { type: Array, default: () => [] },
  initialRecord: { type: Object, default: () => ({}) },
  modalTitle: { type: String, default: '' },
  query: { type: String, default: '' },
})

const emit = defineEmits(['update:visible', 'close', 'save'])
const { t } = useI18n()

const $ComposeAPI = inject('$ComposeAPI')
const $toast = inject('$toast')

const processing = ref(false)
const selectedField = ref(null)
const fields = ref([])
const record = ref(new compose.Record(props.module || {}, {}))

// Initialize
watch(
  () => props.visible,
  isVisible => {
    if (isVisible && props.module) {
      record.value = new compose.Record(props.module, props.initialRecord)

      if (props.selectedFields.length > 0) {
        fields.value = [...props.selectedFields]
      } else {
        fields.value = []
      }
    }
  },
  { immediate: true },
)

const moduleFields = computed(() => {
  if (!props.module || !props.module.fields) return []
  const mFields = [...props.module.fields].sort((a, b) =>
    (a.label || a.name).localeCompare(b.label || b.name),
  )
  const ownedBy = props.module.systemFields
    ? props.module.systemFields().find(f => f.name === 'ownedBy')
    : null
  if (ownedBy) mFields.push(ownedBy)

  return mFields.filter(f => isFieldEditable(f))
})

const availableModuleFields = computed(() => {
  return moduleFields.value.filter(f => !fields.value.includes(f.name))
})

const getFieldSelectorPlaceholder = computed(() => {
  const isAnother = fields.value.length ? 'Another' : ''
  // Use fallback if translations are not present yet
  return (
    t(`block.recordList.bulkRecord.field.add${isAnother}`) ||
    'Add field' + (isAnother ? ' another' : '')
  )
})

function onModalHide() {
  fields.value = []
  if (props.module) {
    record.value = new compose.Record(props.module, {})
  }
  emit('close')
}

function getFieldLabel(f) {
  if (!f) return ''
  return f.label || t(`field.system.${f.name}`, f.name || f.kind)
}

function getField(fieldName) {
  return moduleFields.value.find(f => f.name === fieldName) || {}
}

function addField(fieldName) {
  if (!fieldName) return
  fields.value.push(fieldName)
  selectedField.value = null
}

function getFieldValue(fieldName) {
  if (!record.value) return undefined
  return record.value.values[fieldName]
}

function setFieldValue(fieldName, value) {
  if (!record.value) return
  record.value.setValue(fieldName, value)
}

function onReset() {
  if (props.module) {
    record.value = new compose.Record(props.module, props.initialRecord)
  }
  fields.value = [...props.selectedFields]
}

function isFieldEditable(field) {
  if (!field) return false
  const canUpdateRecordValue = field.canUpdateRecordValue !== false
  if (!canUpdateRecordValue) return false

  if (field.isSystem) {
    if (field.name === 'ownedBy') {
      return true
    }
    return false
  }
  return !field.expressions?.value
}

async function handleBulkUpdate() {
  if (!fields.value.length || props.query === undefined || props.query === null) return

  processing.value = true
  const values = []

  // Serialize fields
  const mockRecord = new compose.Record(props.module, record.value)
  const serialized = mockRecord.serializeValues()

  fields.value.forEach(fieldName => {
    const f = getField(fieldName)
    // Add all values that match the field
    const fieldValues = serialized.filter(v => v.name === fieldName)
    if (fieldValues.length === 0) {
      values.push({ name: fieldName, value: '' })
    } else {
      values.push(...fieldValues)
    }
  })

  try {
    await $ComposeAPI.recordPatch({
      moduleID: props.module.moduleID,
      namespaceID: props.namespace.namespaceID,
      query: props.query,
      values,
    })
    $toast?.toastSuccess?.(
      t('notification.record.bulkRecordUpdateSuccess', 'Bulk update successful'),
    )
    emit('update:visible', false)
    emit('save')
  } catch (error) {
    console.error('Failed to bulk update records', error)
    $toast?.toastDanger?.(
      t('notification.record.bulkRecordUpdateFailed', 'Failed to bulk update records'),
    )
  } finally {
    processing.value = false
  }
}
</script>
