<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="$t('project.llmDialog.header')"
    :style="{ width: '52rem' }"
    :pt="{ footer: { class: 'flex justify-end gap-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      {{ $t('project.llmDialog.blurb') }}
    </p>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <CFormGroup :label="$t('project.llmDialog.provider')" required>
        <Select
          v-model="form.provider"
          :options="providerOptions"
          option-label="label"
          option-value="value"
          fluid
          :placeholder="$t('project.llmDialog.providerPlaceholder')"
        />
      </CFormGroup>

      <CFormGroup :label="$t('project.llmDialog.nameLabel')">
        <InputText
          v-model="form.short"
          fluid
          :placeholder="$t('project.llmDialog.namePlaceholder')"
        />
      </CFormGroup>

      <CFormGroup :label="$t('general.label.handle')">
        <InputText
          v-model="form.handle"
          fluid
          :placeholder="$t('project.llmDialog.handlePlaceholder')"
        />
      </CFormGroup>

      <CFormGroup :label="$t('project.llmDialog.promptURL')" class="md:col-span-2">
        <InputText
          v-model="form.promptURL"
          fluid
          :placeholder="$t('project.llmDialog.promptURLPlaceholder')"
        />
      </CFormGroup>

      <CFormGroup :label="$t('project.llmDialog.apiKey')" required class="md:col-span-2">
        <InputText
          v-model="apiKey"
          fluid
          autocomplete="off"
          :placeholder="$t('project.llmDialog.apiKeyPlaceholder')"
        />
      </CFormGroup>
    </div>

    <p v-if="error" class="text-sm text-red-500 mt-3">{{ error }}</p>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        size="small"
        @click="visible = false"
      />
      <Button
        :label="$t('general.label.create')"
        size="small"
        :loading="saving"
        :disabled="!canSubmit"
        @click="submit"
      />
    </template>
  </Dialog>
</template>

<script setup>
import { system } from '@planetcrust/human-js'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue', 'created'])

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const $SystemAPI = inject('$SystemAPI')
const $toast = inject('$toast')

const providerDefaultURLs = {
  mistral: 'https://api.mistral.ai/v1',
  anthropic: 'https://api.anthropic.com',
}

// Provider names are brand names (untranslated); only the generic "Other" is.
const providerOptions = computed(() => [
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Mistral', value: 'mistral' },
  { label: t('project.llmDialog.providerOther'), value: 'other' },
])

const blank = () => ({ provider: 'anthropic', short: '', handle: '', promptURL: '' })
const form = reactive(blank())
const apiKey = ref('')
const saving = ref(false)
const error = ref('')

const canSubmit = computed(() => !!form.provider && !!apiKey.value.trim())

// Reset the form each time the dialog opens, and prefill the default URL.
watch(visible, v => {
  if (v) {
    Object.assign(form, blank())
    form.promptURL = providerDefaultURLs[form.provider] || ''
    apiKey.value = ''
    error.value = ''
  }
})

// Prefill prompt URL when provider changes (only if empty or another default).
watch(
  () => form.provider,
  p => {
    const defaults = Object.values(providerDefaultURLs)
    if (!form.promptURL || defaults.includes(form.promptURL)) {
      form.promptURL = providerDefaultURLs[p] || ''
    }
  },
)

async function submit() {
  if (!canSubmit.value) return
  saving.value = true
  error.value = ''
  try {
    const model = new system.LlmProvider({
      status: 'active',
      provider: form.provider,
      handle: form.handle || undefined,
      meta: { short: form.short },
      config: { temperature: 0.7, promptURL: form.promptURL },
    })
    const created = await $SystemAPI.llmProviderCreate({
      handle: model.handle,
      provider: model.provider,
      status: model.status,
      meta: model.meta,
      config: model.config,
      apiKey: apiKey.value,
    })
    $toast.toastSuccess(t('project.llmDialog.created'))
    emit('created', created)
    visible.value = false
  } catch (e) {
    error.value = e?.message || t('project.llmDialog.createFailed')
    $toast.toastErrorHandler(t('project.llmDialog.createFailedToast'))(e)
  } finally {
    saving.value = false
  }
}
</script>
