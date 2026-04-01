<template>
  <Dialog
    v-model:visible="showModal"
    :header="discoveryModalTitle"
    :style="{ width: '50vw', maxWidth: '600px' }"
    :breakpoints="{ '1199px': '75vw', '575px': '95vw' }"
    modal
  >
    <div class="flex flex-col gap-4 py-4">
      <FormField name="discoveryPrivateFields" class="flex flex-col gap-2">
        <label class="font-medium text-primary">
          {{ $t('module.edit.discoverySettings.private') }}
        </label>
        <MultiSelect
          v-model="currentFields"
          :options="moduleFieldsList"
          optionLabel="label"
          :placeholder="$t('module.edit.discoverySettings.placeholder')"
          display="chip"
          class="w-full"
        />
        <small class="text-muted-color">
          {{ $t('module.edit.discoverySettings.description') }}
        </small>
      </FormField>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2 px-4 py-3 border-t border-surface bg-surface -mx-6 -mb-6 mt-4">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          @click="showModal = false"
        />
        <Button
          :label="$t('general.label.saveAndClose')"
          @click="onSave"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

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
const currentFields = ref([])
const defaultLang = 'en' // Using explicit string or window default

const discoveryModalTitle = computed(() => {
  const handle = props.module?.handle
  return handle
    ? `${t('module.edit.discoverySettings.title')} (${handle})`
    : t('module.edit.discoverySettings.title')
})

const moduleFieldsList = computed(() => {
  return (props.module?.fields || [])
    .filter((f) => !f.isSystem)
    .map((f) => ({
      name: f.name,
      label: f.label || f.name,
    }))
})

watch(
  () => props.modal,
  (val) => {
    showModal.value = val
    if (val) {
      loadSettings()
    }
  },
  { immediate: true }
)

watch(showModal, (val) => {
  emit('update:modal', val)
})

function loadSettings() {
  const discoveryConfig = props.module.config?.discovery || {}
  const privateConfig = discoveryConfig.private || { result: [] }
  
  const resultItem = privateConfig.result.find((r) => r.lang === defaultLang) || { fields: [] }
  
  // Convert names back into object references for the MultiSelect
  currentFields.value = resultItem.fields
    .map((name) => moduleFieldsList.value.find((f) => f.name === name))
    .filter(Boolean)
}

function onSave() {
  // Save discovery back to meta/config
  const newFieldsNames = currentFields.value.map((f) => f.name)
  
  const discovery = {
    public: { result: [] },
    private: {
      result: [
        {
          lang: defaultLang,
          fields: newFieldsNames,
        },
      ],
    },
    protected: { result: [] },
  }

  emit('save', {
    ...props.module.config,
    discovery,
  })
}
</script>
