<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">{{ recordModule.name }} — {{ $t('module.allRecords.label') }}</span>
  </Teleport>

  <!-- Module not found -->
  <div v-if="!recordModule" class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('general.resourceList.notFound') }}
    </Message>
  </div>

  <!-- Record list — renders RecordListBlock via Grid just like a public page -->
  <div v-else class="flex-1 overflow-auto">
    <Grid :blocks="blocks" :namespace="namespace" :page="syntheticPage" />
  </div>
</template>

<script setup>
import { computed, provide } from 'vue'
import { useRoute } from 'vue-router'
import { useModuleStore } from '@/stores/module'
import Grid from '@/components/PageBlocks/Grid.vue'

defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const moduleStore = useModuleStore()

const moduleID = computed(() => route.params.moduleID)
const recordModule = computed(() => (moduleID.value ? moduleStore.getByID(moduleID.value) : null))

// Provide admin route resolver so RecordListBlock navigates to admin routes
provide('$recordRoutes', {
  view: (modID, recordID) => ({
    name: 'admin.modules.record.view',
    params: { moduleID: modID, recordID },
  }),
  create: (modID, cloneFromID) => ({
    name: 'admin.modules.record.create',
    params: { moduleID: modID },
    ...(cloneFromID ? { query: { cloneFromID } } : {}),
  }),
})

// Single RecordList block spanning the full grid width
const blocks = computed(() => [
  {
    blockID: '_admin_record_list',
    kind: 'RecordList',
    title: '',
    description: '',
    style: { wrap: { kind: 'card' } },
    options: {
      moduleID: moduleID.value,
      fields: [],
      perPage: 20,
      selectable: true,
    },
    xywh: [0, 0, 48, 36],
    meta: { tempID: '_admin_record_list' },
  },
])

const syntheticPage = computed(() => ({
  pageID: '0',
  title: recordModule.value?.name || '',
  moduleID: '0',
  blocks: [],
}))
</script>
