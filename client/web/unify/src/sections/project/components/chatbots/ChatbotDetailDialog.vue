<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '52rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <KindIcon kind="chatbot" size="lg" plain-icon />
        <div class="min-w-0">
          <DialogEyebrow>{{ $t('project.kinds.chatbot.single') }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">{{ draft.name || $t('project.chatbotDetail.unnamed') }}</div>
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
            :placeholder="$t('project.chatbotDetail.namePlaceholder')"
            autofocus
            @keyup.enter="onSave"
          />
          <ValidationMessage :message="submitted ? nameError : ''" />
        </div>
      </CFormGroup>

      <CInputToggleCard
        v-model="draft.enabled"
        :label="$t('project.chatbotDetail.enabled')"
        :description="$t('project.chatbotDetail.enabledHint')"
      />

      <ResourcePermissionsSection :project="project" kind="chatbot" :resource-id="resourceId" />
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="props.resourceId"
        :to="{ name: 'chatbot.edit', params: { chatbotID: props.resourceId } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.chatbotDetail.openEditor')"
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
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import ResourcePermissionsSection from '@/sections/project/components/permissions/ResourcePermissionsSection.vue'
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // Chatbot to edit; always set when the dialog is opened.
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


const chatbot = computed(() =>
  props.resourceId ? store.chatbotsFor(props.project?.projectID).find(c => c.id === props.resourceId) : null,
)

// --- Draft (staged; nothing persists until Save) -------------------------------
const draft = reactive({ name: '', enabled: false })

function initDraft() {
  submitted.value = false
  draft.name = chatbot.value?.name || ''
  draft.enabled = !!chatbot.value?.enabled
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
  draft.name.trim() ? '' : t('project.chatbotDetail.nameRequired'),
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
    await store.updateChatbot(props.project.projectID, props.resourceId, {
      name: draft.name.trim(),
      enabled: draft.enabled,
    })
    emit('saved')
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.chatbotDetail.toastSaveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
