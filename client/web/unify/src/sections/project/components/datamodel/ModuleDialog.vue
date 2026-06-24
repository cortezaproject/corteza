<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '34rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2' } }"
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
          <div class="font-semibold truncate leading-tight">{{ draft.name || $t('project.module.unnamed') }}</div>
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
            :disabled="readonly"
            :invalid="submitted && !!nameError"
            autofocus
          />
          <ValidationMessage :message="submitted ? nameError : ''" />
        </div>
      </CFormGroup>
      <CFormGroup :label="$t('general.label.description')">
        <Textarea v-model="draft.description" rows="3" auto-resize fluid :disabled="readonly" />
      </CFormGroup>
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="moduleEditId && project?.namespaceID"
        :to="{ name: 'admin.modules.edit', params: { slug: project.namespaceID, moduleID: moduleEditId } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.module.openEditor')"
        icon="pi pi-external-link"
        severity="secondary"
        text
        size="small"
      />
      <span v-else />

      <div class="flex gap-2">
        <Button v-if="readonly" :label="$t('general.label.close')" severity="secondary" text size="small" @click="visible = false" />
        <template v-else>
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="visible = false"
          />
          <Button :label="$t('general.label.save')" size="small" :loading="saving" @click="onSave" />
        </template>
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { fieldName } from '@/sections/project/utils/fields'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Existing module to edit; null opens the dialog in create mode.
  moduleId: { type: String, default: null },
  // Open as read-only (e.g. approver view, or a submitted/approved step).
  readonly: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const store = useProjectsStore()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const cfg = kindConfig('module')

const module = computed(() =>
  props.moduleId ? store.resourcesFor(props.project?.id).find(r => r.id === props.moduleId) : null,
)

// --- Draft (staged) edits; nothing persists until Save -----------------------
// The module dialog only owns name + description — fields are added, edited and
// removed inline on the Data Model step.
const draft = reactive({ name: '', description: '' })
// Once a new module is created on Save, remember its id so a re-save (before
// the dialog closes) updates it instead of creating another copy.
const createdId = ref(null)

function initDraft() {
  createdId.value = null
  submitted.value = false
  const m = module.value
  draft.name = m?.name || ''
  draft.description = m?.description || ''
}

// --- Validation ----------------------------------------------------------------
// Checked on Save; invalid states only show after the first attempt so a
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
    if (m.id === (props.moduleId || createdId.value)) return false
    return fieldName(m.name).toLowerCase() === key
  })
  return clash ? t('project.module.nameClash', { name: clash.name }) : ''
})

const isValid = computed(() => !nameError.value)

// The module that exists on the backend (saved already, or created during this
// dialog session) — only then can we link to the full compose module editor.
const moduleEditId = computed(() => props.moduleId || createdId.value)

watch(
  () => [props.modelValue, props.moduleId],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

// --- Commit -------------------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  const pid = props.project.id
  saving.value = true
  try {
    let id = props.moduleId || createdId.value
    if (!id) {
      // createdId guards against duplicates when a later write fails and the
      // user re-saves: the module exists, so the retry only updates it.
      id = await store.addResource(pid, { kind: 'module', name: draft.name })
      createdId.value = id
    }
    if (!id) return
    await store.updateResource(pid, id, { name: draft.name, description: draft.description })
    visible.value = false
  } catch (err) {
    // Stay open so the staged edits aren't lost.
    $toast.toastErrorHandler(t('project.module.toast.saveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
