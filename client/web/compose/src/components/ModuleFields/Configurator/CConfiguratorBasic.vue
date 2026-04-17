<template>
  <div class="flex flex-col gap-3">
    <!-- Basic Flags -->
    <div class="flex flex-col gap-3">
      <div class="flex items-center gap-2">
        <Checkbox
          v-model="field.isRequired"
          inputId="isRequired"
          :binary="true"
          :disabled="!field.cap?.required || showValueExpr"
        />
        <label for="isRequired" class="cursor-pointer">{{ $t('general.label.required') }}</label>
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
          :disabled="field.isRequired || defaultValueEnabled"
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

    <Divider layout="horizontal" v-if="showValueExpr || (showDefaultValue && defaultValueEnabled)" />

    <!-- Value Expression / Default Value Editor -->
    <div v-if="showValueExpr" class="flex flex-col gap-1">
      <label class="font-medium text-primary">{{ $t('field.valueExpr.label') }}</label>
      <small class="text-muted-color">{{ $t('field.valueExpr.description') }}</small>
      <InputText
        v-model="valueExpression"
        class="w-full mt-1"
        :placeholder="$t('field.valueExpr.placeholder')"
      />
    </div>

    <div v-else-if="defaultValueEnabled && showDefaultValue" class="flex flex-col gap-1">
      <label class="font-medium text-primary">{{ $t('field.defaultFieldValue') }}</label>
      <CFieldEditor
        :field="mockField"
        :namespace="namespace"
        :model-value="mockValue"
        @update:model-value="onMockValueChange"
        class="mt-1"
      />
    </div>

    <Divider layout="horizontal" v-if="showValueExpr || (showDefaultValue && defaultValueEnabled)" />

    <!-- Description (View / Edit) -->
    <div class="flex flex-col gap-5 mt-2">
      <div class="flex flex-col gap-2">
        <FloatLabel variant="on">
          <InputText id="desc-default" v-model="field.options.description.view" class="w-full" />
          <label for="desc-default">
            {{ $t(`field.options.description.label.${sameDescription ? 'default' : 'view'}`) }}
          </label>
        </FloatLabel>
      </div>

      <div v-if="!sameDescription" class="flex flex-col gap-2">
        <FloatLabel variant="on">
          <InputText id="desc-edit" v-model="field.options.description.edit" class="w-full" />
          <label for="desc-edit">
            {{ $t('field.options.description.label.edit') }}
          </label>
        </FloatLabel>
      </div>

      <div class="flex items-center gap-2 mt-1">
        <Checkbox v-model="sameDescription" inputId="sameDesc" :binary="true" />
        <label for="sameDesc" class="text-sm cursor-pointer">
          {{ $t('field.options.description.same') }}
        </label>
      </div>
    </div>

    <Divider layout="horizontal" />

    <!-- Hint (View / Edit) -->
    <div class="flex flex-col gap-5 mt-2">
      <div class="flex flex-col gap-2">
        <FloatLabel variant="on">
          <InputText id="hint-default" v-model="field.options.hint.view" class="w-full" />
          <label for="hint-default">
            {{ $t(`field.options.hint.label.${sameHint ? 'default' : 'view'}`) }}
          </label>
        </FloatLabel>
      </div>

      <div v-if="!sameHint" class="flex flex-col gap-2">
        <FloatLabel variant="on">
          <InputText id="hint-edit" v-model="field.options.hint.edit" class="w-full" />
          <label for="hint-edit">
            {{ $t('field.options.hint.label.edit') }}
          </label>
        </FloatLabel>
      </div>

      <div class="flex items-center gap-2 mt-1">
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
import { components } from '@planetcrust/human-vue'

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
const showValueExpr = ref(false)
const valueExpression = computed({
  get: () => field.value.expressions?.value || '',
  set: (val) => {
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
watch(defaultValueEnabled, (val) => {
  if (!val) {
    field.value.defaultValue = []
  } else {
    showValueExpr.value = false
    if (field.value.defaultValue.length === 0) {
      field.value.defaultValue = [{ name: field.value.name, value: '' }]
    }
  }
})

watch(showValueExpr, (val) => {
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
