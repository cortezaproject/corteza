<template>
  <div class="h-full overflow-auto p-6">
    <div class="max-w-3xl mx-auto">
      <CFormGroup label="Connections">
        <template #actions>
          <Button
            v-if="!disabled"
            icon="pi pi-plus"
            label="Add connection"
            severity="secondary"
            size="small"
            @click="pickerOpen = true"
          />
        </template>

        <CFormItemList
          :items="connections"
          item-key="id"
          empty-message="No connections yet."
          :hide-remove="disabled"
          remove-label="Remove connection"
          @select="c => configureResource?.(c.id)"
          @remove="c => store.removeResource(project.id, c.id)"
        >
          <template #default="{ item }">
            <div class="flex items-center gap-3 min-w-0">
              <span
                class="inline-flex items-center justify-center w-8 h-8 rounded-md ring-1 shrink-0"
                :class="[cfg.bg, cfg.ring]"
              >
                <i :class="[connectorIcon(item), cfg.text]" />
              </span>
              <div class="min-w-0">
                <div class="font-medium truncate flex items-center gap-2">
                  <span class="truncate">{{ item.name }}</span>
                  <Tag
                    v-if="sensitivity(item.sensitivity)"
                    :value="sensitivity(item.sensitivity).label"
                    :severity="sensitivity(item.sensitivity).severity"
                    class="!text-[10px] !py-0 shrink-0"
                  />
                </div>
                <div class="text-xs text-muted-color">{{ connectorLabel(item) }}</div>
              </div>
            </div>
          </template>
        </CFormItemList>
      </CFormGroup>
    </div>

    <ConnectorPicker v-model="pickerOpen" @pick="onPick" />
  </div>
</template>

<script setup>
import ConnectorPicker from '@/sections/project/components/connections/ConnectorPicker.vue'
import { connector } from '@/sections/project/config/connectors'
import { kindConfig } from '@/sections/project/config/kinds'
import { sensitivity } from '@/sections/project/config/sensitivity'
import { useProjectsStore } from '@/sections/project/stores/projects'
import { computed, inject, ref } from 'vue'

const props = defineProps({
  project: { type: Object, required: true },
  disabled: { type: Boolean, default: false },
})

const store = useProjectsStore()
const configureResource = inject('configureResource', null)
const createResource = inject('createResource', null)

const cfg = kindConfig('connection')
const connections = computed(() => (props.project.resources || []).filter(r => r.kind === 'connection'))

const connectorLabel = c => connector(c.connector)?.label || 'Connection'
const connectorIcon = c => connector(c.connector)?.icon || cfg.icon

const pickerOpen = ref(false)
function onPick(c) {
  pickerOpen.value = false
  // Open the staged edit dialog for the new connection; Save adds it.
  createResource?.('connection', { connector: c.id, name: c.label })
}
</script>
