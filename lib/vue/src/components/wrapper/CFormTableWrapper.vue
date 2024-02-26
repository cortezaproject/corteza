<template>
  <div class="list-background w-100 p-3">
    <slot></slot>

    <b-button
      v-if="showAddButton"
      variant="primary"
      size="sm"
      :disabled="disableAction"
      @click="handleClick"
    >
      {{ actionLabel }}
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'CFormTableWrapper',

  props: {
    actionLabel: {
      type: String,
      default: '',
    },

    showAddButton: {
      type: Boolean,
      default: false,
    },

    disableAction: {
      type: Boolean,
      default: false,
    },
  },

  methods: {
    handleClick() {
      this.$emit('add-item');
    },
  },
};
</script>

<style scoped>
.list-background {
  background-color: var(--body-bg);
}
</style>
