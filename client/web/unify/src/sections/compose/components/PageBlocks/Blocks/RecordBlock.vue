<template>
  <PageBlock :block="block" :record="activeRecord">
    <div v-if="loading" class="flex items-center justify-center h-full p-5">
      <ProgressSpinner style="width: 28px; height: 28px" />
    </div>

    <div v-else-if="!fieldModule" class="p-5 text-muted-color italic">
      {{ $t('block.record.noModule') }}
    </div>

    <template v-else>
      <!-- Inline edit save/cancel bar -->
      <div
        v-if="hasActiveInlineEdits"
        class="flex items-center gap-1 px-2 py-1 border-b border-surface shrink-0 justify-end"
      >
        <Button
          v-tooltip.bottom="$t('block.recordList.tooltip.saveChanges')"
          icon="pi pi-check"
          text
          size="small"
          severity="primary"
          :loading="localSaving"
          :disabled="localSaving"
          @click="saveInlineEdits"
        />
        <Button
          v-tooltip.bottom="$t('block.recordList.tooltip.discardChanges')"
          icon="pi pi-times"
          text
          size="small"
          severity="secondary"
          :disabled="localSaving"
          @click="cancelInlineEdits"
        />
      </div>

      <div
        v-if="!displayedFields.length"
        class="flex items-center justify-center h-full p-6 text-center text-muted-color italic"
      >
        {{ $t('block.record.noFields') }}
      </div>

      <div v-else ref="fieldContainer" class="p-4 overflow-y-auto flex-1" :class="layoutClass">
        <div
          v-for="field in displayedFields"
          :key="field.fieldID || field.name"
          class="field-item"
          :class="fieldContainerClass"
        >
          <!-- Horizontal layout: label and value side-by-side -->
          <template
            v-if="
              options.horizontalFieldLayoutEnabled && options.recordFieldLayoutOption !== 'noWrap'
            "
          >
            <div class="grid grid-cols-[auto_1fr] gap-x-4 items-start">
              <div class="flex flex-col min-w-32">
                <div class="flex items-center gap-1.5">
                  <label class="text-sm font-medium text-muted-color uppercase tracking-wide">
                    {{ fieldLabel(field) }}
                  </label>
                  <span v-if="field.isRequired && isAnyEditing(field)" class="text-red-500">*</span>
                  <!-- Inline edit button -->
                  <button
                    v-if="showInlineEditButton(field)"
                    class="text-muted-color hover:text-primary transition-colors p-0.5"
                    :title="$t('block.record.inlineEdit.button.title')"
                    @click="startFieldEdit(field)"
                  >
                    <i class="pi pi-pencil text-xs" />
                  </button>
                  <!-- Copy field value button -->
                  <button
                    v-if="showCopyFieldButton(field)"
                    class="text-muted-color hover:text-primary transition-colors p-0.5"
                    :title="$t('block.record.inlineCopy.button.title')"
                    @click="copyFieldValue(field)"
                  >
                    <i class="pi pi-copy text-xs" />
                  </button>
                </div>
                <!-- Field hint -->
                <small
                  v-if="fieldHint(field)"
                  class="text-muted-color"
                  v-tooltip.top="fieldHint(field)"
                >
                  <i class="pi pi-info-circle text-xs" />
                </small>
                <!-- Field description -->
                <small v-if="fieldDescription(field)" class="text-muted-color mt-0.5">
                  {{ fieldDescription(field) }}
                </small>
              </div>
              <div class="field-value text-color min-h-8">
                <template v-if="isFieldEditable(field)">
                  <FormField :name="field.name" v-slot="{ invalid, error }">
                    <CFieldEditor
                      :field="field"
                      :namespace="namespace"
                      :model-value="getFieldValue(field)"
                      @update:model-value="setFieldValue(field, $event)"
                    />
                    <Message v-if="invalid" severity="error" size="small" variant="simple">
                      {{ error?.message }}
                    </Message>
                  </FormField>
                </template>

                <CFieldViewer
                  v-else-if="field.canReadRecordValue !== false"
                  :field="field"
                  :record="activeRecord"
                  :namespace="namespace"
                />

                <span v-else class="text-muted-color italic text-sm">
                  {{ $t('block.field.noPermission') }}
                </span>
              </div>
            </div>
          </template>

          <!-- Default vertical layout -->
          <template v-else>
            <div class="flex items-center gap-1.5 mb-1">
              <label
                v-if="field.kind !== 'Bool' || field.options?.switch || !isEditing"
                class="text-sm font-medium text-muted-color uppercase tracking-wide block"
              >
                {{ fieldLabel(field) }}
              </label>
              <span v-if="field.isRequired && isAnyEditing(field)" class="text-red-500">*</span>
              <!-- Inline edit button -->
              <button
                v-if="showInlineEditButton(field)"
                class="text-muted-color hover:text-primary transition-colors p-0.5"
                :title="$t('block.record.inlineEdit.button.title')"
                @click="startFieldEdit(field)"
              >
                <i class="pi pi-pencil text-xs" />
              </button>
              <!-- Copy field value button -->
              <button
                v-if="showCopyFieldButton(field)"
                class="text-muted-color hover:text-primary transition-colors p-0.5"
                :title="$t('block.record.inlineCopy.button.title')"
                @click="copyFieldValue(field)"
              >
                <i class="pi pi-copy text-xs" />
              </button>
              <!-- Field hint -->
              <span
                v-if="fieldHint(field)"
                v-tooltip.top="fieldHint(field)"
                class="text-muted-color cursor-help"
              >
                <i class="pi pi-info-circle text-xs" />
              </span>
            </div>
            <!-- Field description -->
            <small v-if="fieldDescription(field)" class="text-muted-color block mb-1">
              {{ fieldDescription(field) }}
            </small>

            <div class="field-value text-color min-h-8">
              <!-- Editor -->
              <template v-if="isFieldEditable(field)">
                <FormField :name="field.name" v-slot="{ invalid, error }">
                  <CFieldEditor
                    :field="field"
                    :namespace="namespace"
                    :model-value="getFieldValue(field)"
                    @update:model-value="setFieldValue(field, $event)"
                  />
                  <Message v-if="invalid" severity="error" size="small" variant="simple">
                    {{ error?.message }}
                  </Message>
                </FormField>
              </template>

              <!-- Viewer -->
              <CFieldViewer
                v-else-if="field.canReadRecordValue !== false"
                :field="field"
                :record="activeRecord"
                :namespace="namespace"
              />

              <span v-else class="text-muted-color italic text-sm">
                {{ $t('block.field.noPermission') }}
              </span>
            </div>
          </template>
        </div>
      </div>
    </template>
  </PageBlock>
</template>

<script setup>
import { computed, inject, onBeforeUnmount, provide, reactive, ref, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import { compose } from '@planetcrust/human-js'
const { CFieldViewer, CFieldEditor } = components
import { useModuleStore } from '@planetcrust/human-vue'
import { useRecordStore } from '@planetcrust/human-vue'
import PageBlock from './PageBlock.vue'

const props = defineProps({
  block: {
    type: Object,
    required: true,
  },
  namespace: {
    type: Object,
    default: () => ({}),
  },
  page: {
    type: Object,
    default: () => ({}),
  },
})

const { t } = useI18n()
const route = useRoute()
const moduleStore = useModuleStore()
const recordStore = useRecordStore()
const $ComposeAPI = inject('$ComposeAPI')
const $SystemAPI = inject('$SystemAPI', null)
const $Auth = inject('$Auth', {})
const $toast = inject('$toast', null)

// Inject edit context from RecordView (may be null on non-record pages)
const ctx = inject('recordViewContext', null)

const loading = ref(false)
const localRecord = ref(null)
const referenceRecord = ref(null)
const referenceModule = ref(null)
const fieldContainer = ref(null)

// Per-field inline edit state (only used when NOT on an edit page)
const activeEditFieldNames = ref([])
const localDirtyValues = reactive({})
const localSaving = ref(false)
const hasActiveInlineEdits = computed(() => activeEditFieldNames.value.length > 0)

// Field condition tracking
const hiddenConditions = ref([]) // array of fieldIDs/names that should be hidden
// False until an evaluation has answered for the conditioned fields
const conditionsResolved = ref(false)

// ResizeObserver state
const resizeObserver = ref(null)
const columnWrapClass = ref('')

const options = computed(() => props.block.options || {})

// Builder mode detection
const isBuilder = computed(() => route.name === 'admin.pages.builder')

// Dummy record for builder preview
const builderRecord = ref(null)

// True when the page itself is in edit/create mode (record edit page)
const isOnEditPage = computed(() => !!ctx && ctx.mode.value !== 'view')

// Whether we're in page-level edit or create mode
const isEditing = computed(() => {
  if (isBuilder.value) return true
  return isOnEditPage.value
})

// The page's module
const pageModule = computed(() => {
  const moduleID = props.page?.moduleID
  if (!moduleID) return null
  return moduleStore.getByID(moduleID) || null
})

// Resolve the module: use referenceModule if reference field is set, otherwise page's module
const fieldModule = computed(() => {
  if (options.value.referenceField && referenceModule.value) {
    return referenceModule.value
  }
  return pageModule.value
})

// The active record: use reference record (if reference field), context record (if ctx available), builder record, or local record
const activeRecord = computed(() => {
  // When a reference field is configured, only show the referenced record — never fall back
  // to the page's own record (which would show unrelated data when the selector is empty).
  if (options.value.referenceField) {
    return referenceRecord.value
  }
  if (isBuilder.value) {
    return builderRecord.value
  }
  // When a record view context is present (record page or modal), always use its record.
  // This covers both edit and view modes, including modal where route.params.recordID is absent.
  if (ctx) {
    return ctx.record.value
  }
  return localRecord.value
})

// Provide the resolved record so nested field editors (e.g. a Record field's
// prefilter) can interpolate ${record...}/${recordID}/${ownerID} expressions
provide('$recordContext', activeRecord)

// Determine which fields to display (before condition filtering)
const visibleFields = computed(() => {
  if (!fieldModule.value) return []

  const configuredFields = options.value.fields || []

  if (configuredFields.length === 0) {
    // No fields configured — show all module fields (non-system)
    return fieldModule.value.fields || []
  }

  // Show only configured fields in order
  if (fieldModule.value.filterFields) {
    return fieldModule.value.filterFields(configuredFields)
  }
  // Fallback: filter module fields to only those listed in block options
  const names = configuredFields.map(f => f.name ?? f)
  return (fieldModule.value.fields || []).filter(f => names.includes(f.name))
})

// Apply field conditions filtering
const displayedFields = computed(() => {
  return visibleFields.value.filter(field => canDisplay(field))
})

// Layout class based on block options
const layoutClass = computed(() => {
  const layout = options.value.recordFieldLayoutOption || 'default'
  const classes = {
    default: 'flex flex-col gap-5',
    noWrap: 'flex gap-6',
    wrap: 'flex flex-wrap',
  }
  return classes[layout] || classes.default
})

const fieldContainerClass = computed(() => {
  const layout = options.value.recordFieldLayoutOption || 'default'
  if (layout === 'noWrap') return 'min-w-52'
  if (layout === 'wrap') return columnWrapClass.value
  return ''
})

// --- Field label ---
function fieldLabel(field) {
  if (field.isSystem) {
    return t(`field.system.${field.name}`, field.label || field.name)
  }
  return field.label || field.name
}

// --- Field hint ---
function fieldHint(field) {
  return field.options?.hint?.view || ''
}

// --- Field description ---
function fieldDescription(field) {
  return field.options?.description?.view || ''
}

// --- Field editability ---
function canFieldBeEdited(field) {
  if (!field) return false
  if (field.canReadRecordValue === false) return false
  if (field.canUpdateRecordValue === false) return false
  if (field.isSystem) {
    if (field.name !== 'ownedBy') return false
    const record = activeRecord.value
    const mod = fieldModule.value
    return record?.createdAt
      ? record.canManageOwnerOnRecord !== false
      : mod?.canCreateOwnedRecord !== false
  }
  return !field.expressions?.value
}

function isFieldEditable(field) {
  if (!field) return false

  // Reference record blocks are always read-only; only inline edit applies.
  if (options.value.referenceField) {
    return activeEditFieldNames.value.includes(field.name)
  }

  if (isBuilder.value) return true

  if (isOnEditPage.value) {
    return canFieldBeEdited(field)
  }

  // Local per-field inline edit mode
  return activeEditFieldNames.value.includes(field.name)
}

// Whether required * or bool-label logic should treat field as "being edited"
function isAnyEditing(field) {
  return isEditing.value || activeEditFieldNames.value.includes(field.name)
}

// --- Inline edit ---
function showInlineEditButton(field) {
  if (!options.value.inlineRecordEditEnabled) return false
  if (isOnEditPage.value || isBuilder.value) return false
  if (activeEditFieldNames.value.includes(field.name)) return false
  if (activeRecord.value?.deletedAt) return false
  return canFieldBeEdited(field)
}

function startFieldEdit(field) {
  if (!activeEditFieldNames.value.includes(field.name)) {
    activeEditFieldNames.value = [...activeEditFieldNames.value, field.name]
  }
}

// --- Inline copy ---
function showCopyFieldButton(field) {
  if (!options.value.inlineRecordCopyEnabled) return false
  if (isBuilder.value) return false
  if (activeEditFieldNames.value.includes(field.name)) return false
  if (!activeRecord.value) return false
  if (field.canReadRecordValue === false) return false
  const val = activeRecord.value.values?.[field.name]
  if (val === undefined || val === null) return false
  if (Array.isArray(val) && val.length === 0) return false
  if (val === '') return false
  return true
}

function formatFieldValueForClipboard(field) {
  const r = activeRecord.value
  if (!r) return ''
  const val = r.values?.[field.name]
  if (val === undefined || val === null) return ''
  if (Array.isArray(val)) return val.map(v => (v == null ? '' : String(v))).join('\n')
  return String(val)
}

async function copyFieldValue(field) {
  const text = formatFieldValueForClipboard(field)
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const ta = document.createElement('textarea')
      ta.value = text
      ta.style.position = 'fixed'
      ta.style.opacity = '0'
      document.body.appendChild(ta)
      ta.select()
      document.execCommand('copy')
      document.body.removeChild(ta)
    }
    $toast?.toastSuccess(t('block.record.inlineCopy.success'))
  } catch (e) {
    console.error('Failed to copy field value:', e)
    $toast?.toastErrorHandler(
      t('block.record.inlineCopy.error'),
      t('block.record.inlineCopy.errorSummary'),
    )(e)
  }
}

async function saveInlineEdits() {
  const record = activeRecord.value
  if (!record) return

  localSaving.value = true
  try {
    Object.entries(localDirtyValues).forEach(([fieldName, value]) => {
      record.setValue(fieldName, value)
    })
    const saved = await recordStore.update(record)
    // The response carries what the server computed (value expressions,
    // formatters); handing it to the record view keeps conditions and layout
    // reading the saved record rather than the one we sent.
    if (ctx?.adoptSaved) ctx.adoptSaved(saved)
    else if (!ctx) localRecord.value = saved
    activeEditFieldNames.value = []
    Object.keys(localDirtyValues).forEach(k => delete localDirtyValues[k])
  } catch (e) {
    console.error('Failed to save inline edits:', e)
    $toast?.toastErrorHandler(
      t('block.record.inlineEdit.saveError'),
      t('block.record.inlineEdit.saveErrorSummary'),
    )(e)
  } finally {
    localSaving.value = false
  }
}

function cancelInlineEdits() {
  Object.keys(localDirtyValues).forEach(k => delete localDirtyValues[k])
  activeEditFieldNames.value = []
}

// --- Field value helpers ---
function getFieldValue(field) {
  if (isOnEditPage.value || isBuilder.value) {
    const r = ctx?.record?.value || builderRecord.value
    if (!r) return field.isMulti ? [] : ''
    const val = r.values[field.name]
    return val === undefined || val === null ? (field.isMulti ? [] : '') : val
  }
  // Local inline edit: serve dirty value if present, else record value
  if (field.name in localDirtyValues) {
    return localDirtyValues[field.name]
  }
  const r = activeRecord.value
  if (!r) return field.isMulti ? [] : ''
  const val = r.values[field.name]
  return val === undefined || val === null ? (field.isMulti ? [] : '') : val
}

function setFieldValue(field, value) {
  if (isOnEditPage.value || isBuilder.value) {
    const r = ctx?.record?.value || builderRecord.value
    if (!r) return
    r.setValue(field.name, value)
  } else {
    localDirtyValues[field.name] = value
  }
}

// --- Field conditions ---

// Fields a condition speaks for, by the id canDisplay matches on
const conditionedFieldIDs = computed(() =>
  (options.value.fieldConditions || []).map(({ field }) => field).filter(Boolean),
)

/**
 * A conditioned field is hidden until the first evaluation answers for it:
 * showing it and taking it away once the result lands reads as a glitch, and
 * briefly puts a value on screen the condition means to withhold.
 */
function canDisplay({ fieldID, name }) {
  const id = fieldID && fieldID !== '0' ? fieldID : name
  if (!conditionsResolved.value) return !conditionedFieldIDs.value.includes(id)
  if (hiddenConditions.value.length === 0) return true
  return !hiddenConditions.value.includes(id)
}

let _fieldConditionTimer = null
let _fieldConditionSeq = 0

/**
 * Field conditions follow the record as it is edited — they re-evaluate on
 * every value change, debounced so that typing does not fire a request per
 * keystroke. The first pass runs at once instead: it is what decides whether a
 * conditioned field appears at all, so the debounce would only delay the form.
 */
function evaluateExpressions({ immediate = false } = {}) {
  clearTimeout(_fieldConditionTimer)
  if (immediate) {
    runFieldConditions()
    return
  }
  _fieldConditionTimer = setTimeout(runFieldConditions, 300)
}

async function runFieldConditions() {
  const fieldConditions = options.value.fieldConditions || []
  // Nothing will answer for the conditioned fields on these paths, so release
  // them rather than leaving them hidden on a first-paint that never resolves.
  if (!fieldConditions.length) return resolveConditions()
  // Don't evaluate in builder mode
  if (route.name === 'admin.pages.builder') return resolveConditions()
  if (!$SystemAPI) return resolveConditions()

  const expressions = {}
  const record = activeRecord.value
  const serialized = record?.serialize ? record.serialize() : {}
  const isNew = ctx?.isNew?.value ?? false
  const isEditMode = ctx ? ctx.mode.value !== 'view' && !isNew : false
  const variables = {
    user: $Auth?.user || {},
    record: serialized,
    screen: {
      width: window.innerWidth,
      height: window.innerHeight,
      userAgent: navigator.userAgent,
      breakpoint: (() => {
        const w = window.innerWidth
        if (w >= 1200) return 'lg'
        if (w >= 996) return 'md'
        if (w >= 768) return 'sm'
        if (w >= 480) return 'xs'
        return 'xxs'
      })(),
    },
    isView: !isEditMode && !isNew,
    isCreate: isNew,
    isEdit: isEditMode,
  }

  fieldConditions.forEach(({ field, condition }) => {
    if (field && condition) {
      expressions[field] = condition
    }
  })

  if (Object.keys(expressions).length === 0) return resolveConditions()

  const seq = ++_fieldConditionSeq

  try {
    const res = await $SystemAPI.expressionEvaluate({ variables, expressions })

    // A slower earlier response must not overwrite a newer one, nor clear
    // values based on a stale view of the record
    if (seq !== _fieldConditionSeq) return

    const previousConditions = [...hiddenConditions.value]
    const newHidden = []

    Object.keys(res).forEach(v => {
      if (!res[v]) newHidden.push(v)
    })

    hiddenConditions.value = newHidden
    resolveConditions()

    // Clear values for newly hidden fields
    clearValuesForHiddenFields(previousConditions)
  } catch (e) {
    console.error('Failed to evaluate field conditions:', e)
    // A failed evaluation must not leave the conditioned fields hidden for good
    resolveConditions()
  }
}

function resolveConditions() {
  conditionsResolved.value = true
}

/**
 * Clearing writes to the record every block on the page shares, so it is
 * confined to the modes that lead to a save. In view mode there is nothing to
 * submit and the write would only take the value off the screen — including
 * out of other blocks showing the same field.
 */
function clearValuesForHiddenFields(previousConditions) {
  if (!isOnEditPage.value) return

  const newlyHidden = hiddenConditions.value.filter(id => !previousConditions.includes(id))
  if (newlyHidden.length === 0) return

  const clearAllOnHide = options.value.clearConditionalFieldsOnHide || false
  const fieldConditions = options.value.fieldConditions || []
  const mod = fieldModule.value
  const rec = activeRecord.value

  fieldConditions.forEach(({ field, clearOnHide }) => {
    const shouldClear = clearAllOnHide || clearOnHide
    if (!shouldClear) return
    if (!newlyHidden.includes(field)) return

    const moduleField = mod?.fields?.find(f => f.fieldID === field || f.name === field)
    if (!moduleField || !rec?.values) return

    const fieldName = moduleField.name
    if (moduleField.isMulti) {
      rec.values[fieldName] = []
    } else {
      rec.values[fieldName] = undefined
    }
  })
}

// --- Reference field support ---
async function fetchReferenceModule(moduleID) {
  if (!moduleID) {
    referenceModule.value = null
    return
  }

  try {
    const mod = await moduleStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID,
    })
    referenceModule.value = mod

    if (options.value.referenceField) {
      await loadReferenceRecord()
    }
  } catch (e) {
    console.error('Failed to fetch reference module:', e)
    referenceModule.value = null
  }
}

async function loadReferenceRecord() {
  const mod = referenceModule.value
  if (!mod) return

  const { referenceField } = options.value
  const currentRecord = ctx?.record?.value || localRecord.value
  if (!currentRecord || !pageModule.value) return

  const field = pageModule.value.fields.find(f => f.fieldID === referenceField)
  if (!field) {
    referenceRecord.value = null
    return
  }

  // Multi-value record selectors are not supported for reference
  if (field.isMulti) {
    referenceRecord.value = null
    return
  }

  const recordID = currentRecord.values[field.name]
  if (!recordID) {
    referenceRecord.value = null
    return
  }

  try {
    const rec = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: mod.moduleID,
      recordID,
    })
    referenceRecord.value = rec
  } catch (e) {
    console.error('Failed to load reference record:', e)
    referenceRecord.value = null
  }
}

// --- Record loading (view mode) ---
async function loadRecord() {
  // When a record view context is injected (record page or modal), the record is managed
  // by RecordView — no need to load it here.
  if (ctx) return
  if (!pageModule.value) return

  const recordID = route.params?.recordID
  if (!recordID || recordID === '0') return

  loading.value = true

  try {
    localRecord.value = await recordStore.findByID({
      namespaceID: props.namespace.namespaceID,
      moduleID: pageModule.value.moduleID,
      recordID,
      force: true,
    })
  } catch (e) {
    console.error('Failed to load record for Record block:', e)
    localRecord.value = null
  } finally {
    loading.value = false
  }
}

// --- Responsive wrap layout (ResizeObserver) ---
function initializeResizeObserver(el) {
  if (!el) return
  if (resizeObserver.value) {
    resizeObserver.value.disconnect()
  }

  resizeObserver.value = new ResizeObserver(entries => {
    for (const entry of entries) {
      applyColumnClasses(entry.contentRect.width)
    }
  })

  resizeObserver.value.observe(el)
}

function applyColumnClasses(width) {
  const breakpoints = {
    xs: 576,
    md: 768,
    lg: 992,
    xl: 1200,
  }

  // Tailwind equivalents for Bootstrap column classes
  const columnClasses = {
    xs: 'w-full px-3 mb-4', // col-12
    md: 'w-1/2 px-3 mb-4', // col-6
    lg: 'w-1/3 px-3 mb-4', // col-4
    xl: 'w-1/4 px-3 mb-4', // col-3
  }

  let columnClass
  if (width <= breakpoints.xs) {
    columnClass = columnClasses.xs
  } else if (width <= breakpoints.md) {
    columnClass = columnClasses.md
  } else if (width <= breakpoints.lg) {
    columnClass = columnClasses.lg
  } else {
    columnClass = columnClasses.xl
  }

  columnWrapClass.value = columnClass
}

// --- Watchers ---

// Create dummy record for builder preview. A record needs a module WITH fields
// (compose.Record throws otherwise), so wait until the module's fields are
// loaded/present — guards both a freshly created fieldless module and the brief
// window before fields arrive in the store.
watch(
  () => [isBuilder.value, fieldModule.value],
  ([builder, mod]) => {
    builderRecord.value = builder && mod?.fields?.length ? new compose.Record(mod) : null
  },
  { immediate: true },
)

// Load record in view mode
watch(
  () => [pageModule.value?.moduleID, route.params?.recordID, isEditing.value],
  () => loadRecord(),
  { immediate: true },
)

// Load reference module when referenceModuleID option changes
watch(
  () => options.value.referenceModuleID,
  moduleID => {
    if (moduleID) {
      fetchReferenceModule(moduleID)
    } else {
      referenceModule.value = null
      referenceRecord.value = null
    }
  },
  { immediate: true },
)

// Reload reference record when parent record's reference field value changes
watch(
  () => {
    if (!options.value.referenceField || !pageModule.value) return null
    const rec = ctx?.record?.value || localRecord.value
    if (!rec) return null
    const field = pageModule.value.fields.find(f => f.fieldID === options.value.referenceField)
    return field ? rec.values[field.name] : null
  },
  (newVal, oldVal) => {
    if (newVal !== oldVal && referenceModule.value) {
      loadReferenceRecord()
    }
  },
)

// Evaluate field conditions when record loaded or changes. Arriving at a record
// is the first-paint case, so it skips the debounce; with no record to evaluate
// against there is nothing to wait for and the fields are released.
watch(
  () => activeRecord.value,
  rec => {
    if (rec) {
      evaluateExpressions({ immediate: true })
    } else {
      resolveConditions()
    }
  },
  { immediate: true },
)

// Re-evaluate when record values change (edit mode)
watch(
  () => activeRecord.value?.values,
  () => {
    if (activeRecord.value) {
      evaluateExpressions()
    }
  },
  { deep: true },
)

// Setup ResizeObserver for wrap layout
watch(
  () => [loading.value, fieldModule.value, options.value.recordFieldLayoutOption],
  () => {
    if (options.value.recordFieldLayoutOption === 'wrap' && !loading.value && fieldModule.value) {
      nextTick(() => {
        initializeResizeObserver(fieldContainer.value)
      })
    } else if (resizeObserver.value) {
      resizeObserver.value.disconnect()
      resizeObserver.value = null
      columnWrapClass.value = ''
    }
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  if (resizeObserver.value) {
    resizeObserver.value.disconnect()
  }
})
</script>
