<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-2">
      <h3 class="text-lg font-medium text-primary m-0">
        {{ $t('module.edit.config.uniqueValues.duplicationDetection') }}
      </h3>
      <i
        class="pi pi-info-circle text-orange-400"
        v-tooltip="$t('module.edit.config.uniqueValues.tooltip.performance')"
      />
    </div>

    <!-- Rules Loop -->
    <div
      v-for="(rule, index) in rules"
      :key="index"
      class="border rounded-lg p-4 bg-surface shadow-sm"
    >
      <div class="flex items-center justify-between mb-4 pb-3 border-b border-surface">
        <label class="font-bold text-primary">
          {{ $t('module.edit.config.uniqueValues.uniqueValueConstraint', { index: index + 1 }) }}
        </label>
        <Button
          icon="pi pi-trash"
          severity="danger"
          size="small"
          text
          @click="rules.splice(index, 1)"
        />
      </div>

      <div class="flex flex-wrap items-end justify-between gap-4 mb-4">
        <CFormGroup :name="`rule_${index}_field`" class="flex-grow max-w-sm">
          <Select
            v-model="rule.currentField"
            :placeholder="$t('module.edit.config.uniqueValues.searchFields')"
            :options="filterFieldOptions(rule)"
            optionLabel="name"
            optionValue="name"
            class="w-full"
            @change="updateRuleConstraint(rule)"
          >
            <template #option="{ option }">
              {{ getOptionLabel(option) }}
            </template>
            <template #value="{ value, placeholder }">
              <span v-if="value">{{ getOptionLabel(getField(value)) }}</span>
              <span v-else class="text-muted-color">{{ placeholder }}</span>
            </template>
          </Select>
        </CFormGroup>

        <FormField :name="`rule_${index}_strict`" class="flex items-center gap-3">
          <label class="font-medium text-primary mb-0">
            {{ $t('module.edit.config.uniqueValues.preventRecordsSave') }}
          </label>
          <ToggleSwitch v-model="rule.strict" />
        </FormField>
      </div>

      <p
        v-if="rule.constraints && rule.constraints.length > 1"
        class="text-sm text-muted-color mb-3 mt-0"
      >
        {{ $t('module.edit.config.uniqueValues.allFieldsMustMatch') }}
      </p>

      <DataTable
        v-if="rule.constraints && rule.constraints.length > 0"
        :value="rule.constraints"
        size="small"
        responsiveLayout="scroll"
        class="border border-surface"
      >
        <Column :header="$t('module.edit.config.uniqueValues.field')">
          <template #body="{ data }">
            {{ getOptionLabel(getField(data.attribute)) }}
          </template>
        </Column>

        <Column :header="$t('module.edit.config.uniqueValues.type')">
          <template #body="{ data }">
            {{ getField(data.attribute).kind }}
          </template>
        </Column>

        <Column :header="$t('module.edit.config.uniqueValues.valueModifiers')">
          <template #body="{ data }">
            <Select
              v-model="data.modifier"
              :options="modifierOptions"
              optionLabel="text"
              optionValue="value"
              size="small"
              class="w-full"
            />
          </template>
        </Column>

        <Column :header="$t('module.edit.config.uniqueValues.multiValues')">
          <template #body="{ data }">
            <Select
              v-model="data.multiValue"
              :options="multiValueOptions"
              optionLabel="text"
              optionValue="value"
              :disabled="!getField(data.attribute).isMulti"
              size="small"
              class="w-full"
            />
          </template>
        </Column>

        <Column headerStyle="width: 4rem" bodyClass="text-right">
          <template #body="{ index: cIndex }">
            <Button
              icon="pi pi-times"
              severity="secondary"
              text
              size="small"
              @click="rule.constraints.splice(cIndex, 1)"
            />
          </template>
        </Column>
      </DataTable>
    </div>

    <div>
      <Button
        :label="$t('module.edit.config.uniqueValues.addNewConstraint')"
        icon="pi pi-plus"
        severity="secondary"
        @click="addNewConstraint"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'

const module = inject('moduleDraft')

const { t } = useI18n()

// Safely ensure config paths exist
if (!module.value.config) module.value.config = {}
if (!module.value.config.recordDeDup) {
  module.value.config.recordDeDup = { enabled: false, rules: [] }
}

const rules = computed({
  get() {
    return module.value.config.recordDeDup.rules || []
  },
  set(value) {
    module.value.config.recordDeDup.rules = value
  },
})

const modifierOptions = computed(() => {
  const ruleModifiers = rules.value.reduce((acc, { constraints }) => {
    if (!constraints) return acc
    constraints.forEach(({ modifier }) => {
      if (!acc.includes(modifier)) acc.push(modifier)
    })
    return acc
  }, [])

  return [
    { value: 'ignore-case', text: t('module.edit.config.uniqueValues.ignoreCase') },
    { value: 'fuzzy-match', text: t('module.edit.config.uniqueValues.fuzzyMatch'), legacy: true },
    { value: 'sounds-like', text: t('module.edit.config.uniqueValues.soundsLike'), legacy: true },
    { value: 'case-sensitive', text: t('module.edit.config.uniqueValues.caseSensitive') },
  ].filter(({ value, legacy }) => !legacy || ruleModifiers.includes(value))
})

const multiValueOptions = computed(() => {
  return [
    { value: 'one-of', text: t('module.edit.config.uniqueValues.oneOf') },
    { value: 'equal', text: t('module.edit.config.uniqueValues.equal') },
  ]
})

function addNewConstraint() {
  if (!module.value.config.recordDeDup.rules) {
    module.value.config.recordDeDup.rules = []
  }
  module.value.config.recordDeDup.rules.push({
    name: '',
    strict: true,
    constraints: [],
  })
}

function updateRuleConstraint(rule) {
  const currentFieldName = rule.currentField
  const fieldObj = module.value.fields.find(({ name }) => name === currentFieldName)
  
  if (!fieldObj) {
    rule.currentField = undefined
    return
  }

  if (!rule.constraints) {
    rule.constraints = []
  }

  rule.constraints.push({
    attribute: fieldObj.name,
    modifier: 'case-sensitive',
    multiValue: 'equal',
    type: fieldObj.kind,
    isMulti: fieldObj.isMulti,
  })

  // Reset the selector UI state
  rule.currentField = undefined
}

function filterFieldOptions(rule) {
  const selectedFields = rule.constraints ? rule.constraints.map(({ attribute }) => attribute) : []
  return module.value.fields.filter(({ name }) => !selectedFields.includes(name))
}

function getField(attribute) {
  const field = module.value.fields.find(({ name }) => name === attribute)
  return field || {}
}

function getOptionLabel(fieldObj) {
  if (!fieldObj) return ''
  return fieldObj.label || fieldObj.name || fieldObj.kind || ''
}
</script>
