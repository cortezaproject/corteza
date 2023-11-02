<template>
  <b-col class="d-flex align-items-center flex-wrap gap-1">
    <b-form-radio-group
      v-model="$attrs.value"
      v-bind="$attrs"
      :options="[
        { value: 0, text: $attrs['excluded-label'] || 'excluded' },
        { value: 1, text: $attrs['inclusive-label'] || 'inclusive' },
        { value: 2, text: $attrs['exclusive-label'] || 'exclusive' }
      ]"
      buttons
      button-variant="outline-primary"
      size="sm"
      name="radio-btn-outline"
      v-on="$listeners"
    />
    <span class="text-nowrap">
      {{ $attrs.label }}
    </span>
  </b-col>
</template>

<script>
export default {
  name: 'CResourceListStatusFilter',
}
</script>
