<template>
  <!-- Page title in topbar -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="page">{{ page.title }}</span>
  </Teleport>

  <!-- Page builder button in topbar tools -->
  <Teleport to="#topbar-tools" :defer="true">
    <ButtonGroup v-if="page?.canUpdatePage" class="gap-1">
      <Button
        :label="$t('page.block.general.label.pageBuilder')"
        icon="pi pi-wrench"
        size="small"
        @click="goToBuilder"
      />
      <Button
        v-tooltip.bottom="$t('navigation.editPage')"
        icon="pi pi-pencil"
        size="small"
        @click="goToEditPage"
      />
    </ButtonGroup>
  </Teleport>

  <!-- Loading state -->
  <div v-if="loading" class="flex items-center justify-center h-full">
    <ProgressSpinner style="width: 32px; height: 32px" />
  </div>

  <!-- Page content -->
  <div v-else-if="page && positionedBlocks.length" class="h-full">
    <Grid :blocks="positionedBlocks" :namespace="namespace" :page="page" />
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
import { usePageLayoutStore } from '@/stores/page-layout'
import { usePageStore } from '@/stores/page'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { compose } from '@cortezaproject/corteza-js-next'

defineProps({
  namespace: {
    type: Object,
    required: true,
  },
})

const route = useRoute()
const router = useRouter()
const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()

const loading = ref(false)
const page = ref(null)
const layout = ref(null)

const positionedBlocks = computed(() => {
  if (!page.value || !layout.value) {
    // No layout — fall back to page blocks with their default xywh
    if (page.value?.blocks?.length) {
      return page.value.blocks
    }
    return []
  }

  // Merge layout block positions with page block definitions
  return layout.value.blocks
    .map(layoutBlock => {
      const pageBlock = page.value.blocks.find(b => b.blockID === layoutBlock.blockID)
      if (!pageBlock) return null

      // Clone page block and override xywh from layout
      return compose.PageBlockMaker({
        ...pageBlock,
        xywh: layoutBlock.xywh || pageBlock.xywh,
      })
    })
    .filter(Boolean)
})

function loadPage() {
  const pageID = route.params.pageID
  if (!pageID) return

  loading.value = true

  page.value = pageStore.getByID(pageID) || null

  if (page.value) {
    // Get layouts for this page, pick the first one (no visibility expressions yet)
    const layouts = pageLayoutStore.getByPageID(pageID)
    layout.value = layouts.length > 0 ? layouts[0] : null
  }

  loading.value = false
}

function goToBuilder() {
  if (page.value) {
    router.push({
      name: 'admin.pages.builder',
      params: { pageID: page.value.pageID },
    })
  }
}

function goToEditPage() {
  if (page.value) {
    router.push({
      name: 'admin.pages.edit',
      params: { pageID: page.value.pageID },
    })
  }
}

// Load on mount and when pageID changes
watch(
  () => route.params.pageID,
  () => loadPage(),
  { immediate: true },
)
</script>
