<template>
  <div class="flex flex-col gap-4">
    <!-- Display type -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-muted-color text-sm">
        {{ $t('field.kind.number.displayType.label') }}
      </label>
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <RadioButton inputId="displayNumber" v-model="field.options.display" value="number" @update:model-value="onDisplayChange" />
          <label for="displayNumber" class="cursor-pointer">{{ $t('field.kind.number.displayType.number') }}</label>
        </div>
        <div class="flex items-center gap-2">
          <RadioButton inputId="displayProgress" v-model="field.options.display" value="progress" @update:model-value="onDisplayChange" />
          <label for="displayProgress" class="cursor-pointer">{{ $t('field.kind.number.displayType.progress') }}</label>
        </div>
      </div>
    </div>

    <!-- Number display options -->
    <template v-if="!isProgress">
      <!-- Format Group -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.formatLabel') }}
        </label>
        <InputText
          v-model="field.options.format"
          :placeholder="$t('field.kind.number.formatPlaceholder')"
          class="w-full"
        />
        <small class="text-muted-color">{{ $t('field.kind.number.examplesLabel') }}</small>
      </div>

      <!-- Prefix / Suffix Display Group -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.prefixLabel') }}
          </label>
          <InputText v-model="field.options.prefix" placeholder="$" class="w-full" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.suffixLabel') }}
          </label>
          <InputText v-model="field.options.suffix" placeholder="USD" class="w-full" />
        </div>
      </div>

      <!-- Precision -->
      <div class="flex flex-col gap-2">
        <label class="font-medium text-muted-color text-sm">
          {{ $t('field.kind.number.precisionLabel') }}
        </label>
        <InputNumber
          v-model="field.options.precision"
          :min="0"
          :max="10"
          show-buttons
          class="w-full md:w-1/2"
        />
      </div>
    </template>

    <!-- Progress bar options -->
    <template v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.progress.minimumValue') }}
          </label>
          <InputNumber v-model="field.options.min" show-buttons class="w-full" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="font-medium text-muted-color text-sm">
            {{ $t('field.kind.number.progress.maximumValue') }}
          </label>
          <InputNumber v-model="field.options.max" show-buttons class="w-full" />
        </div>
      </div>
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
          <Checkbox v-model="field.options.showValue" inputId="showValue" :binary="true" />
          <label for="showValue" class="cursor-pointer">{{ $t('field.kind.number.progress.show.value') }}</label>
        </div>
        <div class="flex items-center gap-2">
          <Checkbox v-model="field.options.animated" inputId="animated" :binary="true" />
          <label for="animated" class="cursor-pointer">{{ $t('field.kind.number.progress.animated') }}</label>
        </div>
      </div>
    </template>

    <CConfiguratorMultiDelimiter :field="field" />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import CConfiguratorMultiDelimiter from '../CConfiguratorMultiDelimiter.vue'

const props = defineProps({
  field: {
    type: Object,
    required: true,
  },
})

const isProgress = computed(() => props.field.options?.display === 'progress')

function onDisplayChange(val) {
  props.field.options.display = val
}

onMounted(() => {
  if (!props.field.options.display) {
    props.field.options.display = 'number'
  }
})
</script>
