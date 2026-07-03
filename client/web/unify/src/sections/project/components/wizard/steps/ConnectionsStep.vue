<template>
  <div class="h-full overflow-auto p-4">
    <div class="flex flex-col gap-2">
      <div>
        <Button
          v-if="!disabled"
          icon="pi pi-plus"
          :label="$t('project.connections.add')"
          severity="secondary"
          size="small"
          @click="createResource?.('connection')"
        />
      </div>

      <div class="mt-1">
        <CFormItemList
          :items="connections"
          item-key="id"
          :reveal-on-hover="!disabled"
          :empty-message="$t('project.connections.empty')"
          :hide-remove="disabled"
          :remove-label="$t('project.connections.remove')"
          @select="onSelect"
          @remove="onRemove"
        >
          <template #default="{ item }">
            <CFormItemContent :title="item.name" :subtitle="labelForConnector(item.catalogID)" />
          </template>

          <template #actions>
            <Tag
              :value="$t('project.connections.configured')"
              severity="secondary"
              class="!text-xs shrink-0"
            >
              <template #icon>
                <i class="pi pi-check-circle text-green-500 !text-xs" />
              </template>
            </Tag>
          </template>
        </CFormItemList>
      </div>
    </div>
  </div>
</template>

<script setup>
import { connector } from '@/sections/project/config/connectors'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { useConfirmDelete } from '@planetcrust/human-vue'
import { computed, inject, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const { t } = useI18n()
const $toast = inject('$toast')
const { confirmDelete } = useConfirmDelete()

// Detail/create dialogs are mounted once by the wizard and opened via these
// injected helpers; the step never mounts them itself.
const inspectResource = inject('inspectResource', null)
const createResource = inject('createResource', null)

// The library carries no UI icon; resolve a display label from the static
// catalog by catalogID, falling back to the raw catalogID.
const labelForConnector = id => connector(id)?.label || id || ''

const connections = computed(() => store.connectionsFor(props.project.id))

// Clicking a configured connection opens its detail dialog for editing.
function onSelect(item) {
  if (props.disabled) return
  inspectResource?.('connection', item.id)
}

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
