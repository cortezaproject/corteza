<template>
  <Dialog
    v-model:visible="visible"
    modal
    header="Add LLM provider"
    :style="{ width: '52rem' }"
    :pt="{ footer: { class: 'flex justify-end gap-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      Create a new provider here. It becomes available to all projects and is auto-selected when saved.
    </p>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <CFormGroup label="Provider" required>
        <Select
          v-model="form.provider"
          :options="providerOptions"
          option-label="label"
          option-value="value"
          fluid
          placeholder="Select a provider"
        />
      </CFormGroup>

      <CFormGroup label="Name">
        <InputText v-model="form.short" fluid placeholder="e.g. Production Anthropic" />
      </CFormGroup>

      <CFormGroup label="Handle">
        <InputText v-model="form.handle" fluid placeholder="optional, e.g. anthropic-prod" />
      </CFormGroup>

      <CFormGroup label="Prompt URL" class="md:col-span-2">
        <InputText v-model="form.promptURL" fluid placeholder="https://…" />
      </CFormGroup>

      <CFormGroup label="API Key" required class="md:col-span-2">
        <InputText v-model="apiKey" fluid autocomplete="off" placeholder="Secret key" />
      </CFormGroup>
    </div>

    <p v-if="error" class="text-sm text-red-500 mt-3">{{ error }}</p>

    <template #footer>
      <Button label="Cancel" severity="secondary" outlined size="small" @click="visible = false" />
      <Button label="Create" size="small" :loading="saving" :disabled="!canSubmit" @click="submit" />
    </template>
  </Dialog>
</template>

<script setup>
import { system } from '@planetcrust/human-js'
import { computed, inject, reactive, ref, watch } from 'vue'

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

const providerOptions = [
  { label: 'Anthropic', value: 'anthropic' },
  { label: 'Mistral', value: 'mistral' },
  { label: 'Other', value: 'other' },
]

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
    $toast.toastSuccess('LLM provider created')
    emit('created', created)
    visible.value = false
  } catch (e) {
    error.value = e?.message || 'Failed to create provider'
    $toast.toastErrorHandler('Failed to create LLM provider')(e)
  } finally {
    saving.value = false
  }
}
</script>
