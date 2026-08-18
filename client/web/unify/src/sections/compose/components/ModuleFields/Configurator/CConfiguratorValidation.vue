<template>
  <div class="flex flex-col gap-3">
    <!-- Sanitizers -->
    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <label class="font-medium text-muted-color text-sm uppercase tracking-wide">
          {{ $t('field.sanitizers.label') }}
        </label>
        <Button
          icon="pi pi-plus"
          :label="$t('general.label.add')"
          size="small"
          text
          @click="addSanitizer"
        />
      </div>
      <small class="text-muted-color">{{ $t('field.sanitizers.description') }}</small>

      <div v-for="(_, index) in sanitizers" :key="'san-' + index" class="flex items-center gap-2">
        <CInputExpression
          v-model="sanitizers[index]"
          dialect="expr"
          :scope="sanitizerScope"
          class="flex-1"
          :min-lines="1"
          :placeholder="$t('field.sanitizers.expression.placeholder')"
          @update:model-value="syncSanitizers"
        />
        <Button
          icon="pi pi-trash"
          severity="danger"
          text
          size="small"
          @click="removeSanitizer(index)"
        />
      </div>
    </div>

    <Divider layout="horizontal" />

    <!-- Validators -->
    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <label class="font-medium text-muted-color text-sm uppercase tracking-wide">
          {{ $t('field.validators.label') }}
        </label>
        <Button
          icon="pi pi-plus"
          :label="$t('general.label.add')"
          size="small"
          text
          @click="addValidator"
        />
      </div>
      <small class="text-muted-color">{{ $t('field.validators.description') }}</small>

      <div v-for="(v, index) in validators" :key="'val-' + index" class="flex items-center gap-2">
        <CInputExpression
          v-model="v.test"
          dialect="expr"
          :scope="validatorScope"
          class="flex-1"
          :min-lines="1"
          :placeholder="$t('field.validators.expression.placeholder')"
          @update:model-value="syncValidators"
        />
        <InputText
          v-model="v.error"
          class="flex-1"
          :placeholder="$t('field.validators.error.placeholder')"
          @update:model-value="syncValidators"
        />
        <Button
          icon="pi pi-trash"
          severity="danger"
          text
          size="small"
          @click="removeValidator(index)"
        />
      </div>

      <div v-if="validators.length > 0" class="flex items-center gap-2 mt-1">
        <Checkbox
          v-model="disableDefaultValidators"
          inputId="disableDefaultValidators"
          :binary="true"
        />
        <label for="disableDefaultValidators" class="text-sm cursor-pointer">
          {{ $t('field.validators.disableBuiltIn') }}
        </label>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, inject, onMounted, ref } from 'vue'
import { buildFieldExprScope } from '@planetcrust/human-vue'

const field = inject('fieldDraft')

// The server hands each of these a different scope — a sanitizer sees only the
// value, a validator also sees its previous value and the record's other
// fields (server/compose/service/values/).
const fieldModule = inject('fieldModule', null)
const sanitizerScope = computed(() => buildFieldExprScope('sanitizer'))
const validatorScope = computed(() => buildFieldExprScope('validator', fieldModule?.value || null))

// -- Sanitizers --
const sanitizers = ref([])

function addSanitizer() {
  sanitizers.value.push('')
  syncSanitizers()
}

function removeSanitizer(index) {
  sanitizers.value.splice(index, 1)
  syncSanitizers()
}

function syncSanitizers() {
  if (!field.value.expressions) field.value.expressions = {}
  field.value.expressions.sanitizers = [...sanitizers.value]
}

// -- Validators --
const validators = ref([])
const disableDefaultValidators = computed({
  get: () => field.value.expressions?.disableDefaultValidators || false,
  set: val => {
    if (!field.value.expressions) field.value.expressions = {}
    field.value.expressions.disableDefaultValidators = val
  },
})

function addValidator() {
  validators.value.push({ test: '', error: '' })
  syncValidators()
}

function removeValidator(index) {
  validators.value.splice(index, 1)
  syncValidators()
}

function syncValidators() {
  if (!field.value.expressions) field.value.expressions = {}
  field.value.expressions.validators = validators.value.map(v => ({ ...v }))
}

onMounted(() => {
  if (!field.value.expressions) {
    field.value.expressions = {}
  }
  sanitizers.value = [...(field.value.expressions.sanitizers || [])]
  validators.value = (field.value.expressions.validators || []).map(v => ({ ...v }))
})
</script>
