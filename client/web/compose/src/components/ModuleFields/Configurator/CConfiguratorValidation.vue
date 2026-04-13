<template>
  <div class="flex flex-col gap-3">
    <!-- Sanitizers -->
    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <label class="font-medium">{{ $t('field.sanitizers.label') }}</label>
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
        <InputText
          v-model="sanitizers[index]"
          class="flex-1"
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
        <label class="font-medium">{{ $t('field.validators.label') }}</label>
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
        <InputText
          v-model="v.test"
          class="flex-1"
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
import { computed, onMounted, ref } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

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
  if (!props.field.expressions) props.field.expressions = {}
  props.field.expressions.sanitizers = [...sanitizers.value]
}

// -- Validators --
const validators = ref([])
const disableDefaultValidators = computed({
  get: () => props.field.expressions?.disableDefaultValidators || false,
  set: (val) => {
    if (!props.field.expressions) props.field.expressions = {}
    props.field.expressions.disableDefaultValidators = val
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
  if (!props.field.expressions) props.field.expressions = {}
  props.field.expressions.validators = validators.value.map(v => ({ ...v }))
}

onMounted(() => {
  if (!props.field.expressions) {
    props.field.expressions = {}
  }
  sanitizers.value = [...(props.field.expressions.sanitizers || [])]
  validators.value = (props.field.expressions.validators || []).map(v => ({ ...v }))
})
</script>
