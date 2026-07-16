<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '56rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <KindIcon kind="module" size="lg" plain-icon />
        <div class="min-w-0">
          <DialogEyebrow>{{ $t('project.module.label') }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">{{ draft.name || $t('project.module.unnamed') }}</div>
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-6">
      <!-- Module meta (staged; committed on Save) -->
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
        <CFormGroup :label="$t('general.label.description')">
          <Textarea v-model="draft.description" rows="3" auto-resize fluid />
        </CFormGroup>
      </div>

      <!-- Fields (persist immediately via FieldDialog) -->
      <CFormGroup :label="$t('project.module.fields')">
        <template #actions>
          <Button
            icon="pi pi-plus"
            :label="$t('project.dataModel.addField')"
            severity="secondary"
            size="small"
            @click="createField?.(resourceId)"
          />
        </template>

        <div class="border border-surface rounded-border overflow-hidden mt-1">
          <div v-if="fields.length" class="divide-y divide-surface">
            <div
              v-for="f in fields"
              :key="f.id"
              class="group flex items-center gap-3 px-3 min-h-11 hover:bg-emphasis transition-colors cursor-pointer"
              @click="editField?.(resourceId, f.id)"
            >
              <div class="min-w-0 flex-1 flex items-center gap-2">
                <span class="text-sm truncate">
                  {{ f.name || $t('project.dataModel.untitledField') }}
                </span>
                <FieldKindTag :type="f.type" />
                <template v-if="f.required">
                  <span class="shrink-0 text-muted-color opacity-50 text-[11px]">·</span>
                  <span class="shrink-0 text-[11px] font-medium text-amber-600 dark:text-amber-400">
                    {{ $t('project.dataModel.required') }}
                  </span>
                </template>
                <template v-if="f.multi">
                  <span class="shrink-0 text-muted-color opacity-50 text-[11px]">·</span>
                  <span class="shrink-0 text-[11px] font-medium text-sky-600 dark:text-sky-400">
                    {{ $t('project.dataModel.multiple') }}
                  </span>
                </template>
              </div>
              <Button
                icon="pi pi-trash"
                severity="danger"
                text
                size="small"
                class="opacity-0 focus:opacity-100 group-hover:opacity-100 transition-opacity"
                :aria-label="$t('general.label.remove')"
                :title="$t('general.label.remove')"
                @click.stop="removeField(f)"
              />
            </div>
          </div>
          <div v-else class="px-3 py-6 text-center text-sm text-muted-color italic">
            {{ $t('project.dataModel.emptyFields') }}
          </div>
        </div>
      </CFormGroup>

      <ResourcePermissionsSection :project="project" kind="module" :resource-id="resourceId" />
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="resourceId && project?.hasNamespace"
        :to="{ name: 'admin.modules.edit', params: { slug: project.namespaceID, moduleID: resourceId } }"
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
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="visible = false"
        />
        <Button :label="$t('general.label.save')" size="small" :loading="saving" @click="onSave" />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import FieldKindTag from '@/sections/project/components/datamodel/FieldKindTag.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import ResourcePermissionsSection from '@/sections/project/components/permissions/ResourcePermissionsSection.vue'
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { fieldName } from '@/sections/project/utils/fields'
import { components, useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Module to inspect; this dialog only opens for existing modules (step rows +
  // graph nodes). `resourceId` is the shared detail-dialog contract prop name.
  resourceId: { type: String, default: null },
})

const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

// Field editing is delegated to the shared FieldDialog rendered by the wizard;
// these provides open it (stacked over this dialog) and persist immediately.
const editField = inject('editField', null)
const createField = inject('createField', null)

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const module = computed(() =>
  props.resourceId ? store.resourcesFor(props.project?.projectID).find(r => r.id === props.resourceId) : null,
)
// Live field list — reflects immediate field mutations without closing.
const fields = computed(() => module.value?.fields || [])

// --- Module meta draft (staged; nothing persists until Save) ------------------
const draft = reactive({ name: '', description: '' })

function initDraft() {
  submitted.value = false
  const m = module.value
  draft.name = m?.name || ''
  draft.description = m?.description || ''
}

// --- Validation ---------------------------------------------------------------
const submitted = ref(false)

const nameError = computed(() => {
  const name = draft.name.trim()
  if (!name) return t('project.module.nameRequired')
  // The handle is the slugified title; another module slugging to the same
  // handle would collide, so flag it and make the user rename.
  const key = fieldName(name).toLowerCase()
  const clash = store.resourcesFor(props.project?.projectID).find(m => {
    if (m.kind !== 'module') return false
    if (m.id === props.resourceId) return false
    return fieldName(m.name).toLowerCase() === key
  })
  return clash ? t('project.module.nameClash', { name: clash.name }) : ''
})

const isValid = computed(() => !nameError.value)

watch(
  () => [props.modelValue, props.resourceId],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

// --- Field removal (immediate) ------------------------------------------------
function removeField(f) {
  confirmDelete({
    header: t('project.dataModel.removeField.header'),
    message: t('project.dataModel.removeField.message', {
      name: f.name || t('project.dataModel.untitledField'),
    }),
    onConfirm: () => handleRemoveField(f),
  })
}

async function handleRemoveField(f) {
  try {
    await store.removeField(props.project.projectID, props.resourceId, f.id)
    $toast.toastSuccess(f.name, t('project.dataModel.toast.fieldRemoved'))
  } catch (err) {
    $toast.toastErrorHandler(t('project.dataModel.toast.fieldRemoveFailed'))(err)
  }
}

// --- Commit module meta -------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    await store.updateResource(props.project.projectID, props.resourceId, {
      name: draft.name,
      description: draft.description,
    })
    emit('saved')
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.module.toast.saveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
