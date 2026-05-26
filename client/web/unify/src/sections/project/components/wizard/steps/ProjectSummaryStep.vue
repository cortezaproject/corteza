<template>
  <div class="max-w-3xl mx-auto">
    <GovernanceForm
      :schema="SUMMARY_SCHEMA"
      :model-value="modelValue"
      :disabled="disabled"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </div>
</template>

<script setup>
import GovernanceForm from '@/sections/project/components/wizard/GovernanceForm.vue'
import { SUMMARY_SCHEMA } from '@/sections/project/config/summaryForm'

defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})
defineEmits(['update:modelValue'])
</script>
