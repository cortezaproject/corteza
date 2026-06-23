<template>
  <div class="h-full overflow-auto p-4">
    <CFormGroup :label="$t('project.connections.title')">
      <template #actions>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.connections.add')"
          severity="secondary"
          size="small"
          :loading="preparing"
          @click="pickerOpen = true"
        />
      </template>

      <CFormItemList
        :items="connections"
        item-key="id"
        :empty-message="$t('project.connections.empty')"
        :hide-remove="disabled"
        :remove-label="$t('project.connections.remove')"
        @select="onSelect"
        @remove="onRemove"
      >
        <template #default="{ item }">
          <div class="flex items-center gap-3 min-w-0">
            <span
              class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
              :class="[cfg.bg, cfg.ring]"
            >
              <i :class="[iconForConnector(item.catalogID), cfg.text]" />
            </span>
            <div class="min-w-0">
              <div class="font-medium truncate flex items-center gap-2">
                <span class="truncate">{{ item.name }}</span>
                <Tag :value="$t('project.connections.configured')" severity="secondary" class="!text-xs shrink-0">
                  <template #icon>
                    <i class="pi pi-check-circle text-green-500 !text-xs" />
                  </template>
                </Tag>
              </div>
              <div class="text-xs text-muted-color">{{ labelForConnector(item.catalogID) }}</div>
            </div>
          </div>
        </template>
      </CFormItemList>
    </CFormGroup>

    <ConnectorPicker v-model="pickerOpen" :items="pickerItems" @pick="onPick" />
    <ConfigureConnectionDialog
      v-model="dialogOpen"
      :connection="activeConnection"
      :project-id="project.id"
      :configured="activeConfigured"
    />
  </div>
</template>

<script setup>
import ConfigureConnectionDialog from '@/sections/project/components/connections/ConfigureConnectionDialog.vue'
import ConnectorPicker from '@/sections/project/components/connections/ConnectorPicker.vue'
import { connector } from '@/sections/project/config/connectors'
import { kindConfig } from '@/sections/project/config/kinds'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

function onRemove(c) {
  confirmDelete({
    header: t('project.connections.removeConfirm.header'),
    message: t('project.connections.removeConfirm.message', { name: c.name }),
    onConfirm: async () => {
      try {
        await store.removeConnection(props.project.id, c.id)
      } catch (err) {
        $toast.toastErrorHandler(t('project.connections.toastRemoveFailed'))(err)
      }
    },
  })
}

const cfg = kindConfig('connection')
// The library carries no UI icon; resolve one from the static catalog by
// catalogID, falling back to the generic connection icon.
const iconForConnector = id => connector(id)?.icon || cfg.icon
const labelForConnector = id => connector(id)?.label || id || ''

const connections = computed(() => store.connectionsFor(props.project.id))

// Picker offers the live library, mapped to the picker's item shape and
// filtered to the Resource Management whitelist. With no whitelist (Free mode)
// the full library is offered.
const pickerItems = computed(() => {
  const allowed = store.allowedConnectorIds(props.project.id)
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

const pickerOpen = ref(false)
const dialogOpen = ref(false)
const preparing = ref(false)
const activeConnection = ref(null)
const activeConfigured = ref(null)

// Picking imports the real connection (so we have its auth-field schema), then
// opens the configure dialog to create the configured connection.
async function onPick(item) {
  pickerOpen.value = false
  await openConfigure(item, null)
}

// Clicking a configured connection re-opens the dialog to edit it.
function onSelect(item) {
  if (props.disabled) return
  openConfigure(
    { connectionID: item.connectionID, catalogID: item.catalogID },
    { configuredConnectionID: item.configuredConnectionID, name: item.name, config: item.config },
  )
}

async function openConfigure(item, configured) {
  preparing.value = true
  try {
    activeConnection.value = await store.prepareConnection(item)
    activeConfigured.value = configured
    dialogOpen.value = true
  } catch (err) {
    $toast.toastErrorHandler(t('project.configureConnection.toastImportFailed'))(err)
  } finally {
    preparing.value = false
  }
}

async function refresh(id) {
  try {
    await Promise.all([store.loadConnectionLibrary(), store.loadConnections(id)])
  } catch (err) {
    $toast.toastErrorHandler(t('project.connections.toastLoadFailed'))(err)
  }
}

onMounted(() => refresh(props.project.id))
watch(
  () => props.project.id,
  id => id && refresh(id),
)
</script>
