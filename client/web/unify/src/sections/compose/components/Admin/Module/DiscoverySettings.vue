<template>
  <Dialog
    v-model:visible="showModal"
    :header="discoveryModalTitle"
    :style="{ width: '50vw', maxWidth: '600px' }"
    :breakpoints="{ '1199px': '75vw', '575px': '95vw' }"
    :pt="{
      footer: { class: 'border-t border-surface p-3' },
    }"
    modal
  >
    <div class="flex flex-col gap-4 py-4">
      <CFormGroup>
        <CFieldPicker
          :all-fields="moduleFieldsList"
          :model-value="currentFieldNames"
          list-class="max-h-[32rem]"
          :available-label="$t('field.selector.available')"
          :selected-label="$t('field.selector.selected')"
          :select-all-label="$t('field.selector.selectAll')"
          :unselect-all-label="$t('field.selector.unselectAll')"
          :search-placeholder="$t('general.label.search')"
          :no-items-label="$t('field.no-items-found')"
          @update:model-value="currentFieldNames = $event"
        />
        <template #description>
          {{ $t('module.edit.discoverySettings.description') }}
          If no fields are selected, all fields will be exposed.
        </template>
      </CFormGroup>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          size="small"
          text
          @click="showModal = false"
        />
        <Button :label="$t('general.label.saveAndClose')" size="small" @click="onSave" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { components } from '@planetcrust/human-vue'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CFieldPicker } = components

const props = defineProps({
  modal: {
    type: Boolean,
    default: false,
  },
  module: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['save', 'update:modal'])
const { t } = useI18n()

const showModal = ref(false)
const currentFieldNames = ref([])
const defaultLang = 'en' // Using explicit string or window default

const discoveryModalTitle = computed(() => {
  const handle = props.module?.handle
  return handle
    ? `${t('module.edit.discoverySettings.title')} (${handle})`
    : t('module.edit.discoverySettings.title')
})

const moduleFieldsList = computed(() => {
  return (props.module?.fields || [])
    .filter(f => !f.isSystem)
    .map(f => ({
      name: f.name,
      label: f.label || f.name,
    }))
})

watch(
  () => props.modal,
  val => {
    showModal.value = val
    if (val) {
      loadSettings()
    }
  },
  { immediate: true },
)

watch(showModal, val => {
  emit('update:modal', val)
})

function loadSettings() {
  const discoveryConfig = props.module.config?.discovery || {}
  const privateConfig = discoveryConfig.private || { result: [] }

  const resultItem = privateConfig.result.find(r => r.lang === defaultLang) || { fields: [] }

  currentFieldNames.value = resultItem.fields.filter(name =>
    moduleFieldsList.value.some(f => f.name === name),
  )
}

function onSave() {
  const discovery = {
    public: { result: [] },
    private: {
      result: [
        {
          lang: defaultLang,
          fields: [...currentFieldNames.value],
        },
      ],
    },
    protected: { result: [] },
  }

  emit('save', {
    config: {
      ...(props.module.config || {}),
      discovery,
    },
  })
}
</script>
