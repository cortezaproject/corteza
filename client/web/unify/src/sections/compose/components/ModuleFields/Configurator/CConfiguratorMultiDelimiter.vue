<template>
  <div v-if="field.isMulti" class="flex flex-col gap-2">
    <label class="font-medium text-muted-color text-sm">
      {{ $t('field.options.multiDelimiter.label') }}
    </label>
    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-2">
        <RadioButton
          inputId="delimComma"
          :model-value="delimiterType"
          value="comma"
          @update:model-value="setDelimiter(', ')"
        />
        <label for="delimComma" class="cursor-pointer">
          {{ $t('field.options.multiDelimiter.comma') }}
        </label>
      </div>
      <div class="flex items-center gap-2">
        <RadioButton
          inputId="delimNewline"
          :model-value="delimiterType"
          value="newline"
          @update:model-value="setDelimiter('\n')"
        />
        <label for="delimNewline" class="cursor-pointer">
          {{ $t('field.options.multiDelimiter.newline') }}
        </label>
      </div>
      <div class="flex items-center gap-2">
        <RadioButton
          inputId="delimCustom"
          :model-value="delimiterType"
          value="custom"
          @update:model-value="setDelimiter(customDelimiter || ' ')"
        />
        <label for="delimCustom" class="cursor-pointer">
          {{ $t('field.options.multiDelimiter.custom') }}
        </label>
      </div>
    </div>
    <InputText
      v-if="delimiterType === 'custom'"
      v-model="customDelimiter"
      :placeholder="$t('field.options.multiDelimiter.customPlaceholder')"
      class="w-full md:w-1/2"
      size="small"
      @update:model-value="field.options.multiDelimiter = $event"
    />
  </div>
</template>

<script setup>
import { computed, inject, ref, onMounted } from 'vue'

const field = inject('fieldDraft')

const customDelimiter = ref('')

const delimiterType = computed(() => {
  const d = field.value.options?.multiDelimiter
  if (!d || d === ', ') return 'comma'
  if (d === '\n') return 'newline'
  return 'custom'
})

function setDelimiter(val) {
  if (!field.value.options) field.value.options = {}
  field.value.options.multiDelimiter = val
  if (delimiterType.value === 'custom') {
    customDelimiter.value = val
  }
}

onMounted(() => {
  const d = field.value.options?.multiDelimiter
  if (d && d !== ', ' && d !== '\n') {
    customDelimiter.value = d
  }
})
</script>
