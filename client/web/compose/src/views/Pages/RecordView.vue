<template>
  <!-- Page title in topbar -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="page">{{ page.title }}</span>
  </Teleport>

  <!-- Loading state -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner style="width: 32px; height: 32px" />
  </div>

  <!-- Page content with record context -->
  <div v-else-if="page && positionedBlocks.length" class="flex flex-col h-full">
    <div class="flex-1 overflow-auto p-4">
      <Grid :blocks="positionedBlocks" :namespace="namespace" :page="page" />
    </div>

    <!-- Record Toolbar -->
    <div class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <Button
          :label="$t('general.label.back')"
          icon="pi pi-arrow-left"
          severity="secondary"
          @click="$router.back()"
        />
        <div class="flex gap-2">
          <CInputDelete
            v-if="record?.canDeleteRecord"
            :label="$t('general.label.delete')"
            :message="$t('page.public.record.toolbar.deleteConfirm')"
            :header="page.title"
            :disabled="deleting"
            @confirm="handleDelete"
          />
          <Button
            v-if="record?.canUpdateRecord"
            :label="$t('general.label.edit')"
            icon="pi pi-pencil"
            @click="handleEdit"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- No blocks -->
  <div v-else-if="page" class="flex items-center justify-center h-full">
    <p class="text-muted-color">{{ $t('page.noBlock') }}</p>
  </div>

  <!-- Page not found -->
  <div v-else class="flex items-center justify-center h-full">
    <Message severity="warn" :closable="false">
      {{ $t('page.invalid') }}
    </Message>
  </div>
</template>

<script setup>
import Grid from '@/components/PageBlocks/Grid.vue'
import { useModuleStore } from '@/stores/module'
import { usePageLayoutStore } from '@/stores/page-layout'
import { usePageStore } from '@/stores/page'
import { useRecordStore } from '@/stores/record'
import { components } from '@cortezaproject/corteza-vue-next'
import { computed, inject, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

const { CInputDelete } = components

defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const $toast = inject('$toast')
const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()

const loading = ref(false)
const deleting = ref(false)
const page = ref(null)
const layout = ref(null)
const record = ref(null)

const positionedBlocks = computed(() => {
  if (!page.value || !layout.value) {
    if (page.value?.blocks?.length) {
      return page.value.blocks
    }
    return []
  }

  return layout.value.blocks
    .map(layoutBlock => {
      const pageBlock = page.value.blocks.find(b => b.blockID === layoutBlock.blockID)
      if (!pageBlock) return null

      return {
        ...pageBlock,
        xywh: layoutBlock.xywh || pageBlock.xywh,
      }
    })
    .filter(Boolean)
})

async function loadPage() {
  const pageID = route.params.pageID
  const recordID = route.params.recordID
  if (!pageID) return

  loading.value = true

  page.value = pageStore.getByID(pageID) || null

  if (page.value) {
    const layouts = pageLayoutStore.getByPageID(pageID)
    layout.value = layouts.length > 0 ? layouts[0] : null

    // Load the record for toolbar permission checks
    if (recordID && page.value.moduleID) {
      try {
        const mod = moduleStore.getByID(page.value.moduleID)
        if (mod) {
          record.value = await recordStore.findByID({
            namespaceID: mod.namespaceID,
            moduleID: mod.moduleID,
            recordID,
          })
        }
      } catch (e) {
        console.error('Failed to load record for toolbar:', e)
        record.value = null
      }
    }
  }

  loading.value = false
}

async function handleDelete() {
  if (!record.value || !page.value) return

  deleting.value = true
  try {
    await recordStore.delete({
      namespaceID: record.value.namespaceID,
      moduleID: record.value.moduleID,
      recordID: record.value.recordID,
    })
    $toast.toastSuccess(t('notification.record.deleteSuccess'))
    router.back()
  } catch (e) {
    console.error('Failed to delete record:', e)
    $toast.toastDanger(t('notification.record.deleteFailed'))
  } finally {
    deleting.value = false
  }
}

function handleEdit() {
  // For now, navigate to the same page in edit mode
  // This can be extended when a dedicated record edit route is available
  // For now we'll just emit or toggle an edit state
  // TODO: implement record edit navigation when route is available
}

// Load on mount and when route params change
watch(
  () => [route.params.pageID, route.params.recordID],
  () => loadPage(),
  { immediate: true },
)
</script>
