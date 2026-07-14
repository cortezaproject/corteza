<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '52rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
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
            {{ $t('project.kinds.agent.single') }}
          </div>
          <div class="font-semibold truncate leading-tight">{{ draft.name || $t('project.agentDetail.unnamed') }}</div>
        </div>
      </div>
    </template>

    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('general.label.name')" required>
        <div>
          <InputText
            v-model="draft.name"
            fluid
            :invalid="submitted && !!nameError"
            :placeholder="$t('project.agentDetail.namePlaceholder')"
            autofocus
            @keyup.enter="onSave"
          />
          <ValidationMessage :message="submitted ? nameError : ''" />
        </div>
      </CFormGroup>

      <CFormGroup :label="$t('general.label.description')">
        <Textarea v-model="draft.description" fluid auto-resize rows="3" />
      </CFormGroup>

      <CInputToggleCard
        v-model="draft.active"
        :label="$t('project.agentDetail.active')"
        :description="$t('project.agentDetail.activeHint')"
      />

      <ResourcePermissionsSection :project="project" kind="agent" :resource-id="resourceId" />
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="props.resourceId"
        :to="{ name: 'agentic.edit', params: { agentID: props.resourceId } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.agentDetail.openEditor')"
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
        <Button
          :label="$t('general.label.save')"
          size="small"
          :loading="saving"
          @click="onSave"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import ResourcePermissionsSection from '@/sections/project/components/permissions/ResourcePermissionsSection.vue'
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Agent to edit; always set when the dialog is opened.
  resourceId: { type: String, default: null },
})
const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const cfg = kindConfig('agent')

const agent = computed(() =>
  props.resourceId ? store.agentsFor(props.project?.id).find(a => a.id === props.resourceId) : null,
)

// --- Draft (staged; nothing persists until Save) -------------------------------
const draft = reactive({ name: '', description: '', active: false })

function initDraft() {
  submitted.value = false
  draft.name = agent.value?.name || ''
  draft.description = agent.value?.description || ''
  draft.active = agent.value?.status === 'active'
}

// (Re)seed the form on every open / target change.
watch(
  () => [props.modelValue, props.resourceId],
  () => {
    if (props.modelValue) initDraft()
  },
  { immediate: true },
)

// --- Validation -----------------------------------------------------------------
const submitted = ref(false)

const nameError = computed(() =>
  draft.name.trim() ? '' : t('project.agentDetail.nameRequired'),
)

const isValid = computed(() => !nameError.value)

// --- Save -----------------------------------------------------------------------
const saving = ref(false)

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    await store.updateAgent(props.project.projectID, props.resourceId, {
      name: draft.name.trim(),
      description: draft.description.trim(),
      status: draft.active ? 'active' : 'inactive',
    })
    emit('saved')
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.agentDetail.toastSaveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
