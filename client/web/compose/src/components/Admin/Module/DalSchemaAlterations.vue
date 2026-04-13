<template>
  <Dialog
    v-model:visible="showModal"
    :header="$t('module.edit.schemaAlterations.title')"
    :style="{ width: '80vw', maxWidth: '1200px' }"
    :breakpoints="{ '1199px': '90vw', '575px': '95vw' }"
    :pt="{ content: { class: 'p-0 pb-4' } }"
    modal
  >
    <div v-if="loading" class="flex justify-center items-center py-12">
      <ProgressSpinner />
    </div>

    <div v-else-if="!sortedAlterations.length" class="flex justify-center items-center py-12 text-muted-color">
      {{ $t('module.edit.schemaAlterations.noAlterations') }}
    </div>

    <DataTable
      v-else
      :value="sortedAlterations"
      size="small"
      responsiveLayout="scroll"
      :rowClass="rowClassFunc"
      class="border-t border-surface"
      @row-mouseenter="(e) => dependOnHover = e.data.dependsOn"
      @row-mouseleave="dependOnHover = undefined"
    >
      <Column :header="$t('module.edit.schemaAlterations.columns.alteration')" field="alterationID" />

      <Column :header="$t('module.edit.schemaAlterations.columns.change')" style="max-width: 300px">
        <template #body="{ data }">
          {{ stringifyParams(data.params) }}
        </template>
      </Column>

      <Column :header="$t('module.edit.schemaAlterations.columns.status')" alignHeader="center">
        <template #body="{ data }">
          <div class="flex justify-center">
            <Badge v-if="data.error" severity="danger" :value="data.error" />
            <Badge v-else-if="data.completedAt" severity="success" :value="$t('module.edit.schemaAlterations.resolved')" />
            <Badge v-else-if="data.dependsOn" severity="secondary" :value="$t('module.edit.schemaAlterations.waitingFor', { id: data.dependsOn })" />
          </div>
        </template>
      </Column>

      <Column headerStyle="width: 200px" bodyClass="text-right">
        <template #body="{ data }">
          <ProgressSpinner v-if="data.processing" style="width: 20px; height: 20px;" strokeWidth="4" />
          <div v-else-if="!data.completedAt" class="flex justify-end gap-2">
            <Button
              :label="$t('general.label.resolve')"
              size="small"
              :disabled="!canResolve(data) || processing"
              @click.stop="confirmResolve(data)"
            />
            <Button
              :label="$t('general.label.dismiss')"
              severity="secondary"
              text
              size="small"
              :disabled="!canDismiss(data) || processing"
              @click.stop="confirmDismiss(data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="canResolveAlterations ? $t('general.label.cancel') : $t('general.label.close')"
          severity="secondary"
          text
          :disabled="processing"
          @click="showModal = false"
        />
        <Button
          v-if="canResolveAlterations"
          :label="$t('module.edit.schemaAlterations.resolveAuto')"
          :loading="processing"
          @click="confirmResolve()"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirmDelete } from '@cortezaproject/corteza-vue-next'

const props = defineProps({
  modal: {
    type: Boolean,
    default: false,
  },
  module: {
    type: Object,
    required: true,
  },
  batch: {
    type: Array,
    default: undefined,
  },
})

const emit = defineEmits(['update:modal'])
const { t } = useI18n()
const { confirmDelete: confirm } = useConfirmDelete()
const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const showModal = ref(false)
const loading = ref(false)
const processing = ref(false)
const dependOnHover = ref(undefined)
const alterations = ref([])

watch(
  () => props.modal,
  (val) => {
    showModal.value = val
  },
  { immediate: true }
)

watch(showModal, (val) => {
  emit('update:modal', val)
})

watch(
  () => props.batch,
  (newBatch) => {
    if (newBatch && newBatch.length) {
      load(...newBatch)
    }
  },
  { deep: true, immediate: true }
)

const sortedAlterations = computed(() => {
  return [...alterations.value].sort((a, b) => 
    (a.batchID || '').localeCompare(b.batchID || '') || 
    (a.dependsOn || '').localeCompare(b.dependsOn || '')
  )
})

const canResolveAlterations = computed(() => {
  return sortedAlterations.value.some(canResolve)
})

const rowClassFunc = (data) => {
  return {
    'bg-emphasis': data.alterationID === dependOnHover.value,
  }
}

async function load(...batchID) {
  if (!batchID || (batchID && !batchID.length)) {
    alterations.value = []
    return
  }

  loading.value = true
  try {
    const { set } = await $SystemAPI.dalSchemaAlterationList({ batchID })
    alterations.value = set || []
    if (alterations.value.length) {
      showModal.value = true
    }
  } catch (e) {
    if ($toast && $toast.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.schemaAlterations.notification.load.error'))(e)
  } finally {
    loading.value = false
  }
}

function confirmResolve(alteration) {
  confirm({
    message: t('module.edit.schemaAlterations.confirmResolve', 'Are you sure you want to apply this schema alteration?'),
    header: t('general.label.resolve'),
    onConfirm: () => onResolve(alteration),
  })
}

function confirmDismiss(alteration) {
  confirm({
    message: t('module.edit.schemaAlterations.confirmDismiss', 'Are you sure you want to dismiss this schema alteration?'),
    header: t('general.label.dismiss'),
    onConfirm: () => onDismiss(alteration),
  })
}

async function onDismiss(alteration) {
  processing.value = true
  const list = alteration ? [alteration] : alterations.value
  const alterationID = []

  list.forEach(a => {
    alterationID.push(a.alterationID)
    a.processing = true
  })

  try {
    await $SystemAPI.dalSchemaAlterationDismiss({ alterationID })
    if ($toast && $toast.toastSuccess) $toast.toastSuccess(t('module.edit.schemaAlterations.notification.dismiss.success'))
  } catch (e) {
    if ($toast && $toast.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.schemaAlterations.notification.dismiss.error'))(e)
  } finally {
    list.forEach(a => { a.processing = false })
    load(...(props.batch || []))
    processing.value = false
  }
}

async function onResolve(alteration) {
  processing.value = true
  const list = alteration ? [alteration] : alterations.value
  const alterationID = []

  list.forEach(a => {
    alterationID.push(a.alterationID)
    a.processing = true
  })

  try {
    await $SystemAPI.dalSchemaAlterationApply({ alterationID })
    if ($toast && $toast.toastSuccess) $toast.toastSuccess(t('module.edit.schemaAlterations.notification.resolve.success'))
  } catch (e) {
    if ($toast && $toast.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.schemaAlterations.notification.resolve.error'))(e)
  } finally {
    list.forEach(a => { a.processing = false })
    load(...(props.batch || []))
    processing.value = false
  }
}

function canDismiss(alteration) {
  if (alteration.completedAt) return false
  if (alteration.dependsOn) {
    return alterations.value.some(a => a.alterationID === alteration.dependsOn && !a.completedAt)
  }
  return true
}

function canResolve(alteration) {
  return canDismiss(alteration)
}

function stringifyParams(params) {
  try {
    if (params.attributeAdd) return stringifyAttributeAddParams(params.attributeAdd)
    if (params.attributeDelete) return stringifyAttributeDeleteParams(params.attributeDelete)
    if (params.attributeReType) return stringifyAttributeReTypeParams(params.attributeReType)
    if (params.attributeReEncode) return stringifyAttributeReEncodeParams(params.attributeReEncode)
    if (params.modelAdd) return stringifyModelAddParams(params.modelAdd)
    if (params.modelDelete) return stringifyModelDeleteParams(params.modelDelete)
  } catch(e) {}
  return t('module.edit.schemaAlterations.unknownType', 'Unknown alteration type')
}

function stringifyAttributeAddParams({ attr = {} }) {
  return t('module.edit.schemaAlterations.params.attribute.add', { ident: attr.ident, storeType: attr.store?.type, attrType: attr.type?.type })
}

function stringifyAttributeDeleteParams({ attr = {} }) {
  return t('module.edit.schemaAlterations.params.attribute.delete', { ident: attr.ident, storeType: attr.store?.type })
}

function stringifyAttributeReTypeParams({ attr = {}, to = {} }) {
  return t('module.edit.schemaAlterations.params.attribute.reType', { ident: attr.ident, toType: to.type })
}

function stringifyAttributeReEncodeParams({ attr = {}, to = {} }) {
  return t('module.edit.schemaAlterations.params.attribute.reEncode', { ident: attr.ident, toType: to.type })
}

function stringifyModelAddParams({ model = {} }) {
  return t('module.edit.schemaAlterations.params.model.add', { ident: model.ident })
}

function stringifyModelDeleteParams({ model = {} }) {
  return t('module.edit.schemaAlterations.params.model.delete', { ident: model.ident })
}
</script>
