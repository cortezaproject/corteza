<template>
  <div v-if="module" class="flex flex-col gap-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <CFormGroup
        :label="$t('module.edit.config.dal.connection.label')"
        :description="$t('module.edit.config.dal.connection.description')"
        input-id="connectionID"
      >
        <Select
          id="connectionID"
          v-model="module.config.dal.connectionID"
          :options="connections"
          optionLabel="label"
          optionValue="connectionID"
          :disabled="processing"
          :placeholder="$t('module.edit.config.dal.connection.placeholder')"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup
        :label="$t('module.edit.config.dal.ident.label')"
        :description="$t('module.edit.config.dal.ident.description')"
        input-id="ident"
      >
        <InputText
          id="ident"
          v-model="module.config.dal.ident"
          :placeholder="$t('module.edit.config.dal.ident.placeholder')"
          class="w-full"
        />
      </CFormGroup>
    </div>

    <Divider />

    <CFormGroup
      :label="$t('module.edit.config.dal.module-fields.label')"
      :description="$t('module.edit.config.dal.module-fields.description')"
    >
      <div class="flex flex-col gap-2 bg-surface backdrop-blur-sm border rounded-lg p-4 shadow-sm">
        <DalFieldStoreEncoding
          v-for="{ field, storeIdent, label, isMulti } in moduleFields"
          :key="field"
          :config="moduleFieldEncoding[field] || {}"
          :field="field"
          :label="label"
          :is-multi="isMulti"
          :default-strategy="moduleFieldDefaultEncodingStrategy"
          :store-ident="storeIdent"
          @change="applyModuleFieldStrategyConfig(field, $event)"
        />
      </div>
    </CFormGroup>

    <Divider />

    <CFormGroup
      :label="$t('module.edit.config.dal.system-fields.label')"
      :description="$t('module.edit.config.dal.system-fields.description')"
    >
      <div class="flex flex-col gap-2 bg-surface backdrop-blur-sm border rounded-lg p-4 shadow-sm">
        <div class="flex justify-end">
          <SelectButton
            v-model="selectedGroup"
            :options="optionsGroups"
            optionLabel="text"
            optionValue="value"
            :allow-empty="false"
            @update:modelValue="applySelectedSystemFields"
          />
        </div>

        <DalFieldStoreEncoding
          v-for="{ field, storeIdent, label, disabled } in systemFields"
          :key="field"
          :config="systemFieldEncoding[field] || {}"
          :field="field"
          :label="label"
          :store-ident="storeIdent"
          :allow-omit-strategy="true"
          :disabled="disabled"
          @change="applySystemFieldStrategyConfig(field, $event)"
        />
      </div>
    </CFormGroup>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { moduleFieldStrategyConfig, systemFieldStrategyConfig, types } from './encoding-strategy'
import DalFieldStoreEncoding from './DalFieldStoreEncoding.vue'

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')
const module = inject('moduleDraft')

const PrimaryConnType = 'corteza::system:primary-dal-connection'

const sysFieldsBase = [
  { field: 'id', storeIdent: 'id', disabled: true },
  { field: 'namespaceID', storeIdent: 'rel_namespace', group: 'partition' },
  { field: 'moduleID', storeIdent: 'rel_module', group: 'partition' },
  { field: 'revision', storeIdent: 'revision', group: 'extras' },
  { field: 'meta', storeIdent: 'meta', group: 'extras' },
  { field: 'ownedBy', storeIdent: 'owned_by', group: 'user_reference' },
  { field: 'createdAt', storeIdent: 'created_at', group: 'timestamps' },
  { field: 'createdBy', storeIdent: 'created_by', group: 'user_reference' },
  { field: 'updatedAt', storeIdent: 'updated_at', group: 'timestamps' },
  { field: 'updatedBy', storeIdent: 'updated_by', group: 'user_reference' },
  { field: 'deletedAt', storeIdent: 'deleted_at', group: 'timestamps' },
  { field: 'deletedBy', storeIdent: 'deleted_by', group: 'user_reference' },
]

const systemFields = sysFieldsBase.map(sf => ({ ...sf, label: t(`field.system.${sf.field}`) }))

const optionsGroups = [
  { text: t('module.edit.config.dal.system-fields.grouptypes.all'), value: 'all' },
  { text: t('module.edit.config.dal.system-fields.grouptypes.partition'), value: 'partition' },
  {
    text: t('module.edit.config.dal.system-fields.grouptypes.userReference'),
    value: 'user_reference',
  },
  { text: t('module.edit.config.dal.system-fields.grouptypes.timestamps'), value: 'timestamps' },
  { text: t('module.edit.config.dal.system-fields.grouptypes.extras'), value: 'extras' },
]

const processing = ref(false)
const connections = ref([])

const moduleFields = ref([])
const moduleFieldEncoding = ref({})
const selectedGroup = ref('all')

const initialSystemFieldEncoding =
  module.value?.config?.dal?.systemFieldEncoding &&
  typeof module.value.config.dal.systemFieldEncoding === 'object'
    ? module.value.config.dal.systemFieldEncoding
    : {}

const systemFieldEncoding = ref(
  systemFields.reduce((enc, { field }) => {
    enc[field] = field === 'id' ? {} : initialSystemFieldEncoding[field] || {}
    return enc
  }, {}),
)

const moduleFieldDefaultEncodingStrategy = computed(() => types.JSON)

watch(
  () => module.value?.fields,
  () => {
    if (!module.value?.fields) return

    const list = []
    for (const f of module.value.fields) {
      if (f.isSystem) continue

      const entry = {
        field: f.name,
        label: f.label || f.name,
        storeIdent: f.name,
        isMulti: f.isMulti,
      }

      // JSON encoding stores values inside the "values" column
      const strat = f.config?.dal?.encodingStrategy
      if (!strat || strat[types.JSON]) {
        entry.storeIdent = 'values'
      }

      list.push(entry)
    }
    moduleFields.value = list

    moduleFieldEncoding.value = list.reduce((enc, { field }) => {
      const f = module.value.fields.find(mf => mf.name === field)
      enc[field] = f?.config?.dal?.encodingStrategy || {}
      return enc
    }, {})
  },
  { deep: true, immediate: true },
)

onMounted(() => {
  // recordID is always on; strip any stale omit so it cannot persist as omitted
  const sfe = module.value?.config?.dal?.systemFieldEncoding
  if (sfe && sfe.id) {
    delete sfe.id
  }
  fetchConnections()
})

async function fetchConnections() {
  processing.value = true
  try {
    const { set = [] } = await $SystemAPI.dalConnectionList()
    connections.value = set.map(c => ({
      ...c,
      label: c.meta?.name || c.handle || c.connectionID,
    }))

    const connectionID = module.value.config?.dal?.connectionID
    if (!connectionID || connectionID === '0') {
      const primary = connections.value.find(c => c.type === PrimaryConnType)
      if (!module.value.config) module.value.config = {}
      if (!module.value.config.dal) module.value.config.dal = {}
      module.value.config.dal.connectionID = primary ? primary.connectionID : '0'
    }
  } catch (e) {
    if ($toast?.toastErrorHandler) {
      $toast.toastErrorHandler(t('module.edit.config.dal.connections.fetch-failed'))(e)
    } else {
      console.error('Failed to fetch connections:', e)
    }
  } finally {
    processing.value = false
  }
}

function applyModuleFieldStrategyConfig(field, { strategy, config }) {
  const value = moduleFieldStrategyConfig(strategy, config)

  moduleFieldEncoding.value = { ...moduleFieldEncoding.value, [field]: value }

  const moduleField = module.value.fields.find(mf => mf.name === field)
  if (moduleField) {
    if (!moduleField.config) moduleField.config = {}
    if (!moduleField.config.dal) moduleField.config.dal = {}
    moduleField.config.dal.encodingStrategy = value
  }
}

function applySystemFieldStrategyConfig(field, { strategy, config }) {
  const value = systemFieldStrategyConfig(strategy, config)
  systemFieldEncoding.value = { ...systemFieldEncoding.value, [field]: value }
  persistSystemFieldEncoding()
}

function applySelectedSystemFields(selectedOption) {
  const choice = selectedOption || 'all'
  const next = {}
  for (const { field, group } of systemFields) {
    if (field === 'id') {
      next[field] = {}
      continue
    }
    if (choice === 'all') {
      next[field] = {}
    } else {
      next[field] = group === choice ? {} : { omit: true }
    }
  }
  systemFieldEncoding.value = next
  persistSystemFieldEncoding()
}

function persistSystemFieldEncoding() {
  const filtered = Object.entries(systemFieldEncoding.value).reduce((enc, [f, c]) => {
    if (c === null || (c && Object.keys(c).length)) {
      enc[f] = c
    }
    return enc
  }, {})

  if (!module.value.config) module.value.config = {}
  if (!module.value.config.dal) module.value.config.dal = {}
  module.value.config.dal.systemFieldEncoding = filtered
}
</script>
