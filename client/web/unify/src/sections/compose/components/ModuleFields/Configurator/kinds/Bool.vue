<template>
  <div class="flex flex-col gap-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- True value label -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.bool.checkedValueLabel') }}
        </label>
        <InputText
          v-model="field.options.trueLabel"
          :placeholder="$t('field.kind.bool.checkedValuePlaceholder')"
          class="w-full"
        />
      </div>

      <!-- False value label -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.bool.uncheckedValueLabel') }}
        </label>
        <InputText
          v-model="field.options.falseLabel"
          :placeholder="$t('field.kind.bool.uncheckedValuePlaceholder')"
          class="w-full"
        />
      </div>
    </div>

    <!-- Switch / Checkbox type -->
    <div>
      <label class="font-medium text-muted-color text-sm block mb-3">
        {{ $t('field.kind.bool.toggleTypeLabel') }}
      </label>
      <div class="flex items-center gap-2">
        <Checkbox v-model="field.options.switch" inputId="boolSwitch" :binary="true" />
        <label for="boolSwitch" class="cursor-pointer">
          {{ $t('field.kind.bool.toggleType') }}
        </label>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const field = inject('fieldDraft')
</script>
