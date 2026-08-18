<template>
  <div class="flex flex-col gap-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <CFormGroup :label="$t('field.kind.bool.checkedValueLabel')">
        <InputText
          v-model="field.options.trueLabel"
          :placeholder="$t('field.kind.bool.checkedValuePlaceholder')"
          class="w-full"
        />
      </CFormGroup>

      <CFormGroup :label="$t('field.kind.bool.uncheckedValueLabel')">
        <InputText
          v-model="field.options.falseLabel"
          :placeholder="$t('field.kind.bool.uncheckedValuePlaceholder')"
          class="w-full"
        />
      </CFormGroup>
    </div>

    <CFormGroup :label="$t('field.kind.bool.toggleTypeLabel')">
      <div class="flex items-center gap-2">
        <Checkbox v-model="field.options.switch" inputId="boolSwitch" :binary="true" />
        <label for="boolSwitch" class="cursor-pointer text-sm">
          {{ $t('field.kind.bool.toggleType') }}
        </label>
      </div>
    </CFormGroup>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const field = inject('fieldDraft')
</script>
