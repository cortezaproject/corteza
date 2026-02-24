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
    </div>

    <Divider layout="horizontal" />

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
import { onMounted, ref, watch } from 'vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const sameDescription = ref(true)
const sameHint = ref(true)

// Initialize local state based on field options
onMounted(() => {
  if (!props.field.options.description) {
    props.field.options.description = { view: '', edit: '' }
  }
  if (!props.field.options.hint) {
    props.field.options.hint = { view: '', edit: '' }
  }

  sameDescription.value = typeof props.field.options.description.edit === 'undefined'
  sameHint.value = typeof props.field.options.hint.edit === 'undefined'
})

// Sync states back to field object
watch(sameDescription, val => {
  if (val) {
    props.field.options.description.edit = undefined
  } else {
    props.field.options.description.edit = props.field.options.description.edit || ''
  }
})

watch(sameHint, val => {
  if (val) {
    props.field.options.hint.edit = undefined
  } else {
    props.field.options.hint.edit = props.field.options.hint.edit || ''
  }
})
</script>
