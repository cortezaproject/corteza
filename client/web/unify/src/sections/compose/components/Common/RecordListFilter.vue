<template>
  <div class="inline-flex">
    <Button
      ref="filterBtnRef"
      v-tooltip.bottom="$t('block.recordList.filter.title')"
      :icon="hasActiveFilters ? 'pi pi-filter-fill' : 'pi pi-filter'"
      :severity="hasActiveFilters ? 'primary' : 'secondary'"
      text
      size="small"
      @click="toggle"
    />

    <Popover
      ref="popoverRef"
      :pt="{
        root: { class: 'record-list-filter-popover' },
        content: { class: 'p-0' },
      }"
    >
      <div class="flex flex-col" style="width: min(90vw, 850px); max-height: 60vh">
        <!-- Presets dropdown -->
        <div v-if="presets.length" class="flex items-center gap-2 p-3 border-b">
          <i class="pi pi-bookmark text-muted-color" />
          <Select
            :model-value="null"
            :options="presetsWithDelete"
            option-label="name"
            :placeholder="$t('block.recordList.filter.presets')"
            size="small"
            class="flex-1"
            @change="onPresetSelect($event.value)"
          >
            <template #option="{ option }">
              <div class="flex items-center justify-between w-full gap-2">
                <span>{{ option.name }}</span>
                <Button
                  v-if="option.deletable"
                  icon="pi pi-trash"
                  text
                  rounded
                  severity="danger"
                  class="w-6 h-6 p-0 shrink-0"
                  @click.stop="deletePreset(option.index)"
                />
              </div>
            </template>
          </Select>
        </div>

        <!-- Filter rows -->
        <div class="flex-1 overflow-auto p-3">
          <template v-for="(group, gi) in internalFilter" :key="gi">
            <div v-if="group.filter.length" class="flex flex-col gap-2">
              <div
                v-for="(f, fi) in group.filter"
                :key="`${gi}-${fi}`"
                class="flex items-start gap-2"
              >
                <div
                  class="flex items-start gap-2 flex-1 min-w-0 border border-surface rounded-border p-2"
                >
                  <!-- Field picker -->
                  <Select
                    v-model="f.name"
                    :options="filterableFields"
                    option-label="label"
                    option-value="name"
                    :placeholder="$t('block.recordList.filter.fieldPlaceholder')"
                    filter
                    class="w-52 shrink-0"
                    size="small"
                    @change="onFieldChange(gi, fi)"
                  />

                  <!-- Operator -->
                  <Select
                    v-if="f.name"
                    v-model="f.operator"
                    :options="getOperators(f.kind, getField(f.name))"
                    option-label="text"
                    option-value="value"
                    class="w-40 shrink-0"
                    size="small"
                    @change="onOperatorChange(gi, fi)"
                  />

                  <!-- Value editor — omitted for the operators that test only
                       for the presence of a value -->
                  <template v-if="f.name && !isValuelessOperator(f.operator)">
                    <template v-if="isBetween(f.operator)">
                      <div class="flex flex-col gap-1 flex-1 min-w-0">
                        <CFieldEditor
                          :field="makeBetweenField(f.name, 'start')"
                          :namespace="namespace"
                          :model-value="f.value?.start"
                          size="small"
                          @update:model-value="f.value = { ...f.value, start: $event }"
                        />
                        <span class="text-center text-xs text-muted-color">
                          {{ $t('block.general.label.and') }}
                        </span>
                        <CFieldEditor
                          :field="makeBetweenField(f.name, 'end')"
                          :namespace="namespace"
                          :model-value="f.value?.end"
                          size="small"
                          @update:model-value="f.value = { ...f.value, end: $event }"
                        />
                      </div>
                    </template>
                    <CFieldEditor
                      v-else
                      :field="getEditorField(f.name, f.operator)"
                      :namespace="namespace"
                      :model-value="isMultiOperator(f.operator) ? ensureArray(f.value) : f.value"
                      :allow-empty="isMultiOperator(f.operator)"
                      size="small"
                      class="flex-1 min-w-0"
                      @update:model-value="f.value = $event"
                    />
                  </template>
                </div>

                <!-- Delete row (outside the card) -->
                <Button
                  v-if="f.name"
                  icon="pi pi-trash"
                  text
                  severity="danger"
                  size="small"
                  class="shrink-0 mt-2"
                  @click="deleteFilter(gi, fi)"
                />
              </div>

              <!-- Add condition -->
              <Button
                v-if="group.filter[0]?.name"
                :label="$t('general.label.add')"
                icon="pi pi-plus"
                severity="secondary"
                text
                size="small"
                class="self-start"
                @click="addFilter(gi)"
              />
            </div>

            <!-- Group separator / add group -->
            <div class="flex items-center my-3" v-if="group.filter[0]?.name">
              <Divider class="flex-1" />
              <Select
                v-if="gi < internalFilter.length - 1"
                v-model="internalFilter[gi + 1].groupCondition"
                :options="groupConditionOptions"
                option-label="text"
                option-value="value"
                class="mx-2"
                size="small"
              />
              <Button
                v-else
                icon="pi pi-plus"
                severity="secondary"
                outlined
                rounded
                size="small"
                class="mx-2"
                @click="addGroup"
              />
              <Divider class="flex-1" />
            </div>
          </template>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between gap-2 p-3 border-t">
          <Button
            :label="$t('block.recordList.filter.reset')"
            severity="secondary"
            text
            size="small"
            @click="resetFilter"
          />
          <div class="flex items-center gap-2">
            <Button
              v-if="allowPresetSave && hasValidFilters"
              :label="$t('block.recordList.filterPresets.saveFilterAsPreset')"
              icon="pi pi-bookmark"
              severity="secondary"
              outlined
              size="small"
              @click="showSaveDialog = true"
            />
            <Button
              :label="$t('block.recordList.filter.update')"
              severity="primary"
              size="small"
              @click="onSave"
            />
          </div>
        </div>
      </div>
    </Popover>

    <!-- Save preset dialog -->
    <Dialog
      v-model:visible="showSaveDialog"
      :header="$t('block.recordList.filterPresets.saveFilterAsPreset')"
      modal
      :style="{ width: '400px' }"
    >
      <div class="flex flex-col gap-3">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium text-color">
            {{ $t('block.recordList.filter.name.label') }}
          </label>
          <InputText
            v-model="presetName"
            :placeholder="$t('block.recordList.filter.name.placeholder')"
            autofocus
            @keyup.enter="confirmSavePreset"
          />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="$t('general.label.cancel')"
            severity="secondary"
            text
            size="small"
            @click="showSaveDialog = false"
          />
          <Button
            :label="$t('general.label.save')"
            severity="primary"
            size="small"
            :disabled="!presetName.trim()"
            @click="confirmSavePreset"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@planetcrust/human-vue'
import {
  isBetweenOperator,
  isValuelessOperator,
  IS_EMPTY,
  IS_NOT_EMPTY,
} from '../../lib/record-filter'

const { CFieldEditor } = components
const { t } = useI18n()

const props = defineProps({
  module: { type: Object, required: true },
  namespace: { type: Object, required: true },
  modelValue: { type: Array, default: () => [] },
  allowPresetSave: { type: Boolean, default: false },
  presets: { type: Array, default: () => [] },
})

const emit = defineEmits([
  'update:modelValue',
  'reset',
  'save-preset',
  'delete-preset',
  'load-preset',
])

const popoverRef = ref(null)
const filterBtnRef = ref(null)

// --- Preset state ---
const showSaveDialog = ref(false)
const presetName = ref('')

// Presets with index for dropdown option template
const presetsWithDelete = computed(() => {
  return props.presets.map((p, i) => ({
    ...p,
    index: i,
    deletable: true,
  }))
})

// --- Internal filter state (editable copy) ---
const internalFilter = ref([])

// Sync from modelValue (external → internal)
watch(
  () => props.modelValue,
  raw => {
    internalFilter.value = raw?.length ? JSON.parse(JSON.stringify(raw)) : [createDefaultGroup()]
  },
  { immediate: true, deep: true },
)

// --- Field helpers ---
const filterableFields = computed(() => {
  if (!props.module) return []
  const regular = (props.module.fields || [])
    .filter(f => f.isFilterable && f.canReadRecordValue !== false)
    .sort((a, b) => (a.label || a.name).localeCompare(b.label || b.name))

  const system = (props.module.systemFields?.() || [])
    .filter(f => f.isFilterable !== false)
    .map(f => ({
      ...f,
      label: t(`field.system.${f.name}`, f.label || f.name),
    }))

  return [...regular, ...system]
})

function getField(name) {
  return filterableFields.value.find(f => f.name === name)
}

function getEditorField(name, operator) {
  const f = getField(name)
  if (!f) return { name, kind: 'String' }
  // For IN/NOT IN operators, preserve isMulti so the editor shows multi-select UI
  // For other operators, force single-value editing
  const keepMulti = isMultiOperator(operator) && f.isMulti
  return {
    ...f,
    isMulti: keepMulti,
    // Ensure multi-absorbing types get selectType=multiple so CFieldEditor handles them
    options: keepMulti
      ? { ...f.options, selectType: 'multiple' }
      : {
          ...f.options,
          selectType: f.options?.selectType === 'multiple' ? 'default' : f.options?.selectType,
        },
  }
}

function isMultiOperator(op) {
  return ['IN', 'NOT IN'].includes(op)
}

function ensureArray(val) {
  if (Array.isArray(val)) return val
  if (val == null || val === '') return []
  return [val]
}

function makeBetweenField(name, suffix) {
  const f = getEditorField(name)
  return { ...f, name: `${name}-${suffix}` }
}

// --- Active filter tracking ---
const hasActiveFilters = computed(() => props.modelValue?.some(g => g.filter?.some(f => f.name)))

const hasValidFilters = computed(() => internalFilter.value?.some(g => g.filter?.some(f => f.name)))

// --- Operators ---
const groupConditionOptions = [
  { value: 'OR', text: t('block.recordList.filter.conditions.or') },
  { value: 'AND', text: t('block.recordList.filter.conditions.and') },
]

function getOperators(kind, field) {
  const eq = [
    { value: '=', text: t('block.recordList.filter.operators.equal') },
    { value: '!=', text: t('block.recordList.filter.operators.notEqual') },
  ]
  const containsOps = [
    { value: 'IN', text: t('block.recordList.filter.operators.contains') },
    { value: 'NOT IN', text: t('block.recordList.filter.operators.notContains') },
  ]
  const cmp = [
    { value: '>', text: t('block.recordList.filter.operators.greaterThan') },
    { value: '<', text: t('block.recordList.filter.operators.lessThan') },
  ]
  const like = [
    { value: 'LIKE', text: t('block.recordList.filter.operators.like') },
    { value: 'NOT LIKE', text: t('block.recordList.filter.operators.notLike') },
  ]
  const between = [
    { value: 'BETWEEN', text: t('block.recordList.filter.operators.between') },
    { value: 'NOT BETWEEN', text: t('block.recordList.filter.operators.notBetween') },
  ]
  // Offered for every kind, multi-value fields included — those carry only
  // Contains/Not contains otherwise, and so had no way to ask for the records
  // where nothing has been picked at all. Geometry is the one kind the server
  // refuses to query on ("attribute can not be used in query expression"), but
  // that is true of its every operator, not just these.
  const empty = [
    { value: IS_EMPTY, text: t('block.recordList.filter.operators.isEmpty') },
    { value: IS_NOT_EMPTY, text: t('block.recordList.filter.operators.isNotEmpty') },
  ]

  if (field?.multi || field?.isMulti) return [...containsOps, ...empty]

  switch (kind) {
    case 'Number':
    case 'DateTime':
      return [...eq, ...cmp, ...between, ...empty]
    case 'String':
    case 'Url':
    case 'Email':
      return [...eq, ...like, ...empty]
    default:
      return [...eq, ...empty]
  }
}

function isBetween(op) {
  return isBetweenOperator(op)
}

// --- CRUD helpers ---
function createDefaultFilter(field) {
  return {
    name: field?.name || '',
    operator: field?.isMulti ? 'IN' : '=',
    value: undefined,
    kind: field?.kind || '',
    condition: '',
  }
}

function createDefaultGroup(field, groupCondition = 'OR') {
  return {
    filter: [createDefaultFilter(field)],
    groupCondition,
  }
}

function onFieldChange(gi, fi) {
  const f = internalFilter.value[gi].filter[fi]
  const field = getField(f.name)
  if (field) {
    f.kind = field.kind
    const multi = field.isMulti
    f.operator = multi ? 'IN' : '='
    f.value = multi ? [] : undefined
  }
}

// Switching to "is empty" drops whatever was typed for the previous operator:
// the editor for it is gone, so a value left behind is one nobody can see, edit
// or clear, and it would ride along into a saved preset.
function onOperatorChange(gi, fi) {
  const f = internalFilter.value[gi].filter[fi]
  if (isValuelessOperator(f.operator)) {
    f.value = undefined
  }
}

function addFilter(gi) {
  const firstField = filterableFields.value[0]
  internalFilter.value[gi].filter.push(createDefaultFilter(firstField))
}

function deleteFilter(gi, fi) {
  internalFilter.value[gi].filter.splice(fi, 1)
  if (!internalFilter.value[gi].filter.length) {
    internalFilter.value.splice(gi, 1)
    if (!internalFilter.value.length) {
      internalFilter.value = [createDefaultGroup()]
    }
  }
}

function addGroup() {
  const firstField = filterableFields.value[0]
  internalFilter.value.push(createDefaultGroup(firstField))
}

// --- Actions ---
function toggle(event) {
  popoverRef.value?.toggle(event)
}

function cleanedFilter() {
  return internalFilter.value
    .map(g => ({
      ...g,
      filter: g.filter.filter(f => f.name),
    }))
    .filter(g => g.filter.length)
}

function onSave() {
  emit('update:modelValue', cleanedFilter())
  popoverRef.value?.hide()
}

function resetFilter() {
  internalFilter.value = [createDefaultGroup()]
  emit('update:modelValue', [])
  emit('reset')
  popoverRef.value?.hide()
}

// --- Preset methods ---
function onPresetSelect(preset) {
  if (!preset) return
  if (Array.isArray(preset.filter)) {
    internalFilter.value = JSON.parse(JSON.stringify(preset.filter))
  }
  emit('load-preset', preset)
}

function deletePreset(index) {
  emit('delete-preset', index)
}

function confirmSavePreset() {
  const name = presetName.value.trim()
  if (!name) return

  const filter = cleanedFilter()
  emit('save-preset', { name, filter })
  showSaveDialog.value = false
  presetName.value = ''
}
</script>

<style>
/* Above the sidebar.
 *
 * The sidebar is a PrimeVue Drawer, and its mask is a viewport-wide fixed
 * element at z-index 1101 — even though the drawer itself is only as wide as the
 * sidebar. At 1040 this popover was painted under it and, worse, the mask
 * swallowed the clicks: the popover is 850px wide and gets placed flush against
 * the left edge, so its field picker sat behind the sidebar and could not be
 * opened at all. It has to outrank the mask to be usable. */
.record-list-filter-popover {
  z-index: 1102 !important;
}
</style>
