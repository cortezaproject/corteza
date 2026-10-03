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

      <!-- Rich text is always sanitized; the opt-out only applies to plain text -->
      <div v-if="!field.options.useRichTextEditor" class="flex flex-col gap-1">
        <div class="flex items-center gap-2">
          <Checkbox v-model="field.options.sanitizeXSS" inputId="sanitizeXSS" :binary="true" />
          <label for="sanitizeXSS" class="cursor-pointer">
            {{ $t('field.kind.string.sanitizeXSS.label') }}
          </label>
        </div>
        <small v-if="!field.options.sanitizeXSS" class="text-orange-600">
          {{ $t('field.kind.string.sanitizeXSS.warning') }}
        </small>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

const field = inject('fieldDraft')
</script>
