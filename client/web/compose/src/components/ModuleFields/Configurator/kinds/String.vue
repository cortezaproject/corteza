<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3">
      <div class="flex items-center gap-2">
        <Checkbox v-model="field.options.multiLine" inputId="multiLine" :binary="true" />
        <label for="multiLine" class="cursor-pointer">
          {{ $t('field.kind.string.multiLine') }}
        </label>
      </div>

      <!-- Rich Text Editor option is commonly tied to multiLine or stands alone -->
      <div class="flex items-center gap-2">
        <Checkbox
          v-model="field.options.useRichTextEditor"
          inputId="useRichTextEditor"
          :binary="true"
        />
        <label for="useRichTextEditor" class="cursor-pointer">
          {{ $t('field.kind.string.richText') }}
        </label>
      </div>
    </div>

    <CConfiguratorMultiDelimiter :field="field" />
  </div>
</template>

<script setup>
import CConfiguratorMultiDelimiter from '../CConfiguratorMultiDelimiter.vue'

defineProps({
  field: {
    type: Object,
    required: true,
  },
})
</script>
