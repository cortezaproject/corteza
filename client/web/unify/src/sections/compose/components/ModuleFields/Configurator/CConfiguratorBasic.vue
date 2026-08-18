<template>
  <div class="flex flex-col gap-3">
    <!-- Basic Flags -->
    <div class="flex flex-col gap-3">
      <div class="flex items-center gap-2">
        <Checkbox
          v-model="field.isRequired"
          inputId="isRequired"
          :binary="true"
          :disabled="!field.cap?.required"
        />
        <label for="isRequired" class="cursor-pointer">
          {{ showValueExpr ? $t('field.valueExpr.requiredLabel') : $t('general.label.required') }}
        </label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox
          v-model="field.isMulti"
          inputId="isMulti"
          :binary="true"
          :disabled="!field.cap?.multi"
        />
        <label for="isMulti" class="cursor-pointer">{{ $t('field.label.multi') }}</label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox
          v-model="showValueExpr"
          inputId="showValueExpr"
          :binary="true"
          :disabled="defaultValueEnabled"
        />
        <label for="showValueExpr" class="cursor-pointer">
          {{ $t('field.valueExpr.label') }}
        </label>
      </div>

      <div class="flex items-center gap-2" v-if="showDefaultValue">
        <Checkbox
          v-model="defaultValueEnabled"
          inputId="defaultValueEnabled"
          :binary="true"
          :disabled="showValueExpr"
        />
        <label for="defaultValueEnabled" class="cursor-pointer">
          {{ $t('field.defaultValue') }}
        </label>
      </div>
    </div>

    <Divider
      layout="horizontal"
      v-if="showValueExpr || (showDefaultValue && defaultValueEnabled)"
    />

    <!-- Value Expression / Default Value Editor -->
    <CFormGroup
      v-if="showValueExpr"
      :label="$t('field.valueExpr.label')"
      :description="$t('field.valueExpr.description')"
    >
      <CInputExpression
        v-model="valueExpression"
        dialect="expr"
        :scope="valueExprScope"
        class="w-full mt-1"
        :placeholder="$t('field.valueExpr.placeholder')"
      />
    </CFormGroup>

    <CFormGroup
      v-else-if="defaultValueEnabled && showDefaultValue"
      :label="$t('field.defaultFieldValue')"
    >
      <CFieldEditor
        :field="mockField"
        :namespace="namespace"
        :model-value="mockValue"
        @update:model-value="onMockValueChange"
        class="mt-1"
      />
    </CFormGroup>

    <Divider
      layout="horizontal"
      v-if="showValueExpr || (showDefaultValue && defaultValueEnabled)"
    />

    <!-- Description (View / Edit) -->
    <div class="flex flex-col gap-3">
      <CFormGroup
        input-id="desc-default"
        :label="$t(`field.options.description.label.${sameDescription ? 'default' : 'view'}`)"
      >
        <InputText id="desc-default" v-model="field.options.description.view" class="w-full" />
      </CFormGroup>

      <CFormGroup
        v-if="!sameDescription"
        input-id="desc-edit"
        :label="$t('field.options.description.label.edit')"
      >
        <InputText id="desc-edit" v-model="field.options.description.edit" class="w-full" />
      </CFormGroup>

      <div class="flex items-center gap-2">
        <Checkbox v-model="sameDescription" inputId="sameDesc" :binary="true" />
        <label for="sameDesc" class="text-sm cursor-pointer">
          {{ $t('field.options.description.same') }}
        </label>
      </div>
    </div>

    <Divider layout="horizontal" />

    <!-- Hint (View / Edit) -->
    <div class="flex flex-col gap-3">
      <CFormGroup
        input-id="hint-default"
        :label="$t(`field.options.hint.label.${sameHint ? 'default' : 'view'}`)"
      >
        <InputText id="hint-default" v-model="field.options.hint.view" class="w-full" />
      </CFormGroup>

      <CFormGroup
        v-if="!sameHint"
        input-id="hint-edit"
        :label="$t('field.options.hint.label.edit')"
      >
        <InputText id="hint-edit" v-model="field.options.hint.edit" class="w-full" />
      </CFormGroup>

      <div class="flex items-center gap-2">
        <Checkbox v-model="sameHint" inputId="sameHint" :binary="true" />
        <label for="sameHint" class="text-sm cursor-pointer">
          {{ $t('field.options.hint.same') }}
        </label>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref, watch } from 'vue'
import { compose } from '@planetcrust/human-js'
import { buildFieldExprScope, components } from '@planetcrust/human-vue'

const { CFieldEditor } = components

defineProps({
  namespace: {
    type: Object,
    default: null,
  },
})

const field = inject('fieldDraft')

const sameDescription = ref(true)
const sameHint = ref(true)

// -- Default Value --
const defaultValueEnabled = ref(false)
const showDefaultValue = computed(() => !['File'].includes(field.value.kind))

const mockField = computed(() => {
  const f = compose.ModuleFieldMaker(JSON.parse(JSON.stringify(field.value)))
  f.label = f.label || 'Default value'
  f.name = 'defValField'
  return f
})

const mockValue = computed(() => {
  if (field.value.isMulti) {
    return field.value.defaultValue.map(v => v.value).filter(v => v !== undefined && v !== null)
  }
  return field.value.defaultValue[0]?.value || undefined
})

function onMockValueChange(val) {
  let dv = val
  if (!Array.isArray(dv)) {
    dv = [dv]
  } else if (!dv.length) {
    dv = [undefined]
  }

  field.value.defaultValue = dv.map(v => {
    let valueStr = v
    if (v !== undefined && v.toString) valueStr = v.toString()
    const def = { name: field.value.name }
    if (valueStr) def.value = valueStr
    return def
  })
}

// -- Value Expression --
// Its scope is the record's own fields as bare names, plus `new` and `old`
// (server/compose/service/values/expr.go) — not the same variables a validator
// or a visibility rule sees.
const fieldModule = inject('fieldModule', null)
const valueExprScope = computed(() => buildFieldExprScope('value', fieldModule?.value || null))

const showValueExpr = ref(false)
const valueExpression = computed({
  get: () => field.value.expressions?.value || '',
  set: val => {
    if (!field.value.expressions) field.value.expressions = {}
    field.value.expressions.value = val || undefined
  },
})

// Initialize local state based on field options
onMounted(() => {
  if (!field.value.options.description) {
    field.value.options.description = { view: '', edit: '' }
  }
  if (!field.value.options.hint) {
    field.value.options.hint = { view: '', edit: '' }
  }

  sameDescription.value = typeof field.value.options.description.edit === 'undefined'
  sameHint.value = typeof field.value.options.hint.edit === 'undefined'

  // Init default value
  if (field.value.defaultValue && field.value.defaultValue.length > 0) {
    defaultValueEnabled.value = true
  }

  // Init expressions
  if (!field.value.expressions) {
    field.value.expressions = {}
  }

  if (field.value.expressions.value && field.value.expressions.value.length > 0) {
    showValueExpr.value = true
  }
})

// Sync default value back to field
watch(defaultValueEnabled, val => {
  if (!val) {
    field.value.defaultValue = []
  } else {
    showValueExpr.value = false
    if (field.value.defaultValue.length === 0) {
      field.value.defaultValue = [{ name: field.value.name, value: '' }]
    }
  }
})

watch(showValueExpr, val => {
  if (val) {
    field.value.isRequired = false
    field.value.defaultValue = []
    defaultValueEnabled.value = false
  } else {
    if (field.value.expressions) {
      field.value.expressions.value = undefined
    }
  }
})

// Sync description/hint states back to field object
watch(sameDescription, val => {
  if (val) {
    field.value.options.description.edit = undefined
  } else {
    field.value.options.description.edit = field.value.options.description.edit || ''
  }
})

watch(sameHint, val => {
  if (val) {
    field.value.options.hint.edit = undefined
  } else {
    field.value.options.hint.edit = field.value.options.hint.edit || ''
  }
})
</script>
