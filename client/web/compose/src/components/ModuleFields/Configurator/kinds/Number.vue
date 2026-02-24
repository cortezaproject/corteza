<template>
  <div class="flex flex-col gap-4">
    <!-- Format Group -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-surface-500 text-sm">
        {{ $t('field.kind.number.format') }}
      </label>
      <InputText
        v-model="field.options.format"
        :placeholder="$t('field.kind.number.formatPlaceholder')"
        class="w-full"
      />
      <small class="text-muted-color">{{ $t('field.kind.number.formatDescription') }}</small>
    </div>

    <!-- Prefix / Suffix Display Group -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-2">
        <label class="font-medium text-surface-500 text-sm">
          {{ $t('field.kind.number.prefix') }}
        </label>
        <InputText v-model="field.options.prefix" placeholder="$" class="w-full" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="font-medium text-surface-500 text-sm">
          {{ $t('field.kind.number.suffix') }}
        </label>
        <InputText v-model="field.options.suffix" placeholder="USD" class="w-full" />
      </div>
    </div>

    <!-- Precision -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-surface-500 text-sm">
        {{ $t('field.kind.number.precision') }}
      </label>
      <InputNumber
        v-model="field.options.precision"
        :min="0"
        :max="10"
        show-buttons
        class="w-full md:w-1/2"
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
