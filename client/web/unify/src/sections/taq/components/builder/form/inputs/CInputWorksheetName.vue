<template>
  <div class="flex flex-col gap-1">
    <Select
      :model-value="modelValue"
      :options="tabs"
      option-label="label"
      option-value="label"
      :placeholder="
        loading
          ? t('builder.worksheet.loading', 'Loading worksheets…')
          : t('builder.worksheet.pick', 'Select worksheet')
      "
      :disabled="disabled"
      :loading="loading"
      show-clear
      class="w-full"
      @update:model-value="emit('update:modelValue', $event)"
    />
    <p v-if="error" class="text-xs text-red-500 m-0">{{ error }}</p>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { useSheetTabs } from '@/sections/taq/composables/useSheetTabs'

defineProps({
  modelValue: { type: [String, Number], default: null },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])
const { t } = useI18n()

const { tabs, loading, error } = useSheetTabs()
</script>
