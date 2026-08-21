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

  <!-- Record list — renders RecordListBlock via Grid just like a public page.
       The tiles keep their own height and the table takes the rest, so it is
       drawn as the one block of its own grid. -->
  <div v-else class="flex flex-col h-full">
    <div class="shrink-0">
      <Grid :blocks="blocks.tiles" :namespace="namespace" :page="syntheticPage" />
    </div>
    <div class="flex-1 min-h-0">
      <Grid :blocks="[blocks.list]" :namespace="namespace" :page="syntheticPage" />
    </div>
  </div>
</template>

<script setup>
import { computed, provide } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useModuleStore } from '@planetcrust/human-vue'
import { components } from '@planetcrust/human-vue'
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'
import { adminRecordListBlocks } from '@/sections/compose/lib/record-blocks'

const { CRouterLinkButton } = components

defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const { t } = useI18n()
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

// Metric tiles over the module's whole record set
const blocks = computed(() =>
  adminRecordListBlocks({
    moduleID: moduleID.value,
    idPrefix: '_admin_record',
    metricLabels: {
      total: t('module.allRecords.metric.total'),
      createdRecently: t('module.allRecords.metric.createdRecently'),
      updatedRecently: t('module.allRecords.metric.updatedRecently'),
      ownedByMe: t('module.allRecords.metric.ownedByMe'),
    },
  }),
)

const syntheticPage = computed(() => ({
  pageID: '0',
  title: recordModule.value?.name || '',
  moduleID: '0',
  blocks: [],
}))
</script>
