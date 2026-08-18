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
  <div v-else-if="page" class="flex flex-col items-center justify-center gap-3 h-full">
    <p class="text-muted-color">
      {{ emptyStateMessage }}
    </p>
    <Button
      v-if="hasNoLayouts && page.canUpdatePage"
      :label="$t('page.page-layout.add')"
      icon="pi pi-plus"
      size="small"
      @click="goToEditPage"
    />
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
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { compose, NoID } from '@planetcrust/human-js'
import {
  clearRefusal,
  fetchBlockID,
  refuseOnce,
  usePageVisibility,
} from '@/sections/compose/composables/usePageVisibility'
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
const $toast = inject('$toast', null)
const { t } = useI18n()

const { buildExpressionVariables, determineLayout, evaluateBlocks } = usePageVisibility(
  $SystemAPI,
  $Auth,
)
const { showTranslatorButton } = useResourceTranslations()

const loading = ref(false)
const page = ref(null)
const layout = ref(null)
const invisibleBlockIDs = ref(new Set())
// True once this page has been refused and we stayed anyway (see refuseOnce)
const noLayoutMatched = ref(false)

const pageLayouts = computed(() =>
  page.value ? pageLayoutStore.getByPageID(page.value.pageID) : [],
)

// A page nobody has given a layout yet is unfinished, not withheld — it says so
// and offers the way to finish it, rather than leaving as a no-match does.
const hasNoLayouts = computed(() => !!page.value && pageLayouts.value.length === 0)

const emptyStateMessage = computed(() => {
  if (hasNoLayouts.value) return t('page.noLayouts')
  if (noLayoutMatched.value) return t('notification.page.noMatchingLayout')
  return t('page.noBlock')
})

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
  // No layout, no blocks. A layout is what decides which blocks a viewer sees,
  // so falling back to the page's raw set would show everything precisely when
  // the rules meant to narrow it did not apply.
  if (!page.value || !layout.value) {
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
      await applyLayout()
    }
  } finally {
    loading.value = false
  }
}

/**
 * Picks the layout and settles block visibility under it.
 *
 * A page that has layouts but matches none of them is not a page with nothing
 * on it — it is a page this viewer was not meant to reach, so it says so and
 * leaves rather than rendering an empty grid.
 */
async function applyLayout() {
  const layouts = pageLayoutStore.getByPageID(page.value.pageID)
  const vars = buildExpressionVariables()

  // An explicitly requested layout (?layoutID=, e.g. from a navigation block)
  // wins over automatic selection.
  const requested = route.query.layoutID
  const requestedLayoutID = typeof requested === 'string' ? requested : undefined

  layout.value = await determineLayout(layouts, vars, requestedLayoutID)

  // The request has been spent, so it leaves the URL: keeping it would re-pin
  // this layout on every later resolution, and the address would name a layout
  // that may no longer be the one on screen.
  if (requestedLayoutID) dropLayoutQuery()

  // The builder is exempt — an author editing a layout has to be able to see it
  // whether or not its own condition holds right now.
  if (layouts.length && !layout.value && route.name !== 'admin.pages.builder') {
    noLayoutMatched.value = true
    if (refuseOnce(page.value.pageID)) {
      $toast?.toastWarning(t('notification.page.noMatchingLayout'))
      leaveUnshowablePage()
    }
    return
  }

  noLayoutMatched.value = false
  clearRefusal()

  // Evaluate block visibility after layout is resolved
  if (layoutBlocks.value.length) {
    invisibleBlockIDs.value = await evaluateBlocks(layoutBlocks.value, vars)
  }
}

function dropLayoutQuery() {
  if (!route.query.layoutID) return
  const query = { ...route.query }
  delete query.layoutID
  router.replace({ query })
}

/**
 * Leaves for the namespace's page list rather than through history: the page
 * behind this one can be the namespace's landing page, which redirects to the
 * very page that matched nothing — so going back lands straight on it again.
 */
function leaveUnshowablePage() {
  router.push({ name: 'pages', params: { slug: route.params.slug } })
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

// A navigation block linking to the layout of the page already open changes
// only this parameter. Acting on a truthy value alone keeps the strip that
// follows from reading as a second request.
watch(
  () => route.query.layoutID,
  layoutID => {
    if (layoutID && page.value) applyLayout()
  },
)
</script>
