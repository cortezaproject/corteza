<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '34rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-end gap-2' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <span
          class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
          :class="[cfg.bg, cfg.ring]"
        >
          <i :class="[cfg.icon, cfg.text]" />
        </span>
        <div class="min-w-0">
          <div class="text-[10px] uppercase tracking-wider text-muted-color leading-none mb-0.5">
            {{ $t('project.module.label') }}
          </div>
          <div class="font-semibold truncate leading-tight">{{ draft.name || $t('project.moduleCreate.title') }}</div>
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-5">
      <CFormGroup :label="$t('general.label.name')" required>
        <div>
          <InputText
            v-model="draft.name"
            size="small"
            fluid
            :invalid="submitted && !!nameError"
            autofocus
          />
          <ValidationMessage :message="submitted ? nameError : ''" />
        </div>
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
      <Button :label="$t('project.moduleCreate.create')" size="small" :loading="saving" @click="onSave" />
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { fieldName } from '@/sections/project/utils/fields'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
})

const emit = defineEmits(['update:modelValue', 'created'])

const store = useProjectsStore()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const cfg = kindConfig('module')

// --- Draft (staged) edits; nothing persists until Create ----------------------
// The create dialog is deliberately minimal — a name is all a new module needs;
// description and fields are edited afterwards in the detail dialog.
const draft = reactive({ name: '' })
// Once a new module is created on Save, remember its id so a re-save (before
// the dialog closes) retries against it instead of creating another copy.
const createdId = ref(null)

function initDraft() {
  createdId.value = null
  submitted.value = false
  draft.name = ''
}

// --- Validation ----------------------------------------------------------------
// Checked on Create; invalid states only show after the first attempt so a
// fresh dialog isn't covered in red.
const submitted = ref(false)

const nameError = computed(() => {
  const name = draft.name.trim()
  if (!name) return t('project.module.nameRequired')
  // The handle is the slugified title; a different module slugging to the same
  // handle would collide on create, so flag it here and make the user rename.
  const key = fieldName(name).toLowerCase()
  const clash = store.resourcesFor(props.project?.id).find(m => {
    if (m.kind !== 'module') return false
    if (m.id === createdId.value) return false
    return fieldName(m.name).toLowerCase() === key
  })
  return clash ? t('project.module.nameClash', { name: clash.name }) : ''
})

const isValid = computed(() => !nameError.value)

watch(
  () => props.modelValue,
  open => {
    if (open) initDraft()
  },
  { immediate: true },
)

// --- Commit -------------------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    let id = createdId.value
    if (!id) {
      // createdId guards against duplicates when the create partially fails and
      // the user re-saves: the module exists, so the retry reuses it.
      id = await store.addResource(props.project.projectID, { kind: 'module', name: draft.name.trim() })
      createdId.value = id
    }
    emit('created', id)
    visible.value = false
  } catch (err) {
    // Stay open so the staged name isn't lost.
    $toast.toastErrorHandler(t('project.moduleCreate.toastCreateFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
