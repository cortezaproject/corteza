<template>
  <Dialog
    v-model:visible="visible"
    modal
    :style="{ width: '40rem' }"
    :pt="{ content: { class: '!pt-2' }, footer: { class: 'flex justify-between gap-2 p-3' } }"
  >
    <template #header>
      <div class="flex items-center gap-2.5 min-w-0">
        <KindIcon kind="connection" size="lg" plain-icon />
        <div class="min-w-0">
          <DialogEyebrow>{{ $t('project.kinds.connection.single') }}</DialogEyebrow>
          <div class="font-semibold truncate leading-tight">
            {{ name || entity?.name || $t('project.connectionDetail.unnamed') }}
          </div>
        </div>
      </div>
    </template>

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

    <template #footer>
      <!-- Connections have no standalone builder route yet — keep the footer
           balanced with a spacer so Cancel/Save stay right-aligned. -->
      <span />

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
          :disabled="preparing"
          @click="onSave"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import DialogEyebrow from '@/sections/project/components/DialogEyebrow.vue'
import KindIcon from '@/sections/project/components/KindIcon.vue'
import ValidationMessage from '@/sections/project/components/ValidationMessage.vue'
import { connector } from '@/sections/project/config/connectors'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  project: { type: Object, required: true },
  // The configured connection to edit; read from the store by this id
  // (= configuredConnectionID). Always set when opened.
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

// The configured connection entry from the store — it carries
// configuredConnectionID (NOT configurationID, which is only on the raw API
// response). Keyed by project id.
const entity = computed(() =>
  props.resourceId
    ? store.connectionsFor(props.project?.projectID).find(c => c.id === props.resourceId)
    : null,
)

// The imported base connection (its auth-field schema), read on open so we can
// render the params form. { connectionID, catalogID, derivedParams, meta }.
const connection = ref(null)
const preparing = ref(false)
const name = ref('')
const paramValues = reactive({})
const saving = ref(false)
const submitted = ref(false)

const connectorLabel = computed(
  () =>
    connector(connection.value?.catalogID || entity.value?.catalogID)?.label ||
    connection.value?.meta?.short ||
    connection.value?.handle ||
    entity.value?.name ||
    t('project.connections.title'),
)

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
  name.value.trim() ? '' : t('project.connectionDetail.nameRequired'),
)
const isValid = computed(() => !nameError.value)

// On open, import the base connection to learn its schema, then seed the form
// from the existing configured connection. This mirrors what ConnectionsStep
// used to do before opening the configure dialog — now owned by the dialog.
async function loadContext() {
  const target = entity.value
  submitted.value = false
  name.value = target?.name || ''
  Object.keys(paramValues).forEach(k => delete paramValues[k])
  if (!target) return
  preparing.value = true
  connection.value = null
  try {
    connection.value = await store.prepareConnection({
      connectionID: target.connectionID,
      catalogID: target.catalogID,
    })
    // Re-seed once the schema is known, but don't clobber edits the user has
    // already made: only fill a param if it's still untouched.
    const existing = target.config?.params || []
    for (const p of params.value) {
      if (paramValues[p.name] === undefined) {
        paramValues[p.name] = existing.find(e => e.name === p.name)?.value || p.default || ''
      }
    }
  } catch (err) {
    $toast.toastErrorHandler(t('project.connectionDetail.toastLoadFailed'))(err)
  } finally {
    preparing.value = false
  }
}

watch(
  () => [props.modelValue, props.resourceId],
  () => {
    if (props.modelValue) loadContext()
  },
  { immediate: true },
)

function collectParams() {
  const out = []
  for (const p of params.value) {
    const value = paramValues[p.name] || ''
    if (value) out.push({ scope: p.scope, name: p.name, value })
  }
  return out
}

async function onSave() {
  if (saving.value) return
  submitted.value = true
  if (!isValid.value) return
  saving.value = true
  try {
    await store.saveConnection(props.project.projectID, {
      connection: connection.value,
      configuredConnectionID: entity.value?.configuredConnectionID,
      name: name.value.trim(),
      config: { params: collectParams() },
    })
    emit('saved')
    visible.value = false
  } catch (err) {
    $toast.toastErrorHandler(t('project.connectionDetail.toastSaveFailed'))(err)
  } finally {
    saving.value = false
  }
}
</script>
