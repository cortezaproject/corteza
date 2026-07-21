<template>
  <Dialog
    v-model:visible="visible"
    modal
    :header="phase === 'pick' ? $t('project.connectionCreate.title') : configureHeader"
    :style="{ width: phase === 'pick' ? '60rem' : '40rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <!-- Phase 1: pick a connector from the live library ---------------------- -->
    <template v-if="phase === 'pick'">
      <p class="text-sm text-muted-color mb-3">
        {{ $t('project.connectionCreate.pickBlurb') }}
      </p>

      <IconField class="mb-3">
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="query"
          :placeholder="$t('project.connectorPicker.searchPlaceholder')"
          fluid
          autofocus
        />
      </IconField>

      <div class="max-h-[26rem] overflow-auto p-1">
        <div class="grid grid-cols-2 lg:grid-cols-3 gap-3">
          <button
            v-for="c in filtered"
            :key="c.id"
            type="button"
            class="flex items-start gap-3 p-3 rounded-xl border border-surface hover:border-primary hover:bg-emphasis transition-colors text-left"
            @click="onPick(c)"
          >
            <span
              class="inline-flex items-center justify-center w-9 h-9 rounded-md ring-1 ring-surface bg-emphasis shrink-0"
            >
              <i :class="[c.icon, 'text-lg text-primary']" />
            </span>
            <div class="min-w-0">
              <div class="font-medium text-sm leading-tight">{{ c.label }}</div>
              <div class="text-xs text-muted-color line-clamp-2">{{ c.description }}</div>
            </div>
          </button>
        </div>
        <div v-if="!filtered.length" class="text-sm text-muted-color italic text-center py-6">
          {{ $t('project.connectorPicker.noResults', { query }) }}
        </div>
      </div>
    </template>

    <!-- Phase 2: configure the picked connector ----------------------------- -->
    <template v-else>
      <!-- While the base connection is being imported we don't yet know the
           auth-field schema, so show a spinner in place of the form. -->
      <div v-if="preparing" class="flex justify-center py-10">
        <ProgressSpinner style="width: 2.5rem; height: 2.5rem" />
      </div>

      <template v-else>
        <p class="text-sm text-muted-color mb-4">
          {{ $t('project.configureConnection.blurb', { connector: connectorLabel }) }}
        </p>

        <div class="flex flex-col gap-4">
          <CFormGroup :label="$t('project.configureConnection.name')" required>
            <div>
              <InputText
                v-model="name"
                fluid
                :placeholder="connectorLabel"
                :invalid="submitted && !!nameError"
                autofocus
              />
              <ValidationMessage :message="submitted ? nameError : ''" />
            </div>
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
      </template>
    </template>

    <template #footer>
      <!-- Phase 1 has no back target; keep the footer balanced with a spacer. -->
      <Button
        v-if="phase === 'configure'"
        :label="$t('general.label.back')"
        icon="pi pi-arrow-left"
        severity="secondary"
        text
        size="small"
        :disabled="saving"
        @click="backToPick"
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
          v-if="phase === 'configure'"
          :label="$t('project.connectionCreate.create')"
          size="small"
          :loading="saving"
          :disabled="preparing"
          @click="onCreate"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { connector } from '@/sections/project/config/connectors'
import { kindConfig } from '@/sections/project/config/kinds'
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

// 'pick' = the connector grid, 'configure' = the name + params form.
const phase = ref('pick')
const query = ref('')

const cfg = kindConfig('connection')
// The library carries no UI icon; resolve one from the static catalog by
// catalogID, falling back to the generic connection icon.
const iconForConnector = id => connector(id)?.icon || cfg.icon

// Picker offers the live library, mapped to the picker's item shape and
// filtered to the Resource Management whitelist. With no whitelist declared
// yet, the full library is offered.
const pickerItems = computed(() => {
  const allowed = store.allowedConnectorIds(props.project.projectID)
  return store.connectionLibrary
    .filter(c => !allowed || allowed.has(c.catalogID))
    .map(c => ({
      id: c.catalogID,
      catalogID: c.catalogID,
      connectionID: c.connectionID,
      label: c.label,
      description: c.description,
      icon: iconForConnector(c.catalogID),
    }))
})

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return pickerItems.value
  return pickerItems.value.filter(
    c => c.label.toLowerCase().includes(q) || (c.description || '').toLowerCase().includes(q),
  )
})

// The imported base connection: { connectionID, catalogID, derivedParams, meta }.
const connection = ref(null)
const preparing = ref(false)
const name = ref('')
const paramValues = reactive({})
const saving = ref(false)
const submitted = ref(false)

const connectorLabel = computed(
  () =>
    connector(connection.value?.catalogID)?.label ||
    connection.value?.meta?.short ||
    connection.value?.handle ||
    t('project.connections.title'),
)

const configureHeader = computed(() => t('project.connectionCreate.title'))

// Auth/config fields the connector declares, de-duplicated by name.
const params = computed(() => {
  const seen = new Set()
  return (connection.value?.derivedParams || []).filter(p => {
    if (seen.has(p.name)) return false
    seen.add(p.name)
    return true
  })
})

const nameError = computed(() =>
  name.value.trim() ? '' : t('project.connectionCreate.nameRequired'),
)
const isValid = computed(() => !nameError.value)

// Reset to the picker whenever the dialog opens; ensure the library is loaded
// so the grid has something to show even if the step hasn't fetched it yet.
watch(visible, async open => {
  if (!open) return
  phase.value = 'pick'
  query.value = ''
  submitted.value = false
  connection.value = null
  name.value = ''
  Object.keys(paramValues).forEach(k => delete paramValues[k])
  if (!store.connectionLibrary.length) {
    try {
      await store.loadConnectionLibrary()
    } catch (err) {
      $toast.toastErrorHandler(t('project.connections.toastLoadFailed'))(err)
    }
  }
})

// Picking imports the real connection (so we have its auth-field schema), then
// switches to the configure form. The spinner lives in this dialog, not on the
// button that launched the flow.
async function onPick(item) {
  phase.value = 'configure'
  preparing.value = true
  submitted.value = false
  connection.value = null
  try {
    connection.value = await store.prepareConnection(item)
    // Seed the form once the schema is known: name defaults to the connector
    // label (editable), params to their declared defaults.
    name.value = connectorLabel.value
    Object.keys(paramValues).forEach(k => delete paramValues[k])
    for (const p of params.value) {
      paramValues[p.name] = p.default || ''
    }
  } catch (err) {
    phase.value = 'pick'
    $toast.toastErrorHandler(t('project.configureConnection.toastImportFailed'))(err)
  } finally {
    preparing.value = false
  }
}

function backToPick() {
  phase.value = 'pick'
  submitted.value = false
}

function collectParams() {
  const out = []
  for (const p of params.value) {
    const value = paramValues[p.name] || ''
    if (value) out.push({ scope: p.scope, name: p.name, value })
  }
  return out
}

// Persist the configured connection and close; never auto-open the detail dialog.
async function onCreate() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    const entry = await store.saveConnection(props.project.projectID, {
      connection: connection.value,
      name: name.value.trim(),
      config: { params: collectParams() },
    })
    emit('created', entry?.id)
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.connectionCreate.toastCreateFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
