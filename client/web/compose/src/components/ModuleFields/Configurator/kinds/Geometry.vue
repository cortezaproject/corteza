<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.prefillWithCurrentLocation"
        inputId="prefillLocation"
        :binary="true"
      />
      <label for="prefillLocation" class="cursor-pointer">
        {{ $t('field.kind.geometry.prefillWithCurrentLocation') }}
      </label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.hideCurrentLocationButton"
        inputId="hideLocationBtn"
        :binary="true"
      />
      <label for="hideLocationBtn" class="cursor-pointer">
        {{ $t('field.kind.geometry.hideCurrentLocationButton') }}
      </label>
    </div>
    <div class="flex items-center gap-2">
      <Checkbox
        v-model="field.options.hideGeoSearch"
        inputId="hideGeoSearch"
        :binary="true"
      />
      <label for="hideGeoSearch" class="cursor-pointer">
        {{ $t('field.kind.geometry.hideGeoSearch') }}
      </label>
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
