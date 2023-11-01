<template>
  <div
    v-if="tooltip"
    v-b-tooltip.hover.right="{ variant: tooltipVariant, title: tooltip, container: '#body' }"
    class="d-flex align-items-center ml-1"
  >
    <font-awesome-icon
      :icon="icon"
      :class="iconClass"
    />
  </div>
</template>

<script>

export default {
  props: {
    icon: {
      type: Array,
      default: () => ['far', 'question-circle'],
    },

    iconClass: {
      type: String,
      default: 'text-primary',
    },

    tooltip: {
      type: String,
      default: '',
    },

    tooltipVariant: {
      type: String,
      default: 'dark',
    },
  },
}
</script>
