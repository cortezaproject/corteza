<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="
      configured?.configurationID
        ? $t('project.configureConnection.editTitle')
        : $t('project.configureConnection.createTitle')
    "
    :style="{ width: '40rem' }"
    :pt="{ content: { class: '!pt-2' } }"
  >
    <p class="text-sm text-muted-color mb-4">
      {{ $t('project.configureConnection.blurb', { connector: connectorLabel }) }}
    </p>

    <div class="flex flex-col gap-4">
      <CFormGroup :label="$t('project.configureConnection.name')" required>
        <InputText v-model="name" fluid :placeholder="connectorLabel" />
      </CFormGroup>

      <CFormGroup
        v-for="param in params"
        :key="param.name"
        :label="param.label || param.name"
        :description="param.description || ''"
        :required="param.required"
      >
        <InputText v-model="paramValues[param.name]" fluid />
      </CFormGroup>

      <p v-if="!params.length" class="text-sm text-muted-color italic">
        {{ $t('project.configureConnection.noParams') }}
      </p>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          outlined
          size="small"
          @click="visible = false"
        />
        <Button
          :label="$t('general.label.save')"
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
import { connector } from '@/sections/project/config/connectors'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  // The imported base connection: { connectionID, catalogID, derivedParams, meta }.
  connection: { type: Object, default: null },
  projectId: { type: [String, Number], required: true },
  // When set, edit that configured connection instead of creating a new one.
  configured: { type: Object, default: null },
})
const emit = defineEmits(['update:modelValue', 'saved'])

const store = useProjectsStore()
const { t } = useI18n()
const toast = useToast()

const visible = computed({
  get: () => props.modelValue,
  set: v => emit('update:modelValue', v),
})

const connectorLabel = computed(
  () =>
    connector(props.connection?.catalogID)?.label ||
    props.connection?.meta?.short ||
    props.connection?.handle ||
    t('project.connections.title'),
)

// Auth/config fields the connector declares, de-duplicated by name.
const params = computed(() => {
  const seen = new Set()
  return (props.connection?.derivedParams || []).filter(p => {
    if (seen.has(p.name)) return false
    seen.add(p.name)
    return true
  })
})

const name = ref('')
const paramValues = reactive({})
const saving = ref(false)

// (Re)seed the form whenever the dialog opens — prefill name with the connector
// label (editable), and param values from the existing config when editing.
watch(visible, open => {
  if (!open) return
  name.value = props.configured?.name || connectorLabel.value
  Object.keys(paramValues).forEach(k => delete paramValues[k])
  const existing = props.configured?.config?.params || []
  for (const p of params.value) {
    paramValues[p.name] = existing.find(e => e.name === p.name)?.value || p.default || ''
  }
})

function collectParams() {
  const out = []
  for (const p of params.value) {
    const value = paramValues[p.name] || ''
    if (value) out.push({ scope: p.scope, name: p.name, value })
  }
  return out
}

async function save() {
  saving.value = true
  try {
    const entry = await store.saveConnection(props.projectId, {
      connection: props.connection,
      configuredConnectionID: props.configured?.configuredConnectionID,
      name: name.value.trim(),
      config: { params: collectParams() },
    })
    emit('saved', entry)
    visible.value = false
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: t('project.configureConnection.toastFailed'),
      detail: err.message,
      life: 4000,
    })
  } finally {
    saving.value = false
  }
}
</script>
