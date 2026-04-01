<template>
  <div v-if="module" class="flex flex-col gap-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <FormField name="dal.connectionID" class="flex flex-col gap-2">
        <label for="connectionID" class="font-medium text-primary">
          {{ $t('module.edit.config.dal.connection.label') }}
        </label>
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
        <small class="text-muted-color">
          {{ $t('module.edit.config.dal.connection.description') }}
        </small>
      </FormField>

      <FormField name="dal.ident" class="flex flex-col gap-2">
        <label for="ident" class="font-medium text-primary">
          {{ $t('module.edit.config.dal.ident.label') }}
        </label>
        <InputText
          id="ident"
          v-model="module.config.dal.ident"
          :placeholder="$t('module.edit.config.dal.ident.placeholder')"
          class="w-full"
        />
        <small class="text-muted-color">
          {{ $t('module.edit.config.dal.ident.description') }}
        </small>
      </FormField>
    </div>

    <Divider />

    <FormField name="dal.moduleFields" class="flex flex-col gap-2">
      <label class="font-medium text-primary mb-2">
        {{ $t('module.edit.config.dal.module-fields.label') }}
      </label>
      <small class="text-muted-color mb-4">
        {{ $t('module.edit.config.dal.module-fields.description') }}
      </small>

      <div class="flex flex-col gap-2 bg-surface backdrop-blur-sm border rounded-lg p-4 shadow-sm">
        <DalFieldStoreEncoding
          v-for="({ field, storeIdent, label, isMulti }, i) in moduleFields"
          :key="i"
          :config="moduleFieldEncoding[field] || {}"
          :field="field"
          :label="label"
          :is-multi="isMulti"
          :default-strategy="moduleFieldDefaultEncodingStrategy"
          :store-ident="storeIdent"
          @change="applyModuleFieldStrategyConfig(field, $event)"
        />
      </div>
    </FormField>

    <Divider />

    <FormField name="dal.systemFields" class="flex flex-col gap-2">
      <div class="flex items-center justify-between mb-2">
        <label class="font-medium text-primary">
          {{ $t('module.edit.config.dal.system-fields.label') }}
        </label>
        <SelectButton
          v-model="selectedGroup"
          :options="optionsGroups"
          optionLabel="text"
          optionValue="value"
          @change="applySelectedSystemFields"
        />
      </div>
      <small class="text-muted-color mb-4">
        {{ $t('module.edit.config.dal.system-fields.description') }}
      </small>

      <div class="flex flex-col gap-2 bg-surface backdrop-blur-sm border rounded-lg p-4 shadow-sm">
        <DalFieldStoreEncoding
          v-for="({ field, storeIdent, label, disabled }, i) in systemFields"
          :key="i"
          :config="systemFieldEncoding[field] || {}"
          :field="field"
          :label="label"
          :store-ident="storeIdent"
          :allow-omit-strategy="true"
          :disabled="disabled"
          @change="applySystemFieldStrategyConfig(field, $event)"
        />
      </div>
    </FormField>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { compose } from '@cortezaproject/corteza-js-next'
import { moduleFieldStrategyConfig, systemFieldStrategyConfig, types } from './encoding-strategy'
import DalFieldStoreEncoding from './DalFieldStoreEncoding.vue'

const props = defineProps({
  module: {
    type: Object,
    required: true,
  },
})

const { t } = useI18n()
const $toast = inject('$toast')
const $SystemAPI = inject('$SystemAPI')

const PrimaryConnType = 'corteza::system:primary-dal-connection'

// Initialize default system field encodings safely
const initialSystemFieldEncoding = typeof props.module.config?.dal?.systemFieldEncoding === 'object' && props.module.config.dal.systemFieldEncoding !== null ? props.module.config.dal.systemFieldEncoding : {}

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

const processing = ref(false)
const connections = ref([])

const moduleFields = ref([])
const moduleFieldEncoding = ref({})
const selectedGroup = ref('all')

const systemFieldEncoding = ref(
  systemFields.reduce((enc, { field }) => {
    enc[field] = initialSystemFieldEncoding[field] || {}
    return enc
  }, {})
)

const optionsGroups = [
  { text: t('module.edit.config.dal.system-fields.grouptypes.all'), value: 'all' },
  { text: t('module.edit.config.dal.system-fields.grouptypes.partition'), value: 'partition' },
  { text: t('module.edit.config.dal.system-fields.grouptypes.userReference'), value: 'user_reference' },
  { text: t('module.edit.config.dal.system-fields.grouptypes.timestamps'), value: 'timestamps' },
  { text: t('module.edit.config.dal.system-fields.grouptypes.extras'), value: 'extras' },
]

const moduleFieldDefaultEncodingStrategy = computed(() => types.JSON)

watch(
  () => props.module.fields,
  () => {
    moduleFields.value = []

    for (const f of props.module.fields) {
      if (f.isSystem) continue; // Filter out system fields if they are mixed
      const a = {
        field: f.name,
        label: f.label || f.name,
        storeIdent: f.name,
        isMulti: f.isMulti,
      }

      // In case of a JSON encoding strategy, default to values
      const strat = f.config?.dal?.encodingStrategy
      if (!strat || strat[types.JSON]) {
        a.storeIdent = 'values'
      }

      moduleFields.value.push(a)
    }

    moduleFieldEncoding.value = moduleFields.value.reduce((enc, { field }) => {
      const f = props.module.fields.find(mf => mf.name === field)
      if (f) {
        enc[field] = f.config?.dal?.encodingStrategy || {}
      }
      return enc
    }, {})
  },
  { deep: true, immediate: true }
)

onMounted(() => {
  fetchConnections()
})

onBeforeUnmount(() => {
  processing.value = false
  connections.value = []
  moduleFields.value = []
  moduleFieldEncoding.value = {}
  selectedGroup.value = 'all'
})

async function fetchConnections() {
  processing.value = true
  try {
    const { set = [] } = await $SystemAPI.dalConnectionList()
    connections.value = set.map(c => ({
      ...c,
      label: c.meta?.name || c.handle || c.connectionID,
    }))

    const connectionID = props.module.config?.dal?.connectionID
    if (!connectionID || connectionID === '0') {
      const primaryConnectionID = (connections.value.find(c => c.type === PrimaryConnType) || { connectionID: '0' }).connectionID
      
      // Ensure the config structures exist before mutating
      if (!props.module.config) props.module.config = {}
      if (!props.module.config.dal) props.module.config.dal = {}
      
      props.module.config.dal.connectionID = primaryConnectionID
    }
  } catch (e) {
    if ($toast && $toast.toastErrorHandler) {
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

  // merge new config into existing
  moduleFieldEncoding.value = { ...moduleFieldEncoding.value, [field]: value }

  // update the original config
  const moduleField = props.module.fields.find(mf => mf.name === field)
  if (moduleField) {
    if (!moduleField.config) moduleField.config = {}
    if (!moduleField.config.dal) moduleField.config.dal = {}
    moduleField.config.dal.encodingStrategy = value
  }
}

function applySystemFieldStrategyConfig(field, { strategy, config }) {
  const value = systemFieldStrategyConfig(strategy, config)

  // merge new config into existing
  systemFieldEncoding.value = { ...systemFieldEncoding.value, [field]: value }

  // filter out empty configs and update the original config
  const newEncoding = Object.entries(systemFieldEncoding.value).reduce((enc, [f, c]) => {
    if (c === null || Object.keys(c).length) {
      enc[f] = c
    }
    return enc
  }, {})

  if (!props.module.config) props.module.config = {}
  if (!props.module.config.dal) props.module.config.dal = {}
  props.module.config.dal.systemFieldEncoding = newEncoding
}

function applySelectedSystemFields(event) {
  const selectedOption = event.value || 'all'
  systemFieldEncoding.value = systemFields.reduce((enc, { field, group }) => {
    if (field !== 'id') {
      if (selectedOption === 'all') {
        enc[field] = {}
      } else {
        enc[field] = group === selectedOption ? {} : { omit: true }
      }
    } else {
      enc[field] = systemFieldEncoding.value[field] || {}
    }
    return enc
  }, {})

  // Re-apply to the base module object
  if (!props.module.config) props.module.config = {}
  if (!props.module.config.dal) props.module.config.dal = {}
  props.module.config.dal.systemFieldEncoding = { ...systemFieldEncoding.value }
}
</script>
