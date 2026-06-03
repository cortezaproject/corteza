<template>
  <Dialog
    :visible="visible"
    @update:visible="emit('update:visible', $event)"
    modal
    :header="
      mode === 'create'
        ? $t('list.dialog.create.header')
        : $t('builder.generalConfig.header')
    "
    :style="{ width: '700px' }"
    @hide="resetForm"
  >
    <div class="flex flex-col gap-6">
      <CFormGroup
        :label="mode === 'create' ? $t('list.dialog.create.name') : $t('builder.configSidebar.node')"
        required
        input-id="automation-name"
      >
        <InputText
          id="automation-name"
          v-model="form.name"
          :placeholder="mode === 'create' ? $t('list.dialog.create.namePlaceholder') : ''"
          class="w-full"
          :invalid="!!nameError"
          autofocus
          @keyup.enter="handleSubmit"
        />
        <small v-if="nameError" class="text-red-500">{{ nameError }}</small>
      </CFormGroup>

      <CFormGroup
        :label="$t('list.dialog.create.description')"
        input-id="automation-description"
      >
        <Textarea
          id="automation-description"
          v-model="form.description"
          :placeholder="$t('list.dialog.create.descriptionPlaceholder')"
          class="w-full"
          rows="3"
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('builder.generalConfig.runAs')"
        :description="$t('builder.generalConfig.runAsHint')"
      >
        <CInputUser
          v-model="form.runAs"
          :placeholder="$t('builder.generalConfig.runAsPlaceholder')"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup :label="$t('list.dialog.create.labels')">
        <CInputLabel
          v-model="form.labels"
          :placeholder="$t('list.dialog.create.labelsPlaceholder')"
          :create-label="$t('list.dialog.create.labelsCreateNew')"
          :create-dialog-label="$t('list.dialog.create.labelsDialogCreate')"
          :name-label="$t('list.dialog.create.labelsName')"
          :save-btn-label="$t('list.button.create')"
          :cancel-btn-label="$t('list.button.cancel')"
        />
      </CFormGroup>
    </div>
    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('list.button.cancel')"
          severity="secondary"
          outlined
          size="small"
          @click="emit('update:visible', false)"
        />
        <Button
          :label="mode === 'create' ? $t('list.button.create') : $t('builder.save')"
          size="small"
          :disabled="!form.name.trim() || processing"
          :loading="processing"
          @click="handleSubmit"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { components } from '@planetcrust/human-vue'
import { useAutomationStore } from '@/sections/taq/stores/automation'

const { CInputLabel, CInputUser } = components

const props = defineProps({
  visible: { type: Boolean, default: false },
  mode: { type: String, default: 'create' }, // 'create' or 'edit'
  initialName: { type: String, default: '' },
  initialDescription: { type: String, default: '' },
  initialLabels: { type: Object, default: () => ({}) },
  initialRunAs: { type: [String, Number], default: null },
})

const emit = defineEmits(['update:visible', 'saved'])

const { t } = useI18n()
const router = useRouter()
const $AutomationAPI = inject('$AutomationAPI')
const $Auth = inject('$Auth')
const automationStore = useAutomationStore()

function normalizeRunAs(v) {
  return !v || v === '0' || v === 0 ? null : v
}

const form = ref({
  name: props.initialName,
  description: props.initialDescription,
  labels: {},
  runAs: normalizeRunAs(props.initialRunAs),
})

const nameError = ref('')
const processing = ref(false)

watch(
  () => props.visible,
  val => {
    if (val) {
      form.value.name = props.initialName
      form.value.description = props.initialDescription || ''
      form.value.labels = props.initialLabels ? { ...props.initialLabels } : {}
      form.value.runAs = normalizeRunAs(props.initialRunAs)
    }
  },
)

function resetForm() {
  if (props.mode === 'create') {
    form.value.name = ''
    form.value.description = ''
    form.value.runAs = null
  } else {
    form.value.name = props.initialName
    form.value.description = props.initialDescription || ''
    form.value.labels = props.initialLabels ? { ...props.initialLabels } : {}
    form.value.runAs = normalizeRunAs(props.initialRunAs)
  }
  nameError.value = ''
}

function validateForm() {
  nameError.value = ''
  if (!form.value.name.trim()) {
    nameError.value = t('list.dialog.create.nameRequired')
    return false
  }
  return true
}

async function handleSubmit() {
  if (!validateForm()) return

  processing.value = true

  try {
    if (props.mode === 'create') {
      const created = await automationStore.create($AutomationAPI, {
        meta: {
          short: form.value.name.trim(),
          description: form.value.description.trim() || undefined,
        },
        labels: form.value.labels || {},
        enabled: false,
        triggers: [],
        steps: [],
        paths: [],
        ownedBy: $Auth?.user?.userID,
        runAs: form.value.runAs || undefined,
      })
      emit('update:visible', false)
      router.push(`/taq/builder/${created.automationID}`)
    } else {
      // For editing we emit the saved event and expect parent to handle the actual workflow update,
      // because we need to include all current canvas nodes/edges in the API request, which Builder.vue holds.
      emit('saved', {
        name: form.value.name.trim(),
        description: form.value.description.trim() || undefined,
        labels: form.value.labels || {},
        runAs: form.value.runAs || undefined,
      })
      emit('update:visible', false)
    }
  } catch (e) {
    console.error('Failed to create automation:', e)
  } finally {
    processing.value = false
  }
}
</script>
