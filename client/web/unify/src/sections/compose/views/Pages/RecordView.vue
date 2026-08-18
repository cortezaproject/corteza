<template>
  <!-- Page title in topbar -->
  <Teleport v-if="!inModal" to="#topbar-title" :defer="true">
    <span v-if="page">{{ pageTitle }}</span>
  </Teleport>

  <!-- Admin tools in topbar -->
  <Teleport v-if="!inModal" to="#topbar-tools" :defer="true">
    <ButtonGroup v-if="page?.canUpdatePage" class="gap-1">
      <Button
        v-if="page.isRecordPage"
        :label="$t('page.moduleEdit')"
        icon="pi pi-database"
        size="small"
        @click="goToModuleEdit"
      />
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

  <!-- Loading state. The branch itself is what withholds the page — a load
       rebuilds the record from nothing and the layout it was shown under is
       still the previous one, so there is nothing here worth rendering. Whether
       that wait is worth a spinner is a separate question, and usually no.
       The branch outlives `loading` by whatever the cover has left to run: a
       spinner that leaves the moment the work does is the same blink as one
       that arrives the moment it starts. -->
  <div v-if="loading || pageCover" class="flex items-center justify-center h-full">
    <ProgressSpinner v-if="pageCover" style="width: 32px; height: 32px" />
  </div>

  <!-- Page content with record context -->
  <Form
    ref="formRef"
    v-else-if="page && positionedBlocks.length"
    :resolver="resolver"
    @submit="handleSave"
    class="flex flex-col h-full"
  >
    <div class="flex-1 overflow-auto">
      <Grid
        :blocks="positionedBlocks"
        :namespace="namespace"
        :page="page"
        :record="record"
        :loading="cover"
      />
    </div>

    <!-- Record Toolbar (only on record pages) -->
    <div v-if="page?.isRecordPage" class="shrink-0 border-t border-surface bg-surface">
      <div class="flex items-center justify-between p-3">
        <!-- Left side -->
        <div class="flex gap-2">
          <Button
            v-if="mode === 'view' && layoutButtons.back"
            :label="$t('general.label.back')"
            icon="pi pi-arrow-left"
            severity="secondary"
            @click="handleCancel"
          />
          <Button
            v-else-if="mode !== 'view'"
            :label="$t('general.label.cancel')"
            icon="pi pi-times"
            severity="secondary"
            :disabled="isSaving"
            @click="handleCancel"
          />
        </div>

        <!-- Center: prev/next navigation -->
        <div
          v-if="mode === 'view' && (recordNavigation.prev || recordNavigation.next)"
          class="flex gap-1"
        >
          <Button
            icon="pi pi-chevron-left"
            severity="secondary"
            :disabled="!recordNavigation.prev || navigating !== null"
            :loading="navigating === 'prev' && cover"
            :title="$t('general.recordNavigation.prev')"
            @click="navigateToRecord(recordNavigation.prev, 'prev')"
          />
          <Button
            icon="pi pi-chevron-right"
            severity="secondary"
            :disabled="!recordNavigation.next || navigating !== null"
            :loading="navigating === 'next' && cover"
            :title="$t('general.recordNavigation.next')"
            @click="navigateToRecord(recordNavigation.next, 'next')"
          />
        </div>

        <!-- Right side -->
        <div class="flex gap-2">
          <!-- Delete button (view mode or edit mode for existing records) -->
          <CInputDelete
            v-if="showDeleteButton"
            :label="$t('general.label.delete')"
            :message="$t('page.public.record.toolbar.deleteConfirm')"
            :header="page.title"
            :disabled="deleting || isSaving"
            @confirm="handleDelete"
          />

          <!-- Clone / Save as copy (view mode, existing record) -->
          <Button
            v-if="mode === 'view' && !isNew && record && layoutButtons.clone"
            :label="$t('general.label.saveAsCopy')"
            icon="pi pi-copy"
            severity="secondary"
            @click="handleClone"
          />

          <!-- New record button (view mode) -->
          <Button
            v-if="mode === 'view' && layoutButtons.new"
            :label="$t('general.label.add')"
            icon="pi pi-plus"
            severity="secondary"
            @click="handleNew"
          />

          <!-- Edit button (view mode only) -->
          <Button
            v-if="mode === 'view' && record?.canUpdateRecord && layoutButtons.edit"
            :label="$t('general.label.edit')"
            icon="pi pi-pencil"
            severity="primary"
            @click="handleEdit"
          />

          <!-- Save button (edit/create mode) -->
          <Button
            v-if="mode !== 'view' && layoutButtons.submit"
            type="submit"
            :label="$t('general.label.save')"
            icon="pi pi-check"
            :loading="isSaving"
          />
        </div>
      </div>
    </div>
  </Form>

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
import {
  clearRefusal,
  fetchBlockID,
  refuseOnce,
  usePageVisibility,
} from '@/sections/compose/composables/usePageVisibility'
import { useDeferredBusy } from '@planetcrust/human-vue'
import { useModuleStore } from '@planetcrust/human-vue'
import { usePageLayoutStore } from '@planetcrust/human-vue'
import { usePageStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import { compose, validator, NoID } from '@planetcrust/human-js'
import { evaluatePrefilter, usesRecordVariables } from '@/sections/compose/lib/record-filter'
import { components, useHistoryBack } from '@planetcrust/human-vue'
import { computed, inject, nextTick, onBeforeUnmount, provide, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import {
  mergeAttachmentIDs,
  uploadRecordAttachment,
} from '@/sections/compose/lib/record-attachments'

const { CInputDelete } = components

const props = defineProps({
  namespace: {
    type: Object,
    required: true,
  },
  inModal: {
    type: Boolean,
    default: false,
  },
  modalPageID: {
    type: String,
    default: null,
  },
  modalRecordID: {
    type: String,
    default: null,
  },
})

const emit = defineEmits(['close'])

const route = useRoute()
const router = useRouter()
const historyBack = useHistoryBack()
const { t } = useI18n()
const $toast = inject('$toast')
const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI', null)
const $Auth = inject('$Auth', {})
const $eventBus = inject('$eventBus', null)

const { buildExpressionVariables, determineLayout, evaluateBlocks } = usePageVisibility(
  $SystemAPI,
  $Auth,
)

const pageStore = usePageStore()
const pageLayoutStore = usePageLayoutStore()
const moduleStore = useModuleStore()

const isMultiField = fieldName =>
  !!moduleStore.getByID(page.value?.moduleID)?.fields.find(f => f.name === fieldName)?.isMulti
const recordStore = useRecordStore()

const formRef = ref(null)
const serverErrors = ref({})

const loading = ref(false)

// A load that beats the delay leaves the page area empty for those few frames
// rather than flashing a spinner through them — the same bargain a swap makes.
const pageCover = useDeferredBusy(loading)
const invisibleBlockIDs = ref(new Set())
// True once this page has been refused and we stayed anyway (see refuseOnce)
const noLayoutMatched = ref(false)
// Fields the active layout asks for on top of the module's own required ones
const layoutRequiredFields = ref([])

const recordNavigation = computed(() => {
  const recordID = props.inModal ? props.modalRecordID : route.params.recordID
  if (!recordID || recordID === '0') return {}
  // Access the ref directly so Vue tracks it as a reactive dependency
  const ids = recordStore.paginationRecordIDs
  const idx = ids.indexOf(recordID)
  if (idx === -1) return {}
  return {
    prev: idx > 0 ? ids[idx - 1] : undefined,
    next: idx < ids.length - 1 ? ids[idx + 1] : undefined,
  }
})

function navigateToRecord(targetRecordID, direction) {
  navigating.value = direction
  if (!targetRecordID) return
  if (props.inModal) {
    router.push({ query: { ...route.query, recordID: targetRecordID, edit: undefined } })
  } else {
    router.push({
      name: 'page.record',
      params: { slug: route.params.slug, pageID: route.params.pageID, recordID: targetRecordID },
    })
  }
}
const deleting = ref(false)
const isSaving = ref(false)
const page = ref(null)
const layout = ref(null)
const record = ref(null)

/**
 * The record page's displayed title. A layout may override the page title with
 * its own, interpolated against the open record (`config.useTitle`) — so a
 * layout can title the screen `${record.values.name}` rather than the static
 * page title. Falls back to the page title whenever the override is off, empty,
 * or cannot be evaluated yet.
 */
const pageTitle = computed(() => {
  if (!page.value) return ''

  const { config = {}, meta = {} } = layout.value || {}
  if (!config.useTitle || !meta.title) return page.value.title

  // The record resolves after the layout is picked; a template reading it would
  // throw until then, so keep showing the page title rather than flashing an
  // error or a half-interpolated string.
  if (usesRecordVariables(meta.title) && !record.value) return page.value.title

  try {
    return (
      evaluatePrefilter(meta.title, {
        record: record.value,
        user: $Auth?.user || {},
        recordID: record.value?.recordID || NoID,
        ownerID: record.value?.ownedBy || NoID,
        userID: $Auth?.user?.userID || NoID,
      }) || page.value.title
    )
  } catch {
    return page.value.title
  }
})
const pristineRecord = ref(null)
const navigatingAfterSave = ref(false)

// Mode derived from route or props
const routeMode = computed(() => {
  const recordID = props.inModal ? props.modalRecordID : route.params.recordID
  if (recordID === '0') return 'create'
  if (route.query.edit === '1') return 'edit'
  return 'view'
})

// The mode the page is rendered in. It lags routeMode by one staged evaluation
// so that the form, the toolbar and every block a mode condition governs change
// on the same frame: isView/isCreate/isEdit are condition variables, and
// flipping the mode first put the page on screen under the previous mode's
// answers until the evaluation caught up.
const mode = ref(routeMode.value)

const isNew = computed(() => mode.value === 'create')

// Button visibility from layout config (all enabled by default when no layout)
const layoutButtons = computed(() => {
  const b = layout.value?.config?.buttons ?? {}
  return {
    back: b.back?.enabled ?? true,
    delete: b.delete?.enabled ?? true,
    clone: b.clone?.enabled ?? true,
    new: b.new?.enabled ?? true,
    edit: b.edit?.enabled ?? true,
    submit: b.submit?.enabled ?? true,
  }
})

// Show delete button: in view mode if record has permission, or in edit mode for existing record
const showDeleteButton = computed(() => {
  if (!record.value) return false
  if (!layoutButtons.value.delete) return false
  if (mode.value === 'view') return record.value.canDeleteRecord
  if (mode.value === 'edit') return record.value.canDeleteRecord
  return false
})

// A record swap on the page already on screen: the blocks cover themselves
// rather than the page blanking. `loading` stays for the loads that build the
// page from nothing.
const swapping = ref(false)

/**
 * The one spinner state for a whole transition, deferred.
 *
 * A swap on this machine finishes in about 85ms, so a cover that appeared the
 * moment one started would flash on and straight off again — the very thing the
 * page spinner was doing. It waits instead, and only a swap slow enough to be
 * worth explaining ever draws one.
 */
const cover = useDeferredBusy(swapping)

/**
 * Blocks that answer for conditions of their own — RecordBlock's field
 * conditions — so a transition can settle them before it applies anything.
 *
 * Each returns the call that applies what it worked out. The layout and the
 * block set are only two of the three things a swap changes; a page whose
 * layout and blocks stay the same has nothing else moving, so fields left to
 * catch up afterwards are the whole of what a viewer sees change twice.
 */
const settlers = new Set()

function registerSettler(fn) {
  settlers.add(fn)
  return () => settlers.delete(fn)
}

const pendingByField = reactive(new Map())

// Provide context so child blocks (RecordBlock) can inject it
provide('recordViewContext', {
  mode,
  record,
  isNew,
  isSaving,
  cover,
  registerSettler,

  /**
   * A block that saved the record on its own (RecordBlock's inline edit) hands
   * the response back through here. Both refs move together — leaving
   * pristineRecord behind would make the next view→edit clone from stale
   * values — and the layout is re-picked, since a save can satisfy a different
   * layout condition than the one the record was opened with.
   */
  adoptSaved(saved) {
    pristineRecord.value = saved
    record.value = saved
    resolveLayout()
  },
})

provide('$fileUploadContext', {
  registerPending(fieldName, files) {
    if (files.length > 0) pendingByField.set(fieldName, files)
    else pendingByField.delete(fieldName)
  },
})

// The blocks this page can show, before visibility is applied — the layout's
// selection intersected with the page's definitions. This, not page.blocks, is
// what visibility is evaluated over: a block no layout places is never rendered,
// so its condition would only add an expression that can fail for nothing.
function blocksFor(pg, lay) {
  // No layout, no blocks. A layout is what decides which blocks a viewer sees,
  // so falling back to the page's raw set would show everything precisely when
  // the rules meant to narrow it did not apply.
  if (!pg || !lay) return []

  return lay.blocks
    .map(layoutBlock => {
      const pageBlock = pg.blocks.find(b => b.blockID === layoutBlock.blockID)
      if (!pageBlock) return null

      // PageBlockMaker, not a spread: blocks reach their renderer with the
      // methods their class defines (fetch, reorderViews), not just options.
      return compose.PageBlockMaker({
        ...pageBlock,
        xywh: layoutBlock.xywh || pageBlock.xywh,
      })
    })
    .filter(Boolean)
}

const layoutBlocks = computed(() => blocksFor(page.value, layout.value))

const positionedBlocks = computed(() =>
  // meta.hidden is handled by Grid (tab children must still reach TabsBlock via props.blocks)
  // invisibleBlockIDs are blocks hidden by visibility expressions/roles — remove entirely
  layoutBlocks.value.filter(b => !invisibleBlockIDs.value.has(fetchBlockID(b))),
)

// A page nobody has given a layout yet is unfinished, not withheld — it says so
// and offers the way to finish it, rather than leaving as a no-match does.
const hasNoLayouts = computed(
  () => !!page.value && pageLayoutStore.getByPageID(page.value.pageID).length === 0,
)

const emptyStateMessage = computed(() => {
  if (hasNoLayouts.value) return t('page.noLayouts')
  if (noLayoutMatched.value) return t('notification.page.noMatchingLayout')
  return t('page.noBlock')
})

const navigating = ref(null) // 'prev' | 'next' | null

// Cancels the in-flight record load when we navigate to another record / leave.
let recordLoadAbort = null
function abortRecordLoad() {
  if (recordLoadAbort) {
    recordLoadAbort.abort()
    recordLoadAbort = null
  }
}

/**
 * Picks the layout for the current record and mode.
 *
 * Layout conditions see the record, so this must run *after* the record is
 * resolved. It is deliberately not reactive to field edits: switching layout
 * swaps the rendered block set, which would remount blocks and discard
 * in-progress input. It re-runs only on the transitions that change what the
 * page is about — a different record, or a move between view/edit/create.
 */
let _layoutSeq = 0

/** Picks the layout for a record and mode. Writes nothing. */
async function pickLayout({ forRecord, forMode }) {
  const layouts = pageLayoutStore.getByPageID(page.value.pageID)
  const vars = buildExpressionVariables({
    record: forRecord,
    isRecordPage: true,
    mode: forMode,
  })

  // An explicitly requested layout (?layoutID=, e.g. from a navigation block)
  // wins over automatic selection.
  const requested = props.inModal ? undefined : route.query.layoutID
  const requestedLayoutID = typeof requested === 'string' ? requested : undefined

  const seq = ++_layoutSeq
  return {
    seq,
    layouts,
    requestedLayoutID,
    resolved: await determineLayout(layouts, vars, requestedLayoutID),
  }
}

/** Applies a picked layout. False means nothing more should be applied. */
function commitLayout({ seq, layouts, requestedLayoutID, resolved }) {
  // A superseded resolution (rapid record swap, mode toggle mid-load) must not
  // overwrite the layout picked for the record now on screen
  if (seq !== _layoutSeq) return false

  layout.value = resolved

  // The request has been spent, so it leaves the URL: keeping it would re-pin
  // this layout on every later resolution, and the address would name a layout
  // that may no longer be the one on screen.
  if (requestedLayoutID) dropLayoutQuery()

  // Layouts exist but this record matches none of them: there is no honest
  // block set to fall back to, so say so and leave rather than render the
  // page's raw set, which is everything the layouts were there to narrow.
  // The builder is exempt — an author editing a layout has to be able to see it
  // whether or not its own condition holds right now.
  if (layouts.length && !resolved && route.name !== 'admin.pages.builder') {
    noLayoutMatched.value = true
    if (refuseOnce(page.value.pageID)) {
      $toast?.toastWarning(t('notification.page.noMatchingLayout'))
      if (props.inModal) emit('close')
      else leaveUnshowablePage()
    }
    return false
  }

  noLayoutMatched.value = false
  clearRefusal()
  return true
}

async function resolveLayout() {
  if (!page.value) return

  const forRecord = record.value
  const forMode = mode.value
  const picked = await pickLayout({ forRecord, forMode })
  if (!commitLayout(picked)) return

  // Required fields belong to the layout, so every path that settles one
  // settles them with it.
  commitRequiredFields(await pickRequiredFields(picked.resolved, { forRecord, forMode }))
}

/**
 * Leaves for the namespace's page list rather than through history: the page
 * behind this one can be the namespace's landing page, which redirects to the
 * very page that matched nothing — so going back lands straight on it again.
 */
function leaveUnshowablePage() {
  router.push({ name: 'pages', params: { slug: route.params.slug } })
}

function dropLayoutQuery() {
  if (!route.query.layoutID) return
  const query = { ...route.query }
  delete query.layoutID
  router.replace({ query })
}

async function loadRecord(recordID) {
  if (!page.value || !recordID || recordID === '0') return
  const moduleID = page.value.moduleID
  if (!moduleID) return
  const mod = moduleStore.getByID(moduleID)
  if (!mod) return

  abortRecordLoad()
  const ac = new AbortController()
  recordLoadAbort = ac

  // Swapping to another record of the page already on screen does not blank it.
  // The page keeps its chrome and its geometry, the blocks show a spinner over
  // the record they are still holding, and the new record, its layout and its
  // block conditions are all evaluated before any of it is applied — so the
  // values change once, settled. Blanking to the page spinner is what read as a
  // flash. A first load still gets that spinner, in loadPage().
  swapping.value = true

  try {
    const loaded = await recordStore.findByID({
      namespaceID: mod.namespaceID,
      moduleID: mod.moduleID,
      recordID,
      force: true,
      signal: ac.signal,
    })
    // Fast-swapping to another record skips loadPage(), so the layout has to be
    // re-picked here or the previous record's layout would stay on screen.
    const forMode = routeMode.value
    const apply = await stageTransition({ forRecord: loaded, forMode })
    if (ac.signal.aborted) return

    pristineRecord.value = loaded
    record.value = loaded
    mode.value = forMode
    serverErrors.value = {}
    apply()
    await dropPendingVisibility()
  } catch (e) {
    if (ac.signal.aborted) return
    console.error('Failed to load record:', e)
    record.value = null
  } finally {
    if (recordLoadAbort === ac) recordLoadAbort = null
    if (!ac.signal.aborted) {
      swapping.value = false
      navigating.value = null
    }
  }
}

async function loadPage() {
  const pageID = props.inModal ? props.modalPageID : route.params.pageID
  const recordID = props.inModal ? props.modalRecordID : route.params.recordID
  if (!pageID) return

  abortRecordLoad()
  const ac = new AbortController()
  recordLoadAbort = ac

  loading.value = true
  // Nothing of the previous record survives this load, so the rendered mode
  // catches up with the route at once rather than lagging: what follows builds
  // the record from it.
  mode.value = routeMode.value
  record.value = null
  pristineRecord.value = null
  invisibleBlockIDs.value = new Set()

  try {
    page.value = pageStore.getByID(pageID) || null

    if (page.value) {
      if (!page.value.isRecordPage) {
        router.replace({
          name: 'page',
          params: {
            slug: route.params.slug,
            pageID: page.value.pageID,
          },
        })
        return
      }

      const moduleID = page.value.moduleID
      if (moduleID) {
        const mod = moduleStore.getByID(moduleID)
        if (mod) {
          // A record needs a module WITH fields (compose.Record throws otherwise);
          // a fieldless module renders the page without an initialized record.
          if (mode.value === 'create' && mod.fields?.length) {
            // Handle clone
            if (route.query.cloneFromID) {
              try {
                const source = await recordStore.findByID({
                  namespaceID: mod.namespaceID,
                  moduleID: mod.moduleID,
                  recordID: route.query.cloneFromID,
                  signal: ac.signal,
                })
                const newRec = new compose.Record(mod)
                for (const field of mod.fields) {
                  newRec.setValue(field.name, source.values[field.name])
                }
                // Prefill ownedBy with current user
                newRec.ownedBy = $Auth?.user?.userID || undefined
                record.value = newRec
              } catch (e) {
                if (ac.signal.aborted) return
                console.error('Failed to load source record for clone:', e)
                record.value = new compose.Record(mod, { ownedBy: $Auth?.user?.userID })
              }
            } else {
              record.value = new compose.Record(mod, { ownedBy: $Auth?.user?.userID })
            }

            const refField = route.query.refField
            const refValue = route.query.refValue
            if (record.value && refField && refValue) {
              const field = mod.fields.find(f => f.name === refField)
              if (field) {
                const valueArray = Array.isArray(refValue) ? refValue : [refValue]
                if (field.isMulti) {
                  record.value.setValue(refField, valueArray)
                } else {
                  record.value.setValue(refField, valueArray[0])
                }
              }
            }
          } else if (recordID && recordID !== '0') {
            try {
              const loaded = await recordStore.findByID({
                namespaceID: mod.namespaceID,
                moduleID: mod.moduleID,
                recordID,
                force: true,
                signal: ac.signal,
              })
              pristineRecord.value = loaded
              // In edit mode, start with a clone; mode watch will re-clone on later transitions
              record.value = mode.value === 'edit' ? loaded.clone() : loaded
            } catch (e) {
              if (ac.signal.aborted) return
              console.error('Failed to load record:', e)
              record.value = null
            }
          }
        }
      }
    }

    // Layout conditions read the record, so pick the layout once it is resolved
    if (page.value?.isRecordPage) {
      await resolveLayout()
    }

    // Evaluate block visibility before revealing content (no flash)
    await runBlockVisibility()
  } finally {
    // If this load was superseded/cancelled, leave state to the newer load.
    if (!ac.signal.aborted) {
      loading.value = false
      if (recordLoadAbort === ac) recordLoadAbort = null
    }
  }
}

let _blockVisibilityTimer = null
let _blockVisibilitySeq = 0

/** Evaluates block conditions for a block set, record and mode. Writes nothing. */
async function pickBlockVisibility(blocks, { forRecord, forMode }) {
  clearTimeout(_blockVisibilityTimer)
  const seq = ++_blockVisibilitySeq
  if (!blocks.length) return { seq, invisible: null }

  const vars = buildExpressionVariables({ record: forRecord, isRecordPage: true, mode: forMode })
  return { seq, invisible: await evaluateBlocks(blocks, vars) }
}

function commitBlockVisibility({ seq, invisible }) {
  // A slower earlier response must not overwrite a newer one
  if (seq !== _blockVisibilitySeq || !invisible) return
  invisibleBlockIDs.value = invisible
}

let _requiredFieldsSeq = 0

/**
 * Evaluates the layout's conditional required fields. Writes nothing.
 *
 * A layout may ask for a field the module leaves optional — "a reason is
 * required once the status is rejected". An empty condition asks for it
 * outright. This only ever adds: a field the module already requires stays
 * required on every layout.
 */
async function pickRequiredFields(lay, { forRecord, forMode }) {
  const seq = ++_requiredFieldsSeq
  const rules = (lay?.config?.validation?.requiredFields || []).filter(({ field }) => field)
  if (!rules.length) return { seq, fields: [] }

  const always = rules.filter(({ condition }) => !condition?.trim()).map(({ field }) => field)
  const conditional = rules.filter(({ condition }) => condition?.trim())
  if (!conditional.length || !$SystemAPI) return { seq, fields: always }

  const expressions = {}
  for (const { field, condition } of conditional) expressions[field] = condition

  const variables = buildExpressionVariables({
    record: forRecord,
    isRecordPage: true,
    mode: forMode,
  })

  try {
    const res = await $SystemAPI.expressionEvaluate({ variables, expressions })
    return { seq, fields: [...always, ...Object.keys(res).filter(k => res[k])] }
  } catch (e) {
    // A rule that cannot be evaluated must not block the save: the person is
    // left with a field marked required and no way to see what asked for it.
    console.error('Failed to evaluate layout required fields:', e)
    return { seq, fields: always }
  }
}

function commitRequiredFields({ seq, fields }) {
  if (seq !== _requiredFieldsSeq) return
  layoutRequiredFields.value = fields
}

async function runBlockVisibility() {
  commitBlockVisibility(
    await pickBlockVisibility(layoutBlocks.value, {
      forRecord: record.value,
      forMode: mode.value,
    }),
  )
}

/**
 * Evaluates a whole transition — layout, then the blocks that layout places —
 * and hands back the one call that puts all of it on screen.
 *
 * Applying each answer as it arrives is what read as a flash: the layout lands,
 * the page repaints, and the blocks a condition governs appear or vanish a
 * round-trip later. Staged this way the page changes once, already settled.
 */
async function stageTransition({ forRecord, forMode }) {
  if (!page.value) return () => {}

  const picked = await pickLayout({ forRecord, forMode })
  const visibility = await pickBlockVisibility(blocksFor(page.value, picked.resolved), {
    forRecord,
    forMode,
  })
  const required = await pickRequiredFields(picked.resolved, { forRecord, forMode })
  // A block that fails to settle must not hold up the ones that did, nor leave
  // the page on the record it was showing before.
  const settled = await Promise.all(
    [...settlers].map(settle =>
      settle({ record: forRecord, mode: forMode }).catch(e => {
        console.error('A block failed to settle:', e)
        return () => {}
      }),
    ),
  )

  return () => {
    if (!commitLayout(picked)) return
    commitBlockVisibility(visibility)
    commitRequiredFields(required)
    settled.forEach(apply => apply())
  }
}

/**
 * The values watcher fires on whatever a staged commit just assigned and would
 * re-run the evaluation that commit already made — with the same answer, one
 * round-trip later. The staged answer is the newer one, so its timer goes.
 */
async function dropPendingVisibility() {
  await nextTick()
  clearTimeout(_blockVisibilityTimer)
}

// Debounce rapid changes (e.g. typing in a field)
function evaluateConditions() {
  clearTimeout(_blockVisibilityTimer)
  _blockVisibilityTimer = setTimeout(() => {
    runBlockVisibility()
    runRequiredFields()
  }, 300)
}

async function runRequiredFields() {
  commitRequiredFields(
    await pickRequiredFields(layout.value, { forRecord: record.value, forMode: mode.value }),
  )
}

// Block conditions and the layout's required fields both follow the record as
// it is edited: re-evaluate on value changes. Mode switches and record swaps
// are staged instead, so they do not come through here.
watch(() => record.value?.values, evaluateConditions, { deep: true })

function resolver() {
  const errors = {}

  for (const [fieldName, message] of Object.entries(serverErrors.value)) {
    errors[fieldName] = [{ message }]
  }

  const recordModule = page.value ? moduleStore.getByID(page.value.moduleID) : null
  if (!record.value || !recordModule) return { errors }
  // The layout can only add: a field the module requires is required whichever
  // layout is on screen.
  const requiredByLayout = f =>
    layoutRequiredFields.value.includes(f.fieldID) || layoutRequiredFields.value.includes(f.name)

  for (const field of recordModule.fields) {
    if (field.isRequired || requiredByLayout(field)) {
      const val = record.value.values[field.name]
      if (validator.IsEmpty(val)) {
        errors[field.name] = [{ message: t('field.required-field') }]
      }
    }
  }
  return { errors }
}

async function handleSave({ valid }) {
  if (!valid) {
    $toast.toastWarning(t('general.notification.formErrors'))
    nextTick(() => {
      document
        .querySelector('.p-message-error')
        ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
    })
    return
  }
  if (!record.value || !page.value) return

  isSaving.value = true

  try {
    const namespaceID = props.namespace.namespaceID
    const moduleID = page.value.moduleID
    const recordID = isNew.value ? '' : record.value.recordID

    for (const [fieldName, files] of pendingByField) {
      const ids = await Promise.all(
        files.map(file =>
          uploadRecordAttachment($ComposeAPI, { namespaceID, moduleID, recordID, fieldName, file }),
        ),
      )
      record.value.setValue(
        fieldName,
        mergeAttachmentIDs(record.value.values[fieldName], ids, isMultiField(fieldName)),
      )
    }

    const saved = isNew.value
      ? await recordStore.create(record.value)
      : await recordStore.update(record.value)

    $toast.toastSuccess(
      t(isNew.value ? 'notification.record.createSuccess' : 'notification.record.updateSuccess'),
    )

    // Adopt the response rather than the record we sent: it carries what the
    // server computed, and conditions read those values. Saving always leaves
    // edit/create for view, and that mode change re-picks the layout — so
    // resolving here as well would only spend a request on the older record.
    pristineRecord.value = saved
    record.value = saved

    if (props.inModal) {
      if (isNew.value) {
        // Update query to new recordID instead of '0'; drop the clone/prefill
        // params so the resulting view URL doesn't leak them into later actions.
        const q = { ...route.query }
        delete q.cloneFromID
        delete q.refField
        delete q.refValue
        router.replace({
          query: {
            ...q,
            recordID: saved.recordID,
            edit: undefined,
          },
        })
      } else {
        router.replace({ query: { ...route.query, edit: undefined } })
      }
    } else {
      // Set flag so leave guard allows this programmatic navigation
      navigatingAfterSave.value = true
      router.replace({
        name: 'page.record',
        params: {
          slug: route.params.slug,
          pageID: route.params.pageID,
          recordID: saved.recordID,
        },
      })
    }
  } catch (e) {
    console.error('Failed to save record:', e)
    const details = e?.details ?? []
    const fieldErrors = {}
    for (const detail of details) {
      if (detail.meta?.field) {
        fieldErrors[detail.meta.field] = detail.message
      }
    }
    if (Object.keys(fieldErrors).length > 0) {
      serverErrors.value = fieldErrors
      await nextTick()
      formRef.value?.validate()
    } else {
      $toast.toastErrorHandler(
        t(isNew.value ? 'notification.record.createFailed' : 'notification.record.updateFailed'),
      )(e)
    }
  } finally {
    isSaving.value = false
  }
}

function handleEdit() {
  router.push({ query: { ...route.query, edit: '1' } })
}

function handleClone() {
  if (props.inModal) {
    router.push({
      query: {
        ...route.query,
        recordID: '0',
        cloneFromID: record.value.recordID,
      },
    })
  } else {
    router.push({
      name: 'page.record',
      params: {
        slug: route.params.slug,
        pageID: props.inModal ? props.modalPageID : route.params.pageID,
        recordID: '0',
      },
      query: { cloneFromID: record.value.recordID },
    })
  }
}

function handleNew() {
  if (props.inModal) {
    // Strip clone/prefill params so Add always opens a blank record,
    // even when the current record was reached via clone (cloneFromID lingers in the query).
    const q = { ...route.query }
    delete q.cloneFromID
    delete q.refField
    delete q.refValue
    delete q.edit
    router.push({
      query: {
        ...q,
        recordID: '0',
      },
    })
  } else {
    router.push({
      name: 'page.record',
      params: { slug: route.params.slug, pageID: route.params.pageID, recordID: '0' },
      // Explicit empty query drops any lingering cloneFromID/refField/refValue.
      query: {},
    })
  }
}

function handleCancel() {
  if (props.inModal) {
    if (mode.value === 'view') {
      emit('close')
    } else {
      const q = { ...route.query }
      delete q.edit
      if (isNew.value) {
        emit('close')
      } else {
        router.replace({ query: q })
      }
    }
    return
  }

  if (mode.value === 'edit') {
    // Cancel edit — return to view mode by removing the edit query param
    router.replace({ query: {} })
  } else {
    goBack()
  }
}

function goBack() {
  historyBack({ name: 'pages', params: { slug: route.params.slug } })
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
    if (props.inModal) {
      emit('close')
    } else {
      goBack()
    }
  } catch (e) {
    console.error('Failed to delete record:', e)
    $toast.toastErrorHandler(t('notification.record.deleteFailed'))(e)
  } finally {
    deleting.value = false
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

function goToModuleEdit() {
  if (page.value?.isRecordPage) {
    router.push({
      name: 'admin.modules.edit',
      params: { moduleID: page.value.moduleID },
    })
  }
}

// Guard against navigating away with unsaved changes
onBeforeRouteLeave(() => {
  if (navigatingAfterSave.value) {
    navigatingAfterSave.value = false
    return true
  }
  if (mode.value !== 'view' && !isSaving.value) {
    return window.confirm(t('general.editor.unsavedChanges'))
  }
})

// Handle view↔edit transitions in-place without reloading
watch(routeMode, async (newMode, oldMode) => {
  // A transition that also changes which record is on screen belongs to the
  // load watcher below: it fetches the record and applies the mode with it, and
  // a second staging here would race that one.
  const routeRecordID = props.inModal ? props.modalRecordID : route.params.recordID
  if (routeRecordID !== record.value?.recordID) return

  if (newMode === 'edit' && oldMode === 'view' && pristineRecord.value) {
    record.value = pristineRecord.value.clone()
  } else if (newMode === 'view' && oldMode === 'edit') {
    record.value = pristineRecord.value
  }

  // isView/isCreate/isEdit are layout and block condition variables, so a mode
  // switch changes what the page shows. Evaluated whole and applied with the
  // mode itself, so the form, the toolbar and the blocks change together.
  const apply = await stageTransition({ forRecord: record.value, forMode: newMode })
  mode.value = newMode
  apply()
  await dropPendingVisibility()
})

// Clear server errors as soon as the user edits anything
watch(
  () => record.value?.values,
  () => {
    if (Object.keys(serverErrors.value).length > 0) {
      serverErrors.value = {}
      nextTick(() => formRef.value?.validate())
    }
  },
  { deep: true },
)

// Load on mount and when pageID, recordID, or cloneFromID changes (not on edit-query toggle)
// Separate getters, not one getter returning an array: an array source is
// compared by identity, so it re-fires on every touched dependency — including
// the edit-query toggle this watcher is meant to ignore.
watch(
  [
    () => (props.inModal ? props.modalPageID : route.params.pageID),
    () => (props.inModal ? props.modalRecordID : route.params.recordID),
    () => route.query.cloneFromID,
    () => route.query.refField,
    () => route.query.refValue,
  ],
  ([newPageID, newRecordID, newCloneFromID, newRefField, newRefValue], old) => {
    const [oldPageID, oldRecordID, oldCloneFromID, oldRefField, oldRefValue] = old || []
    // Fast-swap only when moving between two existing records on the same page.
    // Transitions to/from create mode ('0') must go through loadPage() so the
    // blank/clone record is (re)built — loadRecord() no-ops on '0' and would
    // otherwise leave the previously viewed record's values on screen.
    if (
      old &&
      newPageID === oldPageID &&
      newRecordID !== '0' &&
      oldRecordID !== '0' &&
      newCloneFromID === oldCloneFromID &&
      newRefField === oldRefField &&
      newRefValue === oldRefValue &&
      page.value
    ) {
      loadRecord(newRecordID)
    } else {
      loadPage()
    }
  },
  { immediate: true },
)

// A navigation block linking to the layout of the record already open changes
// only this parameter. Acting on a truthy value alone keeps the strip that
// follows from reading as a second request.
watch(
  () => route.query.layoutID,
  layoutID => {
    if (layoutID && page.value && !props.inModal) resolveLayout()
  },
)

const offRefetch = $eventBus?.on('refetch-records', () => loadPage())

onBeforeUnmount(() => {
  abortRecordLoad()
  offRefetch?.()
})
</script>
