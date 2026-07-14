<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="$t('project.roleCreate.title')"
    :style="{ width: '30rem' }"
    :pt="{ footer: { class: 'flex justify-end gap-2' } }"
  >
    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('general.label.name')" required>
        <div>
          <InputText
            v-model="draft.name"
            fluid
            :placeholder="$t('project.roleCreate.namePlaceholder')"
            :invalid="submitted && !!nameError"
            autofocus
            @keyup.enter="onCreate"
          />
          <ValidationMessage :message="submitted ? nameError : ''" />
        </div>
      </CFormGroup>

      <CFormGroup :label="$t('general.label.description')">
        <Textarea
          v-model="draft.description"
          rows="3"
          auto-resize
          fluid
          :placeholder="$t('project.roleCreate.descriptionPlaceholder')"
        />
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
        :label="$t('project.roleCreate.create')"
        size="small"
        :loading="saving"
        @click="onCreate"
      />
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
})

const emit = defineEmits(['update:modelValue', 'created'])

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

// --- Draft ---------------------------------------------------------------------
const draft = reactive({ name: '', description: '' })

// Reseed the (minimal) form on every open.
watch(
  () => props.modelValue,
  open => {
    if (!open) return
    submitted.value = false
    draft.name = ''
    draft.description = ''
  },
)

// --- Validation ----------------------------------------------------------------
const submitted = ref(false)

const nameError = computed(() => (draft.name.trim() ? '' : t('project.roleCreate.nameRequired')))

const isValid = computed(() => !nameError.value)

// --- Create --------------------------------------------------------------------
const saving = ref(false)

async function onCreate() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    const id = await store.addRole(props.project.projectID, {
      name: draft.name.trim(),
      description: draft.description,
    })
    emit('created', id)
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.roleCreate.toastCreateFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
