<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="isEdit ? $t('project.configureAutomation.editTitle') : $t('project.configureAutomation.createTitle')"
    :style="{ width: '40rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      {{ $t('project.configureAutomation.blurb') }}
    </p>

    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('project.configureAutomation.name')" required>
        <InputText
          v-model="name"
          fluid
          :placeholder="$t('project.configureAutomation.namePlaceholder')"
          @keyup.enter="save"
        />
      </CFormGroup>

      <CFormGroup :label="$t('project.configureAutomation.description')">
        <Textarea v-model="description" fluid auto-resize rows="3" />
      </CFormGroup>
    </div>

    <template #footer>
      <CRouterLinkButton
        v-if="isEdit"
        :to="{ name: 'taq.builder-edit', params: { id: automation.id } }"
        target="_blank"
        rel="noopener"
        :label="$t('project.configureAutomation.openBuilder')"
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
          outlined
          size="small"
          @click="visible = false"
        />
        <Button
          :label="isEdit ? $t('general.label.save') : $t('project.configureAutomation.create')"
          size="small"
          :loading="saving"
          :disabled="!name.trim()"
          @click="save"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { useProjectsStore } from '@/sections/project/stores/projects'
import { components } from '@planetcrust/human-vue'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { CRouterLinkButton } = components

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  projectId: { type: [String, Number], required: true },
  // Existing automation to edit; null opens the dialog in create mode.
  automation: { type: Object, default: null },
})
const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const isEdit = computed(() => !!props.automation?.id)

const name = ref('')
const description = ref('')
const saving = ref(false)

// (Re)seed the form whenever the dialog opens.
watch(visible, open => {
  if (!open) return
  name.value = props.automation?.name || ''
  description.value = props.automation?.description || ''
})

// Create just creates the TAQ; editing only updates name/description. The full
// builder is reached via the "Open builder" link (edit mode), never a
// programmatic redirect.
async function save() {
  if (!name.value.trim()) return
  saving.value = true
  try {
    if (isEdit.value) {
      await store.updateAutomation(props.projectId, props.automation.id, {
        name: name.value.trim(),
        description: description.value.trim(),
      })
      emit('saved')
    } else {
      const id = await store.addAutomation(props.projectId, {
        name: name.value.trim(),
        description: description.value.trim(),
      })
      emit('saved', id)
    }
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.configureAutomation.toastFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
