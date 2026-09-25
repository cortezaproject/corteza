<template>
  <CFormGroup :label="t('general.resourceList.status.label')" input-id="resource-status-filter">
    <Select
      input-id="resource-status-filter"
      :model-value="status"
      :options="options"
      option-label="label"
      option-value="value"
      size="small"
      class="w-full"
      @update:model-value="$emit('update:filter', statusFilter($event, states))"
    />
  </CFormGroup>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import CFormGroup from '../input/CFormGroup.vue'
import { statusFilter, statusOf } from '../../composables/useResourceStatus'

const props = defineProps({
  // The list's filter state; only the `states` keys are read
  filter: {
    type: Object,
    required: true,
  },
  // The lifecycle states this resource has, in menu order: deleted, suspended, archived, disabled
  states: {
    type: Array,
    required: true,
  },
})

defineEmits(['update:filter'])

const { t } = useI18n()

const status = computed(() => statusOf(props.filter, props.states))

const options = computed(() =>
  ['active', ...props.states].map(value => ({
    value,
    label: t(`general.resourceList.state.${value}`),
  })),
)
</script>
