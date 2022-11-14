<template>
  <li
    class="d-flex mt-1 mx-1"
  >
    <b-button
      class="flex-grow-1"
      :disabled="!hasPrevPage"
      size="sm"
      @click="$emit('prev')"
    >
      <font-awesome-icon
        :icon="['fas', 'chevron-left']"
      />
    </b-button>
    <b-button
      class="flex-grow-1 ml-1"
      :disabled="!hasNextPage"
      size="sm"
      @click="$emit('next')"
    >
      <font-awesome-icon
        :icon="['fas', 'chevron-right']"
      />
    </b-button>
  </li>
</template>

<script>
// This file is copied over from corteza-webapp-compose; src/components/ModuleFields/Common/Pagination.vue
// Refactor it out when those components get moved into corteza-vue
export default {
  props: {
    hasPrevPage: {
      type: Boolean,
      required: true,
    },

    hasNextPage: {
      type: Boolean,
      required: true,
    },
  },
}
</script>
