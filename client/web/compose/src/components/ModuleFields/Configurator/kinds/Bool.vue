<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <!-- True value label -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-surface-500 text-sm">
        {{ $t('field.kind.bool.trueLabel') }}
      </label>
      <InputText
        v-model="field.options.trueLabel"
        :placeholder="$t('field.kind.bool.truePlaceholder', { default: $t('general.label.yes') })"
        class="w-full"
      />
    </div>

    <!-- False value label -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-surface-500 text-sm">
        {{ $t('field.kind.bool.falseLabel') }}
      </label>
      <InputText
        v-model="field.options.falseLabel"
        :placeholder="$t('field.kind.bool.falsePlaceholder', { default: $t('general.label.no') })"
        class="w-full"
      />
    </div>
  </div>
</template>

<script setup>
defineProps({
  field: {
    type: Object,
    required: true,
  },
})
</script>
