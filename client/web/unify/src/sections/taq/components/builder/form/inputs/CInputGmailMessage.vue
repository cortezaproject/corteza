<template>
  <div class="flex flex-col gap-1">
    <Select
      :model-value="modelValue"
      :options="messages"
      option-label="label"
      option-value="value"
      filter
      :placeholder="
        loading
          ? t('builder.gmailMessage.loading', 'Loading emails…')
          : t('builder.gmailMessage.pick', 'Select an email')
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
import { useGmailMessages } from '@/sections/taq/composables/useGmailMessages'

defineProps({
  modelValue: { type: [String, Number], default: null },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])
const { t } = useI18n()

const { messages, loading, error } = useGmailMessages()
</script>
