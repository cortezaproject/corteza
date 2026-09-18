<template>
  <div class="flex items-center gap-2">
    <Button
      v-tooltip.top="{ value: firstText, showDelay: 500 }"
      icon="pi pi-angle-double-left"
      :aria-label="firstText"
      severity="secondary"
      outlined
      size="small"
      :disabled="!hasPrev"
      data-testid="pager-first"
      @click="emit('first')"
    />
    <Button
      icon="pi pi-angle-left"
      :label="prevLabel || t('general.resourceList.pagination.prev')"
      severity="secondary"
      outlined
      size="small"
      :disabled="!hasPrev"
      data-testid="pager-prev"
      @click="emit('prev')"
    />
    <Button
      icon="pi pi-angle-right"
      icon-pos="right"
      :label="nextLabel || t('general.resourceList.pagination.next')"
      severity="secondary"
      outlined
      size="small"
      :disabled="!hasNext"
      data-testid="pager-next"
      @click="emit('next')"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

/**
 * First / Previous / Next buttons for a cursor-paginated list; First is gated
 * by `hasPrev` too. Labels default to the shared pagination strings.
 */
defineProps({
  hasPrev: {
    type: Boolean,
    default: false,
  },
  hasNext: {
    type: Boolean,
    default: false,
  },
  prevLabel: {
    type: String,
    default: '',
  },
  nextLabel: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['first', 'prev', 'next'])

const { t } = useI18n()

const firstText = computed(() => t('general.resourceList.pagination.first'))
</script>
