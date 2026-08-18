<template>
  <!-- Page title in topbar -->
  <Teleport to="#topbar-title" :defer="true">
    <span v-if="page">{{ pageTitle }}</span>
  </Teleport>

  <!-- Page builder button in topbar tools -->
  <Teleport to="#topbar-tools" :defer="true">
    <ButtonGroup v-if="page?.canUpdatePage || showTranslatorButton" class="gap-1">
      <template v-if="page?.canUpdatePage">
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
      </template>
      <PageTranslator
        v-if="page"
        :page="page"
        :namespace="namespace"
        :layouts="pageLayouts"
        @update:page="page = $event"
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
import Grid from '@/sections/compose/components/PageBlocks/Grid.vue'
import PageTranslator from '@/sections/compose/components/Admin/Page/PageTranslator.vue'
import { usePageLayoutStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { computed, inject, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { compose, NoID } from '@planetcrust/human-js'
import { fetchBlockID, usePageVisibility } from '@/sections/compose/composables/usePageVisibility'
import { useResourceTranslations } from '@/sections/compose/composables/useResourceTranslations'
import { evaluatePrefilter, usesRecordVariables } from '@/sections/compose/lib/record-filter'

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
const $SystemAPI = inject('$SystemAPI', null)
const $Auth = inject('$Auth', null)

const { buildExpressionVariables, determineLayout, evaluateBlocks } = usePageVisibility(
  $SystemAPI,
  $Auth,
)
const { showTranslatorButton } = useResourceTranslations()

const loading = ref(false)
const page = ref(null)
const layout = ref(null)
const invisibleBlockIDs = ref(new Set())

const pageLayouts = computed(() =>
  page.value ? pageLayoutStore.getByPageID(page.value.pageID) : [],
)

/**
 * The page's displayed title. A layout may override the page title with its own,
 * interpolated against the signed-in user (`config.useTitle`) — so a layout can
 * title the screen `Welcome, ${user.name}` rather than the static page title.
 * Mirrors `RecordView.vue`, minus the record: there is none here, so a template
 * reading one cannot be evaluated and the page title stands, as it does
 * whenever the override is off, empty, or the template is malformed.
 */
const pageTitle = computed(() => {
  if (!page.value) return ''

  const { config = {}, meta = {} } = layout.value || {}
  if (!config.useTitle || !meta.title) return page.value.title
  if (usesRecordVariables(meta.title)) return page.value.title

  try {
    return (
      evaluatePrefilter(meta.title, {
        record: undefined,
        user: $Auth?.user || {},
        recordID: NoID,
        ownerID: NoID,
        userID: $Auth?.user?.userID || NoID,
      }) || page.value.title
    )
  } catch {
    return page.value.title
  }
})

// The blocks this page can show, before visibility is applied. This, not
// page.blocks, is what visibility is evaluated over: a block no layout places is
// never rendered, so its condition would only add an expression that can fail
// for nothing.
const layoutBlocks = computed(() => {
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

const positionedBlocks = computed(() =>
  // meta.hidden is handled by Grid (tab children must still reach TabsBlock via props.blocks)
  // invisibleBlockIDs are blocks hidden by visibility expressions/roles — remove entirely
  layoutBlocks.value.filter(b => !invisibleBlockIDs.value.has(fetchBlockID(b))),
)

async function loadPage() {
  const pageID = route.params.pageID
  if (!pageID) return

  loading.value = true
  invisibleBlockIDs.value = new Set()

  try {
    page.value = pageStore.getByID(pageID) || null

    if (page.value) {
      const layouts = pageLayoutStore.getByPageID(pageID)
      const vars = buildExpressionVariables()

      // An explicitly requested layout (?layoutID=, e.g. from a navigation
      // block) wins over automatic selection.
      const requested = route.query.layoutID
      layout.value = await determineLayout(
        layouts,
        vars,
        typeof requested === 'string' ? requested : undefined,
      )

      // Evaluate block visibility after layout is resolved
      if (layoutBlocks.value.length) {
        invisibleBlockIDs.value = await evaluateBlocks(layoutBlocks.value, vars)
      }
    }
  } finally {
    loading.value = false
  }
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
