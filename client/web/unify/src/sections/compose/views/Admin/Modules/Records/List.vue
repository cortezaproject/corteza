<template>
  <!-- Topbar title -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="recordModule">{{ recordModule.name }} — {{ $t('module.allRecords.label') }}</span>
  </Teleport>

  <!-- Topbar navigation -->
  <Teleport to="#topbar-tools" :defer="true">
    <div v-if="recordModule" class="flex gap-2">
      <CRouterLinkButton
        :to="{ name: 'admin.modules.edit', params: { moduleID: recordModule.moduleID } }"
        :label="$t('module.edit.edit')"
        icon="pi pi-pencil"
        size="small"
      />
    </div>
  </Teleport>

  <!-- Module not found -->
  <div v-if="!recordModule" class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('general.resourceList.notFound') }}
    </Message>
  </div>

  <!-- Record list — renders RecordListBlock via Grid just like a public page -->
  <div v-else class="h-full">
    <Grid :blocks="blocks" :namespace="namespace" :page="syntheticPage" />
  </div>
</template>

<script setup>
import { computed, provide } from 'vue'
import { useRoute } from 'vue-router'
import { useModuleStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'

const { CRouterLinkButton } = components

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
      allowExport: true,
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
