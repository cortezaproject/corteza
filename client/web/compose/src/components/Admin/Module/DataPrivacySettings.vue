<template>
  <div v-if="resource && connection" class="flex flex-col gap-6">
    <FormField name="privacy.sensitivityLevel" class="flex flex-col gap-2 max-w-lg">
      <label for="sensitivityLevel" class="font-medium text-primary">
        {{ translations.sensitivity.label }}
      </label>
      <CSensitivityLevelPicker
        v-model="resource.config.privacy.sensitivityLevelID"
        :options="sensitivityLevels"
        :placeholder="translations.sensitivity.placeholder"
        :max-level="maxLevel"
        :disabled="processing"
        class="w-full"
      />
      <small class="text-muted-color">
        {{ translations.sensitivity.description }}
      </small>
    </FormField>

    <FormField name="privacy.usageDisclosure" class="flex flex-col gap-2 max-w-lg">
      <label for="usageDisclosure" class="font-medium text-primary">
        {{ translations.usage.label }}
      </label>
      <Textarea
        id="usageDisclosure"
        v-model="resource.config.privacy.usageDisclosure"
        rows="4"
        autoResize
        class="w-full"
      />
    </FormField>
  </div>
</template>

<script setup>
import { components } from '@cortezaproject/corteza-vue-next'
const { CSensitivityLevelPicker } = components

const props = defineProps({
  resource: {
    type: Object,
    required: true,
  },
  connection: {
    type: Object,
    required: true,
  },
  maxLevel: {
    type: String,
    default: undefined,
  },
  sensitivityLevels: {
    type: Array,
    required: true,
  },
  translations: {
    type: Object,
    required: true,
  },
  processing: {
    type: Boolean,
    default: false,
  },
})

// Ensure config.privacy exists safely
if (!props.resource.config) props.resource.config = {}
if (!props.resource.config.privacy) {
  props.resource.config.privacy = { sensitivityLevelID: '0', usageDisclosure: '' }
}
</script>
