<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="$t('project.automationCreate.title')"
    :style="{ width: '40rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-end gap-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      {{ $t('project.automationCreate.blurb') }}
    </p>

    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('general.label.name')" required>
        <div>
          <InputText
            v-model="draft.name"
            fluid
            :invalid="submitted && !!nameError"
            :placeholder="$t('project.automationCreate.namePlaceholder')"
            autofocus
            @keyup.enter="onCreate(false)"
          />
          <ValidationMessage :message="submitted ? nameError : ''" />
        </div>
      </CFormGroup>

      <CFormGroup :label="$t('general.label.description')">
        <Textarea v-model="draft.description" fluid auto-resize rows="3" />
      </CFormGroup>
    </div>

    <template #footer>
      <Button
        :label="$t('general.label.cancel')"
        severity="secondary"
        text
        size="small"
        @click="visible = false"
      />
      <Button
        :label="$t('project.automationCreate.create')"
        outlined
        size="small"
        :loading="saving && !openingBuilder"
        :disabled="saving"
        @click="onCreate(false)"
      />
      <Button
        :label="$t('project.automationCreate.createAndOpenBuilder')"
        icon="pi pi-external-link"
        icon-pos="right"
        size="small"
        :loading="saving && openingBuilder"
        :disabled="saving"
        @click="onCreate(true)"
      />
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
})
const emit = defineEmits(['update:modelValue', 'created'])

const store = useProjectsStore()
const { t } = useI18n()
const router = useRouter()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

// --- Draft (minimal create form) ------------------------------------------------
const draft = reactive({ name: '', description: '' })

// Reset the form whenever the dialog opens.
watch(visible, open => {
  if (!open) return
  submitted.value = false
  draft.name = ''
  draft.description = ''
})

// --- Validation -----------------------------------------------------------------
const submitted = ref(false)

const nameError = computed(() =>
  draft.name.trim() ? '' : t('project.automationCreate.nameRequired'),
)

const isValid = computed(() => !nameError.value)

// --- Create ---------------------------------------------------------------------
// Create persists the automation and closes; it never auto-opens the detail
// dialog. `openBuilder` additionally opens the new automation in the full TAQ
// builder (new tab, so the wizard stays put).
const saving = ref(false)
// Which create button is in flight, so only it shows a spinner.
const openingBuilder = ref(false)

async function onCreate(openBuilder = false) {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  openingBuilder.value = openBuilder
  // Open the tab synchronously inside the click so it isn't blocked as a popup
  // after the await; we point it at the builder once the automation exists, or
  // close it if the create fails.
  const builderTab = openBuilder ? window.open('', '_blank') : null
  try {
    const id = await store.addAutomation(props.project.projectID, {
      name: draft.name.trim(),
      description: draft.description.trim(),
    })
    if (builderTab) {
      const { href } = router.resolve({ name: 'taq.builder-edit', params: { id } })
      builderTab.location.href = new URL(href, window.location.href).href
    }
    emit('created', id)
    visible.value = false
  } catch (err) {
    builderTab?.close()
    $toast.toastErrorHandler(t('project.automationCreate.toastCreateFailed'))(err)
  } finally {
    saving.value = false
    openingBuilder.value = false
  }
}
</script>
