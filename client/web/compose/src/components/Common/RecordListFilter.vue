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
        <!-- Filter rows -->
        <div class="flex-1 overflow-auto p-3">
          <template v-for="(group, gi) in internalFilter" :key="gi">
            <div v-if="group.filter.length" class="flex flex-col gap-2">
              <div
                v-for="(f, fi) in group.filter"
                :key="`${gi}-${fi}`"
                class="flex items-start gap-2"
              >
                <div class="flex items-start gap-2 flex-1 min-w-0 border border-surface rounded-border p-2">
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
                  />

                  <!-- Value editor -->
                  <template v-if="f.name">
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
              v-if="allowPresetSave"
              :label="$t('block.recordList.filter.addFilterToPreset')"
              severity="secondary"
              outlined
              size="small"
              @click="onSavePreset"
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
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { components } from '@cortezaproject/corteza-vue-next'
import { isBetweenOperator } from '../../lib/record-filter'

const { CFieldEditor } = components
const { t } = useI18n()

const props = defineProps({
  module: { type: Object, required: true },
  namespace: { type: Object, required: true },
  modelValue: { type: Array, default: () => [] },
  allowPresetSave: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'reset', 'save-preset'])

const popoverRef = ref(null)
const filterBtnRef = ref(null)

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
      : { ...f.options, selectType: f.options?.selectType === 'multiple' ? 'default' : f.options?.selectType },
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

const activeFilterCount = computed(() =>
  (props.modelValue || []).reduce(
    (count, g) => count + (g.filter?.filter(f => f.name)?.length || 0),
    0,
  ),
)

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

  if (field?.multi || field?.isMulti) return containsOps

  switch (kind) {
    case 'Number':
    case 'DateTime':
      return [...eq, ...cmp, ...between]
    case 'String':
    case 'Url':
    case 'Email':
      return [...eq, ...like]
    default:
      return eq
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

function onSave() {
  const cleaned = internalFilter.value
    .map(g => ({
      ...g,
      filter: g.filter.filter(f => f.name),
    }))
    .filter(g => g.filter.length)

  emit('update:modelValue', cleaned)
  popoverRef.value?.hide()
}

function onSavePreset() {
  const cleaned = internalFilter.value
    .map(g => ({
      ...g,
      filter: g.filter.filter(f => f.name),
    }))
    .filter(g => g.filter.length)

  emit('save-preset', cleaned)
  popoverRef.value?.hide()
}

function resetFilter() {
  internalFilter.value = [createDefaultGroup()]
  emit('update:modelValue', [])
  emit('reset')
  popoverRef.value?.hide()
}
</script>

<style>
.record-list-filter-popover {
  z-index: 1040 !important;
}
</style>
