<template>
  <small class="text-muted-color">
    <template v-if="isRecordPage">
      {{ $t('block.content.interpolationFootnote') }}
      <code>${record.values.fieldName}</code>
      ,
      <code>${recordID}</code>
      ,
      <code>${ownerID}</code>
      ,
      <code>${userID}</code>
      ,
      <code>${user.name}</code>
    </template>
    <template v-else>
      {{ $t('block.content.interpolationFootnoteNonRecord') }}
      <code>${userID}</code>
      ,
      <code>${user.name}</code>
    </template>
    <template v-if="dependsOnPlacement">
      {{ $t('block.content.interpolationFootnotePlacement') }}
    </template>
  </small>
</template>

<script setup>
defineProps({
  isRecordPage: { type: Boolean, default: false },
  // Set where the configured value is authored away from the page it renders on
  // (charts are namespace-level and can be placed on either kind of page).
  dependsOnPlacement: { type: Boolean, default: false },
})
</script>
