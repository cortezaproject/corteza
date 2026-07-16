<template>
  <Dialog
    v-model:visible="showModal"
    :header="federationModalTitle"
    :style="{ width: '70vw', maxWidth: '900px' }"
    :breakpoints="{ '1199px': '85vw', '575px': '95vw' }"
    :pt="{ content: { class: 'p-0 pb-4' } }"
    modal
  >
    <Tabs value="upstream">
      <TabList class="px-4 pt-2">
        <Tab value="upstream">{{ $t('module.edit.federationSettings.upstream.title') }}</Tab>
        <Tab value="downstream">{{ $t('module.edit.federationSettings.downstream.title') }}</Tab>
      </TabList>

      <TabPanels class="p-0">
        <!-- Upstream Tab -->
        <TabPanel value="upstream" class="flex p-0" style="min-height: 400px; max-height: 60vh;">
          <div class="w-1/3 border-r border-surface overflow-y-auto bg-emphasis">
            <Listbox
              v-model="upstream.active"
              :options="servers"
              optionLabel="name"
              optionValue="nodeID"
              class="border-none w-full"
              listStyle="max-height: 100%"
            >
              <template #option="{ option }">
                <div class="flex flex-col">
                  <span class="font-medium truncate">{{ option.name }}</span>
                  <small class="text-muted-color truncate">{{ option.baseURL }}</small>
                </div>
              </template>
            </Listbox>
          </div>

          <div class="w-2/3 p-4 overflow-y-auto">
            <div v-if="upstream.processing" class="flex justify-center items-center h-full">
              <ProgressSpinner />
            </div>
            
            <div v-else-if="upstream[upstream.active]" class="flex flex-col gap-4">
              <div v-if="upstream[upstream.active].canManageModule" class="flex flex-col gap-4">
                <p class="text-muted-color m-0">
                  {{ $t('module.edit.federationSettings.upstream.description') }}
                </p>

                <FormField name="upstream.copy" class="flex items-center gap-4">
                  <label class="font-medium whitespace-nowrap">
                    {{ $t('module.edit.federationSettings.upstream.copyFrom') }}
                  </label>
                  <Select
                    :key="upstream.active"
                    v-model="upstream[upstream.active].copy"
                    :options="upstream[upstream.active].options"
                    optionLabel="name"
                    optionValue="nodeID"
                    class="w-full max-w-xs"
                    @change="copyUpstreamFrom"
                  />
                </FormField>

                <div class="flex items-center gap-2 pb-2 border-b border-surface">
                  <Checkbox
                    v-model="upstream[upstream.active].allFields"
                    :binary="true"
                    inputId="upstreamAllFields"
                    @change="selectAllFields($event, 'upstream')"
                  />
                  <label for="upstreamAllFields" class="font-bold cursor-pointer select-none">
                    {{ $t('module.edit.federationSettings.upstream.allFields') }}
                  </label>
                </div>

                <div class="flex flex-col gap-2 max-h-64 overflow-y-auto px-1">
                  <div
                    v-for="f in upstream[upstream.active].fields"
                    :key="`${upstream.active}${f.name}`"
                    class="flex items-center gap-2"
                  >
                    <Checkbox
                      v-model="f.value"
                      :binary="true"
                      :inputId="`uf_${upstream.active}_${f.name}`"
                      @change="checkChange($event, 'upstream')"
                    />
                    <label :for="`uf_${upstream.active}_${f.name}`" class="cursor-pointer select-none">
                      {{ f.label }}
                    </label>
                  </div>
                </div>
              </div>

              <div v-else class="flex justify-center items-center h-full text-muted-color">
                {{ $t('module.edit.federationSettings.noPermission') }}
              </div>
            </div>

            <div v-else class="flex justify-center items-center h-full text-muted-color">
              {{ $t('module.edit.federationSettings.noNodes') }}
            </div>
          </div>
        </TabPanel>

        <!-- Downstream Tab -->
        <TabPanel value="downstream" class="flex p-0" style="min-height: 400px; max-height: 60vh;">
          <div class="w-1/3 border-r border-surface overflow-y-auto bg-emphasis">
            <Listbox
              v-model="downstream.active"
              :options="servers"
              optionLabel="name"
              optionValue="nodeID"
              class="border-none w-full"
              listStyle="max-height: 100%"
            >
              <template #option="{ option }">
                <div class="flex flex-col">
                  <span class="font-medium truncate">{{ option.name }}</span>
                  <small class="text-muted-color truncate">{{ option.baseURL }}</small>
                </div>
              </template>
            </Listbox>
          </div>

          <div class="w-2/3 p-4 overflow-y-auto">
            <div v-if="downstream.processing" class="flex justify-center items-center h-full">
              <ProgressSpinner />
            </div>

            <div v-else-if="downstream[downstream.active]" class="flex flex-col gap-4">
              <Select
                :key="downstream.active"
                v-model="downstream[downstream.active].module"
                :options="downstream[downstream.active].options"
                optionLabel="name"
                optionValue="moduleID"
                class="w-1/2"
              />

              <template v-if="downstream[downstream.active].module">
                <div class="flex flex-col gap-3 pb-3 border-b border-surface">
                  <p class="text-muted-color m-0">
                    {{ $t('module.edit.federationSettings.downstream.description') }}
                  </p>
                  <div class="flex items-center gap-2">
                    <Checkbox
                      v-model="downstream[downstream.active].allFields[downstream[downstream.active].module]"
                      :binary="true"
                      inputId="downstreamAllFields"
                      @change="selectAllFields($event, 'downstream')"
                    />
                    <label for="downstreamAllFields" class="font-bold cursor-pointer select-none">
                      {{ $t('module.edit.federationSettings.downstream.allFields') }}
                    </label>
                  </div>
                </div>

                <div class="flex flex-col gap-3 max-h-64 overflow-y-auto px-1">
                  <div
                    v-for="sharedModuleFields in activeSharedModules"
                    :key="`${downstream.active}_${sharedModuleFields.name}`"
                    class="flex items-center justify-between gap-4 p-2 hover:bg-emphasis rounded transition-colors"
                  >
                    <div class="flex items-center gap-2 overflow-hidden">
                      <Checkbox
                        v-model="sharedModuleFields.map"
                        :binary="true"
                        :inputId="`df_${downstream.active}_${sharedModuleFields.name}`"
                        @change="checkChange($event, 'downstream')"
                      />
                      <label :for="`df_${downstream.active}_${sharedModuleFields.name}`" class="truncate cursor-pointer select-none">
                        {{ sharedModuleFields.label }}
                      </label>
                    </div>

                    <Select
                      v-show="sharedModuleFields.map"
                      :key="`${downstream.active}_${sharedModuleFields.name}`"
                      v-model="sharedModuleFields.mapped"
                      :options="transformedModuleFields"
                      optionLabel="label"
                      optionValue="name"
                      class="w-1/2 shrink-0"
                      @change="setUpdated('downstream')"
                    />
                  </div>
                </div>
              </template>
            </div>

            <div v-else class="flex justify-center items-center h-full text-muted-color">
              {{ $t('module.edit.federationSettings.noNodes') }}
            </div>
          </div>
        </TabPanel>
      </TabPanels>
    </Tabs>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          :label="$t('general.label.cancel')"
          severity="secondary"
          text
          size="small"
          @click="showModal = false"
        />
        <Button
          :label="$t('general.label.saveAndClose')"
          size="small"
          @click="handleFederationSettingsSave"
          :loading="saving"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed, inject, ref, watch, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  modal: {
    type: Boolean,
    default: false,
  },
  module: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['change', 'update:modal'])

const { t } = useI18n()
const $FederationAPI = inject('$FederationAPI')
const $toast = inject('$toast')

// State
const showModal = ref(false)
const saving = ref(false)
const servers = ref([])
const moduleFields = ref([])
const sharedModule = ref(null)

const sharedModules = ref({})
const sharedModulesMapped = ref({})
const exposedModules = ref({})
const moduleMappings = ref({})

const downstream = ref({
  active: undefined,
  processing: false,
  enabled: false,
})

const upstream = ref({
  active: undefined,
  processing: false,
  enabled: false,
})

// Computed
const federationModalTitle = computed(() => {
  const handle = props.module?.handle
  return handle
    ? t('module.edit.federationSettings.specificTitle', { handle })
    : t('module.edit.federationSettings.title')
})

const activeSharedModules = computed(() => {
  if (!downstream.value[downstream.value.active]?.module) return []
  return (sharedModulesMapped.value[downstream.value.active] || {})[downstream.value[downstream.value.active].module] || []
})

const transformedModuleMappings = computed(() => {
  const tf = transformFields(moduleFields.value)
  const mm = ((sharedModules.value[downstream.value.active] || {})[sharedModule.value] || {}).fields || []

  return tf.map((el) => {
    el.origin.value = false
    if (mm.find((e) => e.origin.name === el.origin.name)) {
      el.destination.name = el.origin.name
      el.origin.value = true
    }
    return el
  })
})

const transformedModuleFields = computed(() => {
  return [
    { name: null, label: t('module.edit.federationSettings.pickModuleField') },
    ...transformedModuleMappings.value.map((el) => ({
      name: el.origin.name,
      label: el.origin.label,
    })),
  ]
})

// Watchers
watch(
  () => props.modal,
  (val) => {
    showModal.value = val
    if (val && servers.value.length === 0) {
      preload()
    }
  },
  { immediate: true },
)

watch(showModal, (val) => {
  emit('update:modal', val)
  emit('change', val)
})

watch(
  () => props.module.fields,
  (fields) => {
    moduleFields.value = (fields || [])
      .map((f) => ({
        kind: f.kind,
        name: f.name,
        label: f.label || f.name,
        isMulti: f.isMulti,
        value: false,
        map: null,
      }))
      .sort((a, b) => a.label.localeCompare(b.label))
  },
  { immediate: true },
)

watch(
  () => upstream.value.active,
  (nodeID) => getNodeUpstream(nodeID),
)

watch(
  () => downstream.value.active,
  (nodeID) => getNodeDownstream(nodeID),
)

onBeforeUnmount(() => {
  setDefaultValues()
})

// Methods
async function preload() {
  if (!$FederationAPI) {
    console.warn('$FederationAPI not found or not injected.')
    return
  }

  try {
    const { set = [] } = await $FederationAPI.nodeSearch({ status: 'paired' })
    servers.value = set.filter(({ canManageNode }) => canManageNode)
  } catch (e) {
    if ($toast?.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.federationSettings.error.fetch.node'))(e)
  }

  for (const node of servers.value) {
    try {
      await loadExposedModules(node.nodeID)
    } catch (e) {
      if ($toast?.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.federationSettings.error.fetch.exposed'))(e)
    }
    try {
      await loadSharedModules(node.nodeID)
    } catch (e) {
      if ($toast?.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.federationSettings.error.fetch.shared'))(e)
    }
    try {
      await loadModuleMappings(node.nodeID)
    } catch (e) {
      if ($toast?.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.federationSettings.error.fetch.mmap'))(e)
    }
  }

  sharedModulesMapped.value = getSharedModulesMapped()

  if (servers.value.length > 0) {
    downstream.value.active = servers.value[0].nodeID
    upstream.value.active = servers.value[0].nodeID
  }
}

function getSharedModulesMapped() {
  const list = {}

  for (const nodeID in sharedModules.value) {
    list[nodeID] = {}

    for (const sm of sharedModules.value[nodeID]) {
      let f = [...sm.fields].sort((a, b) => a.label.localeCompare(b.label))

      const mappedFields = ((moduleMappings.value[nodeID] || {})[sm.moduleID] || {}).fields || []

      f = f.map((el) => {
        let found = false
        let mapped = (moduleFields.value.find(({ name }) => name === el.name) || {}).name || null

        if (mappedFields && mappedFields.length) {
          const m = mappedFields.find((mf) => el.name === mf.origin?.name)
          mapped = ((m || {}).destination || {}).name || null
          found = !!mapped
        }

        return {
          ...el,
          map: found,
          mapped,
        }
      })

      list[nodeID][sm.moduleID] = f
    }
  }

  return list
}

async function handleFederationSettingsSave() {
  saving.value = true
  try {
    // module mappings (downstream)
    for (const nodeID in sharedModulesMapped.value) {
      for (const moduleID in sharedModulesMapped.value[nodeID]) {
        const crtModule = sharedModules.value[nodeID].find((m) => m.moduleID === moduleID)
        if (!crtModule || !crtModule.updated) continue

        const fields = toModuleMappingFormat(sharedModulesMapped.value[nodeID][moduleID])

        const payload = {
          nodeID,
          moduleID,
          composeModuleID: props.module.moduleID,
          composeNamespaceID: props.module.namespaceID,
          fields,
        }

        try {
          await persistModuleMappings(payload)
          crtModule.updated = false
        } catch (e) {
          if ($toast?.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.federationSettings.error.persist.mmap'))(e)
        }
      }
    }

    const nodes = servers.value.map((s) => s.nodeID)

    // upstream
    for (const nodeID of nodes) {
      if (!upstream.value[nodeID] || !upstream.value[nodeID].updated) continue

      const fields = (upstream.value[nodeID].fields || []).filter((el) => el.value)

      const payload = {
        nodeID,
        moduleID: (exposedModules.value[nodeID] || {}).moduleID,
        composeModuleID: props.module.moduleID,
        composeNamespaceID: props.module.namespaceID,
        name: props.module.name,
        handle: props.module.handle,
        fields,
      }

      try {
        const response = await persistExposedModule(payload)
        upstream.value[nodeID].updated = false
        if (response && response.moduleID) {
          exposedModules.value[nodeID] = response
        }
      } catch (e) {
        if ($toast?.toastErrorHandler) $toast.toastErrorHandler(t('module.edit.federationSettings.error.persist.exposed'))(e)
      }
    }

    showModal.value = false
  } finally {
    saving.value = false
  }
}

function toModuleMappingFormat(fields) {
  return fields
    .filter((el) => el.map)
    .filter((el) => !!el.mapped)
    .map((el) => ({
      origin: { kind: el.kind, name: el.name, label: el.label, isMulti: el.isMulti },
      destination: { kind: el.kind, name: el.mapped, label: el.label, isMulti: el.isMulti },
    }))
}

function transformFields(fields) {
  return fields.map((el) => ({
    origin: { kind: el.kind, name: el.name, label: el.label || 'N/A', isMulti: false },
    destination: { kind: el.kind, name: '', label: '', isMulti: false },
  }))
}

function getNodeUpstream(nodeID) {
  if (!nodeID || upstream.value[nodeID]) return

  upstream.value.processing = true
  const exposedModule = exposedModules.value[nodeID] || {}
  const fields = moduleFields.value.map((f) => ({ ...f, value: false }))
  const exposedFields = exposedModule.fields || []

  exposedFields.forEach(({ name }) => {
    const f = fields.find((field) => field.name === name)
    if (f) f.value = true
  })

  const server = servers.value.find((s) => s.nodeID === nodeID) || {}

  const upstreamNode = {
    options: [
      { nodeID: null, name: t('module.edit.federationSettings.pickServer') },
      ...servers.value.filter((s) => s.nodeID !== nodeID),
    ],
    copy: null,
    allFields: false,
    fields,
    updated: false,
    canManageModule: !!exposedModule.canManageModule || !!server.canCreateModule,
  }

  upstreamNode.allFields = upstreamNode.fields.length > 0 && upstreamNode.fields.every((f) => f.value)

  upstream.value[nodeID] = upstreamNode
  upstream.value.processing = false
}

async function getNodeDownstream(nodeID) {
  if (!nodeID || downstream.value[nodeID]) return

  downstream.value.processing = true
  const fields = moduleFields.value.map((f) => ({ ...f, value: false }))

  const sharedModuleOpts = Object.values(sharedModules.value[nodeID] || {})
    .filter(({ canMapModule }) => canMapModule)
    .map((m) => ({ moduleID: m.moduleID, name: m.name }))

  const preselectedMod = (sharedModules.value[nodeID] || []).find(({ handle }) => handle === props.module.handle)

  const downstreamNode = {
    options: [
      { moduleID: null, name: t('module.edit.federationSettings.pickModule') },
      ...sharedModuleOpts,
    ],
    module: preselectedMod ? preselectedMod.moduleID : null,
    allFields: {},
    fields,
  }

  Object.entries(sharedModulesMapped.value[nodeID] || {}).forEach(([key, value]) => {
    if (value && value.length) {
      downstreamNode.allFields[key] = value.every((f) => f.map)
    } else {
      downstreamNode.allFields[key] = false
    }
  })

  downstream.value[nodeID] = downstreamNode
  downstream.value.processing = false
}

function selectAllFields(event, target) {
  const value = event.checked ?? event
  const active = target === 'upstream' ? upstream.value.active : downstream.value.active

  if (target === 'upstream') {
    upstream.value[active].fields.forEach((f) => (f.value = value))
    upstream.value[active].allFields = value
    upstream.value[active].updated = true
  } else if (target === 'downstream') {
    const activeModule = downstream.value[active].module
    if (!activeModule) return
    sharedModulesMapped.value[active][activeModule]?.forEach((f) => (f.map = value))
    downstream.value[active].allFields[activeModule] = value
    const mod = sharedModules.value[active]?.find((m) => m.moduleID === activeModule)
    if (mod) mod.updated = true
  }
}

function setUpdated(target) {
  const active = target === 'upstream' ? upstream.value.active : downstream.value.active

  if (target === 'upstream') {
    upstream.value[active].updated = true
  } else if (target === 'downstream') {
    const activeModule = downstream.value[active].module
    const mod = sharedModules.value[active]?.find((m) => m.moduleID === activeModule)
    if (mod) mod.updated = true
  }
}

function copyUpstreamFrom(event) {
  const nodeID = event.value || event
  const active = upstream.value.active
  
  upstream.value[active].fields.forEach((f) => {
    let value = false
    if (upstream.value[nodeID]) {
      value = !!upstream.value[nodeID].fields.find(({ name, value: v }) => name === f.name && v)
    } else if (exposedModules.value[nodeID]?.fields) {
      value = !!exposedModules.value[nodeID].fields.find(({ name }) => name === f.name)
    }
    f.value = value
  })
}

function checkChange(event, target) {
  const value = event.checked ?? event
  const active = target === 'upstream' ? upstream.value.active : downstream.value.active

  if (target === 'upstream') {
    upstream.value[active].allFields = value ? upstream.value[active].fields.every((f) => f.value) : false
  } else if (target === 'downstream') {
    const activeModule = downstream.value[active].module
    if (!activeModule) return
    const allMapped = sharedModulesMapped.value[active][activeModule]?.every(({ map }) => map)
    downstream.value[active].allFields[activeModule] = value ? allMapped : false
  }

  setUpdated(target)
}

function persistExposedModule(payload) {
  if (payload.moduleID) {
    return $FederationAPI.manageStructureUpdateExposed(payload)
  }
  return $FederationAPI.manageStructureCreateExposed(payload)
}

function persistModuleMappings(payload) {
  return $FederationAPI.manageStructureCreateMappings(payload)
}

async function loadSharedModules(nodeID) {
  if (sharedModules.value[nodeID]) return

  const data = await $FederationAPI.manageStructureListAll({ nodeID, shared: 1 })
  sharedModules.value[nodeID] = (data || []).map((d) => ({ ...d, updated: false }))
}

async function loadExposedModules(nodeID) {
  if (exposedModules.value[nodeID]) return

  const data = await $FederationAPI.manageStructureListAll({ nodeID, exposed: 1 })
  const exposedModule = (data || []).find(({ composeModuleID }) => composeModuleID === props.module.moduleID)
  if (exposedModule) {
    exposedModules.value[nodeID] = exposedModule
  }
}

async function loadModuleMappings(nodeID) {
  if (moduleMappings.value[nodeID] || !sharedModules.value[nodeID]) return

  const mm = {}
  for (const { moduleID } of sharedModules.value[nodeID]) {
    mm[moduleID] = []
    try {
      const data = await $FederationAPI.manageStructureReadMappings({ nodeID, moduleID, composeModuleID: props.module.moduleID })
      mm[moduleID] = data
    } catch {
      // silent
    }
  }

  moduleMappings.value[nodeID] = mm
}

function setDefaultValues() {
  servers.value = []
  moduleFields.value = []
  sharedModule.value = null
  sharedModules.value = {}
  sharedModulesMapped.value = {}
  exposedModules.value = {}
  moduleMappings.value = {}
  downstream.value = { active: undefined, processing: false, enabled: false }
  upstream.value = { active: undefined, processing: false, enabled: false }
}
</script>
