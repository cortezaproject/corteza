<template>
  <div
    v-if="items.length > 0"
  >
    <!-- Group Header -->
    <div
      v-b-toggle:[collapseId]
      class="group-header d-flex align-items-center justify-content-between text-muted shadow-sm"
      :class="{ 'bg-white border-left border-bottom p-2 pr-3': subgroup, 'bg-light p-2 pr-3': !subgroup }"
      :style="{ top: subgroup ? '36px' : '0', zIndex: subgroup ? 10 : 20 }"
    >
      <div class="d-flex align-items-center gap-1">
        <b-badge
          :variant="subgroup ? 'extra-light' : 'primary'"
          style="font-size: 90%;"
        >
          {{ title }}
        </b-badge>
        <span
          v-if="labels.numberOfResults"
          class="small text-muted"
        >
          {{ labels.numberOfResults(items.length) }}
        </span>
      </div>
      <font-awesome-icon
        :icon="['fas', 'chevron-right']"
        class="chevron-icon ml-2 small"
        :class="{ 'chevron-collapsed': !internalExpanded }"
      />
    </div>
    <!-- Group Content -->
    <b-collapse
      :id="collapseId"
      v-model="internalExpanded"
      visible
      class="border-bottom"
    >
      <slot />
    </b-collapse>
  </div>
</template>

<script>
import { library } from '@fortawesome/fontawesome-svg-core'
import { faChevronRight } from '@fortawesome/free-solid-svg-icons'

library.add(faChevronRight)

export default {
  name: 'ItemGroup',

  props: {
    title: {
      type: String,
      required: true,
    },
    items: {
      type: Array,
      required: true,
    },
    collapseId: {
      type: String,
      required: true,
    },
    expanded: {
      type: Boolean,
      default: true,
    },
    subgroup: {
      type: Boolean,
      default: false,
    },
    labels: {
      type: Object,
      default: () => ({}),
    },
  },

  data () {
    return {
      internalExpanded: this.expanded,
    }
  },

  watch: {
    expanded (val) {
      this.internalExpanded = val
    },
    internalExpanded (val) {
      this.$emit('update:expanded', val)
    },
  },
}
</script>

<style lang="scss" scoped>
.group-header {
  position: sticky;
  cursor: pointer;
  transition: all 0.2s ease;

  .chevron-icon {
    transition: transform 0.2s ease;
    transform: rotate(90deg);
  }
  
  .chevron-collapsed {
    transform: rotate(0deg);
  }

  &:hover {
    .chevron-icon {
      color: var(--primary) !important;
    }
  }
}
</style>
