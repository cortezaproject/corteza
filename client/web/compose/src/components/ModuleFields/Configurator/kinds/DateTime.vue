<template>
  <div class="flex flex-col gap-6">
    <!-- Format -->
    <div class="flex flex-col gap-2">
      <label class="font-medium text-surface-500 text-sm">
        {{ $t('field.kind.dateTime.outputFormat') }}
      </label>
      <InputText v-model="field.options.format" placeholder="YYYY-MM-DD HH:mm" class="w-full" />
      <small class="text-muted-color">
        {{ $t('field.kind.dateTime.outputFormatFootnote') }}
      </small>
    </div>

    <!-- Past / Future Only validations -->
    <div class="flex flex-col gap-3">
      <h4 class="font-semibold text-md m-0">{{ $t('field.kind.dateTime.constraints.label') }}</h4>
      <div class="flex items-center gap-2">
        <Checkbox
          v-model="field.options.onlyPastValues"
          inputId="pastOnly"
          :binary="true"
          :disabled="field.options.onlyFutureValues"
        />
        <label for="pastOnly" class="cursor-pointer">
          {{ $t('field.kind.dateTime.constraints.options.pastValuesOnly') }}
        </label>
      </div>

      <div class="flex items-center gap-2">
        <Checkbox
          v-model="field.options.onlyFutureValues"
          inputId="futureOnly"
          :binary="true"
          :disabled="field.options.onlyPastValues"
        />
        <label for="futureOnly" class="cursor-pointer">
          {{ $t('field.kind.dateTime.constraints.options.futureValuesOnly') }}
        </label>
      </div>
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
