<template>
  <div class="flex flex-col gap-6 max-w-3xl">
    <p class="text-sm text-muted-color">{{ $t('project.deployer.prompt') }}</p>

    <div class="flex flex-col gap-3">
      <div
        v-for="(q, i) in DEPLOYER_QUESTIONS"
        :key="q.key"
        class="flex items-start gap-4 rounded-lg border border-surface p-3"
      >
        <p class="text-sm flex-1 min-w-0">
          <span class="font-medium mr-1">{{ i + 1 }}.</span>
          {{ $t(q.labelKey) }}
        </p>
        <SelectButton
          :model-value="modelValue[q.key]"
          :options="YES_NO"
          option-label="label"
          option-value="value"
          :allow-empty="false"
          :disabled="disabled"
          class="shrink-0"
          @update:model-value="v => update(q.key, v)"
        />
      </div>
    </div>

    <!-- Working outcome — session-local read of the answers above, not a
         backend derivation (that's the separate, untouched FriaRequired
         field the same answers will drive once a backend slice lands). -->
    <div
      class="flex items-start gap-3 rounded-lg border p-4"
      :class="required ? 'border-primary/30 bg-primary/5' : 'border-surface bg-emphasis'"
    >
      <i
        class="pi text-lg mt-0.5"
        :class="required ? 'pi-shield text-primary' : 'pi-info-circle text-muted-color'"
      />
      <div class="text-sm">
        <p class="font-medium">
          {{ required ? $t('fria.determination.required') : $t('fria.determination.notRequired') }}
        </p>
        <p class="text-muted-color mt-1">{{ $t('fria.determination.outcomeHint') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  DEPLOYER_QUESTIONS,
  friaRequiredFrom,
} from '@/sections/project/config/friaDeterminationForm'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: { type: Object, default: () => ({}) },
  disabled: { type: Boolean, default: false },
})
const emit = defineEmits(['update:modelValue'])

const YES_NO = computed(() => [
  { label: t('general.label.yes'), value: true },
  { label: t('general.label.no'), value: false },
])

const required = computed(() => friaRequiredFrom(props.modelValue))

function update(key, value) {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>
