<template>
  <CFormGroup
    v-if="!namespace"
    :label="$t('system.applications.editor.custom.namespace')"
    :description="$t('system.applications.editor.custom.namespaceDescription')"
  >
    <div data-test-id="custom-namespace">
      <CInputNamespace
        v-model="namespaceID"
        :disabled="disabled"
        @update:model-value="onNamespaceChange"
      />
    </div>
  </CFormGroup>

  <CFormGroup
    :label="$t('system.applications.editor.custom.modules')"
    :description="$t('system.applications.editor.custom.modulesDescription')"
  >
    <div data-test-id="custom-modules">
      <CInputModule
        v-model="moduleIDs"
        multiple
        :namespace-i-d="namespaceID"
        :disabled="disabled || !namespaceID"
      />
    </div>
  </CFormGroup>

  <CFormGroup
    :label="$t('system.applications.editor.custom.writes')"
    :description="$t('system.applications.editor.custom.writesDescription')"
  >
    <MultiSelect
      v-model="writeIDs"
      data-test-id="custom-writes"
      :options="writeOptions"
      option-label="label"
      option-value="value"
      display="chip"
      :pt="CHIPS_WRAP"
      class="w-full"
      :placeholder="$t('system.applications.editor.custom.readOnly')"
      :disabled="disabled || !moduleIDs.length"
    />
  </CFormGroup>

  <CFormGroup
    :label="$t('app.declaration.deletes')"
    :description="$t('app.declaration.deletesDescription')"
  >
    <MultiSelect
      v-model="deleteIDs"
      data-test-id="custom-deletes"
      :options="writeOptions"
      option-label="label"
      option-value="value"
      display="chip"
      :pt="CHIPS_WRAP"
      class="w-full"
      :placeholder="$t('app.declaration.noDeletes')"
      :disabled="disabled || !moduleIDs.length"
    />
  </CFormGroup>

  <CFormGroup
    :label="$t('app.declaration.automations')"
    :description="$t('app.declaration.automationsDescription')"
  >
    <MultiSelect
      v-model="automations"
      data-test-id="custom-automations"
      :options="automationOptions"
      option-label="label"
      option-value="value"
      display="chip"
      :pt="CHIPS_WRAP"
      filter
      class="w-full"
      :loading="loadingAutomations"
      :placeholder="$t('app.declaration.noAutomations')"
      :disabled="disabled"
    />
  </CFormGroup>

  <CFormGroup
    :label="$t('app.declaration.chatbots')"
    :description="$t('app.declaration.chatbotsDescription')"
  >
    <MultiSelect
      v-model="chatbots"
      data-test-id="custom-chatbots"
      :options="chatbotOptions"
      option-label="label"
      option-value="value"
      display="chip"
      :pt="CHIPS_WRAP"
      filter
      class="w-full"
      :loading="loadingChatbots"
      :placeholder="$t('app.declaration.noChatbots')"
      :disabled="disabled"
    />
  </CFormGroup>

  <CFormGroup
    :label="$t('app.declaration.origins')"
    :description="$t('app.declaration.originsDescription')"
    class="md:col-span-2"
  >
    <template v-if="!disabled" #actions>
      <Button
        :label="$t('general.label.add')"
        icon="pi pi-plus"
        severity="secondary"
        size="small"
        data-test-id="custom-origins-add"
        @click="origins.push('')"
      />
    </template>
    <CFormList
      v-model="origins"
      data-test-id="custom-origins"
      :columns="[{ label: $t('app.declaration.origin'), width: '1fr' }]"
      :empty-message="$t('app.declaration.noOrigins')"
      :disabled="disabled"
      fit-width
    >
      <template #row="{ index }">
        <InputText
          v-model="origins[index]"
          size="small"
          class="w-full"
          placeholder="https://cdn.jsdelivr.net"
          :disabled="disabled"
        />
      </template>
    </CFormList>
  </CFormGroup>
</template>

<script setup>
// What a custom page may reach — a namespace, the modules it reads and those it
// may change — set with the namespace and module pickers. The pickers hold
// IDs; the value is names (the namespace's slug, module handles), because a
// page carried to another instance is read back by what it declared.
//
// Shared by the custom application editor and the compose Custom block, so
// both declare the same way. Renders its fields as siblings for the parent's
// own layout.
import { computed, inject, onMounted, ref, watch } from 'vue'
import { components } from '@planetcrust/human-vue'

const { CInputModule, CInputNamespace } = components

// Chips wrap onto further lines rather than cutting a long name short.
const CHIPS_WRAP = {
  labelContainer: { class: 'whitespace-normal' },
  label: { class: 'flex flex-wrap gap-1 whitespace-normal break-words' },
  pcChip: { label: { class: 'whitespace-normal break-words' } },
}

const props = defineProps({
  // { namespace, modules, writes }, with the stored namespaceID and moduleIDs
  // when there are any.
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
  // A namespace the page is bound to ({namespaceID, slug}): no picker is shown
  // and the modules are this namespace's.
  namespace: { type: Object, default: null },
})

const emit = defineEmits(['update:modelValue'])

const $ComposeAPI = inject('$ComposeAPI')
const $AutomationAPI = inject('$AutomationAPI', null)
const $SystemAPI = inject('$SystemAPI')

const namespaceID = ref('')
const namespaceSlug = ref('')
const moduleIDs = ref([])
const writeIDs = ref([])
const deleteIDs = ref([])
const modules = ref([])
// Origins are plain text, no picker behind them, one row each; a blank row is
// one still being added and is left out of the value. The server checks them
// when the page is stored.
const origins = ref([])
// Workflows and TAQs the page may run, by handle; only one that has a handle
// and a manual trigger can be named.
const automations = ref([])
const runnable = ref([])
const loadingAutomations = ref(false)

const automationOptions = computed(() => {
  const known = runnable.value.map(a => ({ value: a.handle, label: a.label }))
  // A declared handle this viewer cannot list still shows, so saving keeps it.
  const missing = automations.value
    .filter(h => !known.some(a => a.value === h))
    .map(h => ({ value: h, label: h }))
  return [...known, ...missing]
})

async function loadRunnable() {
  if (!$AutomationAPI) return
  loadingAutomations.value = true
  try {
    const [triggers, taqs] = await Promise.all([
      $AutomationAPI.triggerList({ eventType: 'onManual' }).catch(() => ({})),
      $AutomationAPI.ngAutomationList({ limit: 200 }).catch(() => ({})),
    ])

    const workflowIDs = [...new Set((triggers.set || []).map(t => t.workflowID))]
    const { set: workflows = [] } = workflowIDs.length
      ? await $AutomationAPI.workflowList({ workflowID: workflowIDs }).catch(() => ({}))
      : {}

    const taqSet = Array.isArray(taqs) ? taqs : taqs.set || []
    runnable.value = [
      ...workflows
        .filter(w => w.handle)
        .map(w => ({ handle: w.handle, label: `${w.meta?.name || w.handle} (${w.handle})` })),
      ...taqSet
        .filter(a => a.handle && (a.triggers || []).some(t => t.eventType === 'onManual'))
        .map(a => ({ handle: a.handle, label: `${a.meta?.short || a.handle} (${a.handle})` })),
    ]
  } finally {
    loadingAutomations.value = false
  }
}

// Chatbots the page may open, by handle; one without a handle cannot be named.
const chatbots = ref([])
const knownChatbots = ref([])
const loadingChatbots = ref(false)

const chatbotOptions = computed(() => {
  const known = knownChatbots.value.map(c => ({
    value: c.handle,
    label: `${c.name || c.handle} (${c.handle})`,
  }))
  const missing = chatbots.value
    .filter(h => !known.some(c => c.value === h))
    .map(h => ({ value: h, label: h }))
  return [...known, ...missing]
})

async function loadChatbots() {
  loadingChatbots.value = true
  try {
    const { set = [] } = await $SystemAPI.chatbotList({ limit: 200 }).catch(() => ({}))
    knownChatbots.value = set.filter(c => c.handle)
  } finally {
    loadingChatbots.value = false
  }
}

onMounted(() => {
  loadRunnable()
  loadChatbots()
})

const handleOf = id => modules.value.find(m => m.moduleID === id)?.handle || ''

const writeOptions = computed(() =>
  moduleIDs.value.map(moduleID => ({ value: moduleID, label: handleOf(moduleID) || moduleID })),
)

// The value as the pickers stand, by name. A bound namespace is not part of
// it: the page carries that itself.
function current() {
  return {
    ...(props.namespace ? {} : { namespace: namespaceSlug.value }),
    modules: moduleIDs.value.map(handleOf).filter(Boolean),
    writes: writeIDs.value.map(handleOf).filter(Boolean),
    deletes: deleteIDs.value.map(handleOf).filter(Boolean),
    automations: automations.value,
    chatbots: chatbots.value,
    origins: origins.value.map(o => o.trim()).filter(Boolean),
  }
}

const same = (a = {}, b = {}) =>
  (a.namespace || '') === (b.namespace || '') &&
  (a.modules || []).join() === (b.modules || []).join() &&
  (a.writes || []).join() === (b.writes || []).join() &&
  (a.deletes || []).join() === (b.deletes || []).join() &&
  (a.origins || []).join() === (b.origins || []).join() &&
  (a.automations || []).join() === (b.automations || []).join() &&
  (a.chatbots || []).join() === (b.chatbots || []).join()

// Filling the pickers from the value moves them too; nothing is emitted until
// that is done, so an editor opens clean.
let filling = false

function publish() {
  if (filling) return
  const next = current()
  if (!same(next, props.modelValue)) emit('update:modelValue', next)
}

async function loadModules(ids) {
  const known = new Map(modules.value.map(m => [m.moduleID, m]))
  for (const moduleID of ids) {
    if (known.has(moduleID)) continue
    const module = await $ComposeAPI
      .moduleRead({ namespaceID: namespaceID.value, moduleID })
      .catch(() => null)
    if (module) known.set(moduleID, module)
  }
  modules.value = [...known.values()]
}

// A module the page no longer reads cannot stay one it may change.
watch(moduleIDs, async ids => {
  writeIDs.value = writeIDs.value.filter(id => ids.includes(id))
  deleteIDs.value = deleteIDs.value.filter(id => ids.includes(id))
  await loadModules(ids)
  publish()
})

watch(writeIDs, publish)
watch(deleteIDs, publish)
watch(origins, publish, { deep: true })
watch(automations, publish)
watch(chatbots, publish)

async function onNamespaceChange(id) {
  moduleIDs.value = []
  writeIDs.value = []
  deleteIDs.value = []
  modules.value = []
  namespaceSlug.value = ''

  if (id) {
    const namespace = await $ComposeAPI.namespaceRead({ namespaceID: id }).catch(() => null)
    namespaceSlug.value = namespace?.slug || ''
  }
  publish()
}

// The stored ID first, then the name: a page brought in from another instance
// carries IDs that mean nothing here.
async function fill(meta = {}) {
  filling = true
  try {
    namespaceSlug.value = props.namespace?.slug || meta.namespace || ''
    namespaceID.value = props.namespace?.namespaceID || ''
    moduleIDs.value = []
    writeIDs.value = []
    deleteIDs.value = []
    modules.value = []
    origins.value = [...(meta.origins || [])]
    automations.value = [...(meta.automations || [])]
    chatbots.value = [...(meta.chatbots || [])]

    if (!props.namespace && !meta.namespace) return

    const namespace =
      props.namespace ||
      (meta.namespaceID &&
        (await $ComposeAPI.namespaceRead({ namespaceID: meta.namespaceID }).catch(() => null))) ||
      (await $ComposeAPI.namespaceList({ slug: meta.namespace, limit: 1 }).catch(() => ({})))
        .set?.[0]
    if (!namespace) return

    namespaceID.value = namespace.namespaceID

    const found = []
    for (const handle of meta.modules || []) {
      const moduleID = meta.moduleIDs?.[handle]
      const module = moduleID
        ? await $ComposeAPI
            .moduleRead({ namespaceID: namespace.namespaceID, moduleID })
            .catch(() => null)
        : (
            await $ComposeAPI
              .moduleList({ namespaceID: namespace.namespaceID, handle, limit: 1 })
              .catch(() => ({}))
          ).set?.[0]
      if (module) found.push(module)
    }

    modules.value = found
    moduleIDs.value = found.map(m => m.moduleID)
    writeIDs.value = found.filter(m => (meta.writes || []).includes(m.handle)).map(m => m.moduleID)
    deleteIDs.value = found
      .filter(m => (meta.deletes || []).includes(m.handle))
      .map(m => m.moduleID)
  } finally {
    // The watchers above run after this tick; they must see the fill as done
    // only once they have seen it at all.
    await Promise.resolve()
    filling = false
  }
}

// Refilled only when the value moves from outside, not on its own echo.
watch(
  () => props.modelValue,
  meta => {
    if (!same(meta, current())) fill(meta)
  },
  { immediate: true, deep: true },
)
</script>
