<template>
  <div class="my-1">
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
    <span class="mt-1 ml-2 text-nowrap">
      {{ $attrs.label }}
    </span>
  </div>
</template>

<script>
export default {
  name: 'CResourceListStatusFilter',
}
</script>
